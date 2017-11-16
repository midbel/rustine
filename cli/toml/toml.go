package toml

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"text/scanner"
)

const (
	comma             = ','
	hash              = '#'
	leftAngleBracket  = '['
	rightAngleBracket = ']'
	leftCurlyBrace    = '{'
	rightCurlyBrace   = '}'
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

var ErrInvalidSyntax = errors.New("invalid syntax")

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

func New(r io.Reader) *lexer {
	s := new(scanner.Scanner)
	s.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanFloats | scanner.ScanInts
	return &lexer{Scanner: s}
}

func parseSection(lex *lexer) *section {
	var n string
	switch t := lex.Scan(); t {
	case scanner.Ident:
		n = lex.Text()
	case leftAngleBracket:
		return parseSection(lex)
	}
	lex.Scan()
	return &section{Label: n}
}

func parseOption(lex *lexer) *option {
	o := &option{Label: lex.Text()}
	lex.Scan()
	switch t := lex.Scan(); t {
	case leftAngleBracket, leftCurlyBrace:
		o.Value = parseComposite(lex)
	default:
		o.Value = parseSimple(lex)
	}
	return o
}

func parseComposite(lex *lexer) interface{} {
	switch lex.token {
	case leftAngleBracket:
		vs := make([]interface{}, 0, 10)
		for t := lex.Scan(); t != rightAngleBracket; t = lex.Scan() {
			switch t {
			case comma:
				continue
			case leftAngleBracket, leftCurlyBrace:
				vs = append(vs, parseComposite(lex))
			default:
				vs = append(vs, parseSimple(lex))
			}
		}
		return vs
	case leftCurlyBrace:
		vs := make(map[string]interface{})
		for t := lex.Scan(); t != rightCurlyBrace; t = lex.Scan() {
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
