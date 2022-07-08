package zmapgo

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// WithContext adds a context to a scanner, to make it cancellable and able to use timeout.
func WithContext(ctx context.Context) InitOption {
	return func(s *scanner) error {
		// check ctx already exists
		if s.ctx != nil {
			return errors.New("context is already exists")
		}
		s.ctx = ctx
		return nil
	}
}

// WithBinaryPath sets the zmap binary path for a scanner
func WithBinaryPath(binaryPath string) InitOption {
	return func(s *scanner) error {
		// check binary path already created
		if s.binaryPath != "" {
			return errors.New("binary path is already passed")
		}

		// check binary path exists
		if _, err := os.Stat(binaryPath); errors.Is(err, os.ErrNotExist) {
			return errors.New("given binary path does not exists")
		}

		// check real zmap binary
		out, _ := exec.Command(binaryPath, "--version").Output()

		trimed := strings.Trim(string(out), "\n")
		if !strings.Contains(trimed, "zmap") {
			return errors.New("given binary is not real zmap binary")
		}

		s.binaryPath = binaryPath
		return nil
	}
}
