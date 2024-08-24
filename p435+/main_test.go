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
			"1",
			args{strings.NewReader(`1`)},
			`1
1`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2`)},
			`2
2`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`15`)},
			`12
6`,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`400000`)},
			`393120
192`,
			false,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			// if err := run(tt.args.in, out); (err != nil) != tt.wantErr {
			// 	t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			run(tt.args.in, out)
			if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func trimLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r\n")
	}
	n := len(lines)
	if lines[n-1] == "" {
		lines = lines[:n-1]
	}
	return strings.Join(lines, "\n")
}

var v, maximum int

func Benchmark_solve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v, maximum = solve(400000)
	}
}
