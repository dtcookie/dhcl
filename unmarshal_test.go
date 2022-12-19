package hcl_test

import (
	"context"
	"testing"

	"github.com/dtcookie/assert"
	hcl "github.com/dtcookie/dhcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TestingResourceData map[string]any

func (me TestingResourceData) GetOk(key string) (interface{}, bool) {
	res, ok := me[key]
	// fmt.Println("GetOk", key, "=>", res, ok)
	return res, ok
}

func addr[T any](v T) *T {
	return &v
}

func TestUnmarshalPrimitives(t *testing.T) {

	type Record struct {
		Name               string      `json:"name"`
		StringAlias        StringAlias `json:"string_alias_x"`
		Int                int64
		Intp               *int
		MissingIntp        *int
		StringAliasP       *StringAlias
		StringPointerAlias StringPointerAlias
	}

	var record Record
	assert := assert.New(t)
	assert.Success(hcl.Unmarshal(context.Background(), TestingResourceData{
		"name":                 "name-value",
		"string_alias_x":       "string-alias",
		"string_alias_p":       "string-alias-p",
		"int":                  13,
		"intp":                 14,
		"string_pointer_alias": "string-pointer-alias",
	}, &record))

	assert.Equals(Record{
		Name:               "name-value",
		StringAlias:        StringAlias("string-alias"),
		Int:                13,
		Intp:               addr(14),
		StringAliasP:       addr(StringAlias("string-alias-p")),
		StringPointerAlias: StringPointerAlias(addr("string-pointer-alias")),
	}, record)
}

func TestUnmarshalStringSlice(t *testing.T) {
	type Record struct{ Names []string }

	var record Record

	assert := assert.New(t)
	assert.Success(hcl.Unmarshal(context.Background(), TestingResourceData{"names": []any{"name-value"}}, &record))

	assert.Equals(Record{Names: []string{"name-value"}}, record)
}

func TestUnmarshalStringSlicePointer(t *testing.T) {
	record := struct {
		Names *[]string
	}{}
	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"names": []any{"name-value"},
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(struct {
		Names *[]string
	}{
		Names: addr([]string{"name-value"}),
	}, record)
}

func TestUnmarshalStringPointerSlice(t *testing.T) {
	record := struct {
		Names []*string
	}{}
	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"names": []any{"name-value"},
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(struct {
		Names []*string
	}{
		Names: []*string{addr("name-value")},
	}, record)
}

func TestUnmarshalStringPointerAliasSlice(t *testing.T) {
	record := struct {
		Names []StringPointerAlias
	}{}
	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"names": []any{"name-value"},
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(struct {
		Names []StringPointerAlias
	}{
		Names: []StringPointerAlias{StringPointerAlias(addr("name-value"))},
	}, record)
}

func TestUnmarshalStringPointerSliceAlias(t *testing.T) {
	record := struct {
		Names StringPointerSliceAlias
	}{}
	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"names": []any{"name-value"},
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(struct {
		Names StringPointerSliceAlias
	}{
		Names: StringPointerSliceAlias{StringPointerAlias(addr("name-value"))},
	}, record)
}

func TestStringPointerSliceAliasAlias(t *testing.T) {
	record := struct {
		Names StringPointerSliceAliasAlias
	}{}
	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"names": []any{"name-value"},
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(struct {
		Names StringPointerSliceAliasAlias
	}{
		Names: StringPointerSliceAliasAlias{StringPointerAlias(addr("name-value"))},
	}, record)
}

func TestUnmarshalStringPointerPointerSlice(t *testing.T) {
	record := struct {
		Names []**string
	}{}
	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"names": []any{"name-value"},
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(struct {
		Names []**string
	}{
		Names: []**string{addr(addr("name-value"))},
	}, record)
}

func TestUnmarshalInt64Slice(t *testing.T) {
	record := struct{ Names []int64 }{}
	if err := hcl.Unmarshal(context.Background(), TestingResourceData{"names": []any{42}}, &record); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(struct{ Names []int64 }{Names: []int64{42}}, record)
}

func TestUnmarshalInt64SlicePointer(t *testing.T) {
	record := struct{ Names *[]int64 }{}
	if err := hcl.Unmarshal(context.Background(), TestingResourceData{"names": []any{42}}, &record); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(struct{ Names *[]int64 }{Names: addr([]int64{42})}, record)
}

func TestUnmarshalInt64PointerSlice(t *testing.T) {
	record := struct{ Names []*int64 }{}
	if err := hcl.Unmarshal(context.Background(), TestingResourceData{"names": []any{42}}, &record); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(struct{ Names []*int64 }{Names: []*int64{addr(int64(42))}}, record)
}

func TestUnmarshalInt64PointerAliasSlice(t *testing.T) {
	record := struct{ Names []Int64PointerAlias }{}
	if err := hcl.Unmarshal(context.Background(), TestingResourceData{"names": []any{42}}, &record); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(struct{ Names []Int64PointerAlias }{
		Names: []Int64PointerAlias{Int64PointerAlias(addr(int64(42)))},
	}, record)
}

func TestUnmarshalInt64PointerSliceAlias(t *testing.T) {
	record := struct{ Names Int64PointerSliceAlias }{}
	if err := hcl.Unmarshal(context.Background(), TestingResourceData{"names": []any{42}}, &record); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(struct{ Names Int64PointerSliceAlias }{Names: Int64PointerSliceAlias([]*int64{addr(int64(42))})}, record)
}

func TestInt64PointerSliceAliasAlias(t *testing.T) {
	record := struct {
		Names Int64PointerSliceAliasAlias
	}{}
	if err := hcl.Unmarshal(context.Background(), TestingResourceData{"names": []any{42}}, &record); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(struct {
		Names Int64PointerSliceAliasAlias
	}{
		Names: Int64PointerSliceAliasAlias([]Int64PointerAlias{Int64PointerAlias(addr(int64(42)))}),
	}, record)
}

func TestUnmarshalInt64PointerPointerSlice(t *testing.T) {
	record := struct{ Names []**int64 }{}
	if err := hcl.Unmarshal(context.Background(), TestingResourceData{"names": []any{42}}, &record); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(struct{ Names []**int64 }{Names: []**int64{addr(addr(int64(42)))}}, record)
}

func TestUnmarshalStruct(t *testing.T) {
	record := struct {
		Item struct {
			Name string
		}
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"item.#":      1,
			"item.0.name": "42",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(
		struct {
			Item struct {
				Name string
			}
		}{
			Item: struct {
				Name string
			}{
				Name: "42",
			},
		},
		record,
	)
}

func TestUnmarshalStructPointer(t *testing.T) {
	record := struct {
		Item *struct {
			Name string
		}
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"item.#":      1,
			"item.0.name": "42",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(
		struct {
			Item *struct {
				Name string
			}
		}{
			Item: &struct {
				Name string
			}{
				Name: "42",
			},
		},
		record,
	)
}

func TestUnmarshalNilStructPointer(t *testing.T) {
	record := struct{ Item *struct{ Name string } }{}
	if err := hcl.Unmarshal(context.Background(), TestingResourceData{"item.#": 0}, &record); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(struct{ Item *struct{ Name string } }{}, record)

	record = struct{ Item *struct{ Name string } }{}
	if err := hcl.Unmarshal(context.Background(), TestingResourceData{}, &record); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(struct{ Item *struct{ Name string } }{}, record)
}

func TestUnmarshalStructSlice(t *testing.T) {
	record := struct {
		Items []struct {
			Name string
		}
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"items.#":      1,
			"items.0.name": "42",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(
		struct {
			Items []struct {
				Name string
			}
		}{
			Items: []struct {
				Name string
			}{
				{
					Name: "42",
				},
			},
		},
		record,
	)
}

func TestUnmarshalStructPointerSlice(t *testing.T) {
	record := struct {
		Items []*struct {
			Name string
		}
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"items.#":      1,
			"items.0.name": "42",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(
		struct {
			Items []*struct {
				Name string
			}
		}{
			Items: []*struct {
				Name string
			}{
				{
					Name: "42",
				},
			},
		},
		record,
	)
}

func TestUnmarshalStructSlicePointer(t *testing.T) {
	record := struct {
		Items *[]struct {
			Name string
		}
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"items.#":      1,
			"items.0.name": "42",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	items := []struct {
		Name string
	}{
		{
			Name: "42",
		},
	}
	assert.New(t).Equals(
		struct {
			Items *[]struct {
				Name string
			}
		}{
			Items: &items,
		},
		record,
	)
}

type structPointerAlias *struct {
	Name string
}

func TestUnmarshalStructPointerAliasSlice(t *testing.T) {
	record := struct {
		Items []structPointerAlias
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"items.#":      1,
			"items.0.name": "42",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(
		struct {
			Items []structPointerAlias
		}{
			Items: []structPointerAlias{
				{
					Name: "42",
				},
			},
		},
		record,
	)
}

type structSliceAlias []struct {
	Name string
}

func TestUnmarshalStructSliceAlias(t *testing.T) {
	record := struct {
		Items structSliceAlias
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"items.#":      1,
			"items.0.name": "42",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(
		struct {
			Items structSliceAlias
		}{
			Items: structSliceAlias{
				{
					Name: "42",
				},
			},
		},
		record,
	)

	record = struct {
		Items structSliceAlias
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"items.#":      2,
			"items.0.name": "42",
			"items.1.name": "24",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(
		struct {
			Items structSliceAlias
		}{
			Items: structSliceAlias{
				{
					Name: "42",
				},
				{
					Name: "24",
				},
			},
		},
		record,
	)
}

func TestUnmarshalStructSliceAliasPointer(t *testing.T) {
	record := struct {
		Items *structSliceAlias
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"items.#":      1,
			"items.0.name": "42",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(
		struct {
			Items *structSliceAlias
		}{
			Items: &structSliceAlias{
				{
					Name: "42",
				},
			},
		},
		record,
	)

	record = struct {
		Items *structSliceAlias
	}{}
	if err := hcl.Unmarshal(context.Background(),
		TestingResourceData{
			"items.#":      2,
			"items.0.name": "42",
			"items.1.name": "24",
		},
		&record,
	); err != nil {
		t.Error(err)
	}
	assert.New(t).Equals(
		struct {
			Items *structSliceAlias
		}{
			Items: &structSliceAlias{
				{
					Name: "42",
				},
				{
					Name: "24",
				},
			},
		},
		record,
	)
}

func TestUnmarshalStringSet(t *testing.T) {
	record := struct {
		Names []string `hcl:"foos,unordered"`
	}{}
	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"foos": schema.NewSet(schema.HashString, []any{"name-value"}),
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(struct {
		Names []string `hcl:"foos,unordered"`
	}{
		Names: []string{"name-value"},
	}, record)
}

func TestUnmarshalStructSet(t *testing.T) {
	type Name struct {
		Value string
	}
	type Record struct {
		Names []Name `hcl:"foos,unordered"`
	}
	var record Record

	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"foos.#":         1,
		"foos.666.value": "name-value",
		"foos":           schema.NewSet(func(interface{}) int { return 666 }, []any{map[string]any{"value": "name-value"}}),
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(
		Record{Names: []Name{{Value: "name-value"}}},
		record,
	)
}

func TestUnmarshalStructElemSet(t *testing.T) {
	type Name struct {
		Value string
	}
	type Record struct {
		Names []Name `hcl:"foos,elem=foo,unordered"`
	}
	var record Record

	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"foos.#":               1,
		"foos.0.foo.#":         1,
		"foos.0.foo.666.value": "name-value",
		"foos.0.foo":           schema.NewSet(func(interface{}) int { return 666 }, []any{map[string]any{"value": "name-value"}}),
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(
		Record{Names: []Name{{Value: "name-value"}}},
		record,
	)
}

func TestUnmarshalStructPointerSet(t *testing.T) {
	type Name struct {
		Value string
	}
	type Record struct {
		Names []*Name `hcl:"foos,unordered"`
	}
	var record Record

	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"foos.#":         1,
		"foos.666.value": "name-value",
		"foos":           schema.NewSet(func(interface{}) int { return 666 }, []any{map[string]any{"value": "name-value"}}),
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(
		Record{Names: []*Name{{Value: "name-value"}}},
		record,
	)
}

func TestUnmarshalNilStructPointer2(t *testing.T) {
	type Name struct {
		Value string
	}
	type Record struct {
		Names  *Name `hcl:"foos,unordered"`
		String string
	}
	var record Record

	err := hcl.Unmarshal(context.Background(), TestingResourceData{
		"foos.#": 0,
		"string": "some-string",
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert := assert.New(t)

	assert.Equals(
		Record{String: "some-string"},
		record,
	)

	record = Record{}
	err = hcl.Unmarshal(context.Background(), TestingResourceData{
		"string": "some-string",
	}, &record)
	if err != nil {
		t.Error(err)
	}
	assert.Equals(
		Record{String: "some-string"},
		record,
	)
}
