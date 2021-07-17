package forms

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
)

type errors map[string][]string

//Add error and message to slice
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

//Get the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

//Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This form cannot be blank")
		return false
	}
	return true
}

//MinLength checks the string minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

//IsEmail checks the valid email
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Not valid email")
	}
}
