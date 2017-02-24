package conf

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"reflect"
)

func Load(config interface{}, path string) {
	iterateKeys(reflect.ValueOf(config).Elem(), "")
	writeToFile(config, path)
}

func iterateKeys(val reflect.Value, root string) {
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if valueField.Kind() == reflect.Struct {
			iterateKeys(valueField, root+typeField.Name+".")
		} else {
			for {
				userInput := prompt(getMessage(typeField, root))

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

func getMessage(field reflect.StructField, root string) string {
	message := root + field.Name
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

func writeToFile(config interface{}, path string) {
	data, _ := json.MarshalIndent(config, "", "  ")
	ioutil.WriteFile(path, data, 0644)
}
