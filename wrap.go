package textwrap

import (
	"fmt"
	"iter"
	"slices"
	"strings"
	"unicode/utf8"
)

const DefaultLength = 70

func Shorten(str string, limit int) string {
	if limit <= 0 || limit >= len(str) {
		return str
	}
	next, stop := iter.Pull(Lines(str, limit))
	defer stop()

	str, _ = next()
	return fmt.Sprintf("%s...", strings.TrimSpace(str))
}

func Split(str string, limit int) []string {
	if limit <= 0 || limit >= len(str) {
		return []string{str}
	}
	return slices.Collect(wrap(str, limit))
}

func Wrap(str string, limit int) string {
	if limit <= 0 || limit >= len(str) {
		return str
	}
	var (
		out strings.Builder
		lino int
	)
	for str := range Lines(str, limit) {
		if lino > 0 {
			out.WriteString("\n")
		}
		out.WriteString(str)
	}
	return out.String()
}

func Lines(str string, limit int) iter.Seq[string] {
	return wrap(str, limit)
}

func wrap(str string, limit int) iter.Seq[string] {
	it := func(yield func(string) bool) {
		if len(str) < limit {
			yield(str)
			return
		}
		var ptr int
		for ptr < len(str) {
			ptr = skipCollapsible(str, ptr)
			var (
				cs        int
				lastChar  rune
				done      bool
				prevDelim = -1
			)
			// consumes up to limit characters
			for cs < limit {
				if ptr+cs >= len(str) {
					done = true
					break
				}
				r, z := utf8.DecodeRuneInString(str[ptr+cs:])
				cs += z
				lastChar = r
				if isDelimiter(lastChar) {
					prevDelim = ptr + cs
					dot := consumeDot(str, prevDelim)
					prevDelim += dot
					cs += dot
				}
			}
			// break on delim
			if prevDelim == ptr+cs && isDelimiter(lastChar) {
				if isCollapsible(lastChar) {
					cs -= utf8.RuneLen(lastChar)
				}
				if !yield(str[ptr : ptr+cs]) {
					return
				}
				ptr += cs
				continue
			}
			nextDelim := ptr + cs
			if !done {
				nextDelim = nextDelimiter(str, nextDelim)
			}
			next := ptr + cs
			switch {
			case prevDelim < 0:
				next = nextDelim
			case next-prevDelim <= nextDelim-next:
				next = prevDelim
			default:
				next = nextDelim
			}
			if !yield(str[ptr:next]) {
				return
			}
			ptr = next
		}
	}
	return it
}

func skipCollapsible(str string, ptr int) int {
	// skip collapsible characters
	for ptr < len(str) {
		r, z := utf8.DecodeRuneInString(str[ptr:])
		if !isCollapsible(r) {
			break
		}
		ptr += z
	}
	return ptr
}

func nextDelimiter(str string, ptr int) int {
	var read int
	for ptr+read < len(str) {
		r, z := utf8.DecodeRuneInString(str[ptr+read:])
		read += z
		if r == utf8.RuneError || isDelimiter(r) {
			break
		}
	}
	return ptr + read
}

func consumeDot(str string, ptr int) int {
	var read int
	for ptr+read < len(str) {
		r, z := utf8.DecodeRuneInString(str[ptr+read:])
		if !isDot(r) {
			break
		}
		read += z
	}
	return read
}

func isDot(r rune) bool {
	return r == '.'
}

func isCollapsible(r rune) bool {
	switch r {
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	}
}

func isDelimiter(r rune) bool {
	if isCollapsible(r) {
		return true
	}
	switch r {
	case ';', ',', '.', '!', '?', ':', '(', ')', '[', ']', '{', '}':
		return true
	default:
		return false
	}
}