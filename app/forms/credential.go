package forms

import "github.com/revel/revel"

type (
	Credential struct {
		App      string
		Url      string
		Username string
		Password string
	}
)

func (form *Credential) Validate(v *revel.Validation) {
	v.Required(form.App).Message("App required")
	v.Required(form.Url).Message("Url required")
	v.Required(form.Username).Message("Username required")
	v.Required(form.Password).Message("Password required")
}
