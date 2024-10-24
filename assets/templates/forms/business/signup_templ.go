// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package business

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func SignupForm() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"min-h-screen bg-base-200 flex items-center justify-center p-4\"><div class=\"card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow duration-300 w-full max-w-4xl\"><div class=\"card-body\"><h1 class=\"card-title text-2xl justify-center mb-6\">Opprett Bedriftskonto</h1><form hx-post=\"api/business\" hx-target=\"#toasts\" hx-indicator=\"#loading\"><div class=\"flex flex-col md:flex-row md:space-x-8\"><!-- Left pane: Personal details --><div class=\"w-full md:w-1/2\"><h2 class=\"text-xl font-semibold mb-4\">Kontaktperson</h2><div class=\"flex space-x-4\"><div class=\"w-1/2\"><label for=\"firstName\" class=\"label\">Fornavn</label> <input type=\"text\" id=\"firstName\" name=\"firstName\" required class=\"input input-bordered w-full\"></div><div class=\"w-1/2\"><label for=\"lastName\" class=\"label\">Etternavn</label> <input type=\"text\" id=\"lastName\" name=\"lastName\" required class=\"input input-bordered w-full\"></div></div><label for=\"email\" class=\"label\">Email</label> <input type=\"email\" id=\"email\" name=\"email\" required class=\"input input-bordered w-full\"> <label for=\"phone\" class=\"label\">Telefon</label><div class=\"join\"><div class=\"join-item bg-base-200 px-3 flex items-center\">+</div><input type=\"text\" id=\"countryCode\" name=\"countryCode\" value=\"47\" class=\"join-item input input-bordered w-16\" pattern=\"[0-9]*\" inputmode=\"numeric\"> <input type=\"tel\" id=\"phone\" name=\"phone\" required class=\"join-item input input-bordered flex-1\" placeholder=\"987 65 432\"></div><label for=\"password\" class=\"label\">Passord</label> <input type=\"password\" id=\"password\" name=\"password\" required class=\"input input-bordered w-full\"> <label for=\"confirmPassword\" class=\"label\">Bekreft Passord</label> <input type=\"password\" id=\"confirmPassword\" name=\"confirmPassword\" required class=\"input input-bordered w-full\"></div><!-- Right pane: Business details --><div class=\"w-full md:w-1/2 mt-6 md:mt-0\"><h2 class=\"text-xl font-semibold mb-4\">Bedriftsinformasjon</h2><label for=\"companyName\" class=\"label\">Bedriftsnavn</label> <input type=\"text\" id=\"companyName\" name=\"companyName\" required class=\"input input-bordered w-full\"> <label for=\"orgNumber\" class=\"label\">Organisasjonsnummer</label> <input type=\"text\" id=\"orgNumber\" name=\"orgNumber\" required class=\"input input-bordered w-full\"> <label for=\"businessAddress\" class=\"label\">Forretningsadresse</label> <input type=\"text\" id=\"businessAddress\" name=\"businessAddress\" required class=\"input input-bordered w-full\"><div class=\"flex space-x-4\"><div class=\"w-1/3\"><label for=\"zipCode\" class=\"label\">Postnummer</label> <input type=\"text\" id=\"zipCode\" name=\"zipCode\" required pattern=\"[0-9]*\" inputmode=\"numeric\" class=\"input input-bordered w-full\"></div><div class=\"w-2/3\"><label for=\"city\" class=\"label\">Sted</label> <input type=\"text\" id=\"city\" name=\"city\" required class=\"input input-bordered w-full\"></div></div></div></div><!-- Submit button --><div class=\"mt-6\"><button type=\"submit\" class=\"btn btn-primary w-full\">Opprett Konto</button></div></form><div class=\"mt-6 text-center\"><p class=\"text-sm text-base-content/70\">Har du allerede en konto? <a href=\"/login\" class=\"link link-primary\">Logg inn</a></p></div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
