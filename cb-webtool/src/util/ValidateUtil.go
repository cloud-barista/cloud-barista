package util

import (
	"fmt"
	// "gopkg.in/go-playground/validator.v9"
)

// var validate *validator.Validate

type ValidObject struct {
	Name string `validate:"required"`
}

// Vadidation check
func ValidateString(param string) error {
	vObj := ValidObject{}
	vObj.Name = param

	fmt.Println("param valid ", param)
	// validate = validator.New()
	// err := validate.Struct(vObj)
	// if err != nil {
	//     return err
	// }
	return nil
}
