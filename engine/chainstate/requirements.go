package chainstate

// RequiredIn is the requirements for querying transaction information
type RequiredIn string

const (
	// RequiredInMempool is the transaction in mempool? (minimum requirement for a valid response)
	RequiredInMempool RequiredIn = requiredInMempool

	// RequiredOnChain is the transaction in on-chain? (minimum requirement for a valid response)
	RequiredOnChain RequiredIn = requiredOnChain
)

// ValidRequirement will return valid if the requirement is known
func (c *Client) validRequirement(requirement RequiredIn) bool {
	return requirement == RequiredOnChain || requirement == RequiredInMempool
}

func checkRequirement(requirement RequiredIn, id string, txInfo *TransactionInfo, onChainCondition bool) bool {
	switch requirement {
	case RequiredInMempool:
		return txInfo.ID == id
	case RequiredOnChain:
		return onChainCondition
	default:
		return false
	}
}

func checkRequirementArc(requirement RequiredIn, id string, txInfo *TransactionInfo) bool {
	isConfirmedOnChain := len(txInfo.BlockHash) > 0 && txInfo.TxStatus != ""
	return checkRequirement(requirement, id, txInfo, isConfirmedOnChain)
}

func checkRequirementMapi(requirement RequiredIn, id string, txInfo *TransactionInfo) bool {
	isConfirmedOnChain := len(txInfo.BlockHash) > 0 && txInfo.Confirmations > 0
	return checkRequirement(requirement, id, txInfo, isConfirmedOnChain)
}
