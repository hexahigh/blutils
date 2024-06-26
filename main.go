package main

import (
	"github.com/hexahigh/blutils/cmd"

	_ "github.com/hexahigh/blutils/cmd/bench"
	_ "github.com/hexahigh/blutils/cmd/bitflip"
	_ "github.com/hexahigh/blutils/cmd/report"
	_ "github.com/hexahigh/blutils/cmd/update"
	_ "github.com/hexahigh/blutils/cmd/whatis"
)

func main() {
	cmd.Execute()
}
