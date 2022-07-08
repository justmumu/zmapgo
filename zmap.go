package zmapgo

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type BlockingScanner interface {
	AddOptions(options ...Option) error
	RunBlocking() (results []map[string]interface{}, traces []LogLine, debugs []LogLine, warnings []LogLine, infos []LogLine, fatals []LogLine, err error)
	ListProbeModules() ([]string, error)
	ListOutputModules() ([]string, error)
	ListOutputFields() ([]OutputField, error)
	GetVersion() (string, error)
}

type AsyncScanner interface {
	AddOptions(options ...Option) error
	RunAsync() error
	Wait() error
	GetTraceMessages() []LogLine
	GetDebugMessages() []LogLine
	GetWarningMessages() []LogLine
	GetInfoMessages() []LogLine
	GetFatalMessages() []LogLine
	GetResults() []map[string]interface{}
	ListProbeModules() ([]string, error)
	ListOutputModules() ([]string, error)
	ListOutputFields() ([]OutputField, error)
	GetVersion() (string, error)
}

// InitOptions is initialization option for the Scanner.
// Ex: WithBinaryPath, WithContext..
type InitOption func(*scanner) error

// Options is a function that is used for grouping of Scanner options.
// Option adds or remove zmap command line arguments.
type Option func(*scanner) error

type OutputField struct {
	Name        string
	Type        string
	Explanation string
}

type LogLine struct {
	LogTime time.Time
	LogType string
	Message string
}

// Scanner is represents the zmap scanner.
type scanner struct {
	args       []string
	binaryPath string
	ctx        context.Context

	waiter sync.WaitGroup

	asyncError   error
	asyncTrace   []LogLine
	asyncDebug   []LogLine
	asyncWarning []LogLine
	asyncInfo    []LogLine
	asyncFatal   []LogLine
	asyncResults []map[string]interface{}
}

// Creates new Scanner Interface
func NewBlockingScanner(initOptions ...InitOption) (BlockingScanner, error) {
	sc := &scanner{}

	for _, initOption := range initOptions {
		if err := initOption(sc); err != nil {
			return nil, err
		}
	}

	// After this block binaryPath filled.
	if sc.binaryPath == "" {
		var err error
		sc.binaryPath, err = exec.LookPath("zmap")
		if err != nil {
			return nil, ErrZmapNotInstalled
		}
	}

	// create ctx if not already created
	if sc.ctx == nil {
		sc.ctx = context.Background()
	}

	return sc, nil
}

func NewAsyncScanner(initOptions ...InitOption) (AsyncScanner, error) {
	sc := &scanner{}

	for _, initOption := range initOptions {
		if err := initOption(sc); err != nil {
			return nil, err
		}
	}

	// After this block binaryPath filled.
	if sc.binaryPath == "" {
		var err error
		sc.binaryPath, err = exec.LookPath("zmap")
		if err != nil {
			return nil, ErrZmapNotInstalled
		}
	}

	// create ctx if not already created
	if sc.ctx == nil {
		sc.ctx = context.Background()
	}

	return sc, nil
}

