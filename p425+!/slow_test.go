package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func Test_slowRun(t *testing.T) {
	debugEnable = false
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
			args{strings.NewReader(`4 3
1 3 1
2 4 2
3 4 4`)},
			`4
1 2 3 4 `,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`7 2
1 5 1
3 7 2`)},
			`3
1 2 3 `,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`10 3
1 1 2
1 1 3
1 1 6`)},
			`6
2 3 5 6 8 9 `,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`10 3
1 1 1
2 2 2
3 3 3`)},
			`3
1 2 3`,
			false,
			true,
		},
		{
			"4.1",
			args{strings.NewReader(`74 3
65 65 1
66 66 2
67 67 3`)},
			`3
1 2 3`,
			false,
			true,
		},
		{
			"4.2",
			args{strings.NewReader(`140 3
129 129 1
130 130 2
131 131 3`)},
			`3
1 2 3`,
			false,
			true,
		},
		{
			"5",
			args{strings.NewReader(`100 3
16 21 3
12 59 4
15 64 5`)},
			`7
3 4 5 7 8 9 12`,
			false,
			true,
		},
		{
			"8",
			args{strings.NewReader(`100 4
3 5 1 21 59 1 55 62 1 29 63 2`)},
			`4
1 2 3 4`,
			false,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			if err := slowRun(tt.args.in, out); (err != nil) != tt.wantErr {
				t.Errorf("slowRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); strings.TrimRight(gotOut, " \r\n") != strings.TrimRight(tt.wantOut, " \r\n") {
				t.Errorf("slowRun() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
