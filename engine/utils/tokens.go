package utils

import (
	"regexp"
)

var stasDestinationRegexp = regexp.MustCompile(`^76a914[\da-f]{40}88ac`)

// GetLockingScriptFromSTASLockingScript the the destination lockingScript from a STAS token lockingScript
func GetLockingScriptFromSTASLockingScript(lockingScript string) (string, error) {
	matches := stasDestinationRegexp.FindAllString(lockingScript, -1)
	if len(matches) > 0 {
		p2pkhLockingScript := matches[0]
		if GetDestinationType(p2pkhLockingScript) == ScriptTypePubKeyHash {
			return p2pkhLockingScript, nil
		}
	}

	return "", ErrCouldNotDetermineDestinationOutput
}
