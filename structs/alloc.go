package structs

import (
	"reflect"

	"github.com/0meet1/zero-framework/global"
)

var XreflectMeta = func(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Ptr {
		return t
	}
	meta := t
	for meta.Kind() == reflect.Ptr {
		meta = meta.Elem()
	}
	return meta
}

var Xalloc = func(t reflect.Type) any {
	trydef := func(ndef any) {
		defer func() {
			err := recover()
			if err != nil {
				global.Logger().ErrorS(err.(error))
			}
		}()
		m := reflect.ValueOf(ndef).MethodByName("ThisDef")
		if m.IsValid() {
			m.Call([]reflect.Value{reflect.ValueOf(ndef)})
		}
	}
	n := reflect.New(XreflectMeta(t))
	trydef(n.Interface())
	return n.Interface()
}
