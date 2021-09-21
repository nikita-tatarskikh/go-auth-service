package main

import (
	"fmt"
	"context"
	"log"
	"net/http"
	"github.com/google/uuid"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	//"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type tokensPair struct {
	AccessTokenString string `json:"accessToken"`
	RefreshTokenString string `json:"refreshToken"`
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

	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userId
	rtClaims["refersh_uuid"] = uuid.New().String()
	refreshTokenValue := jwt.NewWithClaims(jwt.SigningMethodHS512, rtClaims)
	tokensPair.RefreshTokenString, err = refreshTokenValue.SignedString(jwtKey)
	
	if err != nil {
		log.Fatal(err)
	}

	//TO-DO сделать сохранение refreshToken в монго

	return tokensPair, err
}

func SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	userId := r.Form.Get("guid")
	payload, err := generateTokensPair(userId)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func Refresh (w http.ResponseWriter, r *http.Request) {
	//implementaion
	// 1) Проверяем валидность токена, если токен валиден, переходм к шагу 2.
	
}


func StoreToken (tokensPair tokensPair ) {
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
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("RefreshTokens").Collection("RefreshTokens")
	
	insertRefreshToken, err := collection.InsertOne(context.TODO(), tokensPair.RefreshTokenString)

	if err != nil {
		log.Fatal(err)
	}
	//tokensPair.RefreshTokenString
	
}

func main() {
	// router := httprouter.New()
	// router.POST("/sign-up/", SignUp)
	// log.Fatal(http.ListenAndServe(":8080", router))
	StoreToken(tokensPair{})
}
