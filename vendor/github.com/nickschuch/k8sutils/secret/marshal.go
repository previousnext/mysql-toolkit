package secret

import (
	"encoding/base64"
	"fmt"
	"reflect"
)

// Marshal struct back to the Secret.
func Marshal(obj interface{}) (map[string][]byte, error) {
	data := make(map[string][]byte)

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

		if val.Field(i).Kind() == reflect.String {
			enc := base64.StdEncoding.EncodeToString([]byte(val.Field(i).String()))
			fmt.Println(tag, "=", enc)
			data[tag] = []byte(enc)
		}
	}

	return data, nil
}
