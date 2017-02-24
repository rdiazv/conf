package conf

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"reflect"
)

func Load(config interface{}, path string) {
	fileConfig := getFileConfig(path)

	iterateKeys(
		reflect.ValueOf(config).Elem(),
		parseFileConfigValue(fileConfig, ""),
		"",
	)

	writeToFile(config, path)
}

func parseFileConfigValue(fileConfig interface{}, key string) map[string]interface{} {
	if fileConfig == nil {
		return nil
	}

	casted := fileConfig.(map[string]interface{})

	if key == "" {
		return casted
	}

	if casted[key] != nil {
		return casted[key].(map[string]interface{})
	}

	return nil
}

func iterateKeys(config reflect.Value, fileConfig map[string]interface{}, root string) {
	for i := 0; i < config.NumField(); i++ {
		valueField := config.Field(i)
		typeField := config.Type().Field(i)

		if valueField.Kind() == reflect.Struct {
			iterateKeys(
				valueField,
				parseFileConfigValue(fileConfig, typeField.Name),
				root+typeField.Name+".",
			)
		} else {
			for {
				var defaultValue string
				var currentValue string

				if fileConfig != nil && fileConfig[typeField.Name] != nil {
					currentValue = getStringValue(
						reflect.ValueOf(fileConfig[typeField.Name]),
					)

					fmt.Println(typeField.Name, fileConfig[typeField.Name])
				}

				if currentValue != "" {
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

	case reflect.Float64:
		value, castError := cast.ToStringE(valueField.Float())

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

func assignValue(field reflect.Value, userInput string) bool {
	var value reflect.Value
	var err error

	if userInput != "" {
		switch field.Kind() {
		case reflect.String:
			value = reflect.ValueOf(userInput)

		case reflect.Int:
			casted, castErr := cast.ToIntE(userInput)
			err = castErr
			value = reflect.ValueOf(casted)

		case reflect.Bool:
			casted, castErr := cast.ToBoolE(userInput)
			err = castErr
			value = reflect.ValueOf(casted)
		}
	}

	if err != nil {
		return false
	}

	if value.IsValid() {
		field.Set(value)
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

func getFileConfig(path string) interface{} {
	var config interface{}

	data, err := ioutil.ReadFile(path)

	if err == nil {
		json.Unmarshal(data, &config)
	}

	return config
}

func writeToFile(config interface{}, path string) {
	data, _ := json.MarshalIndent(config, "", "  ")
	ioutil.WriteFile(path, data, 0644)
}
