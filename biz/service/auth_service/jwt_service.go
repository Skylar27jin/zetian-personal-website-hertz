package auth_service

import (
	"context"
	"fmt"
	"time"
	jwt_pkg "zetian-personal-website-hertz/biz/pkg/JWT"
)


/*
GenerateUserJWT generates a JWT token for a user with given information
now: if now is -1, use current time in the system as the issued at time
username: not null
email: not null
validDuration: if validDuration is -1, 7 * 24 hours is used as the default validDuration


returns a encrypted token string like:
{username: username,
email: email,
iat: now, //iat means issued at time
exp: now + validDuration //exp means expiration time
}

*/
func GenerateUserJWT(ctx context.Context, now int, username string, email string, validDuration int) (string, error) {
	if username == "" || email == "" {
		return "", fmt.Errorf("username and email cannot be empty")
	}
	if now == -1 {
		now = int(time.Now().Unix())
	}
	if validDuration == -1 {
		validDuration = 7 * 24 * 3600 //7 days
	}

	payLoad := map[string]interface{}{
		"username": username,
		"email":    email,
		"iat": 	now,
		"exp":	now + validDuration,
	}

	token, err := jwt_pkg.GenerateJWT(payLoad) 
	if err != nil {
		return "", fmt.Errorf("when generating user JWT: %v", err)
	}

	return token, nil
}

func ParseUserJWT(ctx context.Context, tokenString string) (
	username string,
	email string,
	iat int,
	exp int,
	returnErr error) {

	payload, err := jwt_pkg.ParseJWT(tokenString)
	if err != nil {
		return "", "", -2, -2, fmt.Errorf("when parsing user JWT: %v", err)
	}

	usernameInterface, ok1 := payload["username"]
	emailInterface, ok2 := payload["email"]
	iatInterface, ok3 := payload["iat"]
	expInterface, ok4 := payload["exp"]


	if !ok1 || !ok2 || !ok3 || !ok4 {
		return "", "", -3, -3, fmt.Errorf("invalid JWT: When parsing user JWT: missing fields")
	}

	return usernameInterface.(string), emailInterface.(string), iatInterface.(int), expInterface.(int), nil


}
