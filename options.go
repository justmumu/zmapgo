package zmapgo

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// WithCustomArguments sets custom arguments to give to the zmap binary.
// There should be no reason to use this, unless you are using a custom build
// of zmap or that this repository isn't up to date with the latest options
// of the official zmap release.
func WithCustomArguments(args ...string) Option {
	return func(s *scanner) error {
		s.args = append(s.args, args...)
		return nil
	}
}

////////////////////////////////////////
////// Basic Arguments Section
////////////////////////////////////////

// WithTargets sets the target informations to give to the zmap binary.
// Targets can be ip address or cidr notation
func WithTargets(targets ...string) Option {
	return func(s *scanner) error {
		var realTargets []string
		for _, target := range targets {
			isValidIPAddress := false
			isValidCIDR := false

			// Check target is valid ipv4 address
			ipAddress := net.ParseIP(target)
			if ipAddress.To4() != nil {
				isValidIPAddress = true
			}

			// Check target is valid ipv4 cidr notation
			ip, n, err := net.ParseCIDR(target)
			if err == nil && ip.To4() != nil {
				isValidCIDR = true
			}

			if isValidIPAddress {
				realTargets = append(realTargets, ipAddress.To4().String())
			}

			if isValidCIDR {
				realTargets = append(realTargets, n.String())
			}

			if !isValidCIDR && !isValidIPAddress {
				return fmt.Errorf("given value of %s is not a valid ipv4 ipaddress or ipv4 cidr notation", target)
			}
		}

		s.args = append(s.args, realTargets...)
		return nil
	}
}

// WithTargetPort sets the target port to give to the zmap binary.
// This is required option. And should be used for ones.
// Zmap does not support multiple ports
func WithTargetPort(targetPort string) Option {
	return func(s *scanner) error {
		// check multiple usage
		if err := multiPassChecker(s.args, "--target-port"); err != nil {
			return err
		}

		// check valid port number
		portValue, err := strconv.Atoi(targetPort)
		if err != nil {
			return errors.New("given target port value is not a numeric value")
		}

		if !(portValue >= 0 && portValue <= 65535) {
			return errors.New("target port value must be between 0 and 65535")
		}

		s.args = append(s.args, "--target-port")
		s.args = append(s.args, targetPort)
		return nil
	}
}

// WithOutputFile sets the output file name to give to the zmap binary.
// If you are not passing this option, We will use "-" as value to read from stdout by default.
func WithOutputFile(outputFile string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--output-file"); err != nil {
			return err
		}

		if outputFile != "-" {
			// check path not exists but parent directory exists
			fileInfo, err := os.Stat(outputFile)
			if errors.Is(err, os.ErrNotExist) {
				parent := filepath.Dir(outputFile)
				if _, err = os.Stat(parent); errors.Is(err, os.ErrNotExist) {
					return errors.New("given output file path's parent directory is not exists")
				}
			}
			// check path is a directory
			if fileInfo != nil && fileInfo.IsDir() {
				return errors.New("given output file path is a directory")
			}
		}

		s.args = append(s.args, "--output-file")
		s.args = append(s.args, outputFile)
		return nil
	}
}

// WithBlacklistFile sets the blacklist file name to give to the zmap binary.
// If you are not passing this option, Zmap will use default blacklist file in "/usr/local/etc/zmap/blacklist.conf".
func WithBlacklistFile(blacklistFile string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--blacklist-file"); err != nil {
			return err
		}

		if _, err := os.Stat(blacklistFile); errors.Is(err, os.ErrNotExist) {
			return errors.New("blacklist file is not exists")
		}

		s.args = append(s.args, "--blacklist-file")
		s.args = append(s.args, blacklistFile)
		return nil
	}
}

// WithWhitelistFile sets the whitelist file name to give to the zmap binary.
func WithWhitelistFile(whitelistFile string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--whitelist-file"); err != nil {
			return err
		}

		if _, err := os.Stat(whitelistFile); errors.Is(err, os.ErrNotExist) {
			return errors.New("whitelist file is not exists")
		}

		s.args = append(s.args, "--whitelist-file", whitelistFile)
		s.args = append(s.args, whitelistFile)
		return nil
	}
}

////////////////////////////////////////
////// Scan Options Section
////////////////////////////////////////

// WithRate sets the rate to give to the zmap binary.
// Rate in packet/second (pps)
func WithRate(rate string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--rate"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(rate); err != nil {
			return errors.New("given rate value is not a numeric value")
		}

		s.args = append(s.args, "--rate")
		s.args = append(s.args, rate)
		return nil
	}
}

