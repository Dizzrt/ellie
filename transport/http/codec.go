package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/Dizzrt/ellie/encoding"
	"github.com/Dizzrt/ellie/errors"
	"github.com/gorilla/mux"
)

type Redirector interface {
	Redirect() (string, int)
}

type Request = http.Request
type Flusher = http.Flusher
type ResponseWriter = http.ResponseWriter

type HTTPCodecRequestDecoder = func(*http.Request, any) error
type HTTPCodecResponseEncoder = func(http.ResponseWriter, *http.Request, any) error
type HTTPCodecErrorEncoder = func(http.ResponseWriter, *http.Request, error)

func getCodecByHeaderName(r *http.Request, name string) (encoding.Codec, bool) {
	for _, accept := range r.Header[name] {
		codec := encoding.GetCodec(contentSubType(accept))
		if codec != nil {
			return codec, true
		}
	}

	return encoding.GetCodec("json"), false
}

func DefaultPathParamsDecoder(r *http.Request, v any) error {
	raws := mux.Vars(r)
	vars := make(url.Values, len(raws))
	for k, v := range raws {
		vars[k] = []string{v}
	}

	return BindQueryParams(vars, v)
}

func DefaultQueryParamsDecoder(r *http.Request, v any) error {
	return BindQueryParams(r.URL.Query(), v)
}

func DefaultRequestBodyDecoder(r *http.Request, v any) error {
	codec, ok := getCodecByHeaderName(r, "Content-Type")
	if !ok {
		// TODO return error
		return nil
	}

	data, err := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	if err != nil {
		// TODO return error
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	if err = codec.Unmarshal(data, v); err != nil {
		// TODO return error
		return nil
	}

	return nil
}

func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v any) error {
	if v == nil {
		return nil
	}

	if rd, ok := v.(Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}

	codec, _ := getCodecByHeaderName(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", contentType(codec.Name()))

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	ee := errors.FromError(err)
	codec, _ := getCodecByHeaderName(r, "Accept")

	body, err := codec.Marshal(ee)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType(codec.Name()))
	// TODO w.WriteHeader(int(ee.Code))
	w.Write(body)
}
