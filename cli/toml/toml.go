package toml

import (
	"io"
	"sort"
	"strconv"
	"strings"
	"text/scanner"
)

const (
	dot                = '.'
	comma              = ','
	equal              = '='
	hash               = '#'
	leftSquareBracket  = '['
	rightSquareBracket = ']'
	leftCurlyBracket   = '{'
	rightCurlyBracket  = '}'
)

type section struct {
	Label    string
	Options  []*option
	Sections []*section
}

type option struct {
	Label string
	Value interface{}
}

var ids = map[string]interface{}{
	"true":  true,
	"false": false,
}

type lexer struct {
	*scanner.Scanner
	token rune
}

func (l *lexer) Text() string {
	return l.TokenText()
}

func (l *lexer) Scan() rune {
	l.token = l.Scanner.Scan()
	if l.token == hash {
		p := l.Scanner.Position
		for {
			if l.Position.Line > p.Line {
				break
			}
			l.token = l.Scan()
		}
	}
	return l.token
}

func start(r io.Reader) *lexer {
	s := new(scanner.Scanner)
	s.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanFloats | scanner.ScanInts
	return &lexer{Scanner: s}
}

func parseSection(lex *lexer, s *section, a bool) *section {
	switch lex.token {
	default:
		panic("unexpected token " + scanner.TokenString(lex.token))
	case scanner.Ident:
		sort.Slice(s.Sections, func(i, j int) bool {
			return s.Sections[i].Label < s.Sections[j].Label
		})
		var curr *section
		for sup := s; lex.token != rightSquareBracket; sup = curr {
			if lex.token == dot {
				lex.Scan()
			}
			curr = &section{Label: lex.Text()}
			if lex.Peek() == rightSquareBracket && exists(curr.Label, sup.Sections) && a {
				panic("duplicate section: " + curr.Label)
			}
			sup.Sections = append(sup.Sections, curr)
			if t := lex.Scan(); t == '\n' || t == scanner.EOF {
				panic("invalid syntax")
			}
		}
		if t := lex.Scan(); !a {
			if t != rightSquareBracket {
				panic("invalid syntax: missing " + scanner.TokenString(rightSquareBracket))
			}
			lex.Scan()
		}
		curr.Options = parseOptions(lex)

		return curr
	case leftSquareBracket:
		lex.Scan()
		return parseSection(lex, s, false)
	}
}

func exists(n string, vs []*section) bool {
	ix := sort.Search(len(vs), func(i int) bool {
		return vs[i].Label >= n
	})
	return ix < len(vs) && vs[ix].Label == n
}

func parseOptions(lex *lexer) []*option {
	if lex.token == leftSquareBracket || lex.token == scanner.EOF {
		return nil
	}
	os := make(map[string]struct{})
	vs := make([]*option, 0)
	for {
		o := parseOption(lex)
		if _, ok := os[o.Label]; ok {
			panic("duplicate option: " + o.Label)
		}
		os[o.Label] = struct{}{}
		vs = append(vs, o)
		if t := lex.Scan(); t == leftSquareBracket || t == scanner.EOF {
			break
		}
	}
	return vs
}

func parseOption(lex *lexer) *option {
	o := &option{Label: lex.Text()}
	if t := lex.Scan(); t != equal {
		panic("option: expected: '=' got: " + scanner.TokenString(t))
	}
	if t := lex.Peek(); t == '\n' {
		panic("option: missing value")
	}
	switch t := lex.Scan(); t {
	case leftSquareBracket, leftCurlyBracket:
		o.Value = parseComposite(lex)
	default:
		o.Value = parseSimple(lex)
	}
	return o
}

func parseComposite(lex *lexer) interface{} {
	switch lex.token {
	case leftSquareBracket:
		vs := make([]interface{}, 0, 10)
		for t := lex.Scan(); t != rightSquareBracket; t = lex.Scan() {
			switch t {
			case comma:
				continue
			case leftSquareBracket, leftCurlyBracket:
				vs = append(vs, parseComposite(lex))
			default:
				vs = append(vs, parseSimple(lex))
			}
		}
		return vs
	case leftCurlyBracket:
		vs := make(map[string]interface{})
		for t := lex.Scan(); t != rightCurlyBracket; t = lex.Scan() {
			if t == comma {
				continue
			}
			o := parseOption(lex)
			vs[o.Label] = o.Value
		}
		return vs
	default:
		return nil
	}
}

func parseSimple(lex *lexer) interface{} {
	var v interface{}
	switch t := lex.token; t {
	case scanner.String:
		v = strings.Trim(lex.Text(), "\"")
	case scanner.Int:
		v, _ = strconv.ParseInt(lex.Text(), 0, 64)
	case scanner.Float:
		v, _ = strconv.ParseFloat(lex.Text(), 64)
	case scanner.Ident:
		v = ids[lex.Text()]
	case '-':
		lex.Scan()
		v = parseSimple(lex)
		switch n := v.(type) {
		case int64:
			v = -n
		case float64:
			v = -n
		}
	default:
		v = lex.Text()
	}
	return v
}
