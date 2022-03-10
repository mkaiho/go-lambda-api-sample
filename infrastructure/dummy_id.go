package infrastructure

import "github.com/mkaiho/go-lambda-api-sample/entity"

var _ entity.ID = (*dummyID)(nil)
var _ entity.IDValidator = (*dummyIDValidator)(nil)
var _ entity.IDGenerator = (*dummyIDGenerator)(nil)

/** ID **/
type dummyID string

func (id dummyID) Value() string {
	return string(id)
}

func (id dummyID) IsEmpty() bool {
	return len(id) == 0
}

/** ID Validator **/
func NewDummyIDValidator() entity.IDValidator {
	return &dummyIDValidator{}
}

type dummyIDValidator struct{}

func (v *dummyIDValidator) Validate(value string) error {
	return nil
}

/** ID Generator **/
func NewDummyIDGenerator() entity.IDGenerator {
	return &dummyIDGenerator{}
}

type dummyIDGenerator struct {
	validator dummyIDValidator
}

func (g *dummyIDGenerator) From(value string) (entity.ID, error) {
	if err := g.validator.Validate(value); err != nil {
		return nil, err
	}
	return dummyID(value), nil
}

func (g *dummyIDGenerator) Generate() entity.ID {
	id, _ := g.From("001")
	return id
}
