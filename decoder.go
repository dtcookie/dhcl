package hcl

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type decoder handler

func (me decoder) Decode(ctx context.Context, rd ResourceData, t reflect.Type) (any, error) {
	return me.decode(ctx, rd, t)
}

func setTo(rv reflect.Value, v any) {
	if rv.Type().Kind() == reflect.Pointer {
		referencedValue := reflect.New(rv.Type().Elem()).Elem()
		rv.Set(referencedValue.Addr())
		setTo(referencedValue, v)
		return
	}
	rv.Set(reflect.ValueOf(v).Convert(rv.Type()))
}

func (me decoder) decodePointer(ctx context.Context, rd ResourceData, t reflect.Type) (any, error) {
	sch, err := me.decode(ctx, rd, unref(t))
	if err == nil {
		return sch, err
	}
	if _, ok := err.(UnsupportedTypeError); ok {
		return nil, UnsupportedTypeError{me.Field.Name, t}
	}
	return nil, err
}

func (me decoder) decodeStruct(ctx context.Context, rd ResourceData, t reflect.Type) (any, error) {
	if value, ok := rd.GetOk(me.Property + ".#"); !ok || value == nil || value.(int) == 0 {
		return nil, nil
	}
	targetStructPointer := reflect.New(t).Interface()
	if err := Unmarshal(ctx, &resourceData{parent: rd, prefix: me.Property + ".0"}, targetStructPointer); err != nil {
		return nil, err
	}
	return reflect.ValueOf(targetStructPointer).Elem().Interface(), nil
}

func (me decoder) decodePrimitiveSlice(ctx context.Context, rd ResourceData, t reflect.Type) (any, error) {
	value, ok := rd.GetOk(me.Property)
	if !ok {
		return nil, nil
	}
	targetSlice := reflect.New(t).Elem()
	switch typedValue := value.(type) {
	case []any:
		if me.Unordered {
			fmt.Println("data type is []any, but *schema.Set was expected")
		}
		for _, el := range typedValue {
			targetElem := reflect.New(t.Elem()).Elem()
			setTo(targetElem, el)
			targetSlice = reflect.Append(targetSlice, targetElem)
		}
	case *schema.Set:
		if !me.Unordered {
			fmt.Println("data type is *schema.Set, but []any was expected")
		}
		for _, el := range typedValue.List() {
			targetElem := reflect.New(t.Elem()).Elem()
			setTo(targetElem, el)
			targetSlice = reflect.Append(targetSlice, targetElem)
		}
	}

	return targetSlice.Interface(), nil
}

func (me decoder) decodeStructSlice(ctx context.Context, rd ResourceData, t reflect.Type) (any, error) {
	value, ok := rd.GetOk(me.Property + ".#")
	if !ok || value == nil || value.(int) == 0 {
		return nil, nil
	}
	targetSlice := reflect.New(t).Elem()
	if len(me.Elem) > 0 {
		if me.Unordered {
			v, _ := rd.GetOk(fmt.Sprintf("%s.0.%s", me.Property, me.Elem))
			resourceSet := v.(*schema.Set)
			for _, resource := range resourceSet.List() {
				targetStructPointer := reflect.New(unref(t.Elem())).Interface()
				if err := Unmarshal(ctx, &resourceData{parent: rd, prefix: fmt.Sprintf("%s.0.%s.%d", me.Property, me.Elem, resourceSet.F(resource))}, targetStructPointer); err != nil {
					return nil, err
				}
				targetElem := reflect.New(t.Elem()).Elem()
				setTo(targetElem, reflect.ValueOf(targetStructPointer).Elem().Interface())
				targetSlice = reflect.Append(targetSlice, targetElem)
			}
		} else {
			for idx := 0; idx < value.(int); idx++ {
				targetStructPointer := reflect.New(unref(t.Elem())).Interface()
				if err := Unmarshal(ctx, &resourceData{parent: rd, prefix: fmt.Sprintf("%s.%d", me.Property, idx)}, targetStructPointer); err != nil {
					return nil, err
				}
				targetElem := reflect.New(t.Elem()).Elem()
				setTo(targetElem, reflect.ValueOf(targetStructPointer).Elem().Interface())
				targetSlice = reflect.Append(targetSlice, targetElem)
			}
		}
	} else {
		if me.Unordered {
			v, _ := rd.GetOk(me.Property)
			resourceSet := v.(*schema.Set)
			for _, resource := range resourceSet.List() {
				targetStructPointer := reflect.New(unref(t.Elem())).Interface()
				if err := Unmarshal(ctx, &resourceData{parent: rd, prefix: fmt.Sprintf("%s.%d", me.Property, resourceSet.F(resource))}, targetStructPointer); err != nil {
					return nil, err
				}
				targetElem := reflect.New(t.Elem()).Elem()
				setTo(targetElem, reflect.ValueOf(targetStructPointer).Elem().Interface())
				targetSlice = reflect.Append(targetSlice, targetElem)
			}
		} else {
			for idx := 0; idx < value.(int); idx++ {
				targetStructPointer := reflect.New(unref(t.Elem())).Interface()
				if err := Unmarshal(ctx, &resourceData{parent: rd, prefix: fmt.Sprintf("%s.%d", me.Property, idx)}, targetStructPointer); err != nil {
					return nil, err
				}
				targetElem := reflect.New(t.Elem()).Elem()
				setTo(targetElem, reflect.ValueOf(targetStructPointer).Elem().Interface())
				targetSlice = reflect.Append(targetSlice, targetElem)
			}
		}
	}
	return targetSlice.Interface(), nil
}

func (me decoder) decodeSlice(ctx context.Context, rd ResourceData, t reflect.Type) (any, error) {
	switch elemKind := unref(t.Elem()).Kind(); elemKind {
	case reflect.Float32, reflect.Float64, reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return me.decodePrimitiveSlice(ctx, rd, t)
	case reflect.Struct:
		return me.decodeStructSlice(ctx, rd, t)
	}
	return nil, UnsupportedTypeError{me.Field.Name, t}
}

func (me decoder) decodePrimitive(ctx context.Context, rd ResourceData, t reflect.Type) (any, error) {
	value, ok := rd.GetOk(me.Property)
	if ok {
		rv := reflect.New(t).Elem()
		rv.Set(reflect.ValueOf(value).Convert(t))
		return rv.Interface(), nil
	}
	return nil, nil
}

func (me decoder) decode(ctx context.Context, rd ResourceData, t reflect.Type) (any, error) {
	switch kind := t.Kind(); kind {
	case reflect.Map, reflect.Interface, reflect.Array, reflect.Uintptr, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return nil, UnsupportedTypeError{me.Field.Name, t}
	case reflect.Struct:
		return me.decodeStruct(ctx, rd, t)
	case reflect.Pointer:
		return me.decodePointer(ctx, rd, t)
	case reflect.Slice:
		return me.decodeSlice(ctx, rd, t)
	case reflect.Float32, reflect.Float64, reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return me.decodePrimitive(ctx, rd, t)
	default:
		return nil, UnsupportedTypeError{me.Field.Name, t}
	}
}
