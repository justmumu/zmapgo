package zmapgo

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockingScanner_ZmapNotInstalled(t *testing.T) {
	t.Log("Testing NewBlockingScanner function with zmap binary is not found in $PATH")
	// Copy original env
	oldPathEnv := os.Getenv("PATH")
	_ = os.Setenv("PATH", "")

	s, err := NewBlockingScanner()
	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned if zmap binary is not found in PATH")
	}

	t.Logf("Returned scanner: %s", s)
	if s != nil {
		t.Error("Expected NewBlockingScanner to return a nil scanner if zmap is not found in $PATH")
	}

	// Fix env
	_ = os.Setenv("PATH", oldPathEnv)
}

func TestBlockingScanner_NormalBehavior(t *testing.T) {
	t.Log("Testing NewBlockingScanner function under normal behavior")
	_, err := exec.LookPath("zmap")
	if err != nil {
		panic("zmap is required to run this test")
	}

	scanner, err := NewBlockingScanner()
	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}

	t.Logf("Returned scanner: %s", scanner)
	if scanner == nil {
		t.Error("Expected NewBlockingScanner to return a non-nil scanner while under normal behavior")
	}

}

func TestNewAsyncScanner_ZmapNotInstalled(t *testing.T) {
	t.Log("Testing NewAsyncScanner function with zmap binary is not found in $PATH")
	// Copy original env
	oldPathEnv := os.Getenv("PATH")
	_ = os.Setenv("PATH", "")

	s, err := NewAsyncScanner()
	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned if zmap binary is not found in PATH")
	}

	t.Logf("Returned scanner: %s", s)
	if s != nil {
		t.Error("Expected NewAsyncScanner to return a nil scanner if zmap is not found in $PATH")
	}

	// Fix env
	_ = os.Setenv("PATH", oldPathEnv)
}

func TestAsyncScanner_WrongBinaryPath(t *testing.T) {
	t.Log("Testing NewAsyncScanner function with wrong binary path")

	_, err := NewAsyncScanner(WithBinaryPath("/wrong/binary/path"))
	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed wrong binary path")
	}
}

func TestAsyncScanner_NormalBehavior(t *testing.T) {
	t.Log("Testing NewAsyncScanner function under normal behavior")
	_, err := exec.LookPath("zmap")
	if err != nil {
		panic("zmap is required to run this test")
	}

	scanner, err := NewAsyncScanner()
	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}

	t.Logf("Returned scanner: %s", scanner)
	if scanner == nil {
		t.Error("Expected NewAsyncScanner to return a non-nil scanner while under normal behavior")
	}
}

