package vyos

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// Matches section start `interfaces {`
var rxSection = regexp.MustCompile(`^([\w\-]+) \{$`)

// Matches named section `ethernet eth0 {`
var rxNamedSection = regexp.MustCompile(`(?m)^([\w\-]+) ([\w\-\"\./@:=\+]+) \{`)

// Matches simple key-value pair `duplex auto`
var rxValue = regexp.MustCompile(`^([\w\-]+) "?([^"]+)?"?$`)

// Matches single value (flag) `disable`
var rxFlag = regexp.MustCompile(`^([\w\-]+)$`)

// Matches comments
var rxComment = regexp.MustCompile(`^(\/\*).*(\*\/)`)

type Config map[string]interface{}

func (c *Config) Set(path []string, key string, value interface{}) {
	var v map[string]interface{}
	v = *c
	for _, s := range path {
		if _, ok := v[s]; !ok {
			v[s] = make(map[string]interface{})
		}
		v = v[s].(map[string]interface{})
	}
	v[key] = value
}

func Parse(data []byte) (*Config, error) {
	cfg := &Config{}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	var path []string
	var pathTypes []string
	for i := 0; scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if rxSection.MatchString(line) {
			match := rxSection.FindStringSubmatch(line)[1]
			path = append(path, match)
			pathTypes = append(pathTypes, "section")
		} else if rxNamedSection.MatchString(line) {
			matches := rxNamedSection.FindStringSubmatch(line)
			section, name := matches[1], matches[2]
			if (path)[len(path)-1] != section {
				path = append(path, section)
				pathTypes = append(pathTypes, "named_section")
			}
			path = append(path, name)
			pathTypes = append(pathTypes, "named_section")
		} else if rxValue.MatchString(line) {
			matches := rxValue.FindStringSubmatch(line)
			key, value := matches[1], matches[2]
			cfg.Set(path, key, value)
		} else if rxFlag.MatchString(line) {
			matches := rxFlag.FindStringSubmatch(line)
			flag := matches[1]
			cfg.Set(path, flag, true)
		} else if line == "}" {
			if pathTypes[len(path)-1] == "named_section" {
				pathTypes = pathTypes[:len(path)-1]
				path = path[:len(path)-1]
			}

			pathTypes = pathTypes[:len(path)-1]
			path = path[:len(path)-1]
		} else {
			panic(fmt.Sprintf("Parse error at %d: %q", i, line))
		}
	}

	return cfg, nil
}