func (s *scanner) AddOptions(options ...Option) error {
	for _, option := range options {
		if err := option(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *scanner) RunBlocking() (results []map[string]interface{}, traces []LogLine, debugs []LogLine, warnings []LogLine, infos []LogLine, fatals []LogLine, err error) {
	var (
		stdout, stderr bytes.Buffer
	)

	var (
		dryrunPassed         bool = false
		outputFilePassed     bool = false
		logFilePassed        bool = false
		logDirectoryPassed   bool = false
		verbosityLevelPassed bool = false
		outputFieldsPassed   bool = false
	)

	// Look for --dryrun
	_, err = s.getArgument("--dryrun")
	if err == nil {
		dryrunPassed = true
	}

	// Look for --log-file
	logFilePath, err := s.getArgument("--log-file")
	if err == nil {
		logFilePassed = true
	}

	if logFilePassed {
		logFilePath, err = filepath.Abs(logFilePath)
		if err != nil {
			return nil, traces, debugs, warnings, infos, fatals, err
		}
	}

	// Look for --log-directory
	logDirectoryPath, err := s.getArgument("--log-directory")
	if err == nil {
		logDirectoryPassed = true
	}

	if logDirectoryPassed {
		logDirectoryPath, err = filepath.Abs(logDirectoryPath)
		if err != nil {
			return nil, traces, debugs, warnings, infos, fatals, err
		}
	}

	// Look for --verbosity
	_, err = s.getArgument("--verbosity")
	if err == nil {
		verbosityLevelPassed = true
	}

	if !verbosityLevelPassed {
		optionFunc := WithVerbosity(VerbosityLevel5)
		err = optionFunc(s)
		if err != nil {
			return nil, traces, debugs, warnings, infos, fatals, err
		}
	}

	// look for --output-file
	outputFilePath, err := s.getArgument("--output-file")
	if err == nil {
		outputFilePassed = true
	}

	if outputFilePassed {
		outputFilePath, err = filepath.Abs(outputFilePath)
		if err != nil {
			return nil, traces, debugs, warnings, infos, fatals, err
		}
	}

	// look for --output-fields
	outputFields, err := s.getArgument("--output-fields")
	if err == nil {
		outputFieldsPassed = true
	}

	if !outputFieldsPassed {
		availableOutputFields, err := s.ListOutputFields()
		if err != nil {
			return nil, traces, debugs, warnings, infos, fatals, err
		}
		var newOutputFields []string
		for _, aFields := range availableOutputFields {
			newOutputFields = append(newOutputFields, aFields.Name)
		}
		optionFunc := WithOutputFields(newOutputFields)
		err = optionFunc(s)
		if err != nil {
			return nil, traces, debugs, warnings, infos, fatals, err
		}
	}

	args := s.args

	// Prepare zmap process
	cmd := exec.Command(s.binaryPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run zmap process
	err = cmd.Start()
	if err != nil {
		return nil, traces, debugs, warnings, infos, fatals, err
	}

	// Make a goroutine to notify the select when the scan is done.
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-s.ctx.Done():
		// Context was done before the scan was finished.
		// The process is killed and a timeout error is returned.
		_ = cmd.Process.Kill()
		return nil, traces, debugs, warnings, infos, fatals, ErrScanTimeout
	case <-done:
		// Process zmap is done.
		// Output will be parsed according to passing arguments.

		// Start Log Parsing
		if logFilePassed {
			// Then parse Trace, Debug, Warning, Info and Fatal Message from log file.
			logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, os.ModePerm)
			if err != nil {
				return nil, traces, debugs, warnings, infos, fatals, err
			}
			defer logFile.Close()

			traces, debugs, warnings, infos, fatals, err = s.parseLogs(logFile)
			if err != nil {
				return nil, traces, debugs, warnings, infos, fatals, err
			}
		}

		if logDirectoryPassed {
			// Then parse Trace, Debug, Warning, Info and Fatal Message from latest created file in log directory.
			fileInfos, err := ioutil.ReadDir(logDirectoryPath)
			if err != nil {
				return nil, traces, debugs, warnings, infos, fatals, err
			}
			// Find latest log file
			timeLayout := "zmap-2006-01-02T150405-0700.log"
			var latestFile fs.FileInfo
			for _, fileInfo := range fileInfos {
				if !fileInfo.IsDir() {
					tNew, err := time.Parse(timeLayout, fileInfo.Name())
					if err == nil {
						if latestFile != nil {
							tOld, _ := time.Parse(timeLayout, latestFile.Name())
							if tNew.After(tOld) {
								latestFile = fileInfo
							}
						} else {
							latestFile = fileInfo
						}
					}
				}
			}
			// Now Parse Log file.
			logDirectoryFile, err := os.OpenFile(filepath.Join(logDirectoryPath, latestFile.Name()), os.O_RDONLY, os.ModePerm)
			if err != nil {
				return nil, traces, debugs, warnings, infos, fatals, err
			}
			defer logDirectoryFile.Close()

			traces, debugs, warnings, infos, fatals, err = s.parseLogs(logDirectoryFile)
			if err != nil {
				return nil, traces, debugs, warnings, infos, fatals, err
			}
		}

		if !logDirectoryPassed && !logFilePassed {
			// Then parse Trace, Debug, Warning, Info and Fatal Message from stderr
			traces, debugs, warnings, infos, fatals, err = s.parseLogs(&stderr)
			if err != nil {
				return nil, traces, debugs, warnings, infos, fatals, err
			}
		}
		// End Log Parsing

		// Start Result parsing
		if !dryrunPassed {
			if outputFieldsPassed {
				// User passed output fields.
				// If user requested one field. zmap will not add csv headers.
				outputFieldsSplitted := strings.Split(outputFields, ",")

				if len(outputFieldsSplitted) == 1 {
					// There is no csv headers
					if outputFilePassed {
						// Read from file, but add csv header manually.
						outputFile, err := os.Open(outputFilePath)
						if err != nil {
							return nil, traces, debugs, warnings, infos, fatals, err
						}
						defer outputFile.Close()

						sc := bufio.NewScanner(outputFile)

						for sc.Scan() {
							line := sc.Text()
							result := map[string]interface{}{
								outputFieldsSplitted[0]: line,
							}
							results = append(results, result)
						}
					} else {
						// Read from stdout, but add csv header manually.
						sc := bufio.NewScanner(&stdout)

						for sc.Scan() {
							line := sc.Text()
							result := map[string]interface{}{
								outputFieldsSplitted[0]: line,
							}
							results = append(results, result)
						}
					}
				} else {
					// There is csv headers
					if outputFilePassed {
						// Read from file
						outputFile, err := os.Open(outputFilePath)
						if err != nil {
							return nil, traces, debugs, warnings, infos, fatals, err
						}
						defer outputFile.Close()

						results, err = s.parseCsvOutputFile(outputFile)
						if err != nil {
							return nil, traces, debugs, warnings, infos, fatals, err
						}
					} else {
						// Read from stdout
						// Stdout contains only results.
						results, err = s.parseCsvOutputFile(&stdout)
						if err != nil {
							return nil, traces, debugs, warnings, infos, fatals, err
						}
					}
				}
			} else {
				// We added all output fields by default.
				// So there is csv headers.
				if outputFilePassed {
					outputFile, err := os.Open(outputFilePath)
					if err != nil {
						return nil, traces, debugs, warnings, infos, fatals, err
					}
					defer outputFile.Close()

					results, err = s.parseCsvOutputFile(outputFile)
					if err != nil {
						return nil, traces, debugs, warnings, infos, fatals, err
					}
				} else {
					// Read from stdout
					// Stdout contains only results.
					results, err = s.parseCsvOutputFile(&stdout)
					if err != nil {
						return nil, traces, debugs, warnings, infos, fatals, err
					}
				}
			}
		}
	}
	return
}

func (s *scanner) RunAsync() error {
	s.waiter.Add(1)
	go func() {
		results, traces, debugs, warnings, infos, fatals, err := s.RunBlocking()
		if err != nil {
			s.asyncError = err
			s.waiter.Done()
			return
		}
		s.asyncResults = results
		s.asyncTrace = traces
		s.asyncDebug = debugs
		s.asyncWarning = warnings
		s.asyncInfo = infos
		s.asyncFatal = fatals
		s.waiter.Done()
	}()
	return nil
}

func (s *scanner) Wait() error {
	// Wait until all wait group finished.
	s.waiter.Wait()
	return s.asyncError
}

func (s *scanner) GetTraceMessages() []LogLine {
	return s.asyncTrace
}

func (s *scanner) GetDebugMessages() []LogLine {
	return s.asyncDebug
}

func (s *scanner) GetWarningMessages() []LogLine {
	return s.asyncWarning
}

func (s *scanner) GetInfoMessages() []LogLine {
	return s.asyncInfo
}

func (s *scanner) GetFatalMessages() []LogLine {
	return s.asyncFatal
}

func (s *scanner) GetResults() []map[string]interface{} {
	return s.asyncResults
}

func (s *scanner) ListProbeModules() ([]string, error) {
	returnResults, err := exec.Command(s.binaryPath, "--list-probe-modules").Output()
	if err != nil {
		return nil, err
	}
	splitted := strings.Split(string(returnResults), "\n")

	var results []string
	for _, retRes := range splitted {
		if retRes != "" {
			results = append(results, retRes)
		}
	}
	return results, nil
}

func (s *scanner) ListOutputModules() ([]string, error) {
	returnResults, err := exec.Command(s.binaryPath, "--list-output-modules").Output()
	if err != nil {
		return nil, err
	}
	splitted := strings.Split(string(returnResults), "\n")
	var results []string
	for _, retRes := range splitted {
		if retRes != "" {
			results = append(results, retRes)
		}
	}
	return results, nil
}

func (s *scanner) ListOutputFields() ([]OutputField, error) {
	returnResults, err := exec.Command(s.binaryPath, "--list-output-fields").Output()
	if err != nil {
		return nil, err
	}
	splitted := strings.Split(string(returnResults), "\n")
	var results []OutputField
	for _, line := range splitted {
		if line != "" {
			lineSplitted := strings.Fields(line)
			var result OutputField
			result.Name = lineSplitted[0]
			result.Type = lineSplitted[1]
			for i := 2; i < len(lineSplitted); i++ {
				result.Explanation += fmt.Sprintf("%s ", lineSplitted[i])
			}
			result.Explanation = strings.Trim(result.Explanation, " ")
			results = append(results, result)
		}
	}
	return results, nil
}

func (s *scanner) GetVersion() (string, error) {
	returnResult, err := exec.Command(s.binaryPath, "--version").Output()
	if err != nil {
		return "", err
	}
	trimmedResult := strings.Trim(string(returnResult), "\n")

	if !strings.Contains(trimmedResult, "zmap") {
		return "", errors.New("not a zmap binary")
	}

	var newVersionSlice []string
	versionFields := strings.Fields(trimmedResult)
	for _, versionField := range versionFields {
		if versionField != "zmap" {
			newVersionSlice = append(newVersionSlice, versionField)
		}
	}

	return strings.Join(newVersionSlice, " "), nil
}

func (s *scanner) parseLogs(ioReader io.Reader) (traces []LogLine, debugs []LogLine, warnings []LogLine, infos []LogLine, fatals []LogLine, err error) {
	sc := bufio.NewScanner(ioReader)
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.Contains(line, "[TRACE]"):
			logLine, err := s.parseLogLine(line)
			if err != nil {
				return traces, debugs, warnings, infos, fatals, err
			}
			traces = append(traces, logLine)
		case strings.Contains(line, "[DEBUG]"):
			logLine, err := s.parseLogLine(line)
			if err != nil {
				return traces, debugs, warnings, infos, fatals, err
			}
			debugs = append(debugs, logLine)
		case strings.Contains(line, "[WARN]"):
			logLine, err := s.parseLogLine(line)
			if err != nil {
				return traces, debugs, warnings, infos, fatals, err
			}
			warnings = append(warnings, logLine)
		case strings.Contains(line, "[INFO]"):
			logLine, err := s.parseLogLine(line)
			if err != nil {
				return traces, debugs, warnings, infos, fatals, err
			}
			infos = append(infos, logLine)
		case strings.Contains(line, "[FATAL]"):
			logLine, err := s.parseLogLine(line)
			if err != nil {
				return traces, debugs, warnings, infos, fatals, err
			}
			fatals = append(fatals, logLine)
		}
	}
	return traces, debugs, warnings, infos, fatals, nil
}

