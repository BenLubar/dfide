package raws // import "github.com/BenLubar/dfide/raws"

import (
	"bufio"
	"io"
	"strings"
)

type Writer struct {
	w      *bufio.Writer
	object string
	Indent int
}

func NewWriter(w io.Writer, name, objectType string) (*Writer, error) {
	bw := bufio.NewWriter(w)

	_, err := bw.WriteString(name + "\n\n[OBJECT:" + objectType + "]\n")
	if err != nil {
		return nil, err
	}

	return &Writer{w: bw, object: objectType}, nil
}

func (w *Writer) WriteTag(tag []string) error {
	for _, t := range tag {
		if strings.ContainsAny(t, "[:]\n") {
			return ErrInvalidCharacter
		}
	}

	err := w.w.WriteByte('\n')
	if err != nil {
		return err
	}

	for i := 0; i < w.Indent; i++ {
		err = w.w.WriteByte('\t')
		if err != nil {
			return err
		}
	}

	err = w.w.WriteByte('[')
	if err != nil {
		return err
	}

	_, err = w.w.WriteString(tag[0])
	if err != nil {
		return err
	}

	for _, t := range tag[1:] {
		err = w.w.WriteByte(':')
		if err != nil {
			return err
		}

		_, err = w.w.WriteString(t)
		if err != nil {
			return err
		}
	}

	return w.w.WriteByte(']')
}

func (w *Writer) Flush() error {
	return w.w.Flush()
}
