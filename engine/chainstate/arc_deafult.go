package chainstate

func defaultArcConfig() *broadcastConfig {
	return &broadcastConfig{
		protocol: "arc",
		ArcAPIs: []string{
			"https://arc.taal.com",
		},
	}
}