func (s *scanner) parseLogLine(line string) (LogLine, error) {
	logTimeLayout := "Jan 02 15:04:05.000"
	logSplitted := strings.Split(line, " ")

	logTimeStr := strings.Join(logSplitted[:3], " ")
	logTime, err := time.Parse(logTimeLayout, logTimeStr)
	if err != nil {
		return LogLine{}, err
	}
	logType := strings.Replace(strings.Replace(logSplitted[3], "[", "", -1), "]", "", -1)
	logMessage := strings.Join(logSplitted[4:], " ")

	return LogLine{
		LogTime: logTime,
		LogType: logType,
		Message: logMessage,
	}, nil
}

func (s *scanner) parseCsvOutputFile(ioReader io.Reader) ([]map[string]interface{}, error) {
	reader := csv.NewReader(ioReader)

	var rows []map[string]interface{}
	var header []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]interface{}{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows, nil
}

func (s *scanner) getArgument(argument string) (string, error) {
	var (
		argumentValue string
	)

	for index, arg := range s.args {
		if arg == argument {
			if (index + 1) < len(s.args) {
				argumentValue = s.args[index+1]
				if strings.Contains(argumentValue, "--") {
					return "", nil
				} else {
					return argumentValue, nil
				}
			} else {
				return "", nil
			}
		}
	}
	return "", errors.New("argument not found")
}
