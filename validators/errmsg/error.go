package errmsg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ErrorValidationHandler processes validation errors and returns custom messages
func ErrorValidationHandler(err error) (int, map[string][]string) {
	var (
		errorMessages = make(map[string][]string)
		code          = 400
	)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			var (
				fieldInMsg string
				message    string
				value      = err.Value()
				valueType  = reflect.TypeOf(value)
			)

			fieldInMsg = strings.TrimSpace(strings.ToLower(err.Field()))

			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s harus diisi.", fieldInMsg)
			case "email":
				message = fmt.Sprintf("%s bukan alamat email yang valid.", fieldInMsg)
			case "uuid":
				message = fmt.Sprintf("%s bukan UUID yang valid.", fieldInMsg)
			case "min":
				switch valueType.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
					message = fmt.Sprintf("%s harus minimal %s.", fieldInMsg, err.Param())
				case reflect.String:
					message = fmt.Sprintf("%s harus minimal %s karakter.", fieldInMsg, err.Param())
				case reflect.Slice:
					message = fmt.Sprintf("%s harus minimal %s item.", fieldInMsg, err.Param())
				}
			case "max":
				switch valueType.Kind() {
				case reflect.Int, reflect.Float64:
					message = fmt.Sprintf("%s harus tidak lebih dari %s.", fieldInMsg, err.Param())
				case reflect.String:
					message = fmt.Sprintf("%s harus tidak lebih dari %s karakter.", fieldInMsg, err.Param())
				case reflect.Slice:
					message = fmt.Sprintf("%s harus tidak lebih dari %s item.", fieldInMsg, err.Param())
				}
			case "numeric":
				message = fmt.Sprintf("%s harus angka.", fieldInMsg)
			case "gt":
				message = fmt.Sprintf("%s harus lebih dari %s.", fieldInMsg, err.Param())
			case "gte":
				message = fmt.Sprintf("%s harus lebih dari atau sama dengan %s.", fieldInMsg, err.Param())
			case "lt":
				message = fmt.Sprintf("%s harus kurang dari %s.", fieldInMsg, err.Param())
			case "lte":
				message = fmt.Sprintf("%s harus kurang dari atau sama dengan %s.", fieldInMsg, err.Param())
			case "latitude":
				message = fmt.Sprintf("%s harus latitude yang valid.", fieldInMsg)
			case "longitude":
				message = fmt.Sprintf("%s harus longitude yang valid.", fieldInMsg)
			case "oneof":
				oneOfValues := strings.Split(err.Param(), " ")
				oneOfValues[len(oneOfValues)-1] = "atau " + oneOfValues[len(oneOfValues)-1]
				oneOfValuesStr := strings.Join(oneOfValues, ", ")
				message = fmt.Sprintf("%s harus salah satu dari %s.", fieldInMsg, oneOfValuesStr)
			default:
				message = fmt.Sprintf("%s tidak valid.", fieldInMsg)
			}

			errorMessages[fieldInMsg] = append(errorMessages[fieldInMsg], message)
		}
	}

	return code, errorMessages
}
