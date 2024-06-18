package chainstate

func defaultArcConfig() *broadcastConfig {
	return &broadcastConfig{
		ArcAPIs: []string{
			"https://arc.taal.com",
		},
	}
}
