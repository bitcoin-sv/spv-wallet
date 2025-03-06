package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/operations/operationsmodels"
	"github.com/bitcoin-sv/spv-wallet/lox"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/samber/lo"
)

// OperationsPagedResponse maps a paged result of operations to a response.
func OperationsPagedResponse(operations *models.PagedResult[operationsmodels.Operation]) api.ModelsOperationsSearchResult {
	return api.ModelsOperationsSearchResult{
		Page: api.ModelsSearchPage{
			Size:          operations.PageDescription.Size,
			Number:        operations.PageDescription.Number,
			TotalElements: operations.PageDescription.TotalElements,
			TotalPages:    operations.PageDescription.TotalPages,
		},
		Content: lo.Map(operations.Content, lox.MappingFn(OperationsResponse)),
	}
}

// OperationsResponse maps an operation to a response.
func OperationsResponse(operation *operationsmodels.Operation) api.ModelsOperation {
	return api.ModelsOperation{
		CreatedAt:    operation.CreatedAt,
		Value:        operation.Value,
		TxID:         operation.TxID,
		Type:         api.ModelsOperationType(operation.Type),
		Counterparty: operation.Counterparty,
		TxStatus:     api.ModelsOperationTxStatus(operation.TxStatus),
		BlockHeight:  operation.BlockHeight,
		BlockHash:    operation.BlockHash,
	}
}
