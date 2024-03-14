package apientry

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// 首字母大写
func isExported(name string) bool {
	w, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(w)
}

// 指针类型
func makeValueMaybeNil(typ reflect.Type, v interface{}) reflect.Value {
	if v == nil {
		return reflect.New(typ).Elem()
	} else {
		return reflect.ValueOf(v)
	}
}

// ToLowerCamelCase
// 转成小驼峰
// ASCII only
func ToLowerCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))

	c := s[0]
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	b.WriteByte(c)

	for i := 1; i < len(s); i++ {
		b.WriteByte(s[i])
	}
	return b.String()
}

func CheckInvokeCBFunc(cbFunc HandlerCBFunc, e error, result interface{}) {
	if cbFunc == nil {
		return
	}
	cbFunc(e, result)
}
