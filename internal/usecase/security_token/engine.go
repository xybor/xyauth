package securitytoken

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/xybor/xyauth/internal/logger"
)

type Token interface {
	Unmarshal(payload any) error
}

type Engine struct {
	Issuer     string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// create creates a string that represents for the token with the given
// parameters.
func (e Engine) create(payload any, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": payload,
		"iss": e.Issuer,
		"exp": time.Now().Add(expiration).Unix(),
		"iat": time.Now().Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(e.PrivateKey)
	if err != nil {
		logger.Event("create-token-error").
			Field("claims", claims).Field("error", err).Warning()
		return "", SecurityTokenError.New("can not create token")
	}

	return token, nil
}

// Verify checks if s is a valid token, then parse it to Token t. Parameter t
// must be a pointer to Token struct.
func (e Engine) Verify(s string, t Token) error {
	parsedToken, err := jwt.ParseWithClaims(s, jwt.MapClaims{}, publicKeyFunc)
	if err != nil {
		if !errors.Is(err, SecurityTokenError) {
			logger.Event("parse-token-error").Field("error", err).Debug()
			err = SecurityTokenError.New("can not parse the token")
		}
		return err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return SecurityTokenError.New("invalid token")
	}

	if !claims.VerifyIssuer(issuer, true) {
		return SecurityTokenError.New("token is issued by an unknown issuer")
	}

	payload, ok := claims["sub"]
	if !ok {
		return SecurityTokenError.New("no subject provided")
	}

	if err := t.Unmarshal(payload); err != nil {
		return err
	}

	return nil
}
