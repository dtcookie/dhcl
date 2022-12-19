package hcl

import (
	"fmt"
	"reflect"
)

type UnsupportedTypeError struct {
	Name string
	Type reflect.Type
}

func (me UnsupportedTypeError) Error() string {
	return fmt.Sprintf("The type '%s' (kind '%s') for field %s not is supported", me.Type, me.Type.Kind(), me.Name)
}
