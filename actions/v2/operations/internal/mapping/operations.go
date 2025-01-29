package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/operations/operationsmodels"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// OperationsPagedResponse maps a paged result of operations to a response.
func OperationsPagedResponse(operations *models.PagedResult[operationsmodels.Operation]) response.PageModel[response.Operation] {
	return response.PageModel[response.Operation]{
		Page: response.PageDescription{
			Size:          operations.PageDescription.Size,
			Number:        operations.PageDescription.Number,
			TotalElements: operations.PageDescription.TotalElements,
			TotalPages:    operations.PageDescription.TotalPages,
		},
		Content: utils.MapSlice(operations.Content, OperationsResponse),
	}
}

// OperationsResponse maps an operation to a response.
func OperationsResponse(operation *operationsmodels.Operation) *response.Operation {
	return &response.Operation{
		CreatedAt:    operation.CreatedAt,
		Value:        operation.Value,
		TxID:         operation.TxID,
		Type:         operation.Type,
		Counterparty: operation.Counterparty,
		TxStatus:     operation.TxStatus,
	}
}
