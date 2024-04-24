package engine

import (
	"context"
)

// Delete will delete a model from the Cachestore or Datastore using the provided conditions

func Delete(
	ctx context.Context,
	model ModelInterface,
	conditions map[string]interface{},
) error {
	return model.Client().Datastore().DeleteModel(ctx, model, conditions)
}
