package chainstate

import (
	"errors"
	"strings"
)

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

func groupBroadcastResults(results []*BroadcastResult) *BroadcastResult {
	switch len(results) {
	case 0:
		return nil
	case 1:
		return results[0]
	default:
		return &BroadcastResult{
			Provider: ProviderAll,
			Failure:  groupBroadcastFailures(results),
		}
	}
}

func groupBroadcastFailures(results []*BroadcastResult) *BroadcastFailure {
	// group  failures
	invalidTx := false
	var err error

	for _, r := range results {
		if r.Failure == nil {
			continue
		}
		if r.Failure.InvalidTx {
			invalidTx = true
		}

		err = errors.Join(err, r.Failure.Error)
	}

	if err != nil {
		return &BroadcastFailure{
			InvalidTx: invalidTx,
			Error:     err,
		}
	}

	return nil
}
