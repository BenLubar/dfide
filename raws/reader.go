package raws // import "github.com/BenLubar/dfide/raws"

import (
	"bufio"
	"io"
	"strings"
)

type Reader struct {
	r         *bufio.Reader
	name      string
	object    string
	unreadTag []string
	err       error
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func (r *Reader) init() error {
	if r.err != nil {
		return r.err
	}

	if r.name != "" {
		return nil
	}

	name, err := r.r.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			err = ErrNoRawFileName
		}
		r.err = err
		return err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		r.err = ErrNoRawFileName
		return ErrNoRawFileName
	}

	r.name = name

	tag, err := r.nextTag()
	if err != nil {
		if err == io.EOF || err == ErrIncompleteTag {
			err = ErrNoRawObjectType
		}
		r.err = err
		return err
	}

	if len(tag) != 2 || tag[0] != "OBJECT" {
		r.err = ErrNoRawObjectType
		return ErrNoRawObjectType
	}

	r.object = tag[1]
	return nil
}

func (r *Reader) Name() (string, error) {
	err := r.init()
	return r.name, err
}

func (r *Reader) ObjectType() (string, error) {
	err := r.init()
	return r.object, err
}

func (r *Reader) discardComments() error {
	for {
		ch, _, err := r.r.ReadRune()
		if err != nil {
			r.err = err
			return err
		}

		if ch == '[' {
			if err = r.r.UnreadRune(); err != nil {
				r.err = err
				return err
			}
			return nil
		}
	}
}

func (r *Reader) nextTag() ([]string, error) {
	if r.unreadTag != nil {
		tag := r.unreadTag
		r.unreadTag = nil
		return tag, nil
	}

	if err := r.discardComments(); err != nil {
		return nil, err
	}

	ch, _, err := r.r.ReadRune()
	if err != nil {
		r.err = err
		return nil, err
	}

	if ch != '[' {
		panic("raws: internal error: expected bracket but '" + string(ch) + "' found")
	}

	var tag []string
	var current []rune

	for {
		ch, _, err = r.r.ReadRune()

		if err != nil {
			if err == io.EOF {
				err = ErrIncompleteTag
			}

			r.err = err
			return nil, err
		}

		if ch == ':' {
			tag = append(tag, string(current))
			current = current[:0]
			continue
		}

		if ch == ']' {
			return append(tag, string(current)), nil
		}

		current = append(current, ch)
	}
}
