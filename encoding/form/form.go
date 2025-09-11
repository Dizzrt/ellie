package form

import (
	"net/url"
	"reflect"

	"github.com/Dizzrt/ellie/encoding"
	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
)

var _ encoding.Codec = (*codec)(nil)

type codec struct {
	encoder *form.Encoder
	decoder *form.Decoder
}

const Name = "x-www-form-urlencoded"

var (
	tagName = "json"

	encoder = form.NewEncoder()
	decoder = form.NewDecoder()
)

func init() {
	decoder.SetTagName(tagName)
	encoder.SetTagName(tagName)
	encoding.RegisterCodec(codec{
		encoder: encoder,
		decoder: decoder,
	})
}

func (codec) Name() string {
	return Name
}

func (c codec) Marshal(v any) ([]byte, error) {
	var err error
	var vals url.Values

	if msg, ok := v.(proto.Message); ok {
		vals, err = EncodeValues(msg)
		if err != nil {
			return nil, err
		}
	} else {
		vals, err = c.encoder.Encode(v)
		if err != nil {
			return nil, err
		}
	}

	for k, v := range vals {
		if len(v) == 0 {
			delete(vals, k)
		}
	}

	return []byte(vals.Encode()), nil
}

func (c codec) Unmarshal(data []byte, v any) error {
	vals, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}

		rv = rv.Elem()
	}

	if msg, ok := v.(proto.Message); ok {
		return DecodeValues(msg, vals)
	}

	if msg, ok := rv.Interface().(proto.Message); ok {
		return DecodeValues(msg, vals)
	}

	return c.decoder.Decode(v, vals)
}
