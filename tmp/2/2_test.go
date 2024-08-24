package main

import "testing"

func Test_changeStrings(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := changeStrings(tt.args.s1, tt.args.s2)
			if got != tt.want {
				t.Errorf("changeStrings() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("changeStrings() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Benchmark_changeStrings(b *testing.B) {
	s1, s2 := "Stepik", "Hello"
	for i := 0; i < b.N; i++ {
		s1, s2 = changeStrings(s1, s2)
		_, _ = s1, s2
	}
}
