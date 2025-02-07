package txmodels

// TxEncoder is an interface for encoding transactions into different formats (e.g. BEEF, raw HEX)
type TxEncoder interface {
	ToBEEF() (string, error)
	ToRawHEX() string
}
