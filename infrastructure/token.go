package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/mkaiho/go-lambda-api-sample/entity"
	"github.com/mkaiho/go-lambda-api-sample/util"
)

type TokenManagerConfig interface {
	Issuer() string
	Audience() []string
	ValidityDuration() time.Duration
}

type jwtManagerConfig struct {
	issuer           string
	audience         []string
	validityDuration time.Duration
}

func (g *jwtManagerConfig) Issuer() string {
	return g.issuer
}
func (g *jwtManagerConfig) Audience() []string {
	return g.audience
}
func (g *jwtManagerConfig) ValidityDuration() time.Duration {
	return g.validityDuration
}

func NewTokenManagerConfig(
	issuer string,
	audience []string,
	validityDuration time.Duration,
) TokenManagerConfig {
	return &jwtManagerConfig{
		issuer:           issuer,
		audience:         audience,
		validityDuration: validityDuration,
	}
}

type JWTManager interface {
	ReadRSAPrivatePemFile(path string) (jwk.Key, error)
	ReadJWKSetFile(path string) (jwk.Set, error)
	GenerateToken(subject string) (jwt.Token, error)
	Sign(token jwt.Token, key jwk.Key) ([]byte, error)
	Verify(signed []byte, key jwk.Key) ([]byte, error)
	VerifyWithKeySet(signed []byte, publicKeySet jwk.Set) ([]byte, error)
}

type jwtManager struct {
	issuer           string
	audience         []string
	validityDuration time.Duration
	idgen            entity.IDGenerator
}

func NewJWTManager(config TokenManagerConfig, idgen entity.IDGenerator) JWTManager {
	return &jwtManager{
		idgen:            idgen,
		issuer:           config.Issuer(),
		audience:         config.Audience(),
		validityDuration: config.ValidityDuration(),
	}
}

func (g *jwtManager) ReadRSAPrivatePemFile(path string) (jwk.Key, error) {
	pem, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := jwk.ParseKey(pem, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}
	if key.KeyType() != jwa.RSA {
		return nil, fmt.Errorf("unsupported key type: %s", key.KeyType())
	}
	key.Set(jwk.KeyIDKey, g.idgen.Generate().Value())
	key.Set(jwk.AlgorithmKey, "RS256")
	key.Set(jwk.KeyUsageKey, jwk.ForSignature)

	return key, nil
}

func (g *jwtManager) ReadJWKSetFile(path string) (jwk.Set, error) {
	pem, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keySet, err := jwk.Parse(pem)
	if err != nil {
		return nil, err
	}

	return keySet, nil
}

func (g *jwtManager) GenerateToken(
	subject string,
) (jwt.Token, error) {
	issuedAt := time.Now()
	b := jwt.NewBuilder()
	b = b.JwtID(g.idgen.Generate().Value())
	b = b.Subject(subject)
	b = b.IssuedAt(issuedAt)
	b = b.Issuer(g.issuer)
	b = b.Audience(g.audience)
	b = b.Expiration(issuedAt.Add(g.validityDuration))

	return b.Build()
}

func (g *jwtManager) Sign(
	token jwt.Token,
	key jwk.Key,
) ([]byte, error) {
	headers := jws.NewHeaders()
	km, err := key.AsMap(context.TODO())
	if err != nil {
		return nil, err
	}
	for k, v := range km {
		headers.Set(k, v)
	}
	return jwt.Sign(
		token,
		jwt.WithKey(
			jwa.RS256,
			key,
			jws.WithProtectedHeaders(headers),
		),
	)
}

func (g *jwtManager) Verify(
	signed []byte,
	key jwk.Key,
) ([]byte, error) {
	return jws.Verify(signed, jws.WithKey(jwa.RS256, key))
}

func (g *jwtManager) VerifyWithKeySet(
	signed []byte,
	publicKeySet jwk.Set,
) ([]byte, error) {
	payload, err := jwt.Parse(signed, jwt.WithKeySet(publicKeySet))
	if err != nil {
		return nil, err
	}
	if !g.isValidAudience(payload.Audience()) {
		return nil, errors.New("invalid audience")
	}
	if payload.Issuer() != g.issuer {
		return nil, errors.New("invalid issuer")
	}

	buff, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

func (g *jwtManager) isValidAudience(vs []string) bool {
	util.GetLogger().Info("check udience", "aud", g.audience, "vs", vs)
	for _, aud := range g.audience {
		isMatch := false
		for _, v := range vs {
			if aud == v {
				isMatch = true
				break
			}
		}
		if !isMatch {
			return false
		}
	}

	return true
}
