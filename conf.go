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

		userInput := prompt(getMessage(typeField))

		if userInput == "" {
			userInput = typeField.Tag.Get("default")
		}

		switch valueField.Kind() {
		case reflect.String:
			valueField.SetString(userInput)

		case reflect.Int:
			valueField.SetInt(cast.ToInt64(userInput))
		}
	}
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
