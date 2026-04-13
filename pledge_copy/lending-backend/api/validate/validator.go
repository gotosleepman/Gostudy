package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func BindingValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v
	}
}
