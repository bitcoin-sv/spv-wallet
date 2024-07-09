package notifications

import "time"

type WebhookModel struct {
	URL         string
	TokenHeader string
	TokenValue  string
	BannedTo    *time.Time
}

func (model *WebhookModel) Banned() bool {
	if model.BannedTo == nil {
		return false
	}
	ret := !time.Now().After(*model.BannedTo)
	return ret
}
