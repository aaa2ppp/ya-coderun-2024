package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
		debug   bool
	}{
		{
			"0",
			args{strings.NewReader(`aaa
aaa`)},
			`1
1`,
			false,
			true,
		},
		{
			"1",
			args{strings.NewReader(`bwca
love`)},
			`1
1`,
			false,
			true,
		},
		{
			"1.1",
			args{strings.NewReader(`xxbwca
love`)},
			`2
2 3`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`abab
tat`)},
			`2
1 2`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`shla masha po shosse i sosala sushku
masha`)},
			// 3: `la ma`
			// 6:  masha
			// 10:`a po `
			// 14:` shos`
			`4
3 6 10 14`,
			false,
			true,
		},
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			if err := run(tt.args.in, out); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); strings.TrimRight(gotOut, " \r\n") != strings.TrimRight(tt.wantOut, " \r\n") {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func calcHash(s []byte) uint64 {
	var hash uint64
	for _, c := range s {
		hash = hash*X + uint64(c)
	}
	return hash
}

func calcPow(n int) []uint64 {
	pow := make([]uint64, n+1)
	pow[0] = 1
	for i := 1; i < len(pow); i++ {
		pow[i] = pow[i-1] * X
	}
	return pow
}

func Test_makePatterm(t *testing.T) {
	pow := calcPow(1000)

	type args struct {
		s   []byte
		pow []uint64
	}
	tests := []struct {
		args args
	}{
		{
			args{
				[]byte("abc"),
				pow,
			},
		},
		{
			args{
				[]byte("aaa"),
				pow,
			},
		},
		{
			args{
				[]byte("shla masha po shosse i sosala sushku"),
				pow,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(string(tt.args.s), func(t *testing.T) {
			pattern := makePattern(tt.args.s, tt.args.pow)
			var hash uint64
			for i := range pattern {
				hash += pattern[i].pow * uint64(tt.args.s[pattern[i].pos])
			}
			wantHash := calcHash(tt.args.s)
			if hash != wantHash {
				t.Logf("pattern: %v", pattern)
				t.Errorf("hash = %v, want hash %v", hash, wantHash)
			}
		})
	}
}
