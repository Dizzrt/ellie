package encoding

import (
	"strings"
)

var codecs = make(map[string]Codec)

type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	Name() string
}

func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("can't register a nil codec")
	}

	if codec.Name() == "" {
		panic("can't register a codec with empty name")
	}

	codecType := strings.ToLower(codec.Name())
	codecs[codecType] = codec
}

func GetCodec(codecType string) Codec {
	return codecs[codecType]
}
