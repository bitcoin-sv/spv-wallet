package chainstate

import "github.com/tonicpow/go-minercraft/v2"

func defaultMinecraftConfig() *minercraftConfig {
	miners, _ := minercraft.DefaultMiners()
	apis, _ := minercraft.DefaultMinersAPIs()

	broadcastMiners := []*minercraft.Miner{}
	queryMiners := []*minercraft.Miner{}
	for _, miner := range miners {
		broadcastMiners = append(broadcastMiners, miner)

		if supportsQuerying(miner) {
			queryMiners = append(queryMiners, miner)
		}
	}

	return &minercraftConfig{
		broadcastMiners: broadcastMiners,
		queryMiners:     queryMiners,
		minerAPIs:       apis,
	}
}

func supportsQuerying(mm *minercraft.Miner) bool {
	return mm.Name == minercraft.MinerTaal || mm.Name == minercraft.MinerMempool
}
