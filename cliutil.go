package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/labstack/gommon/color"
	"github.com/mattn/go-isatty"
)

func colorSwitch(clr *color.Color, w io.Writer, fds ...uintptr) {
	clr.Disable()
	if len(fds) > 0 {
		if isatty.IsTerminal(fds[0]) {
			clr.Enable()
		}
	} else if w, ok := w.(*os.File); ok && isatty.IsTerminal(w.Fd()) {
		clr.Enable()
	}
}

// HelpCommandFn implements buildin help command function
func HelpCommandFn(ctx *Context) error {
	var (
		args   = ctx.NativeArgs()
		parent = ctx.Command().Parent()
	)
	if len(args) == 0 {
		ctx.String("%s", parent.Usage(ctx))
		return nil
	}
	var (
		child = parent.Route(args)
		clr   = ctx.Color()
	)
	if child == nil {
		return fmt.Errorf("command %s not found", clr.Yellow(strings.Join(args, " ")))
	}
	ctx.String("%s", child.Usage(ctx))
	return nil
}

// HelpCommand returns a buildin help command
func HelpCommand(desc string) *Command {
	return &Command{
		Name:        "help",
		Desc:        desc,
		CanSubRoute: true,
		NoHook:      true,
		Fn:          HelpCommandFn,
	}
}

// Daemon startup app as a daemon process, success if result from stderr has prefix successPrefix
func Daemon(ctx *Context, successPrefix string) error {
	cmd := exec.Command(os.Args[0], ctx.NativeArgs()...)
	serr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(serr)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	if strings.HasPrefix(line, successPrefix) {
		ctx.String("%s", line)
		cmd.Process.Release()
	} else {
		cmd.Process.Kill()
		line = strings.TrimSuffix(line, "\n")
		return errors.New(line)
	}
	return nil
}

// DaemonResponse output response to stderr
func DaemonResponse(resp string) {
	fmt.Fprintln(os.Stderr, resp)
}

// ReadJSON reads data as a json structure into argv
func ReadJSON(r io.Reader, argv interface{}) error {
	return json.NewDecoder(r).Decode(argv)
}

// ReadJSONFromFile is similar to ReadJSON, but read from file
func ReadJSONFromFile(filename string, argv interface{}) error {
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		err = ReadJSON(file, argv)
	}
	return err
}

// ReadJSONConfigFromFile is similar to ReadJSONFromFile, but allows reading file from where the executable file resides as well
func ReadJSONConfigFromFile(filename string, argv interface{}) error {
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		err = ReadJSON(file, argv)
	} else {
		exe, e := os.Executable()
		if e != nil {
			return e
		}
		// allow self-config .json files to go with the executable file, #40
		file, err = os.Open(filepath.Dir(exe) + string(filepath.Separator) + filename)
		if err == nil {
			defer file.Close()
			err = ReadJSON(file, argv)
		}
	}
	return err
}
