package forms

import (
	"net/url"
	"strings"
)

//Form create a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

//Valid return true if there are no error otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//New initilize a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)

		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}
