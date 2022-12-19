package hcl_test

import (
	"context"
	"testing"

	"github.com/dtcookie/assert"
	hcl "github.com/dtcookie/dhcl"
	"github.com/dtcookie/opt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type StringAlias string
type IntAlias int
type Int8Alias int8
type Int16Alias int16
type Int32Alias int32
type Int64Alias int64
type UIntAlias uint
type UInt8Alias uint8
type UInt16Alias uint16
type UInt32Alias uint32
type UInt64Alias uint64
type Float32Alias float32
type Float64Alias float64
type StringSliceAlias []string
type StringPointerAlias *string
type StringPointerSliceAlias []*string
type StringPointerSliceAliasAlias []StringPointerAlias
type Int64PointerAlias *int64
type Int64PointerSliceAlias []*int64
type Int64PointerSliceAliasAlias []Int64PointerAlias

func TestPrimitives(t *testing.T) {
	assert := assert.New(t)
	record := struct {
		Name    string
		Enabled bool
		Int     int
		Uint    uint
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64
	}{
		"1",
		true,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0.0,
		0.0,
	}

	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"name":     record.Name,
		"enabled":  record.Enabled,
		"int":      record.Int,
		"uint":     int(record.Uint),
		"int_8":    int(record.Int8),
		"int_16":   int(record.Int16),
		"int_32":   int(record.Int32),
		"int_64":   int(record.Int64),
		"uint_8":   int(record.Uint8),
		"uint_16":  int(record.Uint16),
		"uint_32":  int(record.Uint32),
		"uint_64":  int(record.Uint64),
		"float_32": float64(record.Float32),
		"float_64": float64(record.Float64),
	}, m, "TestPrimitives failed")
}

func TestAliases(t *testing.T) {
	assert := assert.New(t)
	record := struct {
		StringAlias  StringAlias
		IntAlias     IntAlias
		Int8Alias    Int8Alias
		Int16Alias   Int16Alias
		Int32Alias   Int32Alias
		Int64Alias   Int64Alias
		UIntAlias    UIntAlias
		UInt8Alias   UInt8Alias
		UInt16Alias  UInt16Alias
		UInt32Alias  UInt32Alias
		UInt64Alias  UInt64Alias
		Float64Alias Float64Alias
		Float32Alias Float32Alias
	}{
		StringAlias:  StringAlias("StringAlias"),
		IntAlias:     IntAlias(25),
		Int8Alias:    Int8Alias(26),
		Int16Alias:   Int16Alias(27),
		Int32Alias:   Int32Alias(28),
		Int64Alias:   Int64Alias(29),
		UIntAlias:    UIntAlias(30),
		UInt8Alias:   UInt8Alias(31),
		UInt16Alias:  UInt16Alias(32),
		UInt32Alias:  UInt32Alias(33),
		UInt64Alias:  UInt64Alias(34),
		Float64Alias: Float64Alias(35.0),
		Float32Alias: Float32Alias(36.0),
	}

	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"string_alias":   string(record.StringAlias),
		"int_alias":      int(record.IntAlias),
		"int_8_alias":    int(record.Int8Alias),
		"int_16_alias":   int(record.Int16Alias),
		"int_32_alias":   int(record.Int32Alias),
		"int_64_alias":   int(record.Int64Alias),
		"uint_alias":     int(record.UIntAlias),
		"uint_8_alias":   int(record.UInt8Alias),
		"uint_16_alias":  int(record.UInt16Alias),
		"uint_32_alias":  int(record.UInt32Alias),
		"uint_64_alias":  int(record.UInt64Alias),
		"float_64_alias": float64(record.Float64Alias),
		"float_32_alias": float64(record.Float32Alias),
	}, m, "TestPrimitives failed")
}

func TestZeroPrimitives(t *testing.T) {
	assert := assert.New(t)
	record := struct {
		ZeroInt     int
		ZeroUint    uint
		ZeroInt8    int8
		ZeroInt16   int16
		ZeroInt32   int32
		ZeroInt64   int64
		ZeroUint8   uint8
		ZeroUint16  uint16
		ZeroUint32  uint32
		ZeroUint64  uint64
		ZeroFloat32 float32
		ZeroFloat64 float64
		ZeroString  string
	}{
		ZeroInt:     0,
		ZeroUint:    0,
		ZeroInt8:    0,
		ZeroInt16:   0,
		ZeroInt32:   0,
		ZeroInt64:   0,
		ZeroUint8:   0,
		ZeroUint16:  0,
		ZeroUint32:  0,
		ZeroUint64:  0,
		ZeroFloat32: 0.0,
		ZeroFloat64: 0.0,
		ZeroString:  "",
	}

	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"zero_int":      int(0),
		"zero_uint":     int(0),
		"zero_int_8":    int(0),
		"zero_int_16":   int(0),
		"zero_int_32":   int(0),
		"zero_int_64":   int(0),
		"zero_uint_8":   int(0),
		"zero_uint_16":  int(0),
		"zero_uint_32":  int(0),
		"zero_uint_64":  int(0),
		"zero_float_32": float64(0.0),
		"zero_float_64": float64(0.0),
		"zero_string":   "",
	}, m, "TestZeroPrimitives failed")
}

