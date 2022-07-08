package zmapgo

import (
	"net"
	"runtime"
	"strconv"
	"testing"
)

func TestWithCustomArguments_NormalBehavior(t *testing.T) {
	t.Log("Testing WithCustomArguments function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithCustomArguments("--dryrun"))
	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithTargets_NotIPv4(t *testing.T) {
	t.Log("Testing WithTargets function with wrong format ipv4 address")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithTargets("192.168.1.256"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed wrong format ipv4 address")
	}
}

func TestWithTargets_WithIPv6(t *testing.T) {
	t.Log("Testing WithTargets function with ipv6 address")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithTargets("FE80:CD00:0000:0CDE:1257:0000:211E:729C"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed ipv6 address")
	}
}

func TestWithTargets_WithWrongCIDR(t *testing.T) {
	t.Log("Testing WithTargets function with wrong cidr notation")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithTargets("192.168.1.1/33"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed wrong cidr notation")
	}
}

func TestWithTargets_NormalBehavior(t *testing.T) {
	t.Log("Testing WithTargets function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithTargets("192.168.1.1", "192.168.1.1/24"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is returned while under normal behavior")
	}
}

func TestWithTargetPort_MultiplePassing(t *testing.T) {
	t.Log("Testing WithTargetPort function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithTargetPort("80"))
	err2 := scanner.AddOptions(WithTargetPort("80"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithTargetPort used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithTargetPort used twice")
	}
}

func TestWithTargetPort_PortValueNotNumeric(t *testing.T) {
	t.Log("Testing WithTargetPort function with non-numeric target port")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithTargetPort("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric string as target port")
	}
}

func TestWithTargetPort_WrongPortValue(t *testing.T) {
	t.Log("Testing WithTargetPort function with passing wrong target port")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	// send bigger then 65535
	err = scanner.AddOptions(WithTargetPort("70000"))

	t.Logf("Returned error when sending bigger than 65535: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed bigger than 65535 as target port")
	}

	scanner, err = NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	// send lower than 0
	err = scanner.AddOptions(WithTargetPort("-1"))

	t.Logf("Returned error when sending lower than 0: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed lower than 0 as target port")
	}
}

