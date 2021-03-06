package main

import (
	"os"

	"github.com/akeylesslabs/cli"
)

type argT struct {
	cli.Helper
	Msg string `edit:"m" usage:"message"`
}

func main() {
	cli.GetEditor = func() (string, error) {
		if editor := os.Getenv("EDITOR"); editor != "" {
			return editor, nil
		}
		return cli.DefaultEditor, nil
	}
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("msg: %s", argv.Msg)
		return nil
	}))
}
