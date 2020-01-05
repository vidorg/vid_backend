package router

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

func matchString(reg string, content string) bool {
	re, err := regexp.Compile(reg)
	if err != nil {
		return true // error reg default to match success
	}
	return re.MatchString(content)
}

func setupRegexParamBinding(tag string) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			return matchString(fl.Param(), fl.Field().String())
		})
	}
}

func setupRegexBinding(tag string, re string) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			return matchString(re, fl.Field().String())
		})
	}
}

func SetupDefinedValidation() {
	setupRegexParamBinding("regexp")
	setupRegexBinding("name", "^[^'`\"\\\\]+$")          // ' ` " \
	setupRegexBinding("pwd", "^[a-zA-Z0-9+\\-*/.=_~]+$") // + - * / . = _ ~
	setupRegexBinding("phone", "^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$")
}
