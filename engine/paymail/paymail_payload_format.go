package paymail

import "fmt"

// PayloadFormat is the format of the paymail payload
type PayloadFormat uint32

// Types of Paymail payload formats
const (
	BasicPaymailPayloadFormat PayloadFormat = iota
	BeefPaymailPayloadFormat
)

func (format PayloadFormat) String() string {
	switch format {
	case BasicPaymailPayloadFormat:
		return "BasicPaymailPayloadFormat"

	case BeefPaymailPayloadFormat:
		return "BeefPaymailPayloadFormat"

	default:
		return fmt.Sprintf("%d", uint32(format))
	}
}
