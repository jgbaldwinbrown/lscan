package lscan

import (
	"io"
	"reflect"
	"unsafe"
	"strings"

	"bufio"
)

type Scanner struct {
	scanner *bufio.Scanner
	splitter Splitter
	linebuf []string
}

func (s *Scanner) Scan() bool {
	return s.scanner.Scan()
}

func (s *Scanner) Line() []string {
	s.linebuf = SplitByFunc(s.linebuf, s.scanner.Text(), s.splitter)
	return s.linebuf
}

func NewScanner(r io.Reader, split Splitter) (*Scanner) {
	s := new(Scanner)
	s.scanner = bufio.NewScanner(r)
	s.scanner.Buffer([]byte{}, 1e12)
	s.splitter = split
	return s
}

type Writer struct {
	writer io.Writer
	writefunc StringFormatter
	sep string
	end string
}

func (w *Writer) WriteLine(line []string) {
	WriteLine(w.writer, line, w.writefunc, w.sep, w.end)
}

func NewWriter(writer io.Writer, wf StringFormatter, sep, end string) *Writer {
	w := new(Writer)
	w.writer = writer
	w.writefunc = wf
	w.sep = sep
	w.end = end
	return w
}

type Splitter func(s string) (first, rest string, done bool)

func SplitByFunc(dest []string, line string, f Splitter) []string {
	dest = dest[:0]
	for {
		first, rest, done := f(line)
		dest = append(dest, first)
		if done {
			break
		}
		line = rest
	}
	return dest
}

func ByByte(b byte) func(string)(first, rest string, done bool){
	return func (s string) (first, rest string, done bool) {
		pos := strings.IndexByte(s, b)
		if pos == -1 {
			return s, s, true
		}
		return s[:pos], s[pos+1:], false
	}
}

func byByteEscapedOnly(sep, escape byte, s string, sind int) (first, rest string, done bool) {
	var b strings.Builder

	if sind != -1 {
		b.Grow(sind+5)
	}

	for i:=0; i<len(s); i++ {
		if s[i] == escape {
			i++
			if i < len(s) {
				b.WriteByte(s[i])
			}
		} else if s[i] == sep {
			return b.String(), s[i+1:], false
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String(), s, true
}

func ByByteEscaped(sep, escape byte) func(string)(first, rest string, done bool) {
	bybyte := ByByte(sep)
	return func (s string) (first, rest string, done bool) {
		sind := strings.IndexByte(s, sep)
		eind := -1
		if sind != -1 {
			eind = strings.IndexByte(s[:sind], escape)
		} else {
			eind = strings.IndexByte(s, escape)
		}

		if eind != -1 && eind < sind {
			return byByteEscapedOnly(sep, escape, s, sind)
		} else {
			return bybyte(s)
		}
	}
}

type StringFormatter func(w io.Writer, s string)

func UnsafeStringBytes(s string) []byte {
	var b []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	h.Cap = len(s)
	h.Len = len(s)
	return b
}

func WriteString(w io.Writer, s string) {
	b := UnsafeStringBytes(s)
	w.Write(b)
}

func WriteEscapedString(sep, escape byte) func(io.Writer, string) {
	return func(w io.Writer, s string) {
		ind := strings.IndexByte(s, sep)
		if ind == -1 {
			WriteString(w, s)
		} else {
			barr := [1]byte{}
			bout := barr[:]
			barr_esc := [1]byte{escape}
			besc := barr_esc[:]
			l := len(s)
			for i := 0; i < l; i++ {
				bout[0] = s[i]
				if bout[0] == sep || bout[0] == escape {
					w.Write(besc)
				}
				w.Write(bout)
			}
		}
	}
}

func WriteLine(w io.Writer, line []string, f StringFormatter, sep string, end string) {
	if len(line) > 0 {
		f(w, line[0])
	}
	for _, s := range line[1:] {
		WriteString(w, sep)
		f(w, s)
	}
	WriteString(w, end)
}
