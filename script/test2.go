package main

import (
	"os"
	"strings"
	"bufio"
	"../lscan"
)

func main() {
	str := "Go h\\\\ang\na salami\\\nI'm a lasagna\\ hog.\nyou dude."
	r := strings.NewReader(str)
	s := bufio.NewScanner(r)
	var line []string
	for s.Scan() {
		line = lscan.SplitByFunc(line, s.Text(), lscan.ByByte(' '))
		lscan.WriteLine(os.Stdout, line, lscan.WriteString, "\t", "\n")
	}

	r = strings.NewReader(str)
	s = bufio.NewScanner(r)
	for s.Scan() {
		line = lscan.SplitByFunc(line, s.Text(), lscan.ByByteEscaped(' ', '\\'))
		lscan.WriteLine(os.Stdout, line, lscan.WriteEscapedString('\t', '\\'), "\t", "\n")
	}
}
