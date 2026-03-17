package main

import "testing"

func TestDummy(t *testing.T) {
	// Test Case 1
	if 1+1 != 2 {
		t.Fatal("math error")
	}

	// Test Case 2
	if 1+1 == 3 {
		t.Fatal("math error")
	}
}