func TestOmittedZeroPrimitives(t *testing.T) {
	assert := assert.New(t)
	record := struct {
		Foo            string
		OmittedInt     int     `json:"OmittedInt,omitempty"`
		OmittedUint    uint    `json:",omitempty"`
		OmittedInt8    int8    `json:",omitempty"`
		OmittedInt16   int16   `json:",omitempty"`
		OmittedInt32   int32   `json:",omitempty"`
		OmittedInt64   int64   `json:",omitempty"`
		OmittedUint8   uint8   `json:",omitempty"`
		OmittedUint16  uint16  `json:",omitempty"`
		OmittedUint32  uint32  `json:",omitempty"`
		OmittedUint64  uint64  `json:",omitempty"`
		OmittedFloat32 float32 `json:",omitempty"`
		OmittedFloat64 float64 `json:",omitempty"`
		OmittedString  string  `json:",omitempty"`
	}{
		Foo:            "foo",
		OmittedInt:     0,
		OmittedUint:    0,
		OmittedInt8:    0,
		OmittedInt16:   0,
		OmittedInt32:   0,
		OmittedInt64:   0,
		OmittedUint8:   0,
		OmittedUint16:  0,
		OmittedUint32:  0,
		OmittedUint64:  0,
		OmittedFloat32: 0.0,
		OmittedFloat64: 0.0,
		OmittedString:  "",
	}

	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{"foo": "foo"}, m, "TestOmittedZeroPrimitives failed")
}

func TestPointers(t *testing.T) {
	assert := assert.New(t)
	record := struct {
		IntP     *int
		UintP    *uint
		Int8P    *int8
		Int16P   *int16
		Int32P   *int32
		Int64P   *int64
		Uint8P   *uint8
		Uint16P  *uint16
		Uint32P  *uint32
		Uint64P  *uint64
		Float32P *float32
		Float64P *float64
		String   *string
	}{
		IntP:     opt.NewInt(37),
		UintP:    opt.NewUint(38),
		Int8P:    opt.NewInt8(39),
		Int16P:   opt.NewInt16(40),
		Int32P:   opt.NewInt32(41),
		Int64P:   opt.NewInt64(42),
		Uint8P:   opt.NewUInt8(43),
		Uint16P:  opt.NewUInt16(44),
		Uint32P:  opt.NewUInt32(45),
		Uint64P:  opt.NewUInt64(46),
		Float32P: opt.NewFloat32(47.0),
		Float64P: opt.NewFloat64(48.0),
		String:   opt.NewString("49"),
	}

	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"int_p":      int(*record.IntP),
		"uint_p":     int(*record.UintP),
		"int_8_p":    int(*record.Int8P),
		"int_16_p":   int(*record.Int16P),
		"int_32_p":   int(*record.Int32P),
		"int_64_p":   int(*record.Int64P),
		"uint_8_p":   int(*record.Uint8P),
		"uint_16_p":  int(*record.Uint16P),
		"uint_32_p":  int(*record.Uint32P),
		"uint_64_p":  int(*record.Uint64P),
		"float_32_p": float64(*record.Float32P),
		"float_64_p": float64(*record.Float64P),
		"string":     "49",
	}, m, "TestPointers failed")
}

func TestZeroPointers(t *testing.T) {
	assert := assert.New(t)
	record := struct {
		ZeroIntP     *int
		ZeroUintP    *uint
		ZeroInt8P    *int8
		ZeroInt16P   *int16
		ZeroInt32P   *int32
		ZeroInt64P   *int64
		ZeroUint8P   *uint8
		ZeroUint16P  *uint16
		ZeroUint32P  *uint32
		ZeroUint64P  *uint64
		ZeroFloat32P *float32
		ZeroFloat64P *float64
		ZeroStringP  *string
	}{
		ZeroIntP:     nil,
		ZeroUintP:    nil,
		ZeroInt8P:    nil,
		ZeroInt16P:   nil,
		ZeroInt32P:   nil,
		ZeroInt64P:   nil,
		ZeroUint8P:   nil,
		ZeroUint16P:  nil,
		ZeroUint32P:  nil,
		ZeroUint64P:  nil,
		ZeroFloat32P: nil,
		ZeroFloat64P: nil,
		ZeroStringP:  nil,
	}

	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"zero_int_p":      nil,
		"zero_uint_p":     nil,
		"zero_int_8_p":    nil,
		"zero_int_16_p":   nil,
		"zero_int_32_p":   nil,
		"zero_int_64_p":   nil,
		"zero_uint_8_p":   nil,
		"zero_uint_16_p":  nil,
		"zero_uint_32_p":  nil,
		"zero_uint_64_p":  nil,
		"zero_float_32_p": nil,
		"zero_float_64_p": nil,
		"zero_string_p":   nil,
	}, m, "TestZeroPointers failed")
}

