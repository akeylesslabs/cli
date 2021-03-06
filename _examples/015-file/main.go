package main

import (
	"os"

	"github.com/akeylesslabs/cli"
	clix "github.com/akeylesslabs/cli/ext"
)

type argT struct {
	Content clix.File `cli:"f,file" usage:"read content from file or stdin"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String(argv.Content.String())
		return nil
	}))
}