func TestRunBlocking(t *testing.T) {
	testLogFilePath := "/tmp/test-log-file.txt"
	testLogDirectoryPath := "/tmp/test-log-directory"
	testOutputFilePath := "/tmp/test-output-file.txt"

	if _, err := os.Stat(testLogDirectoryPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(testLogDirectoryPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	if _, err := exec.LookPath("zmap"); err != nil {
		panic("zmap is required to run those tests")
	}

	tests := []struct {
		testDesc    string
		initOptions []InitOption
		options     []Option

		testTimeout bool

		dontCheckResultButShouldReturned bool
		expectedResult                   []map[string]interface{}
		isErrorExpected                  bool
		isTracesExpected                 bool
		isDebugsExpected                 bool
		isWarningsExpected               bool
		isInfosExpected                  bool
		isFatalsExpected                 bool
	}{
		{
			testDesc: "Context Timeout",
			options: []Option{
				// Don't actually run the test
				WithDryrun(),
			},
			testTimeout:        true,
			isErrorExpected:    true,
			isTracesExpected:   false,
			isDebugsExpected:   false,
			isWarningsExpected: false,
			isInfosExpected:    false,
			isFatalsExpected:   false,
		},
		{
			testDesc: "With Log File",
			options: []Option{
				WithTargets("1.1.1.0/24"),
				WithTargetPort("80"),
				WithLogFile(testLogFilePath),
				WithDryrun(),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
		},
		{
			testDesc: "With Log Directory First",
			options: []Option{
				WithTargets("1.1.1.0/24"),
				WithTargetPort("80"),
				WithLogDirectory(testLogDirectoryPath),
				WithDryrun(),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
		},
		{
			testDesc: "With Log Directory Second",
			options: []Option{
				WithTargets("1.1.1.0/24"),
				WithTargetPort("80"),
				WithLogDirectory(testLogDirectoryPath),
				WithDryrun(),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
		},
		{
			testDesc: "With Output File",
			options: []Option{
				WithTargets("1.1.1.0/24"),
				WithTargetPort("80"),
				WithOutputFile(testOutputFilePath),
				WithDryrun(),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
		},
		{
			testDesc: "With Output Fields",
			options: []Option{
				WithTargets("1.1.1.0/24"),
				WithTargetPort("80"),
				WithOutputFields([]string{"saddr", "daddr"}),
				WithDryrun(),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
		},
		{
			testDesc: "With Verbosity",
			options: []Option{
				WithTargets("1.1.1.0/24"),
				WithTargetPort("80"),
				WithVerbosity(VerbosityLevel3),
				WithDryrun(),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   false,
			isDebugsExpected:   false,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
		},
		{
			testDesc: "With No Dryrun",
			options: []Option{
				WithTargets("1.1.1.1/32"),
				WithTargetPort("80"),
				WithOutputFields([]string{"saddr", "sport"}),
				WithRate("10000"),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
			expectedResult: []map[string]interface{}{
				{
					"saddr": "1.1.1.1",
					"sport": "80",
				},
			},
		},
		{
			testDesc: "With No Dryrun And With Output File And One Output Field",
			options: []Option{
				WithTargets("1.1.1.1/32"),
				WithTargetPort("80"),
				WithOutputFile(testLogFilePath),
				WithOutputFields([]string{"saddr"}),
				WithRate("10000"),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
			expectedResult: []map[string]interface{}{
				{
					"saddr": "1.1.1.1",
				},
			},
		},
		{
			testDesc: "With No Dryrun And One Output Field",
			options: []Option{
				WithTargets("1.1.1.1/32"),
				WithTargetPort("80"),
				WithOutputFields([]string{"saddr"}),
				WithRate("10000"),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
			expectedResult: []map[string]interface{}{
				{
					"saddr": "1.1.1.1",
				},
			},
		},
		{
			testDesc: "With No Dryrun And With Output File And Two Output Field",
			options: []Option{
				WithTargets("1.1.1.1/32"),
				WithTargetPort("80"),
				WithOutputFile(testLogFilePath),
				WithOutputFields([]string{"saddr", "sport"}),
				WithRate("10000"),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
			expectedResult: []map[string]interface{}{
				{
					"saddr": "1.1.1.1",
					"sport": "80",
				},
			},
		},
		{
			testDesc: "With No Dryrun And Two Output Field",
			options: []Option{
				WithTargets("1.1.1.1/32"),
				WithTargetPort("80"),
				WithOutputFields([]string{"saddr", "sport"}),
				WithRate("10000"),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
			expectedResult: []map[string]interface{}{
				{
					"saddr": "1.1.1.1",
					"sport": "80",
				},
			},
		},
		{
			testDesc: "With No Dryrun And No Output Field And With Output File",
			options: []Option{
				WithTargets("1.1.1.1/32"),
				WithTargetPort("80"),
				WithOutputFile(testLogFilePath),
				WithRate("10000"),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
			expectedResult: []map[string]interface{}{
				{
					"saddr": "1.1.1.1",
					"sport": "80",
				},
			},
		},
		{
			testDesc: "With No Dryrun And No Output Field And With No Output File",
			options: []Option{
				WithTargets("1.1.1.1/32"),
				WithTargetPort("80"),
				WithRate("10000"),
				WithCooldownTime("2"),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    true,
			isFatalsExpected:   false,
			expectedResult: []map[string]interface{}{
				{
					"saddr": "1.1.1.1",
					"sport": "80",
				},
			},
		},
		{
			testDesc: "With Fatal Because of Blacklisting",
			options: []Option{
				WithTargets("10.0.0.0/24"),
				WithTargetPort("80"),
				WithRate("10000"),
				WithCooldownTime("2"),
				WithDryrun(),
			},
			testTimeout:        false,
			isErrorExpected:    false,
			isTracesExpected:   true,
			isDebugsExpected:   true,
			isWarningsExpected: true,
			isInfosExpected:    false,
			isFatalsExpected:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.testDesc, func(t *testing.T) {
			if test.testTimeout {
				ctx, cancel := context.WithTimeout(context.Background(), 99*time.Hour)
				test.initOptions = append(test.initOptions, WithContext(ctx))

				go (func() {
					// Cancel context to force timeout
					defer cancel()
					time.Sleep(1 * time.Millisecond)
				})()
			}

			scanner, err := NewBlockingScanner(test.initOptions...)
			if err != nil {
				panic(err) // this is never supposed to err, as we are testing run and not new.
			}

			err = scanner.AddOptions(test.options...)
			if err != nil {
				panic(err) // this is never supposed to err, as we are testing run and not new.
			}

			results, traces, debugs, warnings, infos, fatals, err := scanner.RunBlocking()
			t.Logf("Returned Error: %v", err)
			t.Logf("Returned Traces: %s", traces)
			t.Logf("Returned Debugs: %s", debugs)
			t.Logf("Returned Warnings: %s", warnings)
			t.Logf("Returned Infos: %s", infos)
			t.Logf("Returned Fatals: %s", fatals)
			t.Logf("Returned Results: %s", results)
			if !assert.Equal(t, test.isErrorExpected, err != nil) {
				return
			}

			// Check for traces
			if len(traces) == 0 && test.isTracesExpected == true {
				t.Error("Expected traces is returned. But, got length 0")
			} else if len(traces) != 0 && test.isTracesExpected == false {
				t.Error("Expected traces is empty. But, got length more than zero.")
			}

			// Check for debugs
			if len(debugs) == 0 && test.isDebugsExpected == true {
				t.Error("Expected debugs is returned. But, got length 0")
			} else if len(debugs) != 0 && test.isDebugsExpected == false {
				t.Error("Expected debugs is empty. But, got length more than zero.")
			}

			// Check for warnings
			if len(warnings) == 0 && test.isWarningsExpected == true {
				t.Error("Expected warnings is returned. But, got length 0")
			} else if len(warnings) != 0 && test.isWarningsExpected == false {
				t.Error("Expected warnings is empty. But, got length more than zero.")
			}

			// Check for infos
			if len(infos) == 0 && test.isInfosExpected == true {
				t.Error("Expected infos is returned. But, got length 0")
			} else if len(infos) != 0 && test.isInfosExpected == false {
				t.Error("Expected infos is empty. But, got length more than zero.")
			}

			// Check for fatals
			if len(fatals) == 0 && test.isFatalsExpected == true {
				t.Error("Expected fatals is returned. But, got length 0")
			} else if len(fatals) != 0 && test.isFatalsExpected == false {
				t.Error("Expected fatals is empty. But, got length more than zero.")
			}

			// Check for results
			if results == nil && test.expectedResult != nil {
				t.Error("Expected non-nil result. But, got nil")
			} else if results != nil && test.expectedResult == nil {
				t.Error("Expected nil result, But, got non-nil")
			}

			// CleanUp Section
			// clean test log file if created.
			if _, err := os.Stat(testLogFilePath); !errors.Is(err, os.ErrNotExist) {
				if err := os.Remove(testLogFilePath); err != nil {
					panic(err)
				}
			}

			// clean test output file if created.
			if _, err := os.Stat(testOutputFilePath); !errors.Is(err, os.ErrNotExist) {
				if err := os.Remove(testOutputFilePath); err != nil {
					panic(err)
				}
			}
		})
	}

	// Clean Up
	if _, err := os.Stat(testLogDirectoryPath); !errors.Is(err, os.ErrNotExist) {
		err := os.RemoveAll(testLogDirectoryPath)
		if err != nil {
			panic(err)
		}
	}
}

func TestRunAsync(t *testing.T) {
	if _, err := exec.LookPath("zmap"); err != nil {
		panic("zmap is required to run those tests")
	}

	scanner, err := NewAsyncScanner()
	if err != nil {
		t.Error("Expected that error is not returned while runing TestRunAsync")
		return
	}

	if err := scanner.AddOptions(
		WithTargets("1.1.1.1/32"),
		WithTargetPort("80"),
		WithRate("10000"),
		WithCooldownTime("2"),
		WithDryrun(),
	); err != nil {
		t.Error("Expected that error is not returned while adding option to AsyncScanner")
		return
	}

	if err := scanner.RunAsync(); err != nil {
		t.Error("Expected that error is not returned while running RunAsync function")
		return
	}

	// Block until finished
	if err := scanner.Wait(); err != nil {
		t.Error("Expected that error is not returned while waiting AsyncScanner finished")
	}

	traces := scanner.GetTraceMessages()
	if len(traces) == 0 {
		t.Error("Expected that trace messages returned even under dryrun")
	}

	debugs := scanner.GetDebugMessages()
	if len(debugs) == 0 {
		t.Error("Expected that debugs messages returned even under dryrun")
	}

	warnings := scanner.GetWarningMessages()
	if len(warnings) == 0 {
		t.Error("Expected that warnings messages returned even under dryrun")
	}

	infos := scanner.GetInfoMessages()
	if len(infos) == 0 {
		t.Error("Expected that infos messages returned even under dryrun")
	}

	fatals := scanner.GetFatalMessages()
	if len(fatals) != 0 {
		t.Error("Expected that fatals messages are not returned even under dryrun")
	}

	results := scanner.GetResults()
	if len(results) != 0 {
		t.Error("Expected that results are not returned because of dryrun")
	}
}

func TestBlockingScanner_GetVersion(t *testing.T) {
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Expected that error is not returned while creating BlockingScanner with NewBlockingScanner")
	}

	version, err := scanner.GetVersion()
	if err != nil {
		t.Error("Expected that error is not returned while getting version")
	}

	if version == "" {
		t.Error("Expected that version is returned.")
	}
}
