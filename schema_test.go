package hcl_test

import (
	"strconv"
	"testing"

	"github.com/dtcookie/assert"
	hcl "github.com/dtcookie/dhcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSchemaString(t *testing.T) {
	record := struct {
		Name string
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)
	assert.Equalsf(map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: hcl.NoDocumentationAvailable,
			Required:    true,
		},
	}, sch, "TestSchemaString failed")
}

func TestSchemaStringPointer(t *testing.T) {
	record := struct {
		Name *string
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)
	assert.Equalsf(map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: hcl.NoDocumentationAvailable,
			Required:    true,
		},
	}, sch, "TestSchemaStringPointer failed")
}

func TestSchemaStringOptional(t *testing.T) {
	record := struct {
		Name string `json:",omitempty" doc:"Documentation available"`
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)
	assert.Equalsf(map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Documentation available",
			Optional:    true,
		},
	}, sch, "TestSchemaStringOptional failed")
}

func TestSchemaStruct(t *testing.T) {
	record := struct {
		File struct {
			Name string `doc:"Documentation available"`
		} `doc:"Documentation available"`
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)
	assert.Equalsf(map[string]*schema.Schema{
		"file": {
			Type:        schema.TypeList,
			Description: "Documentation available",
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "Documentation available",
						Required:    true,
					},
				},
			},
		},
	}, sch, "TestSchemaStruct failed")
}

func TestSchemaStringSlice(t *testing.T) {
	record := struct {
		Name []string
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)
	assert.Equalsf(map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeList,
			Description: hcl.NoDocumentationAvailable,
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}, sch, "TestSchemaStringSlice failed")
}

func TestSchemaStringPointerSlice(t *testing.T) {
	record := struct {
		Name []*string
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)
	assert.Equalsf(map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeList,
			Description: hcl.NoDocumentationAvailable,
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}, sch, "TestSchemaStringPointerSlice failed")
}

func TestSchemaStringPointerSlicePointer(t *testing.T) {
	record := struct {
		Name *[]*string
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)
	assert.Equalsf(map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeList,
			Description: hcl.NoDocumentationAvailable,
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}, sch, "TestSchemaStringPointerSlicePointer failed")
}

func TestSchemaStructSlice(t *testing.T) {
	record := struct {
		Files []struct {
			Name string `doc:"Documentation available"`
		} `doc:"Documentation available"`
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)
	assert.Equalsf(map[string]*schema.Schema{
		"files": {
			Type:        schema.TypeList,
			Description: "Documentation available",
			Required:    true,
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "Documentation available",
						Required:    true,
					},
				},
			},
		},
	}, sch, "TestSchemaStructSlice failed")
}

func TestSchemaStructSliceElems(t *testing.T) {
	record := struct {
		Files []struct {
			Name string `doc:"Documentation \u0060 available"`
		} `hcl:"elem=file" doc:"Documentation \u0060 available"`
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)

	assert.Equalsf(map[string]*schema.Schema{
		"files": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "Documentation ` available",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"file": {
						Type:        schema.TypeList,
						Description: "Documentation ` available",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Description: "Documentation ` available",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
	}, sch, "TestSchemaStructSliceElems failed")
}

func TestSchemaUnorderedStringSlice(t *testing.T) {
	record := struct {
		Name []string `hcl:",unordered"`
	}{}

	assert := assert.New(t)

	sch := hcl.Schema(record)

	assert.Equalsf(map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeSet,
			Description: hcl.NoDocumentationAvailable,
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}, sch, "TestSchemaUnorderedStringSlice failed")
}

type Item struct {
	Name    string `doc:"The name of the item"`
	Enabled bool   `hcl:",omitempty" doc:"The item is enabled (\u0060true\u0060) or disabled (\u0060false\u0060). Defaults to \u0060false\u0060."`
	Order   int    `doc:"The item order"`
}
type Record struct {
	Name       string   `doc:"The name of the record"`
	Enabled    bool     `hcl:",omitempty" doc:"The record is enabled (\u0060true\u0060) or disabled (\u0060false\u0060). Defaults to \u0060false\u0060."` //
	Items      []Item   `hcl:"items,unordered,omitempty,elem=item" doc:"The list of items"`
	StringList []string `hcl:"strings,omitempty" doc:"A list of strings"`
}

func (me Record) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the record",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The record is enabled (`true`) or disabled (`false`). Defaults to `false`.",
			Optional:    true,
		},
		"strings": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A list of strings",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"items": {
			Type:        schema.TypeList,
			Description: "The list of items",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"item": {
						Type:        schema.TypeSet,
						Description: "The list of items",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Resource{
							Schema: Item{}.Schema(),
						},
					},
				},
			},
		},
	}
}

func (me Item) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the item",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The item is enabled (`true`) or disabled (`false`). Defaults to `false`.",
			Optional:    true,
		},
		"order": {
			Type:        schema.TypeInt,
			Description: "The item order",
			Required:    true,
		},
	}
}

