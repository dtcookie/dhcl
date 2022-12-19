package hcl

import (
	"reflect"
	"strings"

	"github.com/dtcookie/camel"
)

var NoDocumentationAvailable = "No documentation available"

type handler struct {
	Field         reflect.StructField
	OmitEmpty     bool
	Elem          string
	Property      string
	Unordered     bool
	Sensitive     bool
	Value         reflect.Value
	Documentation string
}

func evalHandlers(refType reflect.Type, v interface{}) []handler {
	handlers := []handler{}
	for idx := 0; idx < refType.NumField(); idx++ {
		field := refType.Field(idx)
		if !field.IsExported() {
			continue
		}
		if field.Anonymous {
			anonFieldValue := reflect.ValueOf(v).Field(idx).Interface()
			handlers = append(handlers, evalHandlers(field.Type, anonFieldValue)...)
			continue
		}
		documentation := NoDocumentationAvailable
		tag := field.Tag
		if tagValue := tag.Get("doc"); len(tagValue) > 0 {
			documentation = tagValue
		}
		sensitive := false
		if tagValue := tag.Get("hcl"); len(tagValue) > 0 {
			tagValues := strings.Split(tagValue, ",")
			propertyName := tagValues[0]
			if propertyName == "-" {
				continue
			} else if propertyName == "" || strings.Contains(propertyName, "=") {
				propertyName = camel.Strip(field.Name)
			}
			omitEmpty := false
			elem := ""
			unordered := false
			for _, tagValue := range tagValues {
				if tagValue == "omitempty" {
					omitEmpty = true
				}
				if tagValue == "sensitive" {
					sensitive = true
				}
				if strings.HasPrefix(tagValue, "elem=") {
					elem = strings.TrimSpace(strings.TrimPrefix(tagValue, "elem="))
				}
				if tagValue == "unordered" {
					unordered = true
				}
			}
			handlers = append(handlers, handler{Field: field, OmitEmpty: omitEmpty, Elem: elem, Property: propertyName, Sensitive: sensitive, Unordered: unordered, Value: reflect.ValueOf(v).Field(idx), Documentation: documentation})
		} else if tagValue := tag.Get("json"); len(tagValue) > 0 {
			tagValues := strings.Split(tagValue, ",")
			propertyName := tagValues[0]
			if propertyName == "-" {
				continue
			} else if propertyName == "" {
				propertyName = field.Name
			}
			omitEmpty := false
			propertyName = camel.Strip(propertyName)
			for _, tagValue := range tagValues {
				if tagValue == "omitempty" {
					omitEmpty = true
				}
			}
			handlers = append(handlers, handler{Field: field, Sensitive: sensitive, OmitEmpty: omitEmpty, Property: propertyName, Value: reflect.ValueOf(v).Field(idx), Documentation: documentation})
		} else {
			handlers = append(handlers, handler{Field: field, Sensitive: sensitive, OmitEmpty: false, Property: camel.Strip(field.Name), Value: reflect.ValueOf(v).Field(idx), Documentation: documentation})
		}
	}
	return handlers
}
