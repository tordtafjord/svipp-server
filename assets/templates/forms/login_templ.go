// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package forms

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func LoginForm() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-base-200 min-h-screen flex items-center justify-center\"><div class=\"card w-96 bg-base-100 shadow-xl hover:shadow-2xl transition-shadow duration-300\"><div class=\"card-body\"><h1 class=\"text-2xl font-bold mb-6 text-center\">Logg inn</h1><form hx-post=\"api/auth\" hx-target=\"#toasts\" hx-indicator=\"#loading\"><div class=\"form-control mb-4\"><label for=\"email\" class=\"label\"><span class=\"label-text\">Email</span></label> <input type=\"email\" id=\"email\" name=\"email\" required class=\"input input-bordered w-full\"></div><div class=\"form-control mb-6\"><label for=\"password\" class=\"label\"><span class=\"label-text\">Passord</span></label> <input type=\"password\" id=\"password\" name=\"password\" required class=\"input input-bordered w-full\"></div><button type=\"submit\" class=\"btn btn-primary w-full\">Logg inn</button></form><div class=\"mt-6 text-center\"><p class=\"text-sm\">Har du ikke en konto? <a href=\"/signup\" class=\"link link-primary\">Opprett konto her</a></p></div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate