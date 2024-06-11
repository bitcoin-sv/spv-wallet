package chainstate

func defaultArcConfig() *broadcastConfig {
	return &broadcastConfig{
		protocol: "arc",
		minerAPIs: []string{
			"https://example.com",
		},
	}
}
