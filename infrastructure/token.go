package infrastructure

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/mkaiho/go-lambda-api-sample/entity"
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
	GenerateToken(subject string) (jwt.Token, error)
	Sign(token jwt.Token, key jwk.Key) ([]byte, error)
	Verify(signed []byte, key jwk.Key) ([]byte, error)
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
