package gexec

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Result struct {
	stdout []byte
}

func (r *Result) Bytes() []byte {
	return r.stdout
}

func (r *Result) String() string {
	return string(r.stdout)
}

func (r *Result) Scan(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.New("cmd: Scan(non-pointer " + rv.Type().String() + ")")
	}
	return decode(string(r.stdout), rv.Elem())
}

func decode(s string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool:
		switch strings.TrimSpace(s) {
		case "1", "t", "T", "true", "TRUE", "True", "y", "Y", "yes", "YES", "Yes", "on", "ON", "On":
			v.SetBool(true)
		case "0", "f", "F", "false", "FALSE", "False", "n", "N", "no", "NO", "No", "off", "OFF", "Off":
			v.SetBool(false)
		}
		return errors.New("cmd: cannot decode " + v.Kind().String() + ": " + s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil || v.OverflowInt(n) {
			return errors.New("cmd: cannot decode " + v.Kind().String() + ": " + s)
		}
		v.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
		if err != nil || v.OverflowUint(n) {
			return errors.New("cmd: cannot decode " + v.Kind().String() + ": " + s)
		}
		v.SetUint(n)

	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(strings.TrimSpace(s), v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			return errors.New("cmd: cannot decode " + v.Kind().String() + ": " + s)
		}
		v.SetFloat(n)

	case reflect.String:
		v.SetString(strings.TrimSpace(s))

	default:
		// reflect.Uintptr, reflect.Complex64, reflect.Complex128,
		// reflect.Array, reflect.Chan, reflect.Func,reflect.Interface,
		// reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct,
		// reflect.UnsafePointer
		return errors.New("cmd: unsupported type: " + v.Kind().String())
	}

	return nil
}
