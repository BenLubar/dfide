package language // import "github.com/BenLubar/dfide/raws/language"

type Tag struct {
	Translation *Translation `raws:"TRANSLATION,union"`
	Symbol      *Symbol      `raws:"SYMBOL,union"`
	Word        *Word        `raws:"WORD,union"`
}
