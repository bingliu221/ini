package ini

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var ErrFormat = errors.New("format error")

type Key string
type Value string
type Section string

type Config map[Section]map[Key]Value

func (v Value) Number() float64 {
	n, _ := strconv.ParseFloat(string(v), 64)
	return n
}

func (v Value) String() string {
	return string(v)
}

func (v Value) Boolean() bool {
	switch string(v) {
	case "Y", "y", "Yes", "yes", "1", "true", "True":
		return true
	}
	return false
}

func Parse(s string) (Config, error) {
	cfg := Config{"": make(map[Key]Value)}
	sectionReg := regexp.MustCompile(`^\[\w+\]$`)
	propertyReg := regexp.MustCompile(`^\s*\w+\s*=`)
	currentSection := cfg[""]
	for n, line := range strings.Split(s, "\n") {
		switch {
		case len(line) == 0:
		case strings.HasPrefix(line, ";"):

		case sectionReg.MatchString(line):
			name := Section(line[1 : len(line)-1])
			cfg[name] = make(map[Key]Value)
			currentSection = cfg[name]

		case propertyReg.MatchString(line):
			parts := strings.SplitN(line, "=", 2)
			key := Key(strings.TrimSpace(parts[0]))
			value := strings.TrimSpace(parts[1])
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") || strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
				value = value[1 : len(value)-1]
			}

			escape := -1
			unqoutedValue := make([]rune, 0)
			for i, c := range value {
				if i == escape {
					switch c {
					case '0':
						c = 0
					case 'a':
						c = '\a'
					case 'b':
						c = '\b'
					case 't':
						c = '\t'
					case 'r':
						c = '\r'
					case 'n':
						c = '\n'
					}
					unqoutedValue = append(unqoutedValue, c)
					continue
				}
				if c == '\\' {
					escape = i + 1
					continue
				}
				unqoutedValue = append(unqoutedValue, c)
			}
			currentSection[key] = Value(string(unqoutedValue))

		default:
			return nil, fmt.Errorf("%v at line %d, %s", ErrFormat, n+1, line)
		}
	}
	return cfg, nil
}
