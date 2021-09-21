package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"


	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ulule/deepcopier"

)

type tokensPair struct {
	AccessTokenString string `json:"accessToken"`
	RefreshTokenString string `bson:"refreshToken" json:"refreshToken"`
}

//Создана отдельная структура для простоты insert в MongoDB
type RefreshTokenDoc struct {
	RefreshTokenString string `bson:"refreshToken"`
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

	StoreToken(*tokensPair)

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

func Refresh(w http.ResponseWriter, r *http.Request) {
	//implementaion
	// 1) Проверяем валидность токена, если токен валиден, переходм к шагу 2.
	
}


func StoreToken(tokensPair tokensPair ) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
	    log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
	    log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	collection := client.Database("RefreshTokens").Collection("RefreshTokens")

	RefreshTokenDoc := &RefreshTokenDoc{}

	deepcopier.Copy(tokensPair).To(RefreshTokenDoc)

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
	log.Fatal(http.ListenAndServe(":8080", router))
}
