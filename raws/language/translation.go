package language // import "github.com/BenLubar/dfide/raws/language"

type Translation struct {
	ID string `raws:"1"`

	Words []TWord `raws:"T_WORD"`
}

type TWord struct {
	English string `raws:"1"`
	Native  string `raws:"2"`
}
