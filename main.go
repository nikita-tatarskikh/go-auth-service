package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func generateTokensPair(userId string) string {
	var jwtKey = []byte("https://www.notion.so/Test-task-Junior-BackDev-215fcddafff2425a8ca7e515e21527e7")
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return err.Error()
	}
	return tokenString
}

func SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	parameter := ps.ByName("guid")
	w.Write([]byte(generateTokensPair(parameter)))
	//implemetation
}

func SignUp2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	fmt.Print(r.Form)
	//w.Write([]byte(generateTokensPair(r.Form["dgdf"])))
	//implemetation
}

func main() {
	router := httprouter.New()
	router.POST("/sign-up/userId=:guid", SignUp)
	router.POST("/sign-up/", SignUp2)
	log.Fatal(http.ListenAndServe(":8080", router))
}
