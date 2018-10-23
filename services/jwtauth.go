package services

import (
	"github.com/gyaan/meta-mask-login/models"
	"time"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

//use it only when you verified user
func GenerateToken(user models.User) (models.TokenEndExpire, error)  {

	var te models.TokenEndExpire
	var expire = time.Now().Add(time.Minute*60).Unix()

    claims := models.Claims{
    	user.PublicAddress,
    	jwt.StandardClaims{
    		ExpiresAt:expire,
    		Issuer:"localhost",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	t, err := token.SignedString([]byte(models.SecretKey))
	if err!= nil{
		return  te, fmt.Errorf("%q",err)
	}
    te.Token = t
    te.Expire = fmt.Sprintf("%q", time.Unix(expire, 0))
    return  te, nil
}
func ValidateAccessToken()  {

}