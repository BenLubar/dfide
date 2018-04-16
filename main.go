package main // import "github.com/BenLubar/dfide"

import "github.com/BenLubar/dfide/gui"

func main() {
	if err := gui.Main(); err != nil {
		panic(err)
	}
}
