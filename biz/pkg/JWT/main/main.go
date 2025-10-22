//this is a simple example of how to use the JWT package
//generate a JWT token
package main

import (
	"fmt"
	JWT "zetian-personal-website-hertz/biz/pkg/JWT"
)
func main() {
	name := "Nina"
	age := 21
	println(name, age)
	payLoad := map[string]interface{}{
		"username": name,
		"age":      age,
	}
	res, err := JWT.GenerateJWT(payLoad)
	
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated JWT:", res)

	decodedCookie, err := JWT.ParseJWT(res)
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed JWT Claims:", decodedCookie)
}
//to get the JWT token from cookie:
// 	rawJWT := string(c.Cookie("JWT"))