package chainmodels

import "fmt"

// ArcError represents an error returned by the ARC API when status code is 4xx.
type ArcError struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail"`
	Instance  string `json:"instance"`
	TxID      string `json:"txid"`
	ExtraInfo string `json:"extraInfo"`
}

// Error returns the error string it's the implementation of the error interface.
func (a *ArcError) Error() string {
	return fmt.Sprintf("ARC error: %s <txID: %s> %s", a.Title, a.TxID, a.Detail)
}

// IsEmpty checks if the error is empty indicating that we could not parse the error response.
func (a *ArcError) IsEmpty() bool {
	return a == nil || a.Status == 0
}
