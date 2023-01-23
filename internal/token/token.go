package token

import (
	"crypto/rsa"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/internal/config"
	"github.com/xybor/xyauth/internal/logger"
)

var issuer = config.GetDefault("oauth.issuer", "xyauth").MustString()

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey
var AccessTokenExpiration = config.GetDefault(
	"oauth.access_token_expiration", 5*time.Minute).MustDuration()
var RefreshTokenExpiration = config.GetDefault(
	"oauth.refresh_token_expiration", 5*time.Minute).MustDuration()

func publicKeyFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, AlgorithmError.Newf("unexpected algorithm %s", t.Header["alg"])
	}
	return publicKey, nil
}

func init() {
	// TODO: ReplaceAll commands will be removed if the PR 156 of godotenv is
	// merged.
	var err error
	key := config.MustGet("OAUTH_PRIVATE_KEY").MustString()
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(strings.ReplaceAll(key, `\n`, "\n")))
	xycond.AssertNil(err)

	key = config.MustGet("OAUTH_PUBLIC_KEY").MustString()
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(strings.ReplaceAll(key, `\n`, "\n")))
	xycond.AssertNil(err)
}

// Create creates a string that represents for the token with the given
// configuration.
func Create(c Config) (string, error) {
	claims := jwt.MapClaims{
		"sub": c.payload,
		"iss": issuer,
		"exp": time.Now().Add(c.expiration).Unix(),
		"iat": time.Now().Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		logger.Event("create-token-error").
			Field("claims", claims).Field("error", err).Warning()
		return "", TokenError.New("can not create token")
	}

	return token, nil
}

type Token interface {
	Unmarshal(payload any) error
}

// Verify checks if s is a valid token, then parse it to Token t. Parameter t
// must be a pointer to Token struct.
func Verify(s string, t Token) error {
	parsedToken, err := jwt.ParseWithClaims(s, jwt.MapClaims{}, publicKeyFunc)
	if err != nil {
		if !errors.Is(err, TokenError) {
			logger.Event("parse-token-error").Field("error", err).Debug()
			err = TokenError.New("can not parse the token")
		}
		return err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return TokenError.New("invalid token")
	}

	if !claims.VerifyIssuer(issuer, true) {
		return TokenError.New("token is issued by an unknown issuer")
	}

	payload, ok := claims["sub"]
	if !ok {
		return TokenError.New("no subject provided")
	}

	if err := t.Unmarshal(payload); err != nil {
		return err
	}

	return nil
}
