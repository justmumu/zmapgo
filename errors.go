package zmapgo

import "errors"

var (
	// ErrNmapNotInstalled means that upon trying to manually locate zmap in the user's path,
	// it was not found. Either use the WithBinaryPath method to set it manually, or make sure that
	// the nmap binary is present in the user's $PATH.
	ErrZmapNotInstalled = errors.New("zmap binary was not found")

	// ErrScanTimeout means that the provided context was done before the scanner finished its scan.
	ErrScanTimeout = errors.New("zmap scan timed out")
)
