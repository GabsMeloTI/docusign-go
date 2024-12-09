package helper

import (
	"reflect"
	"strings"
)

func MapFormToStruct(form map[string][]string, dst interface{}) error {
	dstVal := reflect.ValueOf(dst).Elem()
	dstType := dstVal.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		formTag := field.Tag.Get("form")
		if values, ok := form[formTag]; ok && len(values) > 0 {
			fieldVal := dstVal.FieldByName(field.Name)
			if fieldVal.CanSet() {
				fieldVal.SetString(values[0])
			}
		}
	}
	return nil
}

func ExtractNameFromContentDisposition(header string) string {
	parts := strings.Split(header, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "name=") {
			return strings.Trim(part[5:], `"`)
		}
	}
	return ""
}
