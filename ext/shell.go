package ext

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/akeylesslabs/cli"
)

// InstallBashCompletion install bash_completion
func InstallBashCompletion(root *cli.Command) error {
	if root.Name == "" {
		return fmt.Errorf("root command's name is empty")
	}
	compFilename := "." + root.Name + "_compeltion"
	compFilepath := filepath.Join(os.Getenv("HOME"), compFilename)
	compFile, err := os.OpenFile(compFilepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	buff, err := genBashCompletion(root)
	if err != nil {
		return err
	}
	if _, err := compFile.Write(buff.Bytes()); err != nil {
		return err
	}

	bashrcFiles := []string{
		".bashrc",
		".bash_profile",
	}
	var dstFile *os.File
	for _, bashFilename := range bashrcFiles {
		dstFilepath := filepath.Join(os.Getenv("HOME"), bashFilename)
		dstFile, err = os.OpenFile(dstFilepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err == nil {
			break
		}
	}
	if dstFile == nil {
		return fmt.Errorf("no destination bash file")
	}
	data, err := ioutil.ReadAll(dstFile)
	if err != nil {
		return err
	}

	subslice := fmt.Sprintf("[ -f ~/%s ] && . ~/%s", compFilename, compFilename)
	if !bytes.Contains(data, []byte(subslice)) {
		_, err = dstFile.WriteString("\n#Auto-generated by cli\n" + subslice)
	}

	return err
}

func genBashCompletion(root *cli.Command) (*bytes.Buffer, error) {
	buff := bytes.NewBufferString("")
	t, err := template.New("bash_completion").Parse(shellTemplateText)
	if err != nil {
		return nil, err
	}
	return buff, t.Execute(buff, struct {
		Cli        string
		CompleteFn string
	}{Cli: root.Name, CompleteFn: genCompleteFn(root)})
}

func genCompleteFn(root *cli.Command) string {
	return "TODO=true"
}

const shellTemplateText = `# {{.Cli}} command completion script

COMP_WORDBREAKS=${COMP_WORDBREAKS/=/}
COMP_WORDBREAKS=${COMP_WORDBREAKS/@/}
export COMP_WORDBREAKS

__complete_fn() {
#COMP_CWORD
#COMP_LINE
#COMP_POINT
#COMP_WORDS
#{{.Cli}} completion -- "${COMP_WORDS[@]}"
{{.CompleteFn}}
}

if type complete &>/dev/null; then
  _{{.Cli}}_completion () {
    local si="$IFS"
    IFS=$'\n' COMPREPLY=($(__complete_fn \
                           2>/dev/null)) || return $?
    IFS="$si"
  }
  complete -F _{{.Cli}}_completion {{.Cli}}
elif type compdef &>/dev/null; then
  _{{.Cli}}_completion() {
    si=$IFS
    compadd -- $(COMP_CWORD=$((CURRENT-1)) \
                 COMP_LINE=$BUFFER \
                 COMP_POINT=0 \
                 COMP_WORDS="${words[@]}" \
				 __complete_fn \
                 2>/dev/null)
    IFS=$si
  }
  compdef _{{.Cli}}_completion {{.Cli}}
elif type compctl &>/dev/null; then
  _{{.Cli}}_completion () {
    local cword line point words si
    read -Ac words
    read -cn cword
    let cword-=1
    read -l line
    read -ln point
    si="$IFS"
    IFS=$'\n' reply=($(COMP_CWORD="$cword" \
                       COMP_LINE="$line" \
                       COMP_POINT="$point" \
                       COMP_WORDS="${words[@]}" \
					   __complete_fn \
                       2>/dev/null)) || return $?
    IFS="$si"
  }
  compctl -K _{{.Cli}}_completion {{.Cli}}
fi
`
