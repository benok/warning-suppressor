package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type FilterConfig struct {
	Suppress []string
}

type ConfigDoc struct {
	Config FilterConfig `yaml:"filter-config"`
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
