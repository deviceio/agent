package system

import "github.com/deviceio/hmapi"
import "net/http"

// System ...
type system struct {
	links   map[string]*hmapi.Link
	forms   map[string]*hmapi.Form
	content map[string]*hmapi.Content
}

// HMAPILinks ...
func (t *system) HMAPILinks() map[string]*hmapi.Link {
	return t.links
}

// HMAPIForms ...
func (t *system) HMAPIForms() map[string]*hmapi.Form {
	return t.forms
}

// HMAPIContent ...
func (t *system) HMAPIContent() map[string]*hmapi.Content {
	return t.content
}

// HMAPIHandleForm ...
func (t *system) HMAPIHandleForm(form string, rw http.ResponseWriter, r *http.Request) {
	switch form {
	case "exec":
		t.exec(rw, r)
	}
}

func (t *system) exec(rw http.ResponseWriter, r *http.Request) {

}
