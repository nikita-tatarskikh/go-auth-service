package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ulule/deepcopier"
)

type tokensPair struct {
	AccessTokenString string `json:"accessToken"`
	RefreshTokenString string `bson:"refreshToken" json:"refreshToken"`
}

//Создана отдельная структура для простоты insert в MongoDB и обратки refresh роута.
type RefreshToken struct {
	UserId string `bson:"userId" json:"userId"`
	RefreshTokenString string `bson:"refreshToken" json:"refreshToken"`
}

func generateTokensPair(userId string) (*tokensPair, error) {
	var err error
	var jwtKey = []byte("https://www.notion.so/Test-task-Junior-BackDev-215fcddafff2425a8ca7e515e21527e7")
	tokensPair := &tokensPair{}

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userId
	accessTokenValue := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	tokensPair.AccessTokenString, err = accessTokenValue.SignedString(jwtKey)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Access Token was generated")

	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userId
	rtClaims["refersh_uuid"] = uuid.New().String()
	refreshTokenValue := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	tokensPair.RefreshTokenString, err = refreshTokenValue.SignedString(jwtKey)
	
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Refresh Token was generated")

	StoreToken(*tokensPair, userId)

	return tokensPair, err
}


func SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Started handling SignUp")
	r.ParseForm()
	userId := r.Form.Get("guid")
	payload, err := generateTokensPair(userId)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
	log.Println("User request was handled")
}

func Refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var RefreshToken, mongoSearchResult RefreshToken

	err := json.NewDecoder(r.Body).Decode(&RefreshToken)
	if err != nil {
		log.Fatal(err)
	}

	if (RefreshToken.RefreshTokenString == "") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Key 'refreshToken' not found")
	} else {
		token, err := jwt.Parse(RefreshToken.RefreshTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte("https://www.notion.so/Test-task-Junior-BackDev-215fcddafff2425a8ca7e515e21527e7"), nil
		})
	
		if err != nil {
			log.Println("test0", err)
		}
	
		claims := token.Claims.(jwt.MapClaims);
		RefreshToken.UserId = claims["user_id"].(string)
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.Connect(context.TODO(), clientOptions)
	
		if err != nil {
			log.Fatal(err)
		}
	
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
	
		log.Println("Connected to MongoDB")
	
		collection := client.Database("RefreshTokens").Collection("RefreshTokens")
		err = collection.FindOne(context.TODO(), bson.M{"userId": RefreshToken.UserId}).Decode(&mongoSearchResult)
		log.Println(mongoSearchResult)
		if err!= nil {
			log.Println("test",err)
		}

	
		err = bcrypt.CompareHashAndPassword([]byte(mongoSearchResult.RefreshTokenString), []byte(RefreshToken.RefreshTokenString) )
		log.Println(mongoSearchResult.RefreshTokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Refresh Token is incorrect")
			log.Println("bcrypt error", err)
		} else {
			payload, err := generateTokensPair(RefreshToken.UserId)
				if err != nil {
				log.Printf("An error occurred while processing the request")	
				}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(payload)
			log.Println("Tokens Pair was updated")
		}
	}

	

	

	
	
}


func StoreToken(tokensPair tokensPair, userId string) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
	    log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
	    log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	collection := client.Database("RefreshTokens").Collection("RefreshTokens")

	RefreshTokenDoc := &RefreshToken{}
	deepcopier.Copy(tokensPair).To(RefreshTokenDoc)

	bytes, err := bcrypt.GenerateFromPassword([]byte(RefreshTokenDoc.RefreshTokenString), 4)

	if err != nil {
		log.Fatal(err)
	}
	RefreshTokenDoc.RefreshTokenString = string(bytes)
	RefreshTokenDoc.UserId = userId
	log.Println("Refres Token while saving", RefreshTokenDoc)
	InsertRefrshToken, err := collection.InsertOne(context.TODO(), RefreshTokenDoc)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Refresh Token was stored in MongoDB", InsertRefrshToken.InsertedID)
	
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Disconnected from MongoDB")
}

func main() {
	router := httprouter.New()
	router.POST("/sign-up/", SignUp)
	router.POST("/refresh/", Refresh)
	log.Fatal(http.ListenAndServe(":8080", router))
}

