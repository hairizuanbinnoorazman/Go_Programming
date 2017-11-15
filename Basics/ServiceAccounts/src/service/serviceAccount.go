package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"encoding/json"
	"io/ioutil"
	"time"
	"fmt"
	"log"
)

type serviceAuth struct  {
	Type string `json:"type"`
	ProjectId string `json:"project_id"`
	PrivateKeyId string `json:"private_key_id"`
	PrivateKey string `json:"private_key"`
	ClientId string `json:"client_id"`
	ClientEmail string`json:"client_email"`
	AuthUri string `json:"auth_uri"`
	TokenUri string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl string `json:"client_x509_cert_url"`
}


type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func main(){
	claima := CustomClaims{
		"https://www.googleapis.com/auth/analytics",
		jwt.StandardClaims{
			Issuer:"delete-this-thing@random-project.iam.gserviceaccount.com",
			Audience:"https://www.googleapis.com/oauth2/v4/token",
			IssuedAt:time.Now().Unix(),
			ExpiresAt:time.Now().Unix()+3600,
		},

	}

	//claima2 := jwt.StandardClaims{
	//	Issuer:"delete-this-thing@wolverine-dev-1234.iam.gserviceaccount.com",
	//	Audience:"https://www.googleapis.com/oauth2/v4/token",
	//	IssuedAt:time.Now().Unix(),
	//	ExpiresAt:time.Now().Unix()+3600,
	//}

	authRaw, _ := ioutil.ReadFile("auth.json")
	//fmt.Println(authRaw)
	var serviceauth serviceAuth
	json.Unmarshal(authRaw, &serviceauth)
	fmt.Println(serviceauth.PrivateKey)

	fmt.Println(claima)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claima)
	fmt.Println(claima)

	 lol, err := json.Marshal(claima)
	 log.Println("Printing the token", string(lol))

	claima.Valid()
	fmt.Println(token.Valid)
	fmt.Println(token)
	fmt.Println(token.Signature)

	nyaahaha := []byte(serviceauth.PrivateKey)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(nyaahaha)
	if err != nil{
		log.Println(err)
		log.Println(err.Error())
		log.Fatal("Kill me please")
	}

	ss, err := token.SignedString(key)
	fmt.Println(time.Date(2017, 11, 10, 12, 0, 0, 0, time.UTC).Unix())
	if err != nil {
		log.Println(err)
		log.Println(err.Error())
		log.Fatal("Stop the world")
	}
	log.Println(ss)
}