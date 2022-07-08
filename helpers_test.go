package zmapgo

import (
	"testing"
)

func TestMultiPassChecker_NormalBehavior_ReturnError(t *testing.T) {
	t.Log("Testing multiPassChecker function under normal behavior")

	args := []string{
		"--argument-one",
	}

	err := multiPassChecker(args, "--argument-one")
	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when checking argument already exists")
	}
}

func TestMultiPassChecker_NormalBehavior_NotReturnError(t *testing.T) {
	t.Log("Testing multiPassChecker function under normal behavior")

	args := []string{
		"--argument-one",
	}

	err := multiPassChecker(args, "--argument-two")
	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is returned when checking argument not exists")
	}
}
