package raws // import "github.com/BenLubar/dfide/raws"

import "errors"

var (
	ErrNoRawFileName    = errors.New("raws: no file name (first line was blank)")
	ErrNoRawObjectType  = errors.New("raws: missing OBJECT tag at start of file")
	ErrIncompleteTag    = errors.New("raws: incomplete tag (missing closing bracket)")
	ErrInvalidCharacter = errors.New("raws: invalid character")
)
