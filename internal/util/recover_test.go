package util

import (
	"testing"
)

func TestRecover_NoPanic(t *testing.T) {
	// Arrange.
	funcGotExecuted := false

	// Act.
	panicked := Recover(func() {
		funcGotExecuted = true
	})

	// Assert.
	if !funcGotExecuted {
		t.Fatalf("expected true but got false")
	}
	if panicked {
		t.Fatalf("expected false but got true")
	}
}

func TestRecover_Panic(t *testing.T) {
	// Arrange.
	funcGotExecuted := false

	// Act.
	panicked := Recover(func() {
		funcGotExecuted = true
		panic("panic!!!")
	})

	// Assert.
	if !funcGotExecuted {
		t.Fatalf("expected true but got false")
	}
	if !panicked {
		t.Fatalf("expected true but got false")
	}
}
