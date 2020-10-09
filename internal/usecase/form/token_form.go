package form

import (
	"strconv"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/library/config"
	"github.com/dgrijalva/jwt-go"
)

type TokenForm struct {
	Token string `name:"token" validate:"required"`
}

func (tf TokenForm) ToToken(tokenSettings config.TokenSettings) (*model.Token, error) {
	claim := new(jwt.StandardClaims)

	parser := jwt.Parser{
		ValidMethods: []string{tokenSettings.Alg},
	}
	if _, err := parser.ParseWithClaims(tf.Token, claim, func(*jwt.Token) (interface{}, error) {
		return []byte(tokenSettings.SecretKey), nil
	}); err != nil {
		return nil, err
	}

	userId, err := strconv.ParseUint(claim.Subject, 10, 64)
	if err != nil {
		return nil, err
	}

	token := new(model.Token)
	token.UserID = uint(userId)
	token.Jti = claim.Id
	return token, nil
}
