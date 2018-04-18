package language_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BenLubar/dfide/raws"
	"github.com/BenLubar/dfide/raws/language"
)

func TestLanguage(t *testing.T) {
	t.Parallel()

	files, err := filepath.Glob("../testdata/objects/language_*.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range files {
		name := name // shadow
		base := filepath.Base(name)
		base = base[len("language_") : len(base)-len(".txt")]
		t.Run(base, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(name)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					t.Error(err)
				}
			}()

			r := raws.NewReader(f)

			var l []language.Tag
			if err = r.ParseAll(&l); err != nil {
				t.Error(err)
			}
			// TODO: check contents of l
		})
	}
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	r := raws.NewReader(strings.NewReader("language_empty\n\n[OBJECT:LANGUAGE]\n"))

	var l []language.Tag
	if err := r.ParseAll(&l); err != nil {
		t.Error(err)
	}
	if len(l) != 0 {
		t.Errorf("Expected language raws to be empty, but it contains %d entries.", len(l))
		for i, tag := range l {
			t.Errorf("%d: %#v %#v %#v", i, tag.Symbol, tag.Translation, tag.Word)
		}
	}
}
