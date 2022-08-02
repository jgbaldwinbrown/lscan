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
	s := lscan.NewScanner(in, lscan.ByByteEscaped('\t', '\\'))

	out := bufio.NewWriter(os.Stdout)
	w := lscan.NewWriter(out, lscan.WriteEscapedString('\t', '\\'), "\t", "\n")

	for s.Scan() {
		w.WriteLine(s.Line())
	}
}
