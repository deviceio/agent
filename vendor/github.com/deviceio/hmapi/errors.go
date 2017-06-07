package hmapi

import "fmt"

type FormNotFound struct {
	Resource string
	FormName string
}

func (t *FormNotFound) Error() string {
	return fmt.Sprintf("No such form with name '%v' defined on resource '%v'", t.FormName, t.Resource)
}

type UnsupportedMediaType struct {
	MediaType MediaType
}

func (t *UnsupportedMediaType) Error() string {
	return fmt.Sprintf("%v is not supported", t.MediaType.String())
}
