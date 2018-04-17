package raws

import (
	"io"
	"strings"
	"testing"
)

func TestReader_Valid(t *testing.T) {
	t.Parallel()

	r := NewReader(strings.NewReader(`foo_bar
[OBJECT:FOO]
[FOO:BAR][BAZ]`))
	if err := r.init(); err != nil {
		t.Fatalf("Reader.init: %+v", err)
	}

	if r.name != "foo_bar" {
		t.Errorf("expected name to be foo_bar, but name is %q", r.name)
	}

	if r.object != "FOO" {
		t.Errorf("expected object type to be FOO, but object type is %q", r.object)
	}

	tag, err := r.nextTag()
	if err != nil {
		t.Errorf("Reader.nextTag[0]: %+v", err)
	}

	if len(tag) != 2 {
		t.Errorf("expected 2-element tag but tag length is %d", len(tag))
		t.Logf("tag: %#v", tag)
	} else {
		if tag[0] != "FOO" {
			t.Errorf("expected first element of tag to be FOO, but it is %q", tag[0])
		}
		if tag[1] != "BAR" {
			t.Errorf("expected second element of tag to be BAR, but it is %q", tag[1])
		}
	}

	tag, err = r.nextTag()
	if err != nil {
		t.Errorf("Reader.nextTag[1]: %+v", err)
	}

	if len(tag) != 1 {
		t.Errorf("expected 1-element tag but tag length is %d", len(tag))
		t.Logf("tag: %#v", tag)
	} else if tag[0] != "BAZ" {
		t.Errorf("expected only element of tag to be BAZ, but it is %q", tag[0])
	}

	tag, err = r.nextTag()
	if err != io.EOF {
		t.Errorf("Reader.nextTag[2]: expected io.EOF, but error is %+v", err)
	}

	if tag != nil {
		t.Errorf("Reader.nextTag[2]: returned non-nil tag: %#v", tag)
	}
}

func TestReader_Invalid_Name(t *testing.T) {
	t.Parallel()

	r := NewReader(strings.NewReader(""))
	if err := r.init(); err != ErrNoRawFileName {
		t.Errorf("Reader.init on empty raws file returned unexpected error: %+v", err)
	}

	r = NewReader(strings.NewReader("  \t  "))
	if err := r.init(); err != ErrNoRawFileName {
		t.Errorf("Reader.init on whitespace-only raws file returned unexpected error: %+v", err)
	}

	r = NewReader(strings.NewReader("  \t  \ntest"))
	if err := r.init(); err != ErrNoRawFileName {
		t.Errorf("Reader.init on raws file with whitespace-only first line returned unexpected error: %+v", err)
	}

	r = NewReader(strings.NewReader("test"))
	if err := r.init(); err != ErrNoRawFileName {
		t.Errorf("Reader.init on raws file with no newlines returned unexpected error: %+v", err)
	}
}

func TestReader_Invalid_Type(t *testing.T) {
	t.Parallel()

	r := NewReader(strings.NewReader("foo\n[OBJECTS:FOO]"))
	if err := r.init(); err != ErrNoRawObjectType {
		t.Errorf("Reader.init on raws file with OBJECTS instead of OBJECT returned unexpected error: %+v", err)
	}

	r = NewReader(strings.NewReader("foo\n[OBJECT:FOO:BAR]"))
	if err := r.init(); err != ErrNoRawObjectType {
		t.Errorf("Reader.init on raws file with too many parts to the OBJECT tag returned unexpected error: %+v", err)
	}

	r = NewReader(strings.NewReader("foo\n[OBJECT]"))
	if err := r.init(); err != ErrNoRawObjectType {
		t.Errorf("Reader.init on raws file with too few parts to the OBJECT tag returned unexpected error: %+v", err)
	}

	r = NewReader(strings.NewReader("foo\n[OBJECT:FOO"))
	if err := r.init(); err != ErrNoRawObjectType {
		t.Errorf("Reader.init on raws file with unterminated OBJECT tag returned unexpected error: %+v", err)
	}
}
