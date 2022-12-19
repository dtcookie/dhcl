package hcl

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type encoder handler

func (me encoder) encode(ctx context.Context) (interface{}, error) {
	if me.OmitEmpty {
		switch kind := me.Value.Kind(); kind {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
			if !me.Value.IsValid() || me.Value.IsZero() || me.Value.IsNil() {
				return emptyValue, nil
			}
		default:
			if !me.Value.IsValid() || me.Value.IsZero() {
				return emptyValue, nil
			}
		}
	}

	switch kind := me.Value.Kind(); kind {
	case reflect.Map, reflect.Interface, reflect.Array, reflect.Uintptr, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return nil, UnsupportedTypeError{me.Field.Name, me.Value.Type()}
	case reflect.Struct:
		serialized, err := Marshal(ctx, me.Value.Interface())
		if err != nil {
			return nil, err
		}
		return []interface{}{serialized}, nil
	case reflect.Pointer:
		if !me.Value.IsValid() || me.Value.IsZero() || me.Value.IsNil() {
			return nil, nil
		}
		return encoder{Value: me.Value.Elem(), OmitEmpty: false, Field: me.Field, Unordered: me.Unordered, Property: me.Property}.encode(ctx)
	case reflect.Slice:
		refElemsSlice := reflect.New(serialTypeOf(me.Value.Type())).Elem()
		for idx := 0; idx < me.Value.Len(); idx++ {
			elem, err := encoder{Value: me.Value.Index(idx), OmitEmpty: false, Field: me.Field, Unordered: me.Unordered, Property: me.Property}.encode(ctx)
			if err != nil {
				return nil, err
			}
			elemValue := reflect.ValueOf(elem)
			refElemsSlice = reflect.Append(refElemsSlice, elemValue)
		}
		// special handling for struct slices - we need to re-arrange things here
		// reflect.Slice hints that everything needs to be stored within a []interface{}
		// Slice elements are Structs - which are also represented as []interface{}
		//
		// The correct representation is
		// * if no elem tag was specified
		//   * []interface{}
		//     * elements: map[string]interface{}
		//       * key: property
		//   * if the elem tag was specified
		//     * []interface{}
		//       * exactly ONE element of type map[string]interface{}
		//         * key: property
		//         * exaclty ONE element of type []interface{}
		//           * elements: map[string]interface{}
		//             * key: elem tag
		if refElemsSlice.Type().Elem().Kind() == reflect.Slice {
			sliceOfMaps := []interface{}{}
			for _, sliceContainingOneMap := range refElemsSlice.Interface().([][]interface{}) {
				sliceOfMaps = append(sliceOfMaps, sliceContainingOneMap[0])
			}
			if len(me.Elem) > 0 {
				return []interface{}{map[string]interface{}{
					me.Elem: sliceOfMaps,
				}}, nil
			}
			return sliceOfMaps, nil
		} else if me.Unordered {
			seedSlice := []interface{}{}
			switch refElemsSlice.Type().Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				for _, seedElem := range refElemsSlice.Interface().([]int) {
					seedSlice = append(seedSlice, seedElem)
				}
				return schema.NewSet(schema.HashInt, seedSlice), nil
			case reflect.Float32, reflect.Float64:
				for _, seedElem := range refElemsSlice.Interface().([]float64) {
					seedSlice = append(seedSlice, seedElem)
				}
				return schema.NewSet(func(v interface{}) int { return schema.HashString(fmt.Sprintf("%f", v)) }, seedSlice), nil
			case reflect.String:
				for _, seedElem := range refElemsSlice.Interface().([]string) {
					seedSlice = append(seedSlice, seedElem)
				}
				return schema.NewSet(schema.HashString, seedSlice), nil
			}
		}
		return refElemsSlice.Interface(), nil
	case reflect.Int:
		return me.Value.Convert(reflect.TypeOf(intVar)).Interface().(int), nil
	case reflect.Int8:
		return int(me.Value.Convert(reflect.TypeOf(int8Var)).Interface().(int8)), nil
	case reflect.Int16:
		return int(me.Value.Convert(reflect.TypeOf(int16Var)).Interface().(int16)), nil
	case reflect.Int32:
		return int(me.Value.Convert(reflect.TypeOf(int32Var)).Interface().(int32)), nil
	case reflect.Int64:
		return int(me.Value.Convert(reflect.TypeOf(int64Var)).Interface().(int64)), nil
	case reflect.Uint:
		return int(me.Value.Convert(reflect.TypeOf(uintVar)).Interface().(uint)), nil
	case reflect.Uint8:
		return int(me.Value.Convert(reflect.TypeOf(uint8Var)).Interface().(uint8)), nil
	case reflect.Uint16:
		return int(me.Value.Convert(reflect.TypeOf(uint16Var)).Interface().(uint16)), nil
	case reflect.Uint32:
		return int(me.Value.Convert(reflect.TypeOf(uint32Var)).Interface().(uint32)), nil
	case reflect.Uint64:
		return int(me.Value.Convert(reflect.TypeOf(uint64Var)).Interface().(uint64)), nil
	case reflect.Bool:
		return me.Value.Convert(reflect.TypeOf(boolVar)).Interface().(bool), nil
	case reflect.String:
		return me.Value.Convert(reflect.TypeOf(stringVar)).Interface().(string), nil
	case reflect.Float32:
		return float64(me.Value.Convert(reflect.TypeOf(float32Var)).Interface().(float32)), nil
	case reflect.Float64:
		return me.Value.Convert(reflect.TypeOf(float64Var)).Interface().(float64), nil
	default:
		return nil, UnsupportedTypeError{me.Field.Name, me.Value.Type()}
	}
}

// serialTypeOf translates a type (usually the type of a struct field) into the type that needs to be stored within the structure Terraform expects
// Any kind of integer (uint32, int32, ...) will have to be stored as an int
// Any kind of float will have to be float64
// Structs are getting translated into []interface{} - within a Terraform Schema they are represented as a List of Subresources with one single entry
// Slices and Arrays are translated into []interface{}
// Pointer Types are getting automatically dereferenced
// Currently unsupported types are not getting translated. No error will be thrown
func serialTypeOf(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	switch kind := t.Kind(); kind {
	case reflect.Pointer:
		return serialTypeOf(t.Elem())
	case reflect.Slice, reflect.Array:
		return reflect.SliceOf(serialTypeOf(t.Elem()))
	case reflect.String:
		return reflect.TypeOf(stringVar)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.TypeOf(intVar)
	case reflect.Bool:
		return reflect.TypeOf(boolVar)
	case reflect.Float32, reflect.Float64:
		return reflect.TypeOf(float64Var)
	case reflect.Struct:
		return reflect.TypeOf(sliceValue)
	default:
		return t
	}
}
