package http

import (
	"net/url"

	"github.com/Dizzrt/ellie/encoding"
	"github.com/Dizzrt/ellie/encoding/form"
)

func BindQueryParams(vars url.Values, target any) error {
	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(vars.Encode()), target); err != nil {
		// TODO wrap errors
		return err
	}

	return nil
}
