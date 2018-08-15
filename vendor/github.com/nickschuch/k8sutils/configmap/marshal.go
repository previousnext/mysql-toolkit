package configmap

import (
	"reflect"
	"strconv"
)

// Marshal struct back to the ConfigMap.
func Marshal(obj interface{}) (map[string]string, error) {
	data := make(map[string]string)

	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get(Tag)

		if tag == "" {
			continue
		}

		switch val.Field(i).Kind() {
		case reflect.String:
			data[tag] = val.Field(i).String()

		case reflect.Int64:
			data[tag] = strconv.FormatInt(val.Field(i).Int(), 10)

		case reflect.Bool:
			data[tag] = strconv.FormatBool(val.Field(i).Bool())
		}

	}

	return data, nil
}
