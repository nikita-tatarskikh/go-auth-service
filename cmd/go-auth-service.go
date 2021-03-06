package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/ulule/deepcopier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type tokensPair struct {
	AccessTokenString string `json:"accessToken"`
	RefreshTokenString string `bson:"refreshToken" json:"refreshToken"`
}

// RefreshToken Создана отдельная структура для простоты insert в MongoDB и обратки refresh роута.
type RefreshToken struct {
	UserID string `bson:"userID" json:"userID"`
	RefreshTokenString string `bson:"refreshToken" json:"refreshToken"`
}

func generateTokensPair(userID string) (*tokensPair, error) {
	var err error
	var jwtKey = []byte("https://www.notion.so/Test-task-Junior-BackDev-215fcddafff2425a8ca7e515e21527e7")
	tokensPair := &tokensPair{}

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userID
	atClaims["refresh_uuid"] = uuid.New().String()
	accessTokenValue := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	tokensPair.AccessTokenString, err = accessTokenValue.SignedString(jwtKey)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Access Token was generated")

	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userID
	rtClaims["refresh_uuid"] = uuid.New().String()
	refreshTokenValue := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	tokensPair.RefreshTokenString, err = refreshTokenValue.SignedString(jwtKey)
	
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Refresh Token was generated")

	StoreRefreshToken(*tokensPair, userID)

	return tokensPair, err
}


func SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Started handling SignUp")
	err := r.ParseForm()
	if err != nil {
		return 
	}
	userID := r.Form.Get("guid")
	payload, err := generateTokensPair(userID)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		return 
	}
	log.Println("User request was handled")
}

func Refresh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Started handling refresh token")
	var RefreshToken, mongoSearchResult RefreshToken

	err := json.NewDecoder(r.Body).Decode(&RefreshToken)
	if err != nil {
		log.Fatal(err)
	}

	if RefreshToken.RefreshTokenString == "" {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode("Key 'refreshToken' not found")
		if err != nil {
			return 
		}
	} else {
		token, err := jwt.Parse(RefreshToken.RefreshTokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte("https://www.notion.so/Test-task-Junior-BackDev-215fcddafff2425a8ca7e515e21527e7"), nil
		})
	
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode("Illegal token")
			if err != nil {
				return 
			}
			log.Println("test0", err)
		} else {
			claims := token.Claims.(jwt.MapClaims)
		RefreshToken.UserID = claims["user_id"].(string)
		clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
		client, err := mongo.Connect(context.TODO(), clientOptions)
	
		if err != nil {
			log.Fatal(err)
		}
	
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("in", RefreshToken.RefreshTokenString)
	
		log.Println("Connected to MongoDB")
	
		collection := client.Database("RefreshTokens").Collection("RefreshTokens")
		err = collection.FindOne(context.TODO(), bson.M{"userID": RefreshToken.UserID}).Decode(&mongoSearchResult)
		log.Println("mongo search result", mongoSearchResult)
		if err!= nil {
			log.Println("test",err)
		}
	
		err = bcrypt.CompareHashAndPassword([]byte(mongoSearchResult.RefreshTokenString), []byte(RefreshToken.RefreshTokenString) )
		log.Println(mongoSearchResult.RefreshTokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode("Refresh Token is incorrect")
			if err != nil {
				return 
			}
			log.Println("bcrypt error", err)
		} else {
			if err!= nil {
		 		log.Println("An error occurred while processing the request")
		 	}
			payload, err := generateTokensPair(RefreshToken.UserID)
				if err != nil {
				log.Printf("An error occurred while processing the request")	
				}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(payload)
			if err != nil {
				return 
			}
			log.Println("Tokens Pair was updated")
		}
	}	
		}
	
		
	
}

func StoreRefreshToken(tokensPair tokensPair, userID string) {
	var mongoSearchResult RefreshToken
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
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
	err = deepcopier.Copy(tokensPair).To(RefreshTokenDoc)
	if err != nil {
		return 
	}

	RefreshTokenDoc.UserID = userID

	err = collection.FindOne(context.TODO(), bson.M{"userID": RefreshTokenDoc.UserID}).Decode(&mongoSearchResult)

	if err != nil {
		log.Println(err)
	}

	if mongoSearchResult.RefreshTokenString != "" && mongoSearchResult.UserID != "" {
		bytes, err := bcrypt.GenerateFromPassword([]byte(RefreshTokenDoc.RefreshTokenString), 4)
		if err != nil {
			log.Fatal(err)
		}
		RefreshTokenDoc.RefreshTokenString = string(bytes)
		log.Println("Hashed Token while saving", RefreshTokenDoc)
		UpdateRefreshToken, err := collection.ReplaceOne(context.TODO(),bson.M{"userID": RefreshTokenDoc.UserID}, RefreshTokenDoc)
		log.Println("Hashed Token was updated in MongoDB", UpdateRefreshToken)
			if err != nil {
				log.Fatal(err)
			}


	} else {
		bytes, err := bcrypt.GenerateFromPassword([]byte(RefreshTokenDoc.RefreshTokenString), 4)
		if err != nil {
			log.Fatal(err)
		}
		RefreshTokenDoc.RefreshTokenString = string(bytes)
		log.Println("Hashed Token while saving", RefreshTokenDoc)
		InsertRefreshToken, err := collection.InsertOne(context.TODO(), RefreshTokenDoc)
		log.Println("Hashed Token was stored in MongoDB", InsertRefreshToken.InsertedID)
			if err != nil {
				log.Fatal(err)
			}
	}
	
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

