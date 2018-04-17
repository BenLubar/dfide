package language // import "github.com/BenLubar/dfide/raws/language"

type Symbol struct {
	ID    string   `raws:"1"`
	Words []string `raws:"S_WORD.1"`
}
