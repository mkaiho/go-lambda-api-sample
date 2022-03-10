package entity

/** Email **/
type Email interface {
	Value() string
}

type email struct {
	value string
}

func (em email) Value() string {
	return em.value
}

func (em email) validate() error {
	// TODO: implement
	return nil
}

func NewEmail(value string) (Email, error) {
	em := email{
		value: value,
	}
	if err := em.validate(); err != nil {
		return nil, err
	}
	return em, nil
}
