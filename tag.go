package cli

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	tagCli  = "cli"
	tagPw   = "pw" // password
	tagEdit = "edit"

	tagUsage          = "usage"
	tagDefaut         = "dft"
	tagIgnoreDefault  = "ignoreDft"
	tagIgnoreRequired = "ignoreRequired"
	tagName           = "name"
	tagPrompt         = "prompt"
	tagParser         = "parser"
	tagSep            = "sep" // used to seperate key/value pair of map, default is `=`
	tagHidden         = "hidden"

	dashOne = "-"
	dashTwo = "--"

	sepName = ", "

	defaultSepForKeyValueOfMap = "="
)

type tagProperty struct {
	// is a required flag?
	isRequired bool `cli:"*x" pw:"*y" edit:"*z"`

	// is a force flag?
	isForce bool `cli:"!x" pw:"!y" edit:"!z"`

	// is a password flag?
	isPassword bool `pw:"xxx"`

	// is a edit flag?
	isEdit   bool   `edit:"xxx"`
	editFile string `edit:"FILE:xxx"`

	usage          string            `usage:"usage string"`
	dft            string            `dft:"default value or expression"`
	ignoreDft      bool              `ignoreDft:"true/false, by default false"`
	ignoreRequired bool              `ignoreRequired:"true/false, by default false"`
	name           string            `name:"tag reference name"`
	prompt         string            `prompt:"prompt string"`
	sep            string            `sep:"string for seperate kay/value pair of map"`
	parserCreator  FlagParserCreator `parser:"parser for flag"`
	hidden         bool              `hidden:"hide from help menu"`

	// flag names
	shortNames []string
	longNames  []string
}

func parseTag(fieldName string, structTag reflect.StructTag) (p *tagProperty, isEmpty bool, err error) {
	p = &tagProperty{
		shortNames: []string{},
		longNames:  []string{},
	}
	cliLikeTagCount := 0

	tag := newMultiTag(string(structTag))

	// `cli` TAG
	cli := tag.Get(tagCli)
	if cli != "" {
		cliLikeTagCount++
	}

	// `pw` TAG
	if pw := tag.Get(tagPw); pw != "" {
		p.isPassword = true
		cli = pw
		cliLikeTagCount++
	}

	// `edit` TAG
	if edit := tag.Get(tagEdit); edit != "" {
		// specific filename for editor
		sepIndex := strings.Index(edit, ":")
		if sepIndex > 0 {
			p.editFile = edit[:sepIndex]
			edit = edit[sepIndex+1:]
		}
		p.isEdit = true
		cli = edit
		cliLikeTagCount++
	}

	if cliLikeTagCount > 1 {
		err = errCliTagTooMany
		return
	}

	// `usage` TAG
	p.usage = tag.Get(tagUsage)

	// `dft` TAG
	p.dft = tag.Get(tagDefaut)

	// `ignoreDft` TAG
	if tag := tag.Get(tagIgnoreDefault); tag != "" {
		p.ignoreDft, err = strconv.ParseBool(tag)

	}
	if err != nil {
		return
	}

	// `ignoreRequired` TAG
	if tag := tag.Get(tagIgnoreRequired); tag != "" {
		p.ignoreRequired, err = strconv.ParseBool(tag)
	}
	if err != nil {
		return
	}

	// `name` TAG
	p.name = tag.Get(tagName)

	// `prompt` TAG
	p.prompt = tag.Get(tagPrompt)

	// `parser` TAG
	if parserName := tag.Get(tagParser); parserName != "" {
		if parserCreator, ok := parserCreators[parserName]; ok {
			p.parserCreator = parserCreator
		}
	}

	// `sep` TAG
	p.sep = defaultSepForKeyValueOfMap
	if sep := tag.Get(tagSep); sep != "" {
		p.sep = sep
	}

	if hidden := tag.Get(tagHidden); hidden != "" {
		p.hidden, err = strconv.ParseBool(hidden)
	}
	if err != nil {
		return
	}

	cli = strings.TrimSpace(cli)
	for {
		if strings.HasPrefix(cli, "*") {
			p.isRequired = true
			cli = strings.TrimSpace(strings.TrimPrefix(cli, "*"))
		} else if strings.HasPrefix(cli, "!") {
			p.isForce = true
			cli = strings.TrimSpace(strings.TrimPrefix(cli, "!"))
		} else {
			break
		}
	}

	names := strings.Split(cli, ",")
	isEmpty = true
	for _, name := range names {
		if name = strings.TrimSpace(name); name == dashOne {
			return nil, false, nil
		}
		if len(name) == 0 {
			continue
		} else if len(name) == 1 {
			p.shortNames = append(p.shortNames, dashOne+name)
		} else {
			p.longNames = append(p.longNames, dashTwo+name)
		}
		isEmpty = false
	}
	if isEmpty {
		p.longNames = append(p.longNames, dashTwo+fieldName)
	}
	return
}
