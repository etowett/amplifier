package forms

import "github.com/revel/revel"

type (
	Request struct {
		App     string
		Multi   bool
		Message string
		Count   int
		Times   int
	}
)

func (form *Request) Validate(v *revel.Validation) {
	v.Required(form.App).Message("App required")
	// v.Required(form.Multi).Message("Multi required")
	v.Required(form.Count).Message("Count required")
	v.Required(form.Message).Message("Message required")
	v.Required(form.Times).Message("Times required")
}
