package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// GetValidatorController return validator controller
func GetValidatorController() *validator.Validate {
	return validator.New()
}

// getErrorMessage to define the error when validate field.
// If you want to define or find another specific error to custom,
// please add the case below and also with the message.
//
//	Usage example:
//		errorMessage := getErrorMessage(err, "required")
//		return errorMessage
//
//	Define new error:
//		case "error_something":
//			return "this error occurs because of this case is not fulfilled"
func getErrorMessage(err validator.FieldError, jsonField string) string {
	if jsonField == "" {
		return err.Tag()
	}

	switch err.Tag() {
	case "required", "required_if":
		return fmt.Sprintf("%s is required", jsonField)
	case "min":
		return fmt.Sprintf("%s must be at least %s", jsonField, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", jsonField, err.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", jsonField)
	case "oneof":
		choices := strings.Split(err.Param(), " ")
		choicesStr := ""
		for i, v := range choices {
			choicesStr += fmt.Sprintf(`%v`, v)
			if i != len(choices)-1 {
				choicesStr += ", "
			}
		}
		return fmt.Sprintf("%s valid choices are: %s", jsonField, choicesStr)
	case "datetimeformat":
		return fmt.Sprintf("%s datetime format is YYYY-MM-DD hh:mm", jsonField)
	case "gt":
		minimumLength := err.Param()
		if minimumLength == "0" {
			return fmt.Sprintf("%s can't be empty", jsonField)
		}
		return fmt.Sprintf("Minimum length of %s is %s", jsonField, minimumLength)
	default:
		return fmt.Sprintf("Validation error on field %s", jsonField)
	}
}

// ValidatePayload to validate payload in JSON format
func ValidatePayload(request *http.Request, s interface{}) error {
	err := json.NewDecoder(request.Body).Decode(s)
	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			errorType := err.(*json.UnmarshalTypeError)
			errorFormat := errors.New(fmt.Sprintf("invalid type of %s (expected: %s, got: %s)", errorType.Field, errorType.Type, errorType.Value))
			return errorFormat
		default:
			return errors.New(fmt.Sprintf("payload error: %s", err.Error()))
		}
	}

	errorValidate := ValidateStruct(s)
	if errorValidate != nil {
		return errorValidate
	}

	return nil
}

// ValidateStruct to validate struct using Go Validator (returning map of error: model.Errors)
func ValidateStruct(structObj interface{}) error {
	validatorObj := GetValidatorController()

	// add custom validator to validate field with datetime format (YYYY-MM-DD hh:mm)
	validatorObj.RegisterValidation("datetimeformat", func(fl validator.FieldLevel) bool {
		datetimeStr := fl.Field().String()
		_, err := time.Parse("2006-01-02 15:04", datetimeStr)
		return err == nil
	})

	if err := validatorObj.Struct(structObj); err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorType := err.(validator.ValidationErrors)
			for i := 0; i < len(errorType); {
				errorField := errorType[i]
				jsonField := GetJsonTagInStruct(errorField.Field(), structObj)
				errorMessage := getErrorMessage(errorField, jsonField)
				return errors.New(errorMessage)
			}
		case *validator.InvalidValidationError:
			errorType := err.(*validator.InvalidValidationError)
			log.Println(fmt.Sprintf(`[ERROR] validator.InvalidValidationError: {"type": "%v", "key": "%v", "name": "%v"}`, errorType.Type, errorType.Type.Key(), errorType.Type.Name()))
			return errors.New(fmt.Sprintf("payload error: %s", err.Error()))
		default:
			return errors.New(fmt.Sprintf("payload error: %s", err.Error()))
		}
	}
	return nil
}

// GetJsonTagInStruct to get the JSON tag of struct field
func GetJsonTagInStruct(fieldName string, structOfField any) string {
	res, err := getFieldJSONTagRecursive(reflect.ValueOf(structOfField).Type(), fieldName)
	if err != nil {
		return ""
	}
	return res
}

func getFieldJSONTagRecursive(t reflect.Type, fieldName string) (string, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Check if the field is the one we're looking for
		if field.Name == fieldName {
			return field.Tag.Get("json"), nil
		}

		// If the field is a struct, search recursively
		if field.Type.Kind() == reflect.Struct {
			if tag, err := getFieldJSONTagRecursive(field.Type, fieldName); err == nil {
				return tag, nil
			}
		}
	}
	return "", fmt.Errorf("field %s not found", fieldName)
}
