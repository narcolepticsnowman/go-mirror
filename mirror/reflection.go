package mirror

import (
	"fmt"
	"reflect"
	"strings"
)

type Reflection struct {
	ref     reflect.Value
	execRet []reflect.Value
	err     error
}

func Reflect(val interface{}) *Reflection {
	value := reflect.ValueOf(val)
	if value.Type().Kind() != reflect.Ptr {
		value = reflect.ValueOf(&val)
	}
	return &Reflection{ref: value}
}

func ReflectValue(val reflect.Value) *Reflection {
	return &Reflection{ref: val}
}

func (r *Reflection) UnwrapResult() []interface{} {
	if r.execRet == nil || len(r.execRet) == 0 {
		return nil
	}
	res := make([]interface{}, len(r.execRet))
	for i := range r.execRet {
		res[i] = r.execRet[i].Interface()
	}
	return res
}

func (r *Reflection) Exec(args ...interface{}) *Reflection {
	if r.err != nil {
		return r
	}

	if r.ref.Type().Kind() != reflect.Func {
		r.err = fmt.Errorf("field is not a function")
		return r
	}
	in := make([]reflect.Value, len(args))
	for i, _ := range args {
		in[i] = reflect.ValueOf(args[i])
	}
	out := r.ref.Call(in)
	return &Reflection{execRet: out}
}

func valueEmpty(val reflect.Value) bool {
	return val.Kind() == 0 || (val.Type().Kind() == reflect.Ptr && val.IsNil()) || val.IsZero()
}

//get value at path. i.e. /foo/bar/baz or /baz/bar
func (r *Reflection) GetPath(path string) *Reflection {
	if r.err != nil {
		return r
	}
	value := r.ref
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	for _, part := range strings.Split(path, "/") {
		var nextValue reflect.Value
		if value.Type().Kind() == reflect.Ptr {
			nextValue = value.MethodByName(part)
		}
		if valueEmpty(nextValue) {
			nextValue = reflect.Indirect(value).FieldByName(part)
			if valueEmpty(nextValue) {
				r.err = fmt.Errorf("failed to find field or method at path %s", path)
				return r
			}
		}
		value = nextValue
	}
	return ReflectValue(value)
}

//set value at path. i.e. /foo/bar/baz or /baz/bar
func (r *Reflection) SetPath(path string, newValue interface{}) *Reflection {
	if r.err != nil {
		return r
	}
	value := r.GetPath(path).Value()
	if !value.CanSet() {
		r.err = fmt.Errorf("failed to set field %s, ensure the Reflection was created with a pointer", path)
		return r
	}
	value.Set(reflect.ValueOf(newValue))
	return ReflectValue(value)
}

//the last method call result, if any
func (r *Reflection) Value() reflect.Value {
	return r.ref
}

//the last method call result, if any
func (r *Reflection) Ret() []*Reflection {
	refs := make([]*Reflection, len(r.execRet))
	for i := range r.execRet {
		refs[i] = ReflectValue(r.execRet[i])
	}
	return refs
}

//the last method call result, if any
func (r *Reflection) RetValue() []reflect.Value {
	return r.execRet
}

//returns the current error, if any
func (r *Reflection) Err() error {
	return r.err
}

func (r *Reflection) PanicIfErr() *Reflection {
	if r.err != nil {
		panic(r.err)
	}
	return r
}