func TestUnmarshalComplexItem(t *testing.T) {
	var record Item
	assert := assert.New(t)

	sch := hcl.Schema(record)
	assert.Equals(record.Schema(), sch)
}

func TestUnmarshalComplexRecord(t *testing.T) {
	var record Record
	assert := assert.New(t)

	sch := hcl.Schema(record)

	assert.Equals(record.Schema(), sch)
}

type SchemaDump struct {
	Type       string
	ConfigMode schema.SchemaConfigMode `json:",omitempty"`
	Required   bool                    `json:",omitempty"`
	Optional   bool                    `json:",omitempty"`
	Computed   bool                    `json:",omitempty"`
	ForceNew   bool                    `json:",omitempty"`
	// DiffSuppressFunc SchemaDiffSuppressFunc
	DiffSuppressOnRefresh bool        `json:",omitempty"`
	Default               interface{} `json:",omitempty"`
	// DefaultFunc SchemaDefaultFunc
	Description  string `json:",omitempty"`
	InputDefault string `json:",omitempty"`
	// StateFunc SchemaStateFunc
	Elem     interface{} `json:",omitempty"`
	MaxItems int         `json:",omitempty"`
	MinItems int         `json:",omitempty"`
	// Set SchemaSetFunc
	// ComputedWhen  []string
	ConflictsWith []string `json:",omitempty"`
	ExactlyOneOf  []string `json:",omitempty"`
	AtLeastOneOf  []string `json:",omitempty"`
	RequiredWith  []string `json:",omitempty"`
	Deprecated    string   `json:",omitempty"`
	// ValidateFunc SchemaValidateFunc
	// ValidateDiagFunc SchemaValidateDiagFunc
	Sensitive bool `json:",omitempty"`
}

func resType(t schema.ValueType) string {
	switch t {
	case schema.TypeBool:
		return "schema.TypeBool"
	case schema.TypeInvalid:
		return "schema.TypeInvalid"
	case schema.TypeInt:
		return "schema.TypeInt"
	case schema.TypeFloat:
		return "schema.TypeFloat"
	case schema.TypeString:
		return "schema.TypeString"
	case schema.TypeList:
		return "schema.TypeList"
	case schema.TypeMap:
		return "schema.TypeMap"
	case schema.TypeSet:
		return "schema.TypeSet"
	}
	return strconv.Itoa(int(t))
}

func (me *SchemaDump) Read(res *schema.Schema) *SchemaDump {
	me.Type = resType(res.Type)
	me.ConfigMode = res.ConfigMode
	me.Required = res.Required
	me.Optional = res.Optional
	me.Computed = res.Computed
	me.ForceNew = res.ForceNew
	me.DiffSuppressOnRefresh = res.DiffSuppressOnRefresh
	me.Description = res.Description
	me.InputDefault = res.InputDefault
	me.MaxItems = res.MaxItems
	me.MinItems = res.MinItems
	// me.ComputedWhen = res.ComputedWhen
	me.ConflictsWith = res.ConflictsWith
	me.ExactlyOneOf = res.ExactlyOneOf
	me.AtLeastOneOf = res.AtLeastOneOf
	me.RequiredWith = res.RequiredWith
	me.Deprecated = res.Deprecated
	me.Sensitive = res.Sensitive

	me.Default = res.Default
	if res.Elem != nil {
		switch elem := res.Elem.(type) {
		case *schema.Resource:
			me.Elem = new(ResourceDump).Read(elem)
		case *schema.Schema:
			me.Elem = new(SchemaDump).Read(elem)
		}
	}

	return me
}

type ResourceDump struct {
	Schema        map[string]*SchemaDump
	SchemaVersion int `json:",omitempty"`
	// MigrateState StateMigrateFunc
	// StateUpgraders []StateUpgrader
	// Create CreateFunc
	// Read ReadFunc
	// Update UpdateFunc
	// Delete DeleteFunc
	// Exists ExistsFunc
	// CreateContext CreateContextFunc
	// ReadContext ReadContextFunc
	// UpdateContext UpdateContextFunc
	// DeleteContext DeleteContextFunc
	// CreateWithoutTimeout CreateContextFunc
	// ReadWithoutTimeout ReadContextFunc
	// UpdateWithoutTimeout UpdateContextFunc
	// DeleteWithoutTimeout DeleteContextFunc
	// CustomizeDiff CustomizeDiffFunc
	// Importer *ResourceImporter
	DeprecationMessage string `json:",omitempty"`
	// Timeouts *ResourceTimeout
	Description   string `json:",omitempty"`
	UseJSONNumber bool   `json:",omitempty"`
}

func (me *ResourceDump) Read(res *schema.Resource) *ResourceDump {
	me.DeprecationMessage = res.DeprecationMessage
	me.Description = res.Description
	if res.Schema != nil {
		me.Schema = map[string]*SchemaDump{}
		for k, v := range res.Schema {
			me.Schema[k] = new(SchemaDump).Read(v)
		}
	}

	me.SchemaVersion = res.SchemaVersion
	me.UseJSONNumber = res.UseJSONNumber
	return me
}
