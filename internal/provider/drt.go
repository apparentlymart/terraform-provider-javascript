package provider

import (
	"context"
	"fmt"

	tfsdk "github.com/apparentlymart/terraform-sdk"
	"github.com/apparentlymart/terraform-sdk/tfschema"
	"github.com/dop251/goja"
	"github.com/zclconf/go-cty-goja/ctygoja"
	"github.com/zclconf/go-cty/cty"
)

type drt struct {
	Source string    `cty:"source"`
	Vars   cty.Value `cty:"vars"`

	Result cty.Value `cty:"result"`
}

func javascriptDataResourceType() tfsdk.DataResourceType {
	return tfsdk.NewDataResourceType(&tfsdk.ResourceTypeDef{
		ConfigSchema: &tfschema.BlockType{
			Attributes: map[string]*tfschema.Attribute{
				"source": {
					Type:     cty.String,
					Required: true,
					ValidateFn: func(source string) tfsdk.Diagnostics {
						var diags tfsdk.Diagnostics

						_, err := goja.Compile("", source, true)
						if err != nil {
							diags = diags.Append(tfsdk.Diagnostic{
								Severity: tfsdk.Error,
								Summary:  "JavaScript Syntax Error",
								Detail:   err.Error(),
								Path:     cty.GetAttrPath("source"),
							})
						}

						return diags
					},
				},
				"vars": {
					Type:     cty.DynamicPseudoType,
					Optional: true,
					ValidateFn: func(v cty.Value) tfsdk.Diagnostics {
						var diags tfsdk.Diagnostics
						if v.IsNull() {
							return diags
						}

						ty := v.Type()
						if !(ty.IsObjectType() || ty.IsMapType()) {
							diags = diags.Append(tfsdk.ValidationError(fmt.Errorf("vars must be a mapping describing variables to include in the global scope")))
						}

						return diags
					},
				},
				"result": {Type: cty.DynamicPseudoType, Computed: true},
			},
		},

		ReadFn: func(ctx context.Context, client *client, obj *drt) (*drt, tfsdk.Diagnostics) {
			var diags tfsdk.Diagnostics

			prog, err := goja.Compile("", obj.Source, true)
			if err != nil {
				diags = diags.Append(tfsdk.Diagnostic{
					Severity: tfsdk.Error,
					Summary:  "JavaScript Syntax Error",
					Detail:   err.Error(),
					Path:     cty.GetAttrPath("source"),
				})
				return nil, diags
			}

			js := goja.New()
			installUnderscore(js)

			if !obj.Vars.IsNull() {
				vars := obj.Vars.AsValueMap()
				for name, val := range vars {
					valJS := ctygoja.FromCtyValue(val, js)
					js.Set(name, valJS)
				}
			}

			resultJS, err := js.RunProgram(prog)
			if err != nil {
				diags = diags.Append(tfsdk.Diagnostic{
					Severity: tfsdk.Error,
					Summary:  "JavaScript Runtime Error",
					Detail:   err.Error(),
					Path:     cty.GetAttrPath("source"),
				})
				return nil, diags
			}

			result, err := ctygoja.ToCtyValue(resultJS, js)
			if err != nil {
				diags = diags.Append(tfsdk.Diagnostic{
					Severity: tfsdk.Error,
					Summary:  "Invalid result from JavaScript",
					Detail:   fmt.Sprintf("The program produced a result that isn't compatible with Terraform: %s.", err),
					Path:     cty.GetAttrPath("source"),
				})
				return nil, diags
			}

			obj.Result = result

			return obj, diags
		},
	})
}
