package chainstate

// Network is the supported Bitcoin networks
type Network string

// Supported networks
const (
	MainNet       Network = mainNet // Main public network
	StressTestNet Network = stn     // Stress Test Network (https://bitcoinscaling.io/)
	TestNet       Network = testNet // Test public network
)

// String is the string version of network
func (n Network) String() string {
	return string(n)
}

// Alternate is the alternate string version
func (n Network) Alternate() string {
	switch n {
	case MainNet:
		return mainNetAlt
	case TestNet:
		return testNetAlt
	case StressTestNet:
		return stn
	default:
		return ""
	}
}
