package hcl_test

import (
	"context"
	"hash/crc32"
	"testing"

	"github.com/dtcookie/assert"
	hcl "github.com/dtcookie/dhcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HTTPHeaders []*HTTPHeader

type HTTPHeader struct {
	Name        string  `json:"name" doc:"The name of the HTTP header"`                                                                                                                                                                                                                            // The name of the HTTP header
	Secret      bool    `json:"secret" hcl:"-"`                                                                                                                                                                                                                                                    // The value of this HTTP header is a secret (`true`) or not (`false`).
	Value       *string `json:"value,omitempty" doc:"The value of the HTTP header. May contain an empty value. \u0060secret_value\u0060 and \u0060value\u0060 are mutually exclusive. Only one of those two is allowed to be specified."`                                                          // The value of the HTTP header. May contain an empty value
	SecretValue *string `json:"secretValue,omitempty" hcl:",sensitive,omitempty" doc:"The value of the HTTP header as a sensitive property. May contain an empty value. \u0060secret_value\u0060 and \u0060value\u0060 are mutually exclusive. Only one of those two is allowed to be specified."` // The secret value of the HTTP header. May contain an empty value
}

func hashString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func (me *HTTPHeader) HashCode() int {
	return hashString(me.Name)
}

func (me *HTTPHeader) UnmarshalHCL(ctx context.Context, decoder hcl.ResourceData) error {
	if err := hcl.Unmarshal(ctx, decoder, me); err != nil {
		return err
	}
	return nil
}

type XMatters struct {
	Enabled   bool   `json:"-" hcl:"active" doc:"The configuration is enabled (\u0060true\u0060) or disabled (\u0060false\u0060)"`
	Name      string `json:"-" hcl:"name" doc:"The name of the notification configuration"`
	ProfileID string `json:"-" hcl:"profile" doc:"The ID of the associated alerting profile"`

	URL      string      `json:"url" doc:"The URL of the WebHook endpoint"`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             // The URL of the xMatters webhook
	Insecure bool        `json:"acceptAnyCertificate" hcl:"insecure,omitempty" doc:"Accept any, including self-signed and invalid, SSL certificate (\u0060true\u0060) or only trusted (\u0060false\u0060) certificates"`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                // Accept any SSL certificate (including self-signed and invalid certificates)
	Headers  HTTPHeaders `json:"headers,omitempty" hcl:",omitempty,unordered,elem=header" doc:"A list of the additional HTTP headers"`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  // Additional HTTP headers
	Payload  string      `json:"payload" doc:"The content of the notification message. You can use the following placeholders:  * \u0060{ImpactedEntities}\u0060: Details about the entities impacted by the problem in form of a JSON array.  * \u0060{ImpactedEntity}\u0060: The entity impacted by the problem or *X* impacted entities.  * \u0060{PID}\u0060: The ID of the reported problem.  * \u0060{ProblemDetailsHTML}\u0060: All problem event details, including root cause, as an HTML-formatted string.  * \u0060{ProblemDetailsJSON}\u0060: All problem event details, including root cause, as a JSON object.  * \u0060{ProblemDetailsMarkdown}\u0060: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * \u0060{ProblemDetailsText}\u0060: All problem event details, including root cause, as a text-formatted string.  * \u0060{ProblemID}\u0060: The display number of the reported problem.  * \u0060{ProblemImpact}\u0060: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are \u0060APPLICATION\u0060, \u0060SERVICE\u0060, and \u0060INFRASTRUCTURE\u0060.  * \u0060{ProblemSeverity}\u0060: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are \u0060AVAILABILITY\u0060, \u0060ERROR\u0060, \u0060PERFORMANCE\u0060, \u0060RESOURCE_CONTENTION\u0060, and \u0060CUSTOM_ALERT\u0060.  * \u0060{ProblemTitle}\u0060: A short description of the problem.  * \u0060{ProblemURL}\u0060: The URL of the problem within Dynatrace.  * \u0060{State}\u0060: The state of the problem. Possible values are \u0060OPEN\u0060 and \u0060RESOLVED\u0060.  * \u0060{Tags}\u0060: The list of tags that are defined for all impacted entities, separated by commas"` // The content of the notification message. Type '{' for placeholder suggestions
}

func TestXMatters(t *testing.T) {
	assert := assert.New(t)

	record := XMatters{}
	sch := hcl.Schema(record)

	// dump := &ResourceDump{}

	// dump.Read(&schema.Resource{
	// 	Schema: sch,
	// })
	// data, _ := json.Marshal(dump)
	// t.Log(string(data))

	headersSchema := sch["headers"]
	headersResource := headersSchema.Elem.(*schema.Resource)
	headerSchema := headersResource.Schema["header"]

	expSch := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the notification configuration",
			Required:    true,
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "The configuration is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile",
			Required:    true,
		},

		"url": {
			Type:        schema.TypeString,
			Description: "The URL of the WebHook endpoint",
			Required:    true,
		},
		"insecure": {
			Type:        schema.TypeBool,
			Description: "Accept any, including self-signed and invalid, SSL certificate (`true`) or only trusted (`false`) certificates",
			Optional:    true,
		},
		"headers": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "A list of the additional HTTP headers",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"header": {
						Type:        schema.TypeSet,
						Required:    true,
						MinItems:    1,
						Description: "A list of the additional HTTP headers",
						Set:         headerSchema.Set,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Description: "The name of the HTTP header",
									Required:    true,
								},
								"secret_value": {
									Type:        schema.TypeString,
									Description: "The value of the HTTP header as a sensitive property. May contain an empty value. `secret_value` and `value` are mutually exclusive. Only one of those two is allowed to be specified.",
									Sensitive:   true,
									Optional:    true,
								},
								"value": {
									Type:        schema.TypeString,
									Description: "The value of the HTTP header. May contain an empty value. `secret_value` and `value` are mutually exclusive. Only one of those two is allowed to be specified.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
		},
		"payload": {
			Type:        schema.TypeString,
			Description: "The content of the notification message. You can use the following placeholders:  * `{ImpactedEntities}`: Details about the entities impacted by the problem in form of a JSON array.  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsHTML}`: All problem event details, including root cause, as an HTML-formatted string.  * `{ProblemDetailsJSON}`: All problem event details, including root cause, as a JSON object.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
	}
	assert.Equalsf(
		expSch["headers"].Elem.(*schema.Resource).Schema["header"].Set,
		sch["headers"].Elem.(*schema.Resource).Schema["header"].Set,
		"TestXMatters failed",
	)

	assert.Success(hcl.Unmarshal(context.Background(), TestingResourceData{
		"active":                            true,
		"name":                              "x-matters-name",
		"profile":                           "x-matters-profile",
		"url":                               "x-matters-url",
		"insecure":                          true,
		"payload":                           "x-matters-payload",
		"headers.#":                         1,
		"headers.0.header.#":                1,
		"headers.0.header.3062169091.value": "header-0-value",
		"headers.0.header.3062169091.name":  "header-0-name",
		// "headers.0.header":           schema.NewSet(func(interface{}) int { return 666 }, []any{map[string]any{"name": "header-0-name", "value": "header-0-value"}}),
		"headers.0.header": schema.NewSet(headerSchema.Set, []any{map[string]any{"name": "header-0-name", "value": "header-0-value"}}),
	}, &record))
}
