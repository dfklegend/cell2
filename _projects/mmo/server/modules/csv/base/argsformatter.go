package base

import (
	"reflect"
	"strconv"
)

// IArgsFormatter
// 处理灵活多个参数的情况下，根据某个类型决定了参数类型
// 表格载入之后，调用后处理，将数据格式化成易用使用的数据
type IArgsFormatter interface {
	Format(entry *IArgs)
}

type IArgs struct {
	Arg0     string `csv:"arg0"`
	Arg1     string `csv:"arg1"`
	Arg2     string `csv:"arg2"`
	Arg3     string `csv:"arg3"`
	ExtArgs  string `csv:"extArgs"`
	ArgsImpl any
}

// FormatArgs 按字段顺序，格式化参数
// arg0 -> field(0), arg1 -> field(1), ...
// 最终存在ArgsImpl上
func (s *IArgs) FormatArgs(out any) {
	s.ArgsImpl = s.doFormatArgs(out)
}

func (s *IArgs) doFormatArgs(out any) any {
	valueOf := reflect.ValueOf(out)
	args := []string{s.Arg0, s.Arg1, s.Arg2, s.Arg3}
	size := valueOf.Elem().NumField()

	for i := 0; i < size; i++ {
		if i < len(args) {
			field := valueOf.Elem().Field(i)
			arg := args[i]

			if !field.CanSet() {
				continue
			}

			switch field.Kind() {
			case reflect.String:
				field.SetString(arg)
			case reflect.Bool:
				value, err := strconv.ParseBool(arg)
				if err == nil {
					field.SetBool(value)
				}
			case reflect.Float32:
				value, err := strconv.ParseFloat(arg, 32)
				if err == nil {
					field.SetFloat(value)
				}
			case reflect.Float64:
				value, err := strconv.ParseFloat(arg, 64)
				if err == nil {
					field.SetFloat(value)
				}
			case reflect.Int8:
				value, err := strconv.ParseInt(arg, 10, 8)
				if err == nil {
					field.SetInt(value)
				}
			case reflect.Int16:
				value, err := strconv.ParseInt(arg, 10, 16)
				if err == nil {
					field.SetInt(value)
				}
			case reflect.Int32:
				value, err := strconv.ParseInt(arg, 10, 32)
				if err == nil {
					field.SetInt(value)
				}
			case reflect.Int64, reflect.Int:
				value, err := strconv.ParseInt(arg, 10, 64)
				if err == nil {
					field.SetInt(value)
				}
			}
		}
	}

	return out
}
