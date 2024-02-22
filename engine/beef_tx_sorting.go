package engine

import "github.com/libsv/go-bt/v2"

func kahnTopologicalSortTransactions(transactions []*bt.Tx) []*bt.Tx {
	txByID, incomingEdgesMap, zeroIncomingEdgeQueue := prepareSortStructures(transactions)
	result := make([]*bt.Tx, 0, len(transactions))

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

func prepareSortStructures(dag []*bt.Tx) (txByID map[string]*bt.Tx, incomingEdgesMap map[string]int, zeroIncomingEdgeQueue []string) {
	dagLen := len(dag)
	txByID = make(map[string]*bt.Tx, dagLen)
	incomingEdgesMap = make(map[string]int, dagLen)

	for _, tx := range dag {
		txByID[tx.TxID()] = tx // TODO: perf -> In bt, the TxID is calculated every time we try to get it, which means we hash the tx bytes twice each time. It's expensive operation - try to avoid calculation each time
		incomingEdgesMap[tx.TxID()] = 0
	}

	calculateIncomingEdges(incomingEdgesMap, txByID)
	zeroIncomingEdgeQueue = getTxWithZeroIncomingEdges(incomingEdgesMap)

	return
}

func calculateIncomingEdges(inDegree map[string]int, txByID map[string]*bt.Tx) {
	for _, tx := range txByID {
		for _, input := range tx.Inputs {
			inputUtxoTxID := input.PreviousTxIDStr() // TODO: perf -> In bt, the TxID is calculated every time we try to get it, which means we hash the tx bytes twice each time. It's expensive operation - try to avoid calculation each time
			if _, ok := txByID[inputUtxoTxID]; ok {  // transaction can contains inputs we are not interested in
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

func removeTxFromIncomingEdges(tx *bt.Tx, incomingEdgesMap map[string]int, zeroIncomingEdgeQueue []string) []string {
	for _, input := range tx.Inputs {
		neighborID := input.PreviousTxIDStr() // TODO: perf -> In bt, the TxID is calculated every time we try to get it, which means we hash the tx bytes twice each time. It's expensive operation - try to avoid calculation each time
		incomingEdgesMap[neighborID]--

		if incomingEdgesMap[neighborID] == 0 {
			zeroIncomingEdgeQueue = append(zeroIncomingEdgeQueue, neighborID)
		}
	}

	return zeroIncomingEdgeQueue
}

func reverseInPlace(collection []*bt.Tx) {
	for i, j := 0, len(collection)-1; i < j; i, j = i+1, j-1 {
		collection[i], collection[j] = collection[j], collection[i]
	}
}
