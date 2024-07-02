package chainstate

import "errors"

// ErrMissingBroadcastMiners is when broadcasting miners are missing
var ErrMissingBroadcastMiners = errors.New("missing: broadcasting miners")

// ErrMissingQueryMiners is when query miners are missing
var ErrMissingQueryMiners = errors.New("missing: query miners")
