package entity

/** ID **/
type ID interface {
	Value() string
}

/** ID validator **/
type IDValidator interface {
	Validate(value string) error
}

/** ID generator **/
type IDGenerator interface {
	Generate() ID
	From(value string) (ID, error)
}

/** ID manager **/
func NewIDManager(validator IDValidator, generator IDGenerator) IDManager {
	return &idManager{
		IDValidator: validator,
		IDGenerator: generator,
	}
}

type IDManager interface {
	IDGenerator
	IDValidator
}

type idManager struct {
	IDGenerator
	IDValidator
}