func TestWithTargetPort_NormalBehavior(t *testing.T) {
	t.Log("Testing WithTargetPort function under normal behavior")

	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithTargetPort("80"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithOutputFile_MultiplePassing(t *testing.T) {
	t.Log("Testing WithOutputFile function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithOutputFile("/etc/file.txt"))
	err2 := scanner.AddOptions(WithOutputFile("/etc/file.txt"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithOutputFile used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithOutputFile used twice")
	}
}

func TestWithOutputFile_WithStdoutValue(t *testing.T) {
	t.Log("Testing WithOutputFile function with passing value of '-' to pass output to stdout")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithOutputFile("-"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned when passed value of '-' to WithOutputFile function")
	}
}

func TestWithOutputFile_ParentNotExists(t *testing.T) {
	t.Log("Testing WithOutputFile function with passing parent directory does not exists")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithOutputFile("/non-existing-parent/passwd"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed parent directory does not exists")
	}
}

func TestWithOutputFile_PathIsDirectory(t *testing.T) {
	t.Log("Testing WithOutputFile function with passing directory path")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithOutputFile("/etc"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed directory path")
	}
}

func TestWithOutputFile_NormalBehavior(t *testing.T) {
	t.Log("Testing WithOutputFile function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithOutputFile("/etc/file.txt"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithBlacklistFile_MultiplePassing(t *testing.T) {
	t.Log("Testing WithBlacklistFile function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithBlacklistFile("/etc/passwd"))
	err2 := scanner.AddOptions(WithBlacklistFile("/etc/passwd"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithBlacklistFile used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithBlacklistFile used twice")
	}
}

func TestWithBlacklistFile_FilePathNotExists(t *testing.T) {
	t.Log("Testing WithBlacklistFile function with passing not existing file path as value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithBlacklistFile("/path/to/not/exist/file.txt"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed not existing file path as value")
	}
}

func TestWithBlacklistFile_NormalBehavior(t *testing.T) {
	t.Log("Testing WithBlacklistFile function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithBlacklistFile("/etc/passwd"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithWhitelistFile_MultiplePassing(t *testing.T) {
	t.Log("Testing WithWhitelistFile function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithWhitelistFile("/etc/passwd"))
	err2 := scanner.AddOptions(WithWhitelistFile("/etc/passwd"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithWhitelistFile used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithWhitelistFile used twice")
	}
}

func TestWithWhitelistFile_FilePathNotExists(t *testing.T) {
	t.Log("Testing WithWhitelistFile function with passing not existing file path as value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithWhitelistFile("/path/to/not/exist/file.txt"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed not existing file path as value")
	}
}

func TestWithWhitelistFile_NormalBehavior(t *testing.T) {
	t.Log("Testing WithWhitelistFile function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithWhitelistFile("/etc/passwd"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithRate_MultiplePassing(t *testing.T) {
	t.Log("Testing WithRate function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithRate("10000"))
	err2 := scanner.AddOptions(WithRate("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithRate used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithRate used twice")
	}
}

func TestWithRate_ValueNotNumeric(t *testing.T) {
	t.Log("Testing WithRate function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithRate("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as rate value")
	}
}

func TestWithRate_NormalBehavior(t *testing.T) {
	t.Log("Testing WithRate function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithRate("100000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithBandwidth_MultiplePassing(t *testing.T) {
	t.Log("Testing WithBandwidth function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithBandwidth("10000", UnitBandwidthBps))
	err2 := scanner.AddOptions(WithBandwidth("10000", UnitBandwidthBps))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithBandwidth used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithBandwidth used twice")
	}
}

func TestWithBandwidth_NonNumericBandwidth(t *testing.T) {
	t.Log("Testing WithBandwidth function with non numeric bandwidth value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithBandwidth("non-numeric-value", UnitBandwidthGbps))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as bandwidth value")
	}
}

func TestWithBandwidth_UnsupportedUnit(t *testing.T) {
	t.Log("Testing WithBandwidth function with unsupported unit")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	var unsupportedUnit BandwidthUnit = "unsupported"
	err = scanner.AddOptions(WithBandwidth("100000", unsupportedUnit))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed unsupported unit")
	}
}

func TestWithBandwidth_NormalBehavior(t *testing.T) {
	t.Log("Testing WithBandwidth function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithBandwidth("10", UnitBandwidthGbps))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithMaxTargets_MultiplePassing(t *testing.T) {
	t.Log("Testing WithMaxTargets function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithMaxTargets("10000", false))
	err2 := scanner.AddOptions(WithMaxTargets("10000", false))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithMaxTargets used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithMaxTargets used twice")
	}
}

func TestWithMaxTargets_NonNumericValue(t *testing.T) {
	t.Log("Testing WithMaxTargets function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxTargets("non-numeric-value", false))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as max target value")
	}
}

func TestWithMaxTargets_NormalBehavior(t *testing.T) {
	t.Log("Testing WithMaxTargets function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxTargets("10000", true))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithMaxRuntime_MultiplePassing(t *testing.T) {
	t.Log("Testing WithMaxRuntime function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithMaxRuntime("10000"))
	err2 := scanner.AddOptions(WithMaxRuntime("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithMaxRuntime used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithMaxRuntime used twice")
	}
}

func TestWithMaxRuntime_NonNumericValue(t *testing.T) {
	t.Log("Testing WithMaxRuntime function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxRuntime("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as max target value")
	}
}

func TestWithMaxRuntime_NormalBehavior(t *testing.T) {
	t.Log("Testing WithMaxRuntime function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxRuntime("10000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithMaxResults_MultiplePassing(t *testing.T) {
	t.Log("Testing WithMaxResults function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithMaxResults("10000"))
	err2 := scanner.AddOptions(WithMaxResults("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithMaxResults used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithMaxResults used twice")
	}
}

func TestWithMaxResults_NonNumericValue(t *testing.T) {
	t.Log("Testing WithMaxResults function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxResults("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as max target value")
	}
}

func TestWithMaxResults_NormalBehavior(t *testing.T) {
	t.Log("Testing WithMaxResults function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxResults("10000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithNumberOfProbesPerIP_MultiplePassing(t *testing.T) {
	t.Log("Testing WithNumberOfProbesPerIP function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithNumberOfProbesPerIP("10000"))
	err2 := scanner.AddOptions(WithNumberOfProbesPerIP("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithNumberOfProbesPerIP used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithNumberOfProbesPerIP used twice")
	}
}

func TestWithNumberOfProbesPerIP_NonNumericValue(t *testing.T) {
	t.Log("Testing WithNumberOfProbesPerIP function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithNumberOfProbesPerIP("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as max target value")
	}
}

func TestWithNumberOfProbesPerIP_NormalBehavior(t *testing.T) {
	t.Log("Testing WithNumberOfProbesPerIP function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithNumberOfProbesPerIP("10000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithCooldownTime_MultiplePassing(t *testing.T) {
	t.Log("Testing WithCooldownTime function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithCooldownTime("10000"))
	err2 := scanner.AddOptions(WithCooldownTime("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithCooldownTime used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithCooldownTime used twice")
	}
}

func TestWithCooldownTime_NonNumericValue(t *testing.T) {
	t.Log("Testing WithCooldownTime function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithCooldownTime("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as max target value")
	}
}

func TestWithCooldownTime_NormalBehavior(t *testing.T) {
	t.Log("Testing WithCooldownTime function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithCooldownTime("10000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithSeed_MultiplePassing(t *testing.T) {
	t.Log("Testing WithSeed function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithSeed("10000"))
	err2 := scanner.AddOptions(WithSeed("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithSeed used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithSeed used twice")
	}
}

func TestWithSeed_NonNumericValue(t *testing.T) {
	t.Log("Testing WithSeed function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSeed("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as max target value")
	}
}

func TestWithSeed_NormalBehavior(t *testing.T) {
	t.Log("Testing WithSeed function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSeed("10000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithMaxRetries_MultiplePassing(t *testing.T) {
	t.Log("Testing WithMaxRetries function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithMaxRetries("10000"))
	err2 := scanner.AddOptions(WithMaxRetries("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithMaxRetries used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithMaxRetries used twice")
	}
}

func TestWithMaxRetries_NonNumericValue(t *testing.T) {
	t.Log("Testing WithMaxRetries function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxRetries("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as max target value")
	}
}

func TestWithMaxRetries_NormalBehavior(t *testing.T) {
	t.Log("Testing WithMaxRetries function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxRetries("10000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithDryrun_MultiplePassing(t *testing.T) {
	t.Log("Testing WithDryrun function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithDryrun())
	err2 := scanner.AddOptions(WithDryrun())

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithDryrun used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithDryrun used twice")
	}
}

func TestWithDryrun_NormalBehavior(t *testing.T) {
	t.Log("Testing WithDryrun function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithDryrun())

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithTotalShards_MultiplePassing(t *testing.T) {
	t.Log("Testing WithTotalShards function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithTotalShards("10000"))
	err2 := scanner.AddOptions(WithTotalShards("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithTotalShards used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithTotalShards used twice")
	}
}

func TestWithTotalShards_NonNumericValue(t *testing.T) {
	t.Log("Testing WithTotalShards function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithTotalShards("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as max target value")
	}
}

func TestWithTotalShards_NormalBehavior(t *testing.T) {
	t.Log("Testing WithTotalShards function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithTotalShards("10000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithShardID_MultiplePassing(t *testing.T) {
	t.Log("Testing WithShardID function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithShardID("10000"))
	err2 := scanner.AddOptions(WithShardID("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithShardID used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithShardID used twice")
	}
}

func TestWithShardID_NonNumericValue(t *testing.T) {
	t.Log("Testing WithShardID function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithShardID("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value as max target value")
	}
}

func TestWithShardID_NormalBehavior(t *testing.T) {
	t.Log("Testing WithShardID function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithShardID("10000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior")
	}
}

func TestWithSourcePort_MultiplePassing(t *testing.T) {
	t.Log("Testing WithSourcePort function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithSourcePort("10000"))
	err2 := scanner.AddOptions(WithSourcePort("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithSourcePort used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithSourcePort used twice")
	}
}

func TestWithSourcePort_Range_NotValidRange(t *testing.T) {
	t.Log("Testing WithSourcePort function with not valid range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("1000-1003-1007"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-valid range definition")
	}
}

func TestWithSourcePort_Range_LowerNonNumeric(t *testing.T) {
	t.Log("Testing WithSourcePort function with non numeric value in lower part of range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("nonNumericValue-1003"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non numeric value in lower part of range definition")
	}
}

func TestWithSourcePort_Range_LowerNotValidPortNumber(t *testing.T) {
	t.Log("Testing WithSourcePort function with wrong value in lower part of range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("70000-1003"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed wrong value in lower part of range definition")
	}
}

func TestWithSourcePort_Range_GreaterNonNumeric(t *testing.T) {
	t.Log("Testing WithSourcePort function with non numeric value in greater part of range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("1003-nonNumericValue"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non numeric value in greater part of range definition")
	}
}

func TestWithSourcePort_Range_GreaterNotValidPortNumber(t *testing.T) {
	t.Log("Testing WithSourcePort function with wrong value in greater part of range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("1003-70000"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed wrong value in greater part of range definition")
	}
}

func TestWithSourcePort_Range_LowerAndGreaterEqual(t *testing.T) {
	t.Log("Testing WithSourcePort function with greater and lower part equal in range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("1003-1003"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed greater and lower part equal in range definition")
	}
}

func TestWithSourcePort_Range_LowerGreaterThanGreater(t *testing.T) {
	t.Log("Testing WithSourcePort function with lower value greater than greater part value in range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("1004-1003"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed lower value as greater than greater value")
	}
}

func TestWithSourcePort_Single_NonNumeric(t *testing.T) {
	t.Log("Testing WithSourcePort function with non-numeric value in single port definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("nonNumericValue"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value in single port definition")
	}
}

func TestWithSourcePort_Single_NotValidPortNumber(t *testing.T) {
	t.Log("Testing WithSourcePort function with wrong value in single port definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("70000"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed wrong value in single port definition")
	}
}

func TestWithSorucePort_Range_NormalBehavior(t *testing.T) {
	t.Log("Testing WithSourcePort function under normal behavior in range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("1000-1004"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior in range definition")
	}
}

func TestWithSorucePort_Single_NormalBehavior(t *testing.T) {
	t.Log("Testing WithSourcePort function under normal behavior in single port definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourcePort("10000"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while testing normal behavior in single definition")
	}
}

func TestWithSourceIP_MultiplePassing(t *testing.T) {
	t.Log("Testing WithSourceIP function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithSourceIP("192.168.1.1"))
	err2 := scanner.AddOptions(WithSourceIP("192.168.1.1"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithSourceIP used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithSourceIP used twice")
	}
}

func TestWithSourceIP_Range_NotValidRange(t *testing.T) {
	t.Log("Testing WithSourceIP function with not valid range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourceIP("192.168.1.1-192.168.1.5-192.168.1.8"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-valid range definition")
	}
}

func TestWithSourceIP_Range_NotValidLowerPart(t *testing.T) {
	t.Log("Testing WithSourceIP function with non-valid lower part")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourceIP("nonValidLower-192.168.1.5"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-valid lower part in range definition")
	}
}

func TestWithSourceIP_Range_NotValidGreaterPart(t *testing.T) {
	t.Log("Testing WithSourceIP function with non-valid greater part")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourceIP("192.168.1.5-nonValidGreater"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-valid greater part in range definition")
	}
}

func TestWithSourceIP_Range_NotValidLowerAndGreaterEqual(t *testing.T) {
	t.Log("Testing WithSourceIP function with lower and greater equal in range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourceIP("192.168.1.5-192.168.1.5"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed lower and greater equal in range definition")
	}
}

func TestWithSourceIP_Range_NotValidLowerGreaterThanGreater(t *testing.T) {
	t.Log("Testing WithSourceIP function with lower part greater than greater part in range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourceIP("192.168.1.6-192.168.1.5"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed lower part greater than greater part in range definition")
	}
}

func TestWithSourceIP_Single_NotValidIP(t *testing.T) {
	t.Log("Testing WithSourceIP function with non-valid ip in single definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourceIP("nonValidIP"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-valid ip in single definition")
	}
}

func TestWithSourceIP_Range_NormalBehavior(t *testing.T) {
	t.Log("Testing WithSourceIP function under normal behavior in range definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourceIP("192.168.1.1-192.168.1.5"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior in range definition")
	}
}

func TestWithSourceIP_Single_NormalBehavior(t *testing.T) {
	t.Log("Testing WithSourceIP function under normal behavior in single definition")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourceIP("192.168.1.1"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior in single definition")
	}
}

func TestWithGatewayMAC_MultiplePassing(t *testing.T) {
	t.Log("Testing WithGatewayMAC function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}
	ifas, err := net.Interfaces()
	if err != nil {
		t.Errorf("Error while getting interfaces: %v", err)
	}

	var passMAC string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			passMAC = a
			break
		}
	}

	t.Logf("Choosed MAC Address: %s", passMAC)

	err1 := scanner.AddOptions(WithGatewayMAC(passMAC))
	err2 := scanner.AddOptions(WithGatewayMAC(passMAC))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithGatewayMAC used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithGatewayMAC used twice")
	}
}

func TestWithGatewayMAC_NonValidMAC(t *testing.T) {
	t.Log("Testing WithGatewayMAC function with non-valid mac address")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithGatewayMAC("nonValidMAC"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-valid mac address")
	}
}

func TestWithGatewayMAC_NormalBehavior(t *testing.T) {
	t.Log("Test WithGatewayMAC function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}
	ifas, err := net.Interfaces()
	if err != nil {
		t.Errorf("Error while getting interfaces: %v", err)
	}

	var passMAC string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			passMAC = a
			break
		}
	}
	t.Logf("Choosed MAC Address: %s", passMAC)

	err = scanner.AddOptions(WithGatewayMAC(passMAC))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithSourceMAC_MultiplePassing(t *testing.T) {
	t.Log("Testing WithSourceMAC function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}
	ifas, err := net.Interfaces()
	if err != nil {
		t.Errorf("Error while getting interfaces: %v", err)
	}

	var passMAC string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			passMAC = a
			break
		}
	}
	t.Logf("Choosed MAC Address: %s", passMAC)

	err1 := scanner.AddOptions(WithSourceMAC(passMAC))
	err2 := scanner.AddOptions(WithSourceMAC(passMAC))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithSourceMAC used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithSourceMAC used twice")
	}
}

func TestWithSourceMAC_NonValidMAC(t *testing.T) {
	t.Log("Testing WithSourceMAC function with non-valid mac address")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSourceMAC("nonValidMAC"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-valid mac address")
	}
}

func TestWithSourceMAC_NormalBehavior(t *testing.T) {
	t.Log("Test WithSourceMAC function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}
	ifas, err := net.Interfaces()
	if err != nil {
		t.Errorf("Error while getting interfaces: %v", err)
	}

	var passMAC string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			passMAC = a
			break
		}
	}
	t.Logf("Choosed MAC Address: %s", passMAC)

	err = scanner.AddOptions(WithSourceMAC(passMAC))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithInterface_MultiplePassing(t *testing.T) {
	t.Log("Testing WithInterface function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	ifas, err := net.Interfaces()
	if err != nil {
		t.Errorf("Error while getting interfaces: %v", err)
	}

	var passInterface string
	for _, ifa := range ifas {
		if ifa.Name != "" {
			passInterface = ifa.Name
			break
		}
	}
	t.Logf("Choosed Interface Name: %s", passInterface)

	err1 := scanner.AddOptions(WithInterface(passInterface))
	err2 := scanner.AddOptions(WithInterface(passInterface))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithInterface used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithInterface used twice")
	}
}

func TestWithInterface_NonAvailableInterface(t *testing.T) {
	t.Log("Testing WithInterface function with non-available interface name on system")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithInterface("nonAvailableInterfaceName"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-available interface name")
	}
}

func TestWithInterface_NormalBehavior(t *testing.T) {
	t.Log("Testing WithInterface function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}
	ifas, err := net.Interfaces()
	if err != nil {
		t.Errorf("Error while getting interfaces: %v", err)
	}

	var passInterface string
	for _, ifa := range ifas {
		if ifa.Name != "" {
			passInterface = ifa.Name
			break
		}
	}
	t.Logf("Choosed Interface Name: %s", passInterface)

	err = scanner.AddOptions(WithInterface(passInterface))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithVPN_MultiplePassing(t *testing.T) {
	t.Log("Testing WithVPN function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithVPN())
	err2 := scanner.AddOptions(WithVPN())

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithVPN used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithVPN used twice")
	}
}

func TestWithVPN_NormalBehavior(t *testing.T) {
	t.Log("Testing WithVPN function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithVPN())

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithProbeModule_MultiplePassing(t *testing.T) {
	t.Log("Testing WithProbeModule function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithProbeModule("udp"))
	err2 := scanner.AddOptions(WithProbeModule("udp"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithProbeModule used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithProbeModule used twice")
	}
}

func TestWithProbeModule_NonAvailableModule(t *testing.T) {
	t.Log("Testing WithProbeModule function with non-available probe module name")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithProbeModule("nonAvailableProbeName"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-available probe module name")
	}
}

func TestWithProbeModule_EmptyString(t *testing.T) {
	t.Log("Testing WithProbeModule function with empty probe module name")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithProbeModule(""))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed empty probe module name")
	}
}

func TestWithProbeModule_NormalBehavior(t *testing.T) {
	t.Log("Testing WithProbeModule function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithProbeModule("udp"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithProbeArgs_MultiplePassing(t *testing.T) {
	t.Log("Testing WithProbeArgs function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithProbeArgs("file:packet.pkt"))
	err2 := scanner.AddOptions(WithProbeArgs("file:packet.pkt"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithProbeArgs used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithProbeArgs used twice")
	}
}

func TestWithProbeArgs_NormalBehavior(t *testing.T) {
	t.Log("Testing WithProbeArgs function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithProbeArgs("file:packet.pkt"))
	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithOutputFields_MultiplePassing(t *testing.T) {
	t.Log("Testing WithOutputFields function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	outputFields := []string{"saddr", "daddr"}
	err1 := scanner.AddOptions(WithOutputFields(outputFields))
	err2 := scanner.AddOptions(WithOutputFields(outputFields))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithOutputFields used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithOutputFields used twice")
	}
}

func TestWithOutputFields_NonAvailableModule(t *testing.T) {
	t.Log("Testing WithOutputFields function with non-available probe module name")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	outputFields := []string{"saddr", "nonAvailableOutputField"}
	err = scanner.AddOptions(WithOutputFields(outputFields))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-available probe module name")
	}
}

func TestWithOutputFields_NormalBehavior(t *testing.T) {
	t.Log("Testing WithOutputFields function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	outputFields := []string{"saddr", "daddr"}
	err = scanner.AddOptions(WithOutputFields(outputFields))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithOutputModule_MultiplePassing(t *testing.T) {
	t.Log("Testing WithOutputModule function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithOutputModule("csv"))
	err2 := scanner.AddOptions(WithOutputModule("csv"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithOutputModule used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithOutputModule used twice")
	}
}

func TestWithOutputModule_NonAvailableModule(t *testing.T) {
	t.Log("Testing WithOutputModule function with non-available probe module name")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithOutputModule("nonAvailableProbeName"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-available probe module name")
	}
}

func TestWithOutputModule_EmptyString(t *testing.T) {
	t.Log("Testing WithOutputModule function with empty probe module name")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithOutputModule(""))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed empty probe module name")
	}
}

func TestWithOutputModule_NormalBehavior(t *testing.T) {
	t.Log("Testing WithOutputModule function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithOutputModule("csv"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithOutputArgs_MultiplePassing(t *testing.T) {
	t.Log("Testing WithOutputArgs function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithOutputArgs("args"))
	err2 := scanner.AddOptions(WithOutputArgs("args"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithOutputArgs used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithOutputArgs used twice")
	}
}

func TestWithOutputArgs_NormalBehavior(t *testing.T) {
	t.Log("Testing WithOutputArgs function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithOutputArgs("someargs:someargsvalue"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithOutputFilter_MultiplePassing(t *testing.T) {
	t.Log("Testing WithOutputFilter function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithOutputFilter("filter"))
	err2 := scanner.AddOptions(WithOutputFilter("filter"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithOutputFilter used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithOutputFilter used twice")
	}
}

func TestWithOutputFilter_NormalBehavior(t *testing.T) {
	t.Log("Testing WithOutputFilter function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithOutputFilter("somefilter:somefiltervalue"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithVerbosity_MultiplePassing(t *testing.T) {
	t.Log("Testing WithVerbosity function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithVerbosity(VerbosityLevel5))
	err2 := scanner.AddOptions(WithVerbosity(VerbosityLevel5))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithVerbosity used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithVerbosity used twice")
	}
}

func TestWithVerbosity_WrongLevel(t *testing.T) {
	t.Log("Testing WithVerbosity function with wrong verbosity level")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	var wrongVerbosityLeveL VerbosityLevel = "wrongLevel"
	err = scanner.AddOptions(WithVerbosity(wrongVerbosityLeveL))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed wrong verbosity level")
	}
}

func TestWithVerbosity_NormalBehavior(t *testing.T) {
	t.Log("Testing WithVerbosity function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithVerbosity(VerbosityLevel5))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithLogFile_MultiplePassing(t *testing.T) {
	t.Log("Testing WithLogFile function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithLogFile("/trying.txt"))
	err2 := scanner.AddOptions(WithLogFile("/trying.txt"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithLogFile used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithLogFile used twice")
	}
}

func TestWithLogFile_LogDirectoryPassed(t *testing.T) {
	t.Log("Testing WithLogFile function with --log-directory passed")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(
		WithLogDirectory("/etc"),
		WithLogFile("./log-file.txt"),
	)

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed --log-directory with --log-file")
	}
}

func TestWithLogFile_PathIsADirectory(t *testing.T) {
	t.Log("Testing WithLogFile function with existing directory as value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithLogFile("/etc"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed existing directory as value")
	}
}

func TestWithLogFile_ParentDirNotExists(t *testing.T) {
	t.Log("Testing WithLogFile function with given path's parent directory not exists")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithLogFile("/notexistingparent/deneme.txt"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed not existing parent directory as value")
	}
}

func TestWithLogFile_NormalBehavior(t *testing.T) {
	t.Log("Testing WithLogFile function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithLogFile("/trying.txt"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithLogDirectory_MultiplePassing(t *testing.T) {
	t.Log("Testing WithLogDirectory function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithLogDirectory("/etc"))
	err2 := scanner.AddOptions(WithLogDirectory("/etc"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithLogDirectory used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithLogDirectory used twice")
	}
}

func TestWithLogDirectory_LogFilePassed(t *testing.T) {
	t.Log("Testing WithLogDirectory function with --log-file passed")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(
		WithLogFile("./log-file.txt"),
		WithLogDirectory("/etc"),
	)

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed --log-directory with --log-file")
	}
}

func TestWithLogDirectory_NonExistingPath(t *testing.T) {
	t.Log("Testing WithLogDirectory function with non-existing path")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(
		WithLogDirectory("/non-existing-path"),
	)

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-existing path")
	}
}

func TestWithLogDirectory_PathIsNotDirectory(t *testing.T) {
	t.Log("Testing WithLogDirectory function with path is not directory")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(
		WithLogDirectory("/etc/passwd"),
	)

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed path is not directory")
	}
}

func TestWithLogDirectory_NormalBehavior(t *testing.T) {
	t.Log("Testing WithLogDirectory function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(
		WithLogDirectory("/etc"),
	)

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithMetadataFile_MultiplePassing(t *testing.T) {
	t.Log("Testing WithMetadataFile function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithMetadataFile("/trying.txt"))
	err2 := scanner.AddOptions(WithMetadataFile("/trying.txt"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithMetadataFile used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithMetadataFile used twice")
	}
}

func TestWithMetadataFile_PathIsADirectory(t *testing.T) {
	t.Log("Testing WithMetadataFile function with existing directory as value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMetadataFile("/etc"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed existing directory as value")
	}
}

func TestWithMetadataFile_ParentDirNotExists(t *testing.T) {
	t.Log("Testing WithMetadataFile function with given path's parent directory not exists")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMetadataFile("/notexistingparent/deneme.txt"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed not existing parent directory as value")
	}
}

func TestWithMetadataFile_NormalBehavior(t *testing.T) {
	t.Log("Testing WithMetadataFile function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMetadataFile("/trying.txt"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithStatusUpdatesFile_MultiplePassing(t *testing.T) {
	t.Log("Testing WithStatusUpdatesFile function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithStatusUpdatesFile("/trying.txt"))
	err2 := scanner.AddOptions(WithStatusUpdatesFile("/trying.txt"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithStatusUpdatesFile used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithStatusUpdatesFile used twice")
	}
}

func TestWithStatusUpdatesFile_PathIsADirectory(t *testing.T) {
	t.Log("Testing WithStatusUpdatesFile function with existing directory as value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithStatusUpdatesFile("/etc"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed existing directory as value")
	}
}

func TestWithStatusUpdatesFile_ParentDirNotExists(t *testing.T) {
	t.Log("Testing WithStatusUpdatesFile function with given path's parent directory not exists")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithStatusUpdatesFile("/notexistingparent/deneme.txt"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed not existing parent directory as value")
	}
}

func TestWithStatusUpdatesFile_NormalBehavior(t *testing.T) {
	t.Log("Testing WithStatusUpdatesFile function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithStatusUpdatesFile("/trying.txt"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithQuiet_MultiplePassing(t *testing.T) {
	t.Log("Testing WithQuiet function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithQuiet())
	err2 := scanner.AddOptions(WithQuiet())

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithQuiet used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithQuiet used twice")
	}
}

func TestWithQuiet_NormalBehavior(t *testing.T) {
	t.Log("Testing WithQuiet function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithQuiet())

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is returned while under normal behavior")
	}
}

func TestWithDisableSyslog_MultiplePassing(t *testing.T) {
	t.Log("Testing WithDisableSyslog function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithDisableSyslog())
	err2 := scanner.AddOptions(WithDisableSyslog())

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithDisableSyslog used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithDisableSyslog used twice")
	}
}

func TestWithDisableSyslog_NormalBehavior(t *testing.T) {
	t.Log("Testing WithDisableSyslog function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithDisableSyslog())

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithNotes_MultiplePassing(t *testing.T) {
	t.Log("Testing WithNotes function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithNotes("someNote"))
	err2 := scanner.AddOptions(WithNotes("someNote"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithNotes used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithNotes used twice")
	}
}

func TestWithNotes_NormalBehavior(t *testing.T) {
	t.Log("Testing WithNotes function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithNotes("someNotes"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithUserMetadata_MultiplePassing(t *testing.T) {
	t.Log("Testing WithUserMetadata function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithUserMetadata("metadata"))
	err2 := scanner.AddOptions(WithUserMetadata("metadata"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithUserMetadata used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithUserMetadata used twice")
	}
}

func TestWithUserMetadata_NormalBehavior(t *testing.T) {
	t.Log("Testing WithUserMetadata function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithUserMetadata("metadata"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithConfigFile_MultiplePassing(t *testing.T) {
	t.Log("Testing WithConfigFile function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithConfigFile("/etc/passwd"))
	err2 := scanner.AddOptions(WithConfigFile("/etc/passwd"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithConfigFile used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithConfigFile used twice")
	}
}

func TestWithConfigFile_NonExistingFile(t *testing.T) {
	t.Log("Testing WithConfigFile function with non-existing file path")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithConfigFile("/non-existing-path"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-existing file path")
	}
}

func TestWithConfigFile_WithDirectoryPath(t *testing.T) {
	t.Log("Testing WithConfigFile function with directory path")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithConfigFile("/etc"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed directory path")
	}
}

func TestWithConfigFile_NormalBehavior(t *testing.T) {
	t.Log("Testing WithConfigFile function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithConfigFile("/etc/passwd"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithMaxSendtoFailures_MultiplePassing(t *testing.T) {
	t.Log("Testing WithMaxSendtoFailures function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithMaxSendtoFailures("10000"))
	err2 := scanner.AddOptions(WithMaxSendtoFailures("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithMaxSendtoFailures used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithMaxSendtoFailures used twice")
	}
}

func TestWithMaxSendtoFailures_NonNumericValue(t *testing.T) {
	t.Log("Testing WithMaxSendtoFailures function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxSendtoFailures("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value")
	}
}

func TestWithMaxSendtoFailures_NormalBehavior(t *testing.T) {
	t.Log("Testing WithMaxSendtoFailures function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMaxSendtoFailures("1"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithMinHitrate_MultiplePassing(t *testing.T) {
	t.Log("Testing WithMinHitrate function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithMinHitrate("10000"))
	err2 := scanner.AddOptions(WithMinHitrate("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithMinHitrate used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithMinHitrate used twice")
	}
}

func TestWithMinHitrate_NonDecimalValue(t *testing.T) {
	t.Log("Testing WithMinHitrate function with non-decimal value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMinHitrate("non-decimal-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-decimal value")
	}
}

func TestWithMinHitrate_NormalBehavior_With_Decimal(t *testing.T) {
	t.Log("Testing WithMinHitrate function under normal behavior with decimal")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMinHitrate("1.12"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior with decimal")
	}
}

func TestWithMinHitrate_NormalBehavior_With_Integer(t *testing.T) {
	t.Log("Testing WithMinHitrate function under normal behavior with integer")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithMinHitrate("1"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior with integer")
	}
}

func TestWithSenderThreads_MultiplePassing(t *testing.T) {
	t.Log("Testing WithSenderThreads function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithSenderThreads("10000"))
	err2 := scanner.AddOptions(WithSenderThreads("10000"))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithSenderThreads used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithSenderThreads used twice")
	}
}

func TestWithSenderThreads_NonNumericValue(t *testing.T) {
	t.Log("Testing WithSenderThreads function with non-numeric value")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSenderThreads("non-numeric-value"))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-numeric value")
	}
}

func TestWithSenderThreads_NormalBehavior(t *testing.T) {
	t.Log("Testing WithSenderThreads function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithSenderThreads("1"))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithCores_MultiplePassing(t *testing.T) {
	t.Log("Testing WithCores function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	var cores []string
	for i := 0; i < runtime.NumCPU(); i++ {
		cores = append(cores, strconv.Itoa(i))
	}

	err1 := scanner.AddOptions(WithCores(cores))
	err2 := scanner.AddOptions(WithCores(cores))

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithCores used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithCores used twice")
	}
}

func TestWithCores_NonExistingCoreIndex(t *testing.T) {
	t.Log("Testing WithCores function with non-existing core index")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	cores := []string{
		"0",
		"non-existing-core-index",
	}
	err = scanner.AddOptions(WithCores(cores))

	t.Logf("Returned Error: %v", err)
	if err == nil {
		t.Error("Expected that error is returned when passed non-existing core index")
	}
}

func TestWithCores_NormalBehavior(t *testing.T) {
	t.Log("Testing WithCores function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	var cores []string
	for i := 0; i < runtime.NumCPU(); i++ {
		cores = append(cores, strconv.Itoa(i))
	}
	err = scanner.AddOptions(WithCores(cores))

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}

func TestWithIgnoreInvalidHosts_MultiplePassing(t *testing.T) {
	t.Log("Testing WithIgnoreInvalidHosts function with passing multiple time")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err1 := scanner.AddOptions(WithIgnoreInvalidHosts())
	err2 := scanner.AddOptions(WithIgnoreInvalidHosts())

	t.Logf("Returned First Error: %v", err1)
	if err1 != nil {
		t.Error("Expected that error is not returned when WithIgnoreInvalidHosts used once")
	}

	t.Logf("Returned Second Error: %v", err2)
	if err2 == nil {
		t.Error("Expected that error is returned when WithIgnoreInvalidHosts used twice")
	}
}

func TestWithIgnoreInvalidHosts_NormalBehavior(t *testing.T) {
	t.Log("Testing WithIgnoreInvalidHosts function under normal behavior")
	scanner, err := NewBlockingScanner()
	if err != nil {
		t.Error("Cannot create zmapgo scanner to test")
	}

	err = scanner.AddOptions(WithIgnoreInvalidHosts())

	t.Logf("Returned Error: %v", err)
	if err != nil {
		t.Error("Expected that error is not returned while under normal behavior")
	}
}
