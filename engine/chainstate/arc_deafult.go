package chainstate

func defaultArcConfig() *broadcastConfig {
	return &broadcastConfig{
		protocol: "arc",
		minerAPIs: []string{
			"http://arc.com",
		},
	}
}
