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
validDuration: it is in second. if validDuration is -1, 7 * 24 hours is used as the default validDuration


returns a encrypted token string like:
{
id: userID,
username: username,
email: email,
iat: now, //iat means issued at time
exp: now + validDuration //exp means expiration time
}

*/
func GenerateUserJWT(ctx context.Context, now int64,  id int64, username string, email string, validDuration int64) (string, error) {
	if username == "" || email == "" {
		return "", fmt.Errorf("username and email cannot be empty")
	}
	if now == -1 {
		now = time.Now().Unix()
	}
	if validDuration == -1 {
		validDuration = 7 * 24 * 3600 //7 days
	}

	payLoad := map[string]interface{}{
		"username": username,
		"email":    email,
		"iat": 	now,
		"exp":	now + validDuration,
		"id": id,
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
	iat int64,
	exp int64,
	id int64,
	returnErr error) {

	payload, err := jwt_pkg.ParseJWT(tokenString)
	if err != nil {
		return "", "", -2, -2, -2, fmt.Errorf("when parsing user JWT: %v", err)
	}

	usernameInterface, ok1 := payload["username"]
	emailInterface, ok2 := payload["email"]
	iatInterface, ok3 := payload["iat"]
	expInterface, ok4 := payload["exp"]
	idInterface, ok5 := payload["id"]


	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return "", "", -3, -3, -3, fmt.Errorf("invalid JWT: When parsing user JWT: missing fields")
	}

	iat64, err := safelyConvertToInt64(iatInterface)
	if err != nil {
		return "", "", -3, -3, -3, err
	}
	exp64, err := safelyConvertToInt64(expInterface)
	if err != nil {
		return "", "", -3, -3, -3, err
	}
	id_got, err := safelyConvertToInt64(idInterface)
	if err != nil {
		return "", "", -3, -3, -3, err
	}

	return usernameInterface.(string), emailInterface.(string), iat64, exp64, id_got, nil

}

/*
GenerateVeriEmailJWT generates a JWT token for client's cookie if the user verified the email.
This JWT is used for user to execute action related to this email, including: change the password, bind email, unbind email, change password, etc.
now: if now = -1, use current time in the system as the issued at time
email: not null
validDuration: in second, if valudation duration = -1, 15 min is used as the default valid duration

return a JWT:
payload:
{
"email": email,
"exp": now + validDuration in unix
}
*/
func GenerateVeriEmailJWT(ctx context.Context, now int64, email, purpose string, validDuration int64) (string, error) {
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}
	if now == -1 {
		now = time.Now().Unix()
	}
	if validDuration == -1 {
		validDuration = 3 * 60 //default 3min
	}
	payLoad := map[string]interface{}{
		"email": email,
		"exp": now + validDuration,
		"purpose": purpose,
	}


	token, err := jwt_pkg.GenerateJWT(payLoad)
	if err != nil {
		return "", fmt.Errorf("when generating veriEmail JWT: %v", err)
	}

	return token, nil


}

/*
ParseVeriEmailJWT parses the veriEmailJWT from client's cookie
take in JWT like:
{
"email": email,
"exp": now + validDuration in unix
}

returns email and exp

*/
func ParseVeriEmailJWT(ctx context.Context, tokenString string) (string, int64, string, error) {
	payload, err := jwt_pkg.ParseJWT(tokenString)
	if err != nil {
		return "", -1, "", fmt.Errorf("when parsing VeriEmail JWT: %v", err)
	}

	emailInterface, ok1 := payload["email"]
	expInterface, ok2 := payload["exp"]
	purposeInterface, ok3 := payload["purpose"]


	if !ok1 || !ok2 || !ok3 {
		return "", -2, "", fmt.Errorf("invalid JWT: When parsing user JWT: missing fields")
	}

	exp64, err := safelyConvertToInt64(expInterface)
	if err != nil {
		return "", -3, "", err
	}

	return emailInterface.(string), exp64, purposeInterface.(string), nil

}

// safelyConvertToInt64 converts interface{} to int64 safely, supporting both int64 and float64
func safelyConvertToInt64(v interface{}) (int64, error) {
	switch val := v.(type) {
	case float64:
		return int64(val), nil
	case int64:
		return val, nil
	default:
		return 0, fmt.Errorf("unexpected type for numeric field: %T", v)
	}
}