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
		t.Run(strings.TrimPrefix(strings.TrimSuffix(filepath.Base(name), ".txt"), "language_"), func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(name)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			r := raws.NewReader(f)

			var l []language.Tag
			if err = r.ParseAll(&l); err != nil {
				t.Error(err)
			}
			// TODO: check contents of l
		})
	}
}
