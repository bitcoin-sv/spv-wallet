package chainstate

import (
	"fmt"
	"strings"
	"sync"
)

// struct handles communication with the client - returns first successful broadcast
type broadcastStatus struct {
	mu          *sync.Mutex
	complete    bool
	success     bool
	syncChannel chan string
}

func newBroadcastStatus(synchChannel chan string) *broadcastStatus {
	return &broadcastStatus{complete: false, syncChannel: synchChannel, mu: &sync.Mutex{}}
}

func (g *broadcastStatus) tryCompleteWithSuccess(fastestProvider string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.complete {
		g.complete = true
		g.success = true

		g.syncChannel <- fastestProvider
		close(g.syncChannel)
	}

	// g.mu.Unlock() is done by defer
}

func (g *broadcastStatus) dispose() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.complete {
		g.complete = true
		close(g.syncChannel) // have to close the channel here to not block client
	}

	// g.mu.Unlock() is done by defer
}

// result of single broadcast to provider
type broadcastResult struct {
	isError  bool
	err      error
	provider string
}

func newErrorResult(err error, provider string) broadcastResult {
	return broadcastResult{isError: true, err: err, provider: provider}
}

func newSuccessResult(provider string) broadcastResult {
	return broadcastResult{isError: false, provider: provider}
}

// doesErrorContain will look at a string for a list of strings
func doesErrorContain(err string, messages []string) bool {
	lower := strings.ToLower(err)
	for _, str := range messages {
		if strings.Contains(lower, str) {
			return true
		}
	}
	return false
}

func debugLog(c ClientInterface, txID, msg string) {
	c.DebugLog(fmt.Sprintf("[txID: %s]: %s", txID, msg))
}
