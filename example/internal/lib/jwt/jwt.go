package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	day  = time.Hour * 24
	week = day * 7

	accessTokenExpiresIn  = time.Hour
	refreshTokenExpiresIn = week * 2

	idKey = "id"
)

var (
	ErrExpiredToken = errors.New("token is expired")

	signingMethod = jwt.SigningMethodHS256
	signingSecret = []byte("secret")
)

type Token struct {
	Value     string
	ExpiresAt time.Time
}

func (t Token) Expiration() time.Duration {
	return time.Until(t.ExpiresAt)
}

func (t Token) IsEmpty() bool {
	return t.Value == ""
}

type TokenSet struct {
	Access  Token
	Refresh Token
}

func (ts TokenSet) IsEmpty() bool {
	return ts.Access.IsEmpty() && ts.Refresh.IsEmpty()
}

func New(id string) (TokenSet, error) {
	access, err := generateToken(id, accessTokenExpiresIn)
	if err != nil {
		return TokenSet{}, err
	}
	refresh, err := generateToken(id, refreshTokenExpiresIn)
	if err != nil {
		return TokenSet{}, err
	}

	return TokenSet{Access: access, Refresh: refresh}, nil
}

func generateToken(id string, expiresIn time.Duration) (Token, error) {
	claims := jwt.MapClaims{}
	claims[idKey] = id
	exp := time.Now().Add(expiresIn)
	claims["exp"] = exp.Unix()

	jwtToken := jwt.NewWithClaims(signingMethod, claims)
	value, err := jwtToken.SignedString(signingSecret)
	if err != nil {
		return Token{}, fmt.Errorf("failed to sign token jwt: %w", err)
	}
	return Token{Value: value, ExpiresAt: exp}, nil
}

func Verify(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("failed to parse jwt: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid jwt")
}

func ExtractID(tokenStr string) (string, error) {
	claims, err := Verify(tokenStr)
	if err != nil {
		return "", fmt.Errorf("failed to verify jwt: %w", err)
	}
	id, ok := claims[idKey].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract id from jwt")
	}
	return id, nil
}
