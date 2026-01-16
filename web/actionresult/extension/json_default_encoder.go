package extension

import (
	"github.com/bytedance/sonic/encoder"
	"io"
)

type DefaultJsonEncoder struct {
}

func (jsonEncoder DefaultJsonEncoder) Encode(w io.Writer, data interface{}) error {
	encoder := encoder.NewStreamEncoder(w)
	return encoder.Encode(&data)
}
