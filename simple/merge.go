package simple

import (
	"database/sql"
	"errors"
	"reflect"
)

var (
	ErrInvalidCopyDestination = errors.New("copy destination is invalid")
	ErrInvalidCopyFrom        = errors.New("copy from is invalid")
)

func Merge(toInterface, fromInterface interface{}) (err error) {
	var (
		from = indirect(reflect.ValueOf(fromInterface))
		to   = indirect(reflect.ValueOf(toInterface))
	)

	if !to.CanAddr() {
		return ErrInvalidCopyDestination
	}

	if !from.IsValid() {
		return ErrInvalidCopyFrom
	}

	fromType, _ := indirectType(from.Type())
	toType, _ := indirectType(to.Type())

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}
	var dest, source reflect.Value
	source = indirect(from)
	dest = indirect(to)

	fromTypeFields := deepFields(fromType)
	for _, field := range fromTypeFields {
		srcFieldName, destFieldName := field.Name, field.Name

		if fromField := source.FieldByName(srcFieldName); fromField.IsValid() {
			destFieldNotSet := false
			if f, ok := dest.Type().FieldByName(destFieldName); ok {
				for idx := range f.Index {
					destField := dest.FieldByIndex(f.Index[:idx+1])

					if destField.Kind() != reflect.Ptr {
						continue
					}
					if !destField.IsNil() {
						continue
					}
					if !destField.CanSet() {
						destFieldNotSet = true
						break
					}

					newValue := reflect.New(destField.Type().Elem())
					destField.Set(newValue)
				}
			}

			if destFieldNotSet {
				break
			}

			toField := dest.FieldByName(destFieldName)
			if toField.IsValid() {

				if toField.CanSet() {
					if !set(toField, fromField, true) {
						if err := Merge(toField.Addr().Interface(), fromField.Interface()); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return
}

func set(to, from reflect.Value, deepCopy bool) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {

			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if deepCopy {
			toKind := to.Kind()
			if toKind == reflect.Interface && to.IsNil() {
				if reflect.TypeOf(from.Interface()) != nil {
					to.Set(reflect.New(reflect.TypeOf(from.Interface())).Elem())
					toKind = reflect.TypeOf(to.Interface()).Kind()
				}
			}
			if toKind == reflect.Struct || toKind == reflect.Map || toKind == reflect.Slice {
				return false
			}
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if toScanner, ok := to.Addr().Interface().(sql.Scanner); ok {

			if from.Kind() == reflect.Ptr {
				if from.IsNil() {
					return true
				}
				from = indirect(from)
			}
			err := toScanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem(), deepCopy)
		} else {
			return false
		}
	}
	return true
}

//----------------------------------------------------------------------------------------------------------------------

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) (_ reflect.Type, isPtr bool) {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
		isPtr = true
	}
	return reflectType, isPtr
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	if reflectType, _ = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		fields := make([]reflect.StructField, 0, reflectType.NumField())

		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}

		return fields
	}

	return nil
}

//----------------------------------------------------------------------------------------------------------------------
