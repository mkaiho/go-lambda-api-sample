package infrastructure

import (
	"strings"

	uuidlib "github.com/google/uuid"
	"github.com/mkaiho/go-lambda-api-sample/entity"
)

var _ entity.ID = (*uuid)(nil)
var _ entity.IDValidator = (*uuidValidator)(nil)
var _ entity.IDGenerator = (*uuidGenerator)(nil)

/** ID **/
type uuid struct {
	value string
}

func (id *uuid) Value() string {
	return string(id.value)
}

func (id *uuid) IsEmpty() bool {
	return id == nil || len(id.value) == 0
}

/** ID Validator **/
func NewUUIDValidator() entity.IDValidator {
	return &uuidValidator{}
}

type uuidValidator struct{}

func (v *uuidValidator) Validate(value string) error {
	_, err := uuidlib.Parse(value)
	return err
}

/** ID Generator **/
func NewUUIDGenerator() entity.IDGenerator {
	return &uuidGenerator{}
}

type uuidGenerator struct {
	validator uuidValidator
}

func (g *uuidGenerator) From(value string) (entity.ID, error) {
	if err := g.validator.Validate(value); err != nil {
		return nil, err
	}
	return &uuid{
		value: strings.ToLower(value),
	}, nil
}

func (g *uuidGenerator) Generate() entity.ID {
	return entity.ID(&uuid{
		value: uuidlib.NewString(),
	})
}