// WithBandwidth sets the bandwidth to give to zmap binary.
// It supports B, K, M and G suffixes.
// This option overrides --rate flag.
func WithBandwidth(bandwidth string, unit BandwidthUnit) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--bandwidth"); err != nil {
			return err
		}

		// check bandwidth is numeric
		if _, err := strconv.Atoi(bandwidth); err != nil {
			return errors.New("given bandwidth value is not a numeric value")
		}

		allBandwidthsSupported := []BandwidthUnit{
			UnitBandwidthBps,
			UnitBandwidthKbps,
			UnitBandwidthMbps,
			UnitBandwidthGbps,
		}

		isSupportedUnit := false
		for _, supUnit := range allBandwidthsSupported {
			if supUnit == unit {
				isSupportedUnit = true
				break
			}
		}

		if !isSupportedUnit {
			return errors.New("given unit is unsupported. Supported units: (B,K,M,G)")
		}

		var realValue string = bandwidth

		if unit != UnitBandwidthBps {
			realValue += string(unit)
		}

		s.args = append(s.args, "--bandwidth")
		s.args = append(s.args, realValue)
		return nil
	}
}

// WithMaxTargets sets the max targets to give to zmap binary.
func WithMaxTargets(maxTarget string, isPercentage bool) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--max-targets"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(maxTarget); err != nil {
			return errors.New("given max target value is not a numeric value")
		}

		if isPercentage {
			maxTarget += "%"
		}

		s.args = append(s.args, "--max-targets")
		s.args = append(s.args, maxTarget)
		return nil
	}
}

// WithMaxRuntime sets the max runtime to give to zmap binary.
func WithMaxRuntime(maxRuntime string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--max-runtime"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(maxRuntime); err != nil {
			return errors.New("given max runtime value is not a numeric value")
		}

		s.args = append(s.args, "--max-runtime")
		s.args = append(s.args, maxRuntime)
		return nil
	}
}

// WithMaxResults set the max results to give to zmap binary.
func WithMaxResults(maxResults string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--max-results"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(maxResults); err != nil {
			return errors.New("given max results value is not a numeric value")
		}

		s.args = append(s.args, "--max-results")
		s.args = append(s.args, maxResults)
		return nil
	}
}

// WithNumberOfProbesPerIP set the probes to give to zmap binary.
func WithNumberOfProbesPerIP(numberOfProbes string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--probes"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(numberOfProbes); err != nil {
			return errors.New("given number of probes value is not a numeric value")
		}

		s.args = append(s.args, "--probes")
		s.args = append(s.args, numberOfProbes)
		return nil
	}
}

// WithCooldownTime sets the cooldown to give to zmap binary.
// How long to continue receiving after sending last probe  (default=`8')
func WithCooldownTime(cooldown string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--cooldown-time"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(cooldown); err != nil {
			return errors.New("given cooldown value is not a numeric value")
		}

		s.args = append(s.args, "--cooldown-time")
		s.args = append(s.args, cooldown)
		return nil
	}
}

// WithSeed sets the seed to give to zmap binary.
// Seed used to select address permutation
func WithSeed(seed string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--seed"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(seed); err != nil {
			return errors.New("given seed value is not a numeric value")
		}

		s.args = append(s.args, "--seed")
		s.args = append(s.args, seed)
		return nil
	}
}

// WithMaxRetries sets the retries to give to zmap binary.
// Max number of times to try to send packet if send fails  (default=`10')
func WithMaxRetries(maxRetries string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--retries"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(maxRetries); err != nil {
			return errors.New("given max retries value is not a numeric value")
		}

		s.args = append(s.args, "--retries")
		s.args = append(s.args, maxRetries)
		return nil
	}
}

// WithDryrun sets the dryrun to give to zmap binary.
// Don't actually send packets
func WithDryrun() Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--dryrun"); err != nil {
			return err
		}

		s.args = append(s.args, "--dryrun")
		return nil
	}
}

// WithTotalShards sets the shards to give to zmap binary.
// Set the total number of shards  (default=`1')
func WithTotalShards(shards string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--shards"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(shards); err != nil {
			return errors.New("given shards value is not a numeric value")
		}

		s.args = append(s.args, "--shards")
		s.args = append(s.args, shards)
		return nil
	}
}

// WithShardID sets the shard to give to zmap binary.
// Set which shard this scan is (0 indexed) (default=`0')
func WithShardID(shardID string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--shard"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(shardID); err != nil {
			return errors.New("given shardID value is not a numeric value")
		}

		s.args = append(s.args, "--shard")
		s.args = append(s.args, shardID)
		return nil
	}
}

////////////////////////////////////////
////// Network Options Section
////////////////////////////////////////

// WithSourcePort sets the source port to give to zmap binary.
// Source port(s) for scan packets
// Can be one port (50000) or port range (Ex: 50000-50010)
func WithSourcePort(sourcePort string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--source-port"); err != nil {
			return err
		}

		var realValue string
		if strings.Contains(sourcePort, "-") {
			// if contain "-" it is a range
			splitted := strings.Split(sourcePort, "-")
			if len(splitted) != 2 {
				return errors.New("wrong port range definition")
			}

			// check port value of lower part in range definition
			portLower, err := strconv.Atoi(splitted[0])
			if err != nil {
				return errors.New("port number of lower part in range definition is not a numeric value")
			}

			if !(portLower >= 0 && portLower <= 65535) {
				return errors.New("port number of lower part in range definition must be between 0 and 65535")
			}

			// check port value of greater part in range definition
			portGreater, err := strconv.Atoi(splitted[1])
			if err != nil {
				return errors.New("port number of greater part in range definition is not a numeric value")
			}

			if !(portGreater >= 0 && portGreater <= 65535) {
				return errors.New("port number of greater part in range definition must be between 0 and 65535")
			}

			if portLower == portGreater {
				return errors.New("port number of lower and greater part in range definition cannot be equal")
			}

			if portLower > portGreater {
				return errors.New("port number of lower part in range definition cannot be greater than greater part of port number")
			}

			realValue = fmt.Sprintf("%d-%d", portLower, portGreater)

		} else {
			// if not contain "-" it is single value
			portValue, err := strconv.Atoi(sourcePort)
			if err != nil {
				return errors.New("given port number is not a numeric value")
			}

			if !(portValue >= 0 && portValue <= 65535) {
				return errors.New("given port number must be between 0 and 65535")
			}

			realValue = fmt.Sprintf("%d", portValue)
		}

		s.args = append(s.args, "--source-port")
		s.args = append(s.args, realValue)
		return nil
	}
}

// WithSourceIP sets the source ip to give to zmap binary.
// Source address(es) for scan packets
// Can be one ip (192.168.1.1) or ip range (Ex: 192.168.1.1-192.168.1.5)
func WithSourceIP(sourceIP string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--source-ip"); err != nil {
			return err
		}

		var realValue string

		if strings.Contains(sourceIP, "-") {
			// if ip address contains "-" character, it is range definition
			splitted := strings.Split(sourceIP, "-")
			if len(splitted) != 2 {
				return errors.New("wrong ip range definition")
			}
			lowerIP := net.ParseIP(splitted[0])
			if lowerIP.To4() == nil {
				return errors.New("lower part of ip in range definition is not valid ipv4 address")
			}

			greaterIP := net.ParseIP(splitted[1])
			if greaterIP.To4() == nil {
				return errors.New("greater part of ip in range definition is not valid ipv4 address")
			}

			if lowerIP.To4().Equal(greaterIP.To4()) {
				return errors.New("lower part and greater part cannot be equal in range definition")
			}

			if bytes.Compare(lowerIP.To4(), greaterIP.To4()) == 1 {
				return errors.New("lower part cannot be greater than greater part in range definition")
			}

			realValue = fmt.Sprintf("%s-%s", lowerIP.To4().String(), greaterIP.To4().String())
		} else {
			// if not contain "-" character, it is single ip definition
			ipAddress := net.ParseIP(sourceIP)
			if ipAddress.To4() == nil {
				return errors.New("given value is not valid ipv4 address")
			}
			realValue = ipAddress.To4().String()
		}

		s.args = append(s.args, "--source-ip")
		s.args = append(s.args, realValue)
		return nil
	}
}

// WithGatewayMAC sets the gateway mac to give to zmap binary.
// Specify gateway MAC address
func WithGatewayMAC(gatewayMAC string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--gateway-mac"); err != nil {
			return err
		}

		if _, err := net.ParseMAC(gatewayMAC); err != nil {
			return errors.New("given value is not valid mac address")
		}

		s.args = append(s.args, "--gateway-mac")
		s.args = append(s.args, gatewayMAC)
		return nil
	}
}

// WithSourceMAC sets the source mac to give to zmap binary.
// Source MAC address
func WithSourceMAC(sourceMAC string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--source-mac"); err != nil {
			return err
		}

		if _, err := net.ParseMAC(sourceMAC); err != nil {
			return errors.New("given value is not valid mac address")
		}

		s.args = append(s.args, "--source-mac")
		s.args = append(s.args, sourceMAC)
		return nil
	}
}

// WithInterface sets the interface to give to zmap binary.
// Specify network interface to use
func WithInterface(ifa string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--interface"); err != nil {
			return err
		}

		ifas, err := net.Interfaces()
		if err != nil {
			return errors.New("cannot get available interfaces")
		}

		isInterfaceExists := false
		for _, i := range ifas {
			if ifa == i.Name {
				isInterfaceExists = true
				break
			}
		}

		if !isInterfaceExists {
			return errors.New("given interface name is not available on the system")
		}

		s.args = append(s.args, "--interface")
		s.args = append(s.args, ifa)
		return nil
	}
}

// WithVPN sets the vpn to give to zmap binary.
// Sends IP packets instead of Ethernet (for VPNs)
func WithVPN() Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--vpn"); err != nil {
			return err
		}

		s.args = append(s.args, "--vpn")
		return nil
	}
}

////////////////////////////////////////
////// Probe Module Section
////////////////////////////////////////

// WithProbeModule sets the probe module to give to zmap binary.
// Select probe module  (default=`tcp_synscan')
func WithProbeModule(probeModule string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--probe-module"); err != nil {
			return err
		}

		availableModules, err := s.ListProbeModules()
		if err != nil {
			return err
		}

		isAvailableModule := false
		for _, availableModule := range availableModules {
			if probeModule == availableModule {
				isAvailableModule = true
			}
		}

		if !isAvailableModule {
			return errors.New("given probe module is not in available probe modules")
		}

		s.args = append(s.args, "--probe-module")
		s.args = append(s.args, probeModule)
		return nil
	}
}

// WithProbeArgs sets the probe args to give to zmap binary.
// Arguments to pass to probe module
func WithProbeArgs(probeArgs string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--probe-args"); err != nil {
			return err
		}

		s.args = append(s.args, "--probe-args")
		s.args = append(s.args, probeArgs)
		return nil
	}
}

////////////////////////////////////////
////// Data Output Section
////////////////////////////////////////

// WithOutputFields sets the output fields to give to zmap binary.
// Fields that should be output in result set
func WithOutputFields(fields []string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--output-fields"); err != nil {
			return err
		}

		availableFields, err := s.ListOutputFields()
		if err != nil {
			return err
		}

		for _, field := range fields {
			found := false
			for _, aField := range availableFields {
				if field == aField.Name {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("given field %s is not in available fields", field)
			}
		}

		realValue := strings.Join(fields, ",")

		s.args = append(s.args, "--output-fields")
		s.args = append(s.args, realValue)
		return nil
	}
}

// WithOutputModule sets the output module to give to zmap binary.
// Select output module  (default=`default')
func WithOutputModule(outputModule string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--output-module"); err != nil {
			return err
		}

		availableModules, err := s.ListOutputModules()
		if err != nil {
			return err
		}

		isAvailableModule := false
		for _, aModules := range availableModules {
			if aModules == outputModule {
				isAvailableModule = true
				break
			}
		}

		if !isAvailableModule {
			return errors.New("given output module is not in available output modules")
		}

		s.args = append(s.args, "--output-module")
		s.args = append(s.args, outputModule)
		return nil
	}
}

// WithOutputArgs sets the output args to give to zmap binary.
// Arguments to pass to output module
func WithOutputArgs(outputArgs string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--output-args"); err != nil {
			return err
		}

		s.args = append(s.args, "--output-args")
		s.args = append(s.args, outputArgs)
		return nil
	}
}

// WithOutputFilter sets the output filter to give to zmap binary.
// Specify a filter over the response fields to
// limit what responses get sent to the output
// module
func WithOutputFilter(outputFilter string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--output-filter"); err != nil {
			return err
		}

		s.args = append(s.args, "--output-filter")
		s.args = append(s.args, outputFilter)
		return nil
	}
}

////////////////////////////////////////
////// Logging and Metadata Section
////////////////////////////////////////

