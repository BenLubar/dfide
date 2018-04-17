package raws // import "github.com/BenLubar/dfide/raws"

import (
	"strconv"
)

// Codepage 437, from https://dwarffortresswiki.org/Character_table
var cp437 = []rune("\x00☺☻♥♦♣♠•◘○◙♂♀♪♬☼" +
	"►◄↕‼¶§▬↨↑↓→←∟↔▲▼" +
	" !\"#$%&'()*+,-./" +
	"0123456789:;<=>?" +
	"@ABCDEFGHIJKLMNO" +
	"PQRSTUVWXYZ[\\]^_" +
	"`abcdefghijklmno" +
	"pqrstuvwxyz{|}~⌂" +
	"ÇüéâäàåçêëèïîìÄÅ" +
	"ÉæÆôöòûùÿÖÜ¢£¥₧ƒ" +
	"áíóúñÑªº¿⌐¬½¼¡«»" +
	"░▒▓│┤╡╢╖╕╣║╗╝╜╛┐" +
	"└┴┬├─┼╞╟╚╔╩╦╠═╬╧" +
	"╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀" +
	"αßΓπΣσµτΦΘΩδ∞φε∩" +
	"≡±≥≤⌠⌡÷≈°∙·√ⁿ²■\u00A0")

var cp437Rev = func() map[rune]byte {
	m := make(map[rune]byte)
	for i, r := range cp437 {
		m[r] = byte(i)
	}
	return m
}()

func ToChar(s string) (rune, error) {
	if len(s) > 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		r := []rune(s[1 : len(s)-1])
		if len(r) != 1 {
			return 0, ErrInvalidCharacter
		}
		if _, ok := cp437Rev[r[0]]; !ok {
			return 0, ErrInvalidCharacter
		}
		return r[0], nil
	}

	c, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, ErrInvalidCharacter
	}
	return cp437[c], nil
}

func FromChar(r rune) (string, error) {
	if ' ' <= r && r <= '~' && r != '\'' && r != '[' && r != ']' && r != ':' {
		return string([]rune{'\'', r, '\''}), nil
	}

	c, ok := cp437Rev[r]
	if !ok {
		return "", ErrInvalidCharacter
	}

	return strconv.FormatUint(uint64(c), 10), nil
}
