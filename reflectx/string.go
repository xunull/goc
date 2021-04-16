package reflectx

import (
	"github.com/xunull/goc/commonx"
	"reflect"
	"strconv"
)

func StringToCommon(str string, t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Int:
		i, err := strconv.Atoi(str)
		commonx.CheckErrOrFatal(err)
		return reflect.ValueOf(i)
	case reflect.Int8:
		i, err := strconv.Atoi(str)
		commonx.CheckErrOrFatal(err)
		return reflect.ValueOf(int8(i))
	case reflect.Int16:
		i, err := strconv.Atoi(str)
		commonx.CheckErrOrFatal(err)
		return reflect.ValueOf(int16(i))
	case reflect.Int32:
		i, err := strconv.Atoi(str)
		commonx.CheckErrOrFatal(err)
		return reflect.ValueOf(int32(i))
	case reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		commonx.CheckErrOrFatal(err)
		return reflect.ValueOf(i)
	case reflect.String:
		return reflect.ValueOf(str)
	case reflect.Bool:
		i, err := strconv.ParseBool(str)
		commonx.CheckErrOrFatal(err)
		return reflect.ValueOf(i)
	case reflect.Float32:
		i, err := strconv.ParseFloat(str, 32)
		commonx.CheckErrOrFatal(err)
		return reflect.ValueOf(i)
	case reflect.Float64:
		i, err := strconv.ParseFloat(str, 64)
		commonx.CheckErrOrFatal(err)
		return reflect.ValueOf(i)
	default:
		return reflect.ValueOf(nil)
	}

}