// WithVerbosity sets the verbosity to give to zmap binary.
// Level of log detail (0-5)  (default=`3')
func WithVerbosity(verbosityLevel VerbosityLevel) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--verbosity"); err != nil {
			return err
		}

		availableVerbosityLevels := []VerbosityLevel{
			VerbosityLevel1,
			VerbosityLevel2,
			VerbosityLevel3,
			VerbosityLevel4,
			VerbosityLevel5,
		}

		isAvailableVerbosity := false
		for _, aVerbosity := range availableVerbosityLevels {
			if verbosityLevel == aVerbosity {
				isAvailableVerbosity = true
				break
			}
		}

		if !isAvailableVerbosity {
			return errors.New("given verbosity level is not in available verbosity levels")
		}

		s.args = append(s.args, "--verbosity")
		s.args = append(s.args, string(verbosityLevel))
		return nil
	}
}

// WithLogFile sets the log file to give to zmap binary.
// Write log entries to file
func WithLogFile(logFile string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--log-file"); err != nil {
			return err
		}

		// check --log-directory already passed. If passed, not permit.
		isLogDirectoryPassed := false
		for _, arg := range s.args {
			if arg == "--log-directory" {
				isLogDirectoryPassed = true
			}
		}

		if isLogDirectoryPassed {
			return errors.New("log-file and log-directory cannot specified simultaneously")
		}
		// check passed path is a directory. If directory, not permit.
		fileInfo, err := os.Stat(logFile)
		if err == nil && fileInfo.IsDir() {
			return errors.New("given log file is a directory")
		}
		// check parent directory is exists. If not, not permit.
		if err != nil && errors.Is(err, os.ErrNotExist) {
			parent := filepath.Dir(logFile)
			if _, err := os.Stat(parent); errors.Is(err, os.ErrNotExist) {
				return errors.New("given log-file's parent directory is not exists")
			}
		}

		s.args = append(s.args, "--log-file")
		s.args = append(s.args, logFile)
		return nil
	}
}

// WithLogDirectory sets the log directory to give to zmap binary.
// Write log entries to a timestamped file in this directory
func WithLogDirectory(logDirectory string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--log-directory"); err != nil {
			return err
		}

		// check --log-file already passed. If passed, not permit
		isLogFilePassed := false
		for _, arg := range s.args {
			if arg == "--log-file" {
				isLogFilePassed = true
			}
		}

		if isLogFilePassed {
			return errors.New("log-file and log-directory cannot specified simultaneously")
		}
		// check passed directory already exists. If not exists, not permit.
		fileInfo, err := os.Stat(logDirectory)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			return errors.New("given log-directory path is not exists")
		}
		// check passed path. Is really directory. If not, not permit.
		if !fileInfo.IsDir() {
			return errors.New("given log-directory path is not a directory")
		}

		s.args = append(s.args, "--log-directory")
		s.args = append(s.args, logDirectory)
		return nil
	}
}

// WithMetadataFile sets the metadata file to give to zmap binary.
// Output file for scan metadata (JSON)
func WithMetadataFile(metadataFile string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--metadata-file"); err != nil {
			return err
		}

		// check is directory. If directory, not permit.
		fileInfo, err := os.Stat(metadataFile)
		if err == nil && fileInfo.IsDir() {
			return errors.New("given metadata file path is a directory")
		}
		// check parent directory is exists. If not, not permit.
		if err != nil && errors.Is(err, os.ErrNotExist) {
			parent := filepath.Dir(metadataFile)
			if _, err := os.Stat(parent); errors.Is(err, os.ErrNotExist) {
				return errors.New("given metadata file's parent directory is not exists")
			}
		}

		s.args = append(s.args, "--metadata-file")
		s.args = append(s.args, metadataFile)
		return nil
	}
}

// WithStatusUpdatesFile sets the status updates file to give to zmap binary.
// Write scan progress updates to CSV file
func WithStatusUpdatesFile(statusUpdateFile string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--status-updates-file"); err != nil {
			return err
		}

		// check is directory. If directory, not permit.
		fileInfo, err := os.Stat(statusUpdateFile)
		if err == nil && fileInfo.IsDir() {
			return errors.New("given status updates file path is a directory")
		}
		// check parent directory is exists. If not, not permit.
		if err != nil && errors.Is(err, os.ErrNotExist) {
			parent := filepath.Dir(statusUpdateFile)
			if _, err := os.Stat(parent); errors.Is(err, os.ErrNotExist) {
				return errors.New("given status updates file's parent directory is not exists")
			}
		}
		s.args = append(s.args, "--status-updates-file")
		s.args = append(s.args, statusUpdateFile)
		return nil
	}
}

