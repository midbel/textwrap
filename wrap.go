package textwrap

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const DefaultLength = 70

var Default Wrapper

func init() {
	Default = New()
}

func Shorten(str string, n int) string {
	if n <= 0 || n >= len(str) {
		return str
	}
	str, _, _ = advance(str, n)
	return fmt.Sprintf("%s...", str)
}

type Wrapper struct {
	limit       int
	Indent      string
	MergeBlanks bool
	MergeLines  bool
	Carriage    bool
}

func NewLimit(n int) (Wrapper, error) {
	var w Wrapper
	if n <= 0 {
		return w, fmt.Errorf("limit should be positive")
	}
	w.limit = n
	return w, nil
}

func New() Wrapper {
	w, _ := NewLimit(DefaultLength)
	return w
}

func Wrap(str string) string {
	return Default.Wrap(str)
}

func WrapN(str string, limit int) string {
	d := Default
	d.limit = limit
	return d.Wrap(str)
}

func (w Wrapper) Wrap(str string) string {
	if w.limit <= 0 || len(str) <= w.limit {
		return str
	}
	return w.wrapN(str)
}

func (w Wrapper) wrapN(str string) string {
	var (
		ws  strings.Builder
		ptr int
	)
	for i := 0; ptr < len(str); i++ {
		next, x, addnl := advance(str[ptr:], w.limit)
		if i > 0 && ptr < len(str) && x > 1 {
			if w.Carriage {
				ws.WriteRune(cr)
			}
			ws.WriteRune(nl)
		}
		ptr += x
		ws.WriteString(w.Indent)
		ws.WriteString(next)
		if addnl && len(next) > 0 {
			if w.Carriage {
				ws.WriteRune(cr)
			}
			ws.WriteRune(nl)
		}
		if x == 0 || ptr >= len(str) {
			break
		}
	}
	return ws.String()
}

func advance(str string, limit int) (string, int, bool) {
	if len(str) == 0 {
		return "", 0, false
	}
	var (
		prev int
		curr int
		step int
		last rune
		ws   strings.Builder
	)
	curr += skip(str[curr:], isBlank)
	for {
		last, step = peek(str[curr:], &ws)
		if last == utf8.RuneError {
			return ws.String(), curr, false
		}
		curr += step
		if isNL(last) {
			step = skip(str[curr:], isNL)
			return ws.String(), curr + step, step > 0 && len(str[curr+step:]) > 0
		}
		ws.WriteRune(last)
		if n, ok := canBreak(curr, prev, limit); ok {
			curr = n
			break
		}
		prev = curr
	}
	if str = ws.String(); len(str) > curr {
		str = str[:curr]
	}
	return str, curr, false
}

func peek(str string, ws *strings.Builder) (rune, int) {
	var (
		curr int
		size int
		last rune
	)
	for {
		last, size = next(str[curr:])
		curr += size
		if last == utf8.RuneError {
			return last, curr
		}
		if last == backslash {
			r, x := next(str[curr:])
			if isNL(r) {
				curr += x
				continue
			}
		}
		if isDelimiter(last) {
			break
		}
		ws.WriteRune(last)
	}
	return last, curr
}

func next(str string) (rune, int) {
	return utf8.DecodeRuneInString(str)
}

func canBreak(curr, prev, limit int) (int, bool) {
	if curr < limit {
		return 0, false
	}
	delta := limit - prev
	if diff := curr - limit; delta < diff {
		curr = prev
	}
	return curr, true
}

func skip(str string, fn func(rune) bool) int {
	var n int
	for {
		r, z := utf8.DecodeRuneInString(str[n:])
		if !fn(r) {
			break
		}
		n += z
	}
	return n
}

const (
	space     = ' '
	tab       = '\t'
	nl        = '\n'
	cr        = '\r'
	comma     = ','
	dot       = '.'
	question  = '?'
	bang      = '!'
	colon     = ':'
	semicolon = ';'
	backslash = '\\'
)

func isDelimiter(r rune) bool {
	return isNL(r) || isPunct(r) || isBlank(r) || r == utf8.RuneError
}

func isPunct(r rune) bool {
	return r == comma || r == dot || r == question || r == bang || r == semicolon || r == colon
}

func isBlank(r rune) bool {
	return r == space || r == tab
}

func isNL(r rune) bool {
	return r == nl
}
