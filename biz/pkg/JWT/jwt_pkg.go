package jwt_pkg

import (
	"errors"
	"fmt"
	"zetian-personal-website-hertz/biz/config"

	"github.com/golang-jwt/jwt/v5"
)

// ParseJWT parses and validates JWT (as a string, not []byte), returns payload claims
func ParseJWT(tokenString string) (map[string]interface{}, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // 确保签名方法是 HMAC
        //ensure the sigining method is HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.GetGeneralConfig().JWT_Secret_Key), nil
    })
    if err != nil {
        return nil, err
    }

    // token.Valid 表示签名和过期时间都检查通过
    // token.Valid indicates that the signature and expiration time are valid
    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    // 从 token 中提取 payload（claims）
    //get payLoad from token
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        payload := make(map[string]interface{})
        for k, v := range claims {
            payload[k] = v
        }
        return payload, nil
    }

    return nil, errors.New("failed to parse claims")
}

/*
GenerateJWT generates a JWT token with given payLoad
returned token is a string, not []byte
*/
func GenerateJWT(payLoad map[string]interface{}) (string, error) {

    claims := jwt.MapClaims{}
    for key, value := range payLoad {
        claims[key] = value
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.GetGeneralConfig().JWT_Secret_Key))
}
