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
		debug bool
	}{
		{
			"1",
			args{strings.NewReader(`5 4
1 2
2 3
3 4
4 5`)},
			`6`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`4 3
1 2
2 3
3 4`)},
			`3`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`4 3
1 2
2 3
2 4`)},
			`3`,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`4 4
1 2
1 3
2 3
3 4`)},
			`0`,
			false,
			true,
		},
		{
			"5",
			args{strings.NewReader(`11 12
1 2
1 3
1 4
1 5
4 5
5 6
6 7
7 8
7 9
7 10
8 9
10 11
`)},
			`7`,
			false,
			true,
		},
		{
			"6",
			args{strings.NewReader(`6 7
1 2
1 3
1 4
2 3
2 4
2 5
5 6
`)},
			`1`,
			false,
			true,
		},
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