func TestOmittedZeroPointers(t *testing.T) {
	assert := assert.New(t)
	record := struct {
		Foo             string
		OmittedIntP     *int     `json:",omitempty"`
		OmittedUintP    *uint    `json:",omitempty"`
		OmittedInt8P    *int8    `json:",omitempty"`
		OmittedInt16P   *int16   `json:",omitempty"`
		OmittedInt32P   *int32   `json:",omitempty"`
		OmittedInt64P   *int64   `json:",omitempty"`
		OmittedUint8P   *uint8   `json:",omitempty"`
		OmittedUint16P  *uint16  `json:",omitempty"`
		OmittedUint32P  *uint32  `json:",omitempty"`
		OmittedUint64P  *uint64  `json:",omitempty"`
		OmittedFloat32P *float32 `json:",omitempty"`
		OmittedFloat64P *float64 `json:",omitempty"`
		OmittedStringP  *string  `json:",omitempty"`
	}{
		Foo:             "foo",
		OmittedIntP:     nil,
		OmittedUintP:    nil,
		OmittedInt8P:    nil,
		OmittedInt16P:   nil,
		OmittedInt32P:   nil,
		OmittedInt64P:   nil,
		OmittedUint8P:   nil,
		OmittedUint16P:  nil,
		OmittedUint32P:  nil,
		OmittedUint64P:  nil,
		OmittedFloat32P: nil,
		OmittedFloat64P: nil,
		OmittedStringP:  nil,
	}

	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"foo": "foo",
	}, m, "TestOmittedZeroPointers failed")
}

func TestUnexported(t *testing.T) {
	assert := assert.New(t)
	record := struct {
		unexported  string
		Unexported2 string `json:"-"`
		Name        string `json:"name" hcl:"name"`
		Enabled     bool   `json:"enabled"`
	}{
		"",
		"",
		"name",
		true,
	}

	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"name":    record.Name,
		"enabled": record.Enabled,
	}, m, "TestUnexported failed")
}

type Queue struct {
}

func TestSlices(t *testing.T) {
	assert := assert.New(t)

	record := struct {
		StringSlice      []string
		StringAliasSlice []StringAlias
		StringSliceAlias StringSliceAlias

		IntSlice []int `hcl:",unordered"`
	}{
		StringSlice:      []string{"s1"},
		StringAliasSlice: []StringAlias{StringAlias("s2")},
		StringSliceAlias: StringSliceAlias{"s3"},
		IntSlice:         []int{3, 7, 9},
	}
	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"string_slice":       []string{"s1"},
		"string_alias_slice": []string{"s2"},
		"string_slice_alias": []string{"s3"},
		"int_slice":          schema.NewSet(schema.HashInt, []interface{}{3, 7, 9}),
	}, m, "TestSlices failed")
}

func TestTags(t *testing.T) {
	assert := assert.New(t)

	record := struct {
		AStringA string `json:"AStringA"`
		BStringB string `json:"BStringB,omitempty"`
		CStringC string `json:"CStringC,omitempty"`
		DStringD string `json:",omitempty"`
		StringA  string `hcl:"string_a"`
		StringB  string `hcl:"string_b,omitempty"`
		StringC  string `hcl:"string_c,omitempty"`
		StringD  string `hcl:",omitempty"`
		StringE  string `json:"StringE,omitempty" hcl:",omitempty"`
	}{
		AStringA: "AStringA",
		BStringB: "",
		CStringC: "CStringC",
		DStringD: "DStringD",
		StringA:  "StringA",
		StringB:  "",
		StringC:  "StringC",
		StringD:  "StringD",
		StringE:  "StringE",
	}
	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"astring_a": "AStringA",
		"cstring_c": "CStringC",
		"dstring_d": "DStringD",
		"string_a":  "StringA",
		"string_c":  "StringC",
		"string_d":  "StringD",
		"string_e":  "StringE",
	}, m, "TestTags failed")
}

