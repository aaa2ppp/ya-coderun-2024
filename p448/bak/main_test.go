package main

import (
	"bytes"
	"io"
	"strconv"
	_ "strconv"
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
			"1",
			args{strings.NewReader(`3
			1 10
			a 12
			5 b`)},
			`325`,
			false,
			true,
		},
		{
			"1.1",
			args{strings.NewReader(`3
			1 10
			1 a
			5 b`)},
			strconv.Itoa(func() int {
				var res int
				for a := 1; a <= 10; a++ {
					for b := 1; b <= a; b++ {
						for c := 5; c <= b; c++ {
							res += 1
							res %= modulo
						}
					}
				}
				return res
			}()),
			false,
			true,
		},
		{
			"1.2",
			args{strings.NewReader(`3
			1 10
			a 12
			5 a`)},
			strconv.Itoa(func() int {
				var res int
				for a := 1; a <= 10; a++ { // 5..10
					for b := a; b <= 12; b++ {
						for c := 5; c <= a; c++ {
							res += 1
							res %= modulo
						}
					}
				}
				return res
			}()),
			false,
			true,
		},
		{
			"1.3",
			args{strings.NewReader(`3
			1 10
			a 8
			5 a`)},
			strconv.Itoa(func() int {
				var res int
				for a := 1; a <= 10; a++ { // 5..8
					for b := a; b <= 8; b++ {
						for c := 5; c <= a; c++ {
							res += 1
							res %= modulo
						}
					}
				}
				return res
			}()),
			false,
			true,
		},
		{
			"1.4",
			args{strings.NewReader(`3
			1 10
			7 a
			5 a`)},
			strconv.Itoa(func() int {
				var res int
				for a := 1; a <= 10; a++ { // 5..8
					for b := 7; b <= a; b++ {
						for c := 5; c <= a; c++ {
							res += 1
							res %= modulo
						}
					}
				}
				return res
			}()),
			false,
			true,
		},
		// {
		// 	"2",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	false,
		// 	true,
		// },
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
			if gotOut := out.String(); strings.TrimRight(gotOut, "\r\n") != strings.TrimRight(tt.wantOut, "\r\n") {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
