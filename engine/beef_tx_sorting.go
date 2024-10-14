package engine

import trx "github.com/bitcoin-sv/go-sdk/transaction"

func kahnTopologicalSortTransactions(transactions []*trx.Transaction) []*trx.Transaction {
	txByID, incomingEdgesMap, zeroIncomingEdgeQueue := prepareSortStructures(transactions)
	result := make([]*trx.Transaction, 0, len(transactions))

	for len(zeroIncomingEdgeQueue) > 0 {
		txID := zeroIncomingEdgeQueue[0]
		zeroIncomingEdgeQueue = zeroIncomingEdgeQueue[1:]

		tx := txByID[txID]
		result = append(result, tx)

		zeroIncomingEdgeQueue = removeTxFromIncomingEdges(tx, incomingEdgesMap, zeroIncomingEdgeQueue)
	}

	reverseInPlace(result)
	return result
}

func prepareSortStructures(dag []*trx.Transaction) (txByID map[string]*trx.Transaction, incomingEdgesMap map[string]int, zeroIncomingEdgeQueue []string) {
	dagLen := len(dag)
	txByID = make(map[string]*trx.Transaction, dagLen)
	incomingEdgesMap = make(map[string]int, dagLen)

	for _, tx := range dag {
		txByID[tx.TxID().String()] = tx // TODO: perf -> In GO-SDK, the TxID is calculated every time we try to get it, which means we hash the tx bytes twice each time. It's expensive operation - try to avoid calculation each time
		incomingEdgesMap[tx.TxID().String()] = 0
	}

	calculateIncomingEdges(incomingEdgesMap, txByID)
	zeroIncomingEdgeQueue = getTxWithZeroIncomingEdges(incomingEdgesMap)

	return
}

func calculateIncomingEdges(inDegree map[string]int, txByID map[string]*trx.Transaction) {
	for _, tx := range txByID {
		for _, input := range tx.Inputs {
			inputUtxoTxID := input.SourceTXID.String() // TODO: perf -> In GO-SDK, the TxID is calculated every time we try to get it, which means we hash the tx bytes twice each time. It's expensive operation - try to avoid calculation each time
			if _, ok := txByID[inputUtxoTxID]; ok {    // transaction can contains inputs we are not interested in
				inDegree[inputUtxoTxID]++
			}
		}
	}
}

func getTxWithZeroIncomingEdges(incomingEdgesMap map[string]int) []string {
	zeroIncomingEdgeQueue := make([]string, 0, len(incomingEdgesMap))

	for txID, edgeNum := range incomingEdgesMap {
		if edgeNum == 0 {
			zeroIncomingEdgeQueue = append(zeroIncomingEdgeQueue, txID)
		}
	}

	return zeroIncomingEdgeQueue
}

func removeTxFromIncomingEdges(tx *trx.Transaction, incomingEdgesMap map[string]int, zeroIncomingEdgeQueue []string) []string {
	for _, input := range tx.Inputs {
		neighborID := input.SourceTXID.String() // TODO: perf -> In GO-SDK, the TxID is calculated every time we try to get it, which means we hash the tx bytes twice each time. It's expensive operation - try to avoid calculation each time
		incomingEdgesMap[neighborID]--

		if incomingEdgesMap[neighborID] == 0 {
			zeroIncomingEdgeQueue = append(zeroIncomingEdgeQueue, neighborID)
		}
	}

	return zeroIncomingEdgeQueue
}

func reverseInPlace(collection []*trx.Transaction) {
	for i, j := 0, len(collection)-1; i < j; i, j = i+1, j-1 {
		collection[i], collection[j] = collection[j], collection[i]
	}
}
