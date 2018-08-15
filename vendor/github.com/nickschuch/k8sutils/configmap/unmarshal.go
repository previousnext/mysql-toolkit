package configmap

import (
	"fmt"
	"reflect"
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

// Unmarshal a ConfigMap to a struct.
func Unmarshal(cfg *corev1.ConfigMap, obj interface{}) error {
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

		if d, ok := cfg.Data[tag]; ok {
			switch val.Field(i).Kind() {
			case reflect.String:
				val.Field(i).SetString(d)

			case reflect.Int64:
				v, err := strconv.ParseInt(d, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to unmarshal %s with value %s to type int: %s", tag, d, err)
				}
				val.Field(i).SetInt(v)

			case reflect.Bool:
				v, err := strconv.ParseBool(d)
				if err != nil {
					return fmt.Errorf("failed to unmarshal %s with value %s to type bool: %s", tag, d, err)
				}
				val.Field(i).SetBool(v)
			}
		}
	}

	return nil
}
