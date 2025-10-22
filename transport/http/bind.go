package http

import (
	"net/url"

	"github.com/dizzrt/ellie/encoding"
	"github.com/dizzrt/ellie/encoding/form"
)

func BindQueryParams(vars url.Values, target any) error {
	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(vars.Encode()), target); err != nil {
		// TODO wrap errors
		return err
	}

	return nil
}
