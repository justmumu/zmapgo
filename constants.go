package zmapgo

type BandwidthUnit string

var (
	UnitBandwidthBps  BandwidthUnit = "B"
	UnitBandwidthKbps BandwidthUnit = "K"
	UnitBandwidthMbps BandwidthUnit = "M"
	UnitBandwidthGbps BandwidthUnit = "G"
)

type VerbosityLevel string

var (
	VerbosityLevel1 VerbosityLevel = "1"
	VerbosityLevel2 VerbosityLevel = "2"
	VerbosityLevel3 VerbosityLevel = "3"
	VerbosityLevel4 VerbosityLevel = "4"
	VerbosityLevel5 VerbosityLevel = "5"
)
