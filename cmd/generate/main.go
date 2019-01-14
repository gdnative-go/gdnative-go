package main

import (
	"gitlab.com/pimpam-games-studio/gdnative-go/cmd/generate/gdnative"
	"gitlab.com/pimpam-games-studio/gdnative-go/cmd/generate/types"
)

func main() {

	gdnative.Generate()
	types.Generate()
}