func TestStructs(t *testing.T) {
	assert := assert.New(t)

	record := struct {
		Foo string
		Bar struct {
			String string
		} `hcl:"bar_bar"`
	}{
		Foo: "foo",
		Bar: struct {
			String string
		}{
			String: "string",
		},
	}
	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"foo": "foo",
		"bar_bar": []interface{}{map[string]interface{}{
			"string": "string",
		}},
	}, m, "TestStructs failed")
}

func TestStructPointers(t *testing.T) {
	assert := assert.New(t)

	record := struct {
		Foo string
		Bar *struct {
			String string
		} `hcl:"bar_bar"`
	}{
		Foo: "foo",
		Bar: &struct {
			String string
		}{
			String: "string",
		},
	}
	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"foo": "foo",
		"bar_bar": []interface{}{map[string]interface{}{
			"string": "string",
		}},
	}, m, "TestStructPointers failed")
}

func TestZeroStructPointers(t *testing.T) {
	assert := assert.New(t)

	record := struct {
		Foo string
		Bar *struct {
			String string
		} `hcl:"bar_bar,omitempty"`
	}{
		Foo: "foo",
		Bar: nil,
	}
	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"foo": "foo",
	}, m, "TestZeroStructPointers failed")
}

func TestStructSlices(t *testing.T) {
	assert := assert.New(t)

	record := struct {
		Foo string
		Bar []struct {
			String string
		} `hcl:"bar_bar,omitempty,elem=bar_instance"`
	}{
		Foo: "foo",
		Bar: []struct {
			String string
		}{
			{
				String: "string-a",
			},
			{
				String: "string-b",
			},
		},
	}
	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"foo": "foo",
		"bar_bar": []interface{}{
			map[string]interface{}{
				"bar_instance": []interface{}{
					map[string]interface{}{
						"string": "string-a",
					},
					map[string]interface{}{
						"string": "string-b",
					},
				},
			},
		},
	}, m, "TestStructSlices failed")
}

func TestStructPointerSlices(t *testing.T) {
	assert := assert.New(t)

	record := struct {
		Foo string
		Bar []*struct {
			String string
		} `hcl:"bar_bar,omitempty,elem=bar_instance"`
	}{
		Foo: "foo",
		Bar: []*struct {
			String string
		}{
			{
				String: "string-a",
			},
			{
				String: "string-b",
			},
		},
	}
	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"foo": "foo",
		"bar_bar": []interface{}{
			map[string]interface{}{
				"bar_instance": []interface{}{
					map[string]interface{}{
						"string": "string-a",
					},
					map[string]interface{}{
						"string": "string-b",
					},
				},
			},
		},
	}, m, "TestStructPointerSlices failed")
}

type sampleStruct struct {
	String string
}

type sampleStructSlice []sampleStruct

func TestStructSliceAliases(t *testing.T) {
	assert := assert.New(t)

	record := struct {
		Foo string
		Bar sampleStructSlice `hcl:"bar_bar,omitempty,elem=bar_instance"`
	}{
		Foo: "foo",
		Bar: sampleStructSlice{
			{
				String: "string-a",
			},
			{
				String: "string-b",
			},
		},
	}
	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"foo": "foo",
		"bar_bar": []interface{}{
			map[string]interface{}{
				"bar_instance": []interface{}{
					map[string]interface{}{
						"string": "string-a",
					},
					map[string]interface{}{
						"string": "string-b",
					},
				},
			},
		},
	}, m, "TestStructSliceAliases failed")
}

func TestStructSlicesWithoutElem(t *testing.T) {
	assert := assert.New(t)

	record := struct {
		Foo string
		Bar sampleStructSlice `hcl:"bar_bar,omitempty"`
	}{
		Foo: "foo",
		Bar: sampleStructSlice{
			{
				String: "string-a",
			},
			{
				String: "string-b",
			},
		},
	}
	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"foo": "foo",
		"bar_bar": []interface{}{
			map[string]interface{}{
				"string": "string-a",
			},
			map[string]interface{}{
				"string": "string-b",
			},
		},
	}, m, "TestStructSlicesWithoutElem failed")
}

type Base struct {
	Global   string `json:"global"`
	Property string `json:"base_property" hcl:"base_property"`
}

type Derived struct {
	Base
	Property string `json:"derived_property" hcl:"derived_property"`
}

func TestAnonymousFields(t *testing.T) {
	record := Derived{
		Base: Base{
			Global:   "global",
			Property: "base",
		},
		Property: "derived",
	}
	assert := assert.New(t)

	m, err := hcl.Marshal(context.Background(), record)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equalsf(map[string]interface{}{
		"global":           "global",
		"base_property":    "base",
		"derived_property": "derived",
	}, m, "TestAnonymousFields failed")
}
