package main

import (
	"os"
)

func main() {
	s1, s2 := "Stepik", "Hello"
	s1, s2 = changeStrings(s1, s2)
	
	var buf []byte
	buf = append(buf, s1...)
	buf = append(buf, ' ')
	buf = append(buf, s2...)
	buf = append(buf, '\n')
	os.Stdout.Write(buf)
}

//go:noinline
func changeStrings(s1, s2 string) (string, string) {
	return s2, s1
}
