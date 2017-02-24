package conf

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"reflect"
)

func Load(config interface{}, path string) {
	ok := readFromFile(config, path)
	iterateKeys(reflect.ValueOf(config).Elem(), "", ok)
	writeToFile(config, path)
}

func iterateKeys(val reflect.Value, root string, fromFile bool) {
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if valueField.Kind() == reflect.Struct {
			iterateKeys(valueField, root+typeField.Name+".", fromFile)
		} else {
			for {
				var defaultValue string
				currentValue := getStringValue(valueField)

				if fromFile && currentValue != "" {
					defaultValue = currentValue
				} else {
					defaultValue = typeField.Tag.Get("default")
				}

				required := typeField.Tag.Get("required") == "true"
				userInput := prompt(getMessage(root+typeField.Name, defaultValue, required))

				if userInput == "" {
					userInput = defaultValue
				}

				if userInput == "" && required {
					fmt.Println("Required field.")
					continue
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

func getStringValue(valueField reflect.Value) string {
	switch valueField.Kind() {
	case reflect.String:
		return valueField.String()

	case reflect.Int:
		value, castError := cast.ToStringE(valueField.Int())

		if castError != nil {
			return ""
		}

		return value

	case reflect.Bool:
		value, castError := cast.ToStringE(valueField.Bool())

		if castError != nil {
			return ""
		}

		return value
	}

	return ""
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

func getMessage(field string, defaultValue string, required bool) string {
	message := field

	if defaultValue != "" {
		message = message + " [" + defaultValue + "]"
	}

	if required {
		message = message + "*"
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

func readFromFile(config interface{}, path string) bool {
	data, err := ioutil.ReadFile(path)

	if err == nil {
		json.Unmarshal(data, &config)
		return true
	}

	return false
}

func writeToFile(config interface{}, path string) {
	data, _ := json.MarshalIndent(config, "", "  ")
	ioutil.WriteFile(path, data, 0644)
}
