package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type ColorStr string
type Pattern string
type PatternToColor map[Pattern]ColorStr
type RegExToColor struct {
	Regex *regexp.Regexp
	Color color.Attribute
}

type Colorize_ struct {
	Maps []PatternToColor `yaml:"colorize,omitempty"`
}

type FilterConfig struct {
	Suppress []string
	Colorize Colorize_
}

type ConfigDoc struct {
	Config FilterConfig `yaml:"filter-config"`
}

// https://www.ribice.ba/golang-yaml-string-map/
func (cl *Colorize_) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&cl.Maps)
}

func (c *FilterConfig) Parse(data []byte) error {
	var doc ConfigDoc
	err := yaml.Unmarshal(data, &doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unmarshal error: %s\n", err.Error())
	} else {
		*c = doc.Config
	}
	return err
}

func (c *FilterConfig) Load(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
	}
	err = c.Parse(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *FilterConfig) SuppressRegExp() (regexp.Regexp, error) {
	var a []string
	var pat string

	sups := c.Suppress
	if len(sups) > 0 {
		for _, s := range sups {
			a = append(a, "("+s+")")
		}
		pat = fmt.Sprintf("^.*(%s).*$", strings.Join(a, "|"))
	} else {
		pat = ".^" // invalid pattern
	}
	r, err := regexp.Compile(pat)
	return *r, err
}

func (s ColorStr) ToColor() color.Attribute {
	c_str := strings.ToLower(string(s))
	switch c_str {
	case "red":
		return color.FgHiRed
	case "green":
		return color.FgHiGreen
	case "yellow":
		return color.FgHiYellow
	case "blue":
		return color.FgHiBlue
	case "magenta":
		return color.FgHiMagenta
	case "cyan":
		return color.FgHiCyan
	case "white":
		return color.FgHiWhite
	default:
		return color.Reset
	}
}

func (c *FilterConfig) RegExToColors() []RegExToColor {
	var result []RegExToColor
	pat_to_cols := c.Colorize.Maps
	for _, p2c := range pat_to_cols {
		for pat, col := range p2c {
			var r2c RegExToColor
			r2c.Color = col.ToColor()
			r2c.Regex, _ = regexp.Compile(string(pat))
			result = append(result, r2c)
		}
	}
	return result
}