// WithQuiet sets the quiet to give to zmap binary.
// Do not print status updates
func WithQuiet() Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--quiet"); err != nil {
			return err
		}

		s.args = append(s.args, "--quiet")
		return nil
	}
}

// WithDisableSyslog sets the disable syslog to give to zmap binary.
// Disables logging messages to syslog
func WithDisableSyslog() Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--disable-syslog"); err != nil {
			return err
		}

		s.args = append(s.args, "--disable-syslog")
		return nil
	}
}

// WithNotes sets the notes to give to zmap binary.
// Inject user-specified notes into scan metadata
func WithNotes(notes string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--notes"); err != nil {
			return err
		}

		s.args = append(s.args, "--notes")
		s.args = append(s.args, notes)
		return nil
	}
}

// WithUserMetadata sets the user metadata to give to zmap binary.
// Inject user-specified JSON metadata into scan metadata
func WithUserMetadata(userMetadata string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--user-metadata"); err != nil {
			return err
		}

		s.args = append(s.args, "--user-metadata")
		s.args = append(s.args, userMetadata)
		return nil
	}
}

////////////////////////////////////////
////// Additional Options Section
////////////////////////////////////////

// WithConfigFile sets the config file to give to zmap binary.
// Read a configuration file, which can specify any of these options
// (default=`/usr/local/etc/zmap/zmap.conf')
func WithConfigFile(configFile string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--config"); err != nil {
			return err
		}

		fileInfo, err := os.Stat(configFile)
		// Check path is exists
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("config file path is not exists")
		}

		// check path is directory
		if fileInfo.IsDir() {
			return errors.New("config file path is a directory")
		}

		s.args = append(s.args, "--config")
		s.args = append(s.args, configFile)
		return nil
	}
}

// WithMaxSendtoFailures sets the max sendto failures to give to zmap binary.
// Maximum NIC sendto failures before scan is aborted  (default=`-1')
func WithMaxSendtoFailures(maxSendtoFailures string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--max-sendto-failures"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(maxSendtoFailures); err != nil {
			return errors.New("max send failures is not a valid numeric value")
		}

		s.args = append(s.args, "--max-sendto-failures")
		s.args = append(s.args, maxSendtoFailures)
		return nil
	}
}

// WithMinHitrate sets the min hitrate to give to zmap binary.
// Minimum hitrate that scan can hit before scan is aborted  (default=`0.0')
func WithMinHitrate(minHitrate string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--min-hitrate"); err != nil {
			return err
		}

		if _, err := decimal.NewFromString(minHitrate); err != nil {
			return errors.New("min hitrate is not a valid decimal number")
		}
		s.args = append(s.args, "--min-hitrate")
		s.args = append(s.args, minHitrate)
		return nil
	}
}

// WithSenderThreads sets the sender threads to give to zmap binary.
// Threads used to send packets  (default=`1')
func WithSenderThreads(senderThreads string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--sender-threads"); err != nil {
			return err
		}

		if _, err := strconv.Atoi(senderThreads); err != nil {
			return errors.New("sender threads is not a valid numeric value")
		}

		s.args = append(s.args, "--sender-threads")
		s.args = append(s.args, senderThreads)
		return nil
	}
}

// WithCores sets the cores to give to zmap binary.
// Comma-separated list of cores to pin to
func WithCores(cores []string) Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--cores"); err != nil {
			return err
		}

		for _, core := range cores {
			isAvailableIndex := false
			for i := 0; i < runtime.NumCPU(); i++ {
				if strconv.Itoa(i) == core {
					isAvailableIndex = true
					break
				}
			}

			if !isAvailableIndex {
				return fmt.Errorf("core index \"%s\" is not available on system", core)
			}
		}

		s.args = append(s.args, "--cores")
		s.args = append(s.args, strings.Join(cores, ","))
		return nil
	}
}

// WithIgnoreInvalidHosts sets the ignore invalid hosts to give to zmap binary.
// Ignore invalid hosts in whitelist/blacklist file
func WithIgnoreInvalidHosts() Option {
	return func(s *scanner) error {
		if err := multiPassChecker(s.args, "--ignore-invalid-hosts"); err != nil {
			return err
		}

		s.args = append(s.args, "--ignore-invalid-hosts")
		return nil
	}
}
