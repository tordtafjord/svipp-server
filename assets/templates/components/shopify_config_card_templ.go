// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func ShopifyConfigCard() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-col items-center p-4\"><div id=\"storeConfig\" class=\"card bg-base-100 w-96 shadow-xl my-2\"><div class=\"card-body\"><h2 class=\"card-title\">Store Name</h2><p><span class=\"font-semibold\">Address:</span> address 22, zip</p><div class=\"grid grid-cols-[auto,1fr] gap-x-2\"><span class=\"font-semibold\">Idag's Hentevindu:</span> <span>10:00 - 20:00</span> <span class=\"font-semibold\">Imorgen's Hentevindu:</span> <span>10:00 - 20:00</span></div><div class=\"card-actions justify-end\"><!-- You can add a button here if needed, or remove this div if not --></div></div></div><div id=\"storeConfig\" class=\"card bg-base-100 w-96 shadow-xl my-2\"><div class=\"card-body\"><h2 class=\"card-title\">Store Name</h2><p><span class=\"font-semibold\">Address:</span> address 22, zip</p><div class=\"grid grid-cols-[auto,1fr] gap-x-2\"><span class=\"font-semibold\">Idag's Hentevindu:</span> <span>10:00 - 20:00</span> <span class=\"font-semibold\">Imorgen's Hentevindu:</span> <span>10:00 - 20:00</span></div><div class=\"card-actions justify-end\"><!-- You can add a button here if needed, or remove this div if not --></div></div></div><a hx-get=\"/create-shopify-config\" hx-target=\"#mainContent\" hx-swap=\"innerHTML\" hx-push-url=\"true\" hx-trigger=\"click\" class=\"btn btn-primary w-96 my-2\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-6 w-6 mr-2\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M12 4v16m8-8H4\"></path></svg> Ny Shopify Integrasjon</a></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
