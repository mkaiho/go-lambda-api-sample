package infrastructure

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/mkaiho/go-lambda-api-sample/entity"
	ulidlib "github.com/oklog/ulid/v2"
)

var _ entity.ID = (*ulid)(nil)
var _ entity.IDValidator = (*ulidValidator)(nil)
var _ entity.IDGenerator = (*ulidGenerator)(nil)

/** ID **/
type ulid struct {
	value string
}

func (id *ulid) Value() string {
	return string(id.value)
}

func (id *ulid) IsEmpty() bool {
	return id == nil || len(id.value) == 0
}

/** ID Validator **/
func NewULIDValidator() entity.IDValidator {
	return &ulidValidator{}
}

type ulidValidator struct{}

func (v *ulidValidator) Validate(value string) error {
	_, err := ulidlib.ParseStrict(value)
	return err
}

/** ID Generator **/
func NewULIDGenerator() entity.IDGenerator {
	t := time.Now()
	pool := sync.Pool{
		New: func() interface{} {
			return ulidlib.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
		},
	}
	return &ulidGenerator{
		pool: pool,
	}
}

type ulidGenerator struct {
	pool      sync.Pool
	validator ulidValidator
}

func (g *ulidGenerator) From(value string) (entity.ID, error) {
	if err := g.validator.Validate(value); err != nil {
		return nil, err
	}
	return &ulid{
		value: strings.ToLower(value),
	}, nil
}

func (g *ulidGenerator) Generate() entity.ID {
	t := time.Now()
	entropy := g.pool.Get().(*ulidlib.MonotonicEntropy)
	defer g.pool.Put(entropy)
	id := ulidlib.MustNew(ulidlib.Timestamp(t), entropy)
	return entity.ID(&ulid{
		value: strings.ToLower(id.String()),
	})
}
