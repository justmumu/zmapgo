package zmapgo

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestWithContext_MultiplePassing(t *testing.T) {
	t.Log("Testing WithContext function with multiple passing")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err := NewBlockingScanner(
		WithContext(ctx),
		WithContext(ctx),
	)

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed WithContext more than one")
	}
}

func TestWithBinaryPath_MultiplePassing(t *testing.T) {
	t.Log("Testing WithBinaryPath function with multiple passing")
	out, err := exec.Command("which", "zmap").Output()
	if err != nil {
		t.Errorf("Error while finding zmap binary: %v", err)
	}
	zmapBinary := strings.Trim(string(out), "\n")
	t.Logf("Found zmap binary path: %s", zmapBinary)

	_, err = NewBlockingScanner(
		WithBinaryPath(zmapBinary),
		WithBinaryPath(zmapBinary),
	)

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed WithBinaryPath more than one")
	}
}

func TestWithBinaryPath_PathNotExists(t *testing.T) {
	t.Log("Testing WithBinaryPath function with non-existing path")
	_, err := NewBlockingScanner(WithBinaryPath("/path/to/not/exists"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-existing path")
	}
}

func TestWithBinaryPath_WrongBinary(t *testing.T) {
	t.Log("Testing WithBinaryPath function with passing a non-zmap binary path")
	// find go binary
	out, err := exec.Command("which", "go").Output()
	if err != nil {
		t.Errorf("Error while finding go binary: %s", err)
	}
	goBinary := strings.Trim(string(out), "\n")
	t.Logf("Found go binary path: %s", goBinary)

	_, err = NewBlockingScanner(WithBinaryPath(goBinary))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error returned when passed non-zmap binary path")
	}
}

func TestWithBinaryPath_NormalBehavior(t *testing.T) {
	t.Log("Testing WithBinaryPath function under normal behavior")
	// find zmap binary
	out, err := exec.Command("which", "zmap").Output()
	if err != nil {
		t.Errorf("Error while finding zmap binary: %v", err)
	}
	zmapBinary := strings.Trim(string(out), "\n")
	t.Logf("Found zmap binary path: %s", zmapBinary)

	_, err = NewBlockingScanner(WithBinaryPath(zmapBinary))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned when passed real zmap binary path")
	}
}
