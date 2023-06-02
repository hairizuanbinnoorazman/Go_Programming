package main

import "testing"

func FuzzPositiveNum(f *testing.F) {
	// seedNum := 90
	// f.Add(seedNum)
	f.Fuzz(func(t *testing.T, a int) { // fuzz target
		if a != PositiveNum(a) {
			t.Fail()
		}
	})
}

func TestLol(t *testing.T) {
	type args struct {
		a int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PositiveNum(tt.args.a); got != tt.want {
				t.Errorf("Lol() = %v, want %v", got, tt.want)
			}
		})
	}
}
