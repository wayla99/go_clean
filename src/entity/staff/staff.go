package staff

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidStaff = errors.New("invalid staff")
)

type Staff struct {
	Id        string
	FirstName string `validate:"omitempty,min=2,max=255"`
	LastName  string `validate:"omitempty,min=2,max=255"`
	Email     string `validate:"omitempty,email"`
}

func (s *Staff) Validate() error {
	validate := validator.New()

	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidStaff, err.Error())
	}

	return nil
}
