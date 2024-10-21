// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package business

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func GenerateApiKeyForm() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex justify-center min-h-screen\"><div class=\"card bg-base-100 shadow-2xl max-w-lg\"><div class=\"card-body\"><h2 class=\"card-title mb-4\">Ny Shopify Konfigurasjon</h2><form class=\"space-y-4\" hx-post=\"api/shopify-api\" hx-target=\"#toasts\" hx-indicator=\"#loading\"><div class=\"form-control\"><label for=\"address\" class=\"label\">Adresse</label> <input type=\"text\" id=\"address\" name=\"address\" required class=\"input input-bordered w-full\"></div><div class=\"flex space-x-4\"><div class=\"form-control w-1/3\"><label for=\"zipCode\" class=\"label\">Postnummer</label> <input type=\"text\" id=\"zipCode\" name=\"zipCode\" required pattern=\"[0-9]*\" inputmode=\"numeric\" class=\"input input-bordered w-full\"></div><div class=\"form-control w-2/3\"><label for=\"city\" class=\"label\">Sted</label> <input type=\"text\" id=\"city\" name=\"city\" required class=\"input input-bordered w-full\"></div></div><div class=\"form-control\"><label class=\"label\" for=\"pickupInstructions\">Hente-Instruksjoner</label> <textarea id=\"pickupInstructions\" class=\"textarea textarea-bordered h-24\"></textarea></div><h3 class=\"card-title mb-4\">Vindu for henting av varer</h3>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, day := range []string{"Mandag", "Tirsdag", "Onsdag", "Torsdag", "Fredag", "Lørdag", "Søndag"} {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"form-control\"><label class=\"label\" for=\"{day}-start\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(day)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `assets/templates/forms/business/generate_api_key.templ`, Line: 35, Col: 56}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</label><div class=\"flex space-x-2\"><input type=\"time\" id=\"{day}Start\" class=\"input input-bordered w-1/2\"> <input type=\"time\" id=\"{day}End\" class=\"input input-bordered w-1/2\"></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card-actions justify-end w-full\"><button class=\"btn btn-primary w-full\">Generer Shopify API-nøkkel</button></div></form></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
