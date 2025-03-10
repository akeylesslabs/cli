package main

import (
	"io"
	"os"
	"strings"

	"github.com/akeylesslabs/cli"
	clix "github.com/akeylesslabs/cli/ext"
)

type argT struct {
	Reader *clix.Reader `cli:"r,reader" usage:"read from file, stdin, http or any io.Reader"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		data, err := io.ReadAll(argv.Reader)
		argv.Reader.Close()
		if err != nil {
			return err
		}
		ctx.String("read from file(or http, stdin): %s\n", string(data))
		ctx.String("filename: %s, isStdin=%v\n", argv.Reader.Name(), argv.Reader.IsStdin())

		// Replace the reader
		argv.Reader.SetReader(strings.NewReader("string reader"))
		data, err = io.ReadAll(argv.Reader)
		if err != nil {
			return err
		}
		ctx.String("reade from reader: %s\n", string(data))
		ctx.String("filename: %s, isStdin=%v\n", argv.Reader.Name(), argv.Reader.IsStdin())
		return nil
	}))
}
