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
			`1`,
			false,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2`)},
			`1`,
			false,
			true,
		},
		{
			"3",
			args{strings.NewReader(`3`)},
			`2`,
			false,
			true,
		},
		{
			"4",
			args{strings.NewReader(`4`)},
			`2`,
			false,
			true,
		},
		{
			"5",
			args{strings.NewReader(`5`)},
			`4`,
			false,
			true,
		},
		{
			"6",
			args{strings.NewReader(`6`)},
			`2`,
			false,
			true,
		},
		{
			"7",
			args{strings.NewReader(`7`)},
			`6`,
			false,
			true,
		},
		{
			"8",
			args{strings.NewReader(`8`)},
			`4`,
			false,
			true,
		},
		{
			"9",
			args{strings.NewReader(`9`)},
			`6`,
			false,
			true,
		},
		{
			"10",
			args{strings.NewReader(`10`)},
			`4`,
			false,
			true,
		},
		{
			"54",
			args{strings.NewReader(`54`)},
			`18`,
			false,
			true,
		},
		{
			"96",
			args{strings.NewReader(`96`)},
			`32`,
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
