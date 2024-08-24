package main

import "testing"

func Test_changeStrings(t *testing.T) {
	type args struct {
		s1 *string
		s2 *string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changeStrings(tt.args.s1, tt.args.s2)
		})
	}
}

func Benchmark_changeStrings(b *testing.B) {
	s1, s2 := "Stepik", "Hello"
	for i := 0; i < b.N; i++ {
		changeStrings(&s1, &s2)
		_, _ = s1, s2
	}
}
