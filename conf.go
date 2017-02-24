package conf

import (
  "reflect"
  "fmt"
)

func Load(config interface{}, path string) {
  val := reflect.ValueOf(config).Elem()

  for i := 0; i < val.NumField(); i++ {
    valueField := val.Field(i)
    typeField := val.Type().Field(i)
    // tag := typeField.Tag

    value := prompt(getMessage(typeField))

    if value == "" {
      value = typeField.Tag.Get("default")
    }

    valueField.SetString(value)

    // fmt.Println(typeField.Name, tag.Get("default"), tag.Get("required"))

    // fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("required"))
  }
}

func getMessage(field reflect.StructField) string {
  message := field.Name
  defaultValue := field.Tag.Get("default")

  if (defaultValue != "") {
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
