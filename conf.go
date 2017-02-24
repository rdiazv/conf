package conf

import (
	"fmt"
	"github.com/spf13/cast"
	"reflect"
)

func Load(config interface{}, path string) {
	val := reflect.ValueOf(config).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		for {
			userInput := prompt(getMessage(typeField))

			if userInput == "" {
				userInput = typeField.Tag.Get("default")
			}

			ok := assignValue(valueField, userInput)

			if ok {
				break
			}

			fmt.Printf("Invalid %s value.\n", valueField.Kind())
		}
	}
}

func assignValue(valueField reflect.Value, userInput string) bool {
	switch valueField.Kind() {
	case reflect.String:
		valueField.SetString(userInput)

	case reflect.Int:
		value, castError := cast.ToInt64E(userInput)

		if castError != nil {
			return false
		}

		valueField.SetInt(value)

	case reflect.Bool:
		value, castError := cast.ToBoolE(userInput)

		if castError != nil {
			return false
		}

		valueField.SetBool(value)
	}

	return true
}

func getMessage(field reflect.StructField) string {
	message := field.Name
	defaultValue := field.Tag.Get("default")

	if defaultValue != "" {
		message = message + " (" + defaultValue + ")"
	}

	message = message + ": "

	return message
}

func prompt(message string) string {
	var response string

	fmt.Print(message)
	fmt.Scanln(&response)

	return response
}
