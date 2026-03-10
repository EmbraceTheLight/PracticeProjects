package main

import "testing"

func TestDummy(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("math error")
	}
}
