package main

import (
	"fmt"
	"io"
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

type SyntaxError struct {
	msg   string
	token rune
}

func (s SyntaxError) Error() string {
	if s.token == 0 {
		return s.msg
	}
	return fmt.Sprintf("%s (%s)", s.msg, scanner.TokenString(s.token))
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

func New(r io.Reader) *lexer {
	s := new(scanner.Scanner)
	s.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanFloats | scanner.ScanInts
	return &lexer{Scanner: s}
}

func debug(lex *lexer) {
	for t := lex.Scan(); t != scanner.EOF; t = lex.Scan() {
		switch t {
		case scanner.Ident:
			log.Printf(">> %+v", parseOption(lex))
		case '[':
			log.Printf("> %+v", parseSection(lex))
		default:
			break
		}
	}
}

func parseSection(lex *lexer) *section {
	var n string
	switch lex.token {
	case scanner.Ident:
		ns := []string{lex.Text()}
		for t := lex.Scan(); t != rightSquareBracket; t = lex.Scan() {
			ns = append(ns, lex.Text())
		}
		n = strings.Join(ns, "")
	case leftSquareBracket:
		lex.Scan()
		return parseSection(lex)
	default:
		panic("section: unexpected token " + scanner.TokenString(lex.token))
	}
	for t := lex.Scan(); t == rightSquareBracket; t = lex.Scan() {
	}
	return &section{Label: n, Options: parseOptions(lex)}
}

func parseOptions(lex *lexer) []*option {
	if lex.token == leftSquareBracket {
		return nil
	}
	var os []*option
	for {
		os = append(os, parseOption(lex))
		if t := lex.Scan(); t == leftSquareBracket || t == scanner.EOF {
			break
		}
	}
	return os
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
