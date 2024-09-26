package chainmodels

// TXStatus is the status of the transaction
type TXStatus string

// List of statuses available here: https://github.com/bitcoin-sv/arc
const (
	// Unknown status means that transaction has been sent to metamorph, but no processing has taken place. This should never be the case, unless something goes wrong.
	Unknown TXStatus = "UNKNOWN" // 0
	// Queued status means that transaction has been queued for processing.
	Queued TXStatus = "QUEUED" // 1
	// Received status means that transaction has been properly received by the metamorph processor.
	Received TXStatus = "RECEIVED" // 2
	// Stored status means that transaction has been stored in the metamorph store. This should ensure the transaction will be processed and retried if not picked up immediately by a mining node.
	Stored TXStatus = "STORED" // 3
	// AnnouncedToNetwork status means that transaction has been announced (INV message) to the Bitcoin network.
	AnnouncedToNetwork TXStatus = "ANNOUNCED_TO_NETWORK" // 4
	// RequestedByNetwork status means that transaction has been requested from metamorph by a Bitcoin node.
	RequestedByNetwork TXStatus = "REQUESTED_BY_NETWORK" // 5
	// SentToNetwork status means that transaction has been sent to at least 1 Bitcoin node.
	SentToNetwork TXStatus = "SENT_TO_NETWORK" // 6
	// AcceptedByNetwork status means that transaction has been accepted by a connected Bitcoin node on the ZMQ interface. If metamorph is not connected to ZQM, this status will never by set.
	AcceptedByNetwork TXStatus = "ACCEPTED_BY_NETWORK" // 7
	// SeenOnNetwork status means that transaction has been seen on the Bitcoin network and propagated to other nodes. This status is set when metamorph receives an INV message for the transaction from another node than it was sent to.
	SeenOnNetwork TXStatus = "SEEN_ON_NETWORK" // 8
	// Mined status means that transaction has been mined into a block by a mining node.
	Mined TXStatus = "MINED" // 9
	// SeenInOrphanMempool means that transaction has been sent to at least 1 Bitcoin node but parent transaction was not found.
	SeenInOrphanMempool TXStatus = "SEEN_IN_ORPHAN_MEMPOOL" // 10
	// Confirmed status means that transaction is marked as confirmed when it is in a block with 100 blocks built on top of that block.
	Confirmed TXStatus = "CONFIRMED" // 108
	// Rejected status means that transaction has been rejected by the Bitcoin network.
	Rejected TXStatus = "REJECTED" // 109
)
