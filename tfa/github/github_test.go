package github_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"reflect"
)

func TestGithub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GithubSuite")
}

var _ = Describe("GithubSuite", func() {
	It("", func() {
		arr := [...]interface{}{
			"123","123",
		}

		val := reflect.ValueOf(arr)

		clearField(val)
	})
})

func clearField(field reflect.Value) {
	if !field.CanAddr() {
		return
	}

	fieldType := field.Type()
	if fieldType.Kind() == reflect.Ptr {
		field.SetPointer(nil)
		return
	}
	switch fieldType.Kind() {
	case reflect.Bool:
		field.SetBool(false)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		field.SetFloat(0)
	case reflect.String:
		field.SetString("")
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		field.SetInt(0)
	case reflect.Struct:
		if fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Interface {
			for i := 0; i < field.Elem().NumField(); i++ {
				clearField(field.Elem().Field(i))
			}
		}
	case reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			clearField(field.Field(i))
		}
	case reflect.Map:
		for _, k := range field.MapKeys() {
			originalValue := field.MapIndex(k)
			clearField(originalValue)
		}
	default:
	}
}