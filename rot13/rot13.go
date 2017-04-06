package rot13

import (
	"fmt"
	"io"
	"strings"
)

const inputChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const outputChars = "NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm"

type EncodingError byte

func (e EncodingError) Error() string {
	return fmt.Sprintf("Cannot code %s in rot13", string(e))
}

func Decode(b byte) (byte, error) {
	i := strings.IndexByte(outputChars, b)
	if i == -1 {
		return byte(0), EncodingError(b)
	}
	return inputChars[i], nil
}

func Encode(b byte) (byte, error) {
	i := strings.IndexByte(inputChars, b)
	if i == -1 {
		return byte(0), EncodingError(b)
	}
	return outputChars[i], nil
}

type Reader struct {
	r io.Reader
}

func (r *Reader) Read(p []byte) (int, error) {
	t := make([]byte, len(p))
	n, err_read := r.r.Read(t)
	if err_read != nil {
		return n, err_read
	}
	for i := 0; i < n; i++ {
		var err_enc error
		p[i], err_enc = Encode(t[i])
		if err_enc != nil {
			p[i] = t[i]
		}
	}
	return n, nil
}

type Writer struct {
	w io.Writer
}

func (w *Writer) Write(p []byte) (int, error) {
	es := make([]byte, len(p))
	for i, b := range p {
		eb, err := Encode(b)
		if err == nil {
			es[i] = eb
		} else {
			es[i] = b
		}
	}
	return w.w.Write(es)
}
