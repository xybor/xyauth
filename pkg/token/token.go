package token

import (
	"crypto/rsa"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"github.com/xybor-x/xycond"
	"github.com/xybor/xyauth/internal/models"
	"github.com/xybor/xyauth/internal/utils"
)

var issuer = utils.GetConfig().GetDefault("oauth.issuer", "xyauth").MustString()

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey
var AccessTokenExpiration = utils.GetConfig().GetDefault(
	"oauth.access_token_expiration", 5*time.Minute).MustDuration()
var RefreshTokenExpiration = utils.GetConfig().GetDefault(
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
	key := utils.GetConfig().MustGet("OAUTH_PRIVATE_KEY").MustString()
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(strings.ReplaceAll(key, `\n`, "\n")))
	xycond.AssertNil(err)

	key = utils.GetConfig().MustGet("OAUTH_PUBLIC_KEY").MustString()
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(strings.ReplaceAll(key, `\n`, "\n")))
	xycond.AssertNil(err)
}

type Token interface {
	Unmarshal(payload any) error
}

type AccessToken struct {
	models.User
}

func (t *AccessToken) Unmarshal(payload any) error {
	if err := mapstructure.Decode(payload, &t.User); err != nil {
		return TokenError.New("can not parse to access token")
	}
	return nil
}

type RefreshToken struct {
	Email string
}

func (t *RefreshToken) Unmarshal(payload any) error {
	if err := mapstructure.Decode(payload, t); err != nil {
		return TokenError.New("can not parse to refresh token")
	}
	return nil
}

func Create(c Config) (string, error) {
	claims := jwt.MapClaims{
		"sub": c.payload,
		"iss": issuer,
		"exp": time.Now().Add(c.expiration).Unix(),
		"iat": time.Now().Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		utils.GetLogger().Event("create-token-error").
			Field("claims", claims).Field("error", err).Warning()
		return "", TokenError.New("can not create token")
	}

	return token, nil
}

func Verify(s string, t Token) error {
	parsedToken, err := jwt.ParseWithClaims(s, jwt.MapClaims{}, publicKeyFunc)
	if err != nil {
		if !errors.Is(err, TokenError) {
			utils.GetLogger().Event("parse-token-error").Field("error", err).Debug()
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
