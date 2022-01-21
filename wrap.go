package textwrap

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	DefaultLength    = 70
	DefaultThreshold = 2
)

func Shorten(str string, n int) string {
	if n <= 0 || n >= len(str) {
		return str
	}
	str, _ = advance(str, n)
	return fmt.Sprintf("%s...", str)
}

func Wrap(str string) string {
	return WrapN(str, DefaultLength)
}

func WrapN(str string, n int) string {
	if n <= 0 {
		return str
	}
	var (
		ws  strings.Builder
		ptr int
	)
	for i := 0; ptr < len(str); i++ {
		if i > 0 {
			ws.WriteRune(nl)
		}
		next, x := advance(str[ptr:], n)
		if x == 0 {
			break
		}
		ws.WriteString(strings.TrimSpace(next))
		ptr += x
	}
	return ws.String()
}

func advance(str string, n int) (string, int) {
	if len(str) == 0 {
		return "", 0
	}
	var (
		curr int
		prev int
		ws   strings.Builder
	)
	for {
		r, z := utf8.DecodeRuneInString(str[curr:])
		if r != utf8.RuneError {
			curr += z
		}
		if isNL(r) {
			ws.WriteRune(nl)
			curr += skip(str[curr:], isNL)
			break
		}
		if isDelimiter(r) {
			if isBlank(r) {
				curr += skip(str[curr:], isBlank)
			}
			if z, ok := canBreak(ws.Len(), prev, n); ok {
				curr = z
				if isPunct(r) {
					ws.WriteRune(r)
				}
				break
			}
			prev = curr
		}
		if r == utf8.RuneError {
			break
		}
		ws.WriteRune(r)
	}
	str = ws.String()
	if z := len(str); z > curr {
		str = str[:curr]
	}
	return str, curr
}

func canBreak(curr, prev, limit int) (int, bool) {
	if curr == limit {
		return curr, true
	}
	if curr > limit && curr-limit < DefaultThreshold {
		return curr, true
	}
	if curr > limit && curr-limit > DefaultThreshold {
		return prev, true
	}
	return 0, false
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
	comma     = ','
	dot       = '.'
	question  = '?'
	bang      = '!'
	colon     = ':'
	semicolon = ';'
)

func isDelimiter(r rune) bool {
	return isPunct(r) || isBlank(r) || r == utf8.RuneError
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
