package chainstate

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/tonicpow/go-minercraft/v2"
	"github.com/tonicpow/go-minercraft/v2/apis/mapi"
)

func (c *Client) minercraftInit(ctx context.Context) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		defer txn.StartSegment("start_minercraft").End()
	}
	mi := &minercraftInitializer{client: c, ctx: ctx, minersWithFee: make(minerToFeeMap)}

	if err := mi.newClient(); err != nil {
		return err
	}

	if err := mi.validateMiners(); err != nil {
		return err
	}

	if c.isFeeQuotesEnabled() {
		c.options.config.feeUnit = mi.lowestFee()
	}

	return nil
}

type minercraftInitializer struct {
	client        *Client
	ctx           context.Context
	minersWithFee minerToFeeMap
	lock          sync.Mutex
}

type (
	minerID       string
	minerToFeeMap map[minerID]utils.FeeUnit
)

func (i *minercraftInitializer) defaultMinercraftOptions() (opts *minercraft.ClientOptions) {
	c := i.client
	opts = minercraft.DefaultClientOptions()
	if len(c.options.userAgent) > 0 {
		opts.UserAgent = c.options.userAgent
	}
	return
}

func (i *minercraftInitializer) newClient() (err error) {
	c := i.client

	if c.Minercraft() == nil {
		var optionalMiners []*minercraft.Miner
		var loadedMiners []string

		// Loop all broadcast miners and append to the list of miners
		for _, broadcastMiner := range c.options.config.minercraftConfig.broadcastMiners {
			if !utils.StringInSlice(broadcastMiner.MinerID, loadedMiners) {
				optionalMiners = append(optionalMiners, broadcastMiner)
				loadedMiners = append(loadedMiners, broadcastMiner.MinerID)
			}
		}

		// Loop all query miners and append to the list of miners
		for _, queryMiner := range c.options.config.minercraftConfig.queryMiners {
			if !utils.StringInSlice(queryMiner.MinerID, loadedMiners) {
				optionalMiners = append(optionalMiners, queryMiner)
				loadedMiners = append(loadedMiners, queryMiner.MinerID)
			}
		}
		c.options.config.minercraft, err = minercraft.NewClient(
			i.defaultMinercraftOptions(),
			c.HTTPClient(),
			c.options.config.minercraftConfig.apiType,
			optionalMiners,
			c.options.config.minercraftConfig.minerAPIs,
		)
	}
	return
}

// validateMiners will check if miner is reachable by requesting its FeeQuote
// If there was on error on FeeQuote(), the miner will be deleted from miners list
// If usage of MapiFeeQuotes is enabled and miner is reachable, miner's fee unit will be updated with MAPI fee quotes
// If FeeQuote returns some quote, but fee is not presented in it, it means that miner is valid but we can't use it's feequote
func (i *minercraftInitializer) validateMiners() error {
	ctxWithCancel, cancel := context.WithTimeout(i.ctx, 5*time.Second)
	defer cancel()

	c := i.client
	var wg sync.WaitGroup

	for _, miner := range c.options.config.minercraftConfig.broadcastMiners {
		wg.Add(1)
		currentMiner := miner
		go func() {
			defer wg.Done()
			feeUnit, err := i.getFeeQuote(ctxWithCancel, currentMiner)
			if err != nil {
				c.options.logger.Warn().Msgf("No FeeQuote response from miner %s. Reason: %s", currentMiner.Name, err)
				return
			}
			i.addToMinersWithFee(currentMiner, feeUnit)
		}()
	}
	wg.Wait()

	i.deleteUnreachableMiners()

	switch {
	case len(c.options.config.minercraftConfig.broadcastMiners) == 0:
		return ErrMissingBroadcastMiners
	case len(c.options.config.minercraftConfig.queryMiners) == 0:
		return ErrMissingQueryMiners
	default:
		return nil
	}
}

func (i *minercraftInitializer) getFeeQuote(ctx context.Context, miner *minercraft.Miner) (*utils.FeeUnit, error) {
	c := i.client

	apiType := c.Minercraft().APIType()

	if apiType == minercraft.Arc {
		return nil, fmt.Errorf("we no longer support ARC with Minercraft. (%s)", miner.Name)
	}

	quote, err := c.Minercraft().FeeQuote(ctx, miner)
	if err != nil {
		return nil, fmt.Errorf("no FeeQuote response from miner %s. Reason: %s", miner.Name, err)
	}

	btFee := quote.Quote.GetFee(mapi.FeeTypeData)
	if btFee == nil {
		return nil, fmt.Errorf("fee is missing in %s's FeeQuote response", miner.Name)
	}

	feeUnit := &utils.FeeUnit{
		Satoshis: btFee.MiningFee.Satoshis,
		Bytes:    btFee.MiningFee.Bytes,
	}
	return feeUnit, nil
}

func (i *minercraftInitializer) addToMinersWithFee(miner *minercraft.Miner, feeUnit *utils.FeeUnit) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.minersWithFee[minerID(miner.MinerID)] = *feeUnit
}

// deleteUnreachableMiners deletes miners which can't be reachable from config
func (i *minercraftInitializer) deleteUnreachableMiners() {
	c := i.client
	validMiners := []*minercraft.Miner{}
	for _, miner := range c.options.config.minercraftConfig.broadcastMiners {
		_, ok := i.minersWithFee[minerID(miner.MinerID)]
		if ok {
			validMiners = append(validMiners, miner)
		}
	}
	c.options.config.minercraftConfig.broadcastMiners = validMiners
}

// lowestFees takes the lowest fees among all miners and sets them as the feeUnit for future transactions
func (i *minercraftInitializer) lowestFee() *utils.FeeUnit {
	fees := make([]utils.FeeUnit, 0)
	for _, fee := range i.minersWithFee {
		fees = append(fees, fee)
	}
	lowest := utils.LowestFee(fees, i.client.options.config.feeUnit)
	return lowest
}
