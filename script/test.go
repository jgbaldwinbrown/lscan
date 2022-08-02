package main

import (
	"os"
	"bufio"
	"github.com/jgbaldwinbrown/lscan/pkg"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start().Stop()

	in, err := os.Open("../input.txt")
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(in)

	out := bufio.NewWriter(os.Stdout)
	var line []string

	split := lscan.ByByteEscaped('\t', '\\')
	write := lscan.WriteEscapedString('\t', '\\')
	for s.Scan() {
		line = lscan.SplitByFunc(line, s.Text(), split)
		lscan.WriteLine(out, line, write, "\t", "\n")
	}
}
