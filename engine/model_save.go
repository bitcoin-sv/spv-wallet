package engine

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/spverrors"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/pkg/errors"
)

// Save will save the model(s) into the Datastore
func Save(ctx context.Context, model ModelInterface) (err error) {
	// Check for a client
	c := model.Client()
	if c == nil {
		return spverrors.ErrMissingClient
	}

	// Check for a datastore
	ds := c.Datastore()
	if ds == nil {
		return spverrors.ErrDatastoreRequired
	}
	// Create new Datastore transaction
	// @siggi: we need this to be in a callback context for Mongo
	// NOTE: a DB error is not being returned from here
	return ds.NewTx(ctx, func(tx *datastore.Transaction) (err error) {
		parentBeforeHook := _beforeHook(model)
		if err = parentBeforeHook(ctx); err != nil {
			return _closeTxWithError(tx, err)
		}

		// Set the record's timestamps
		model.SetRecordTime(model.IsNew())

		// Start the list of models to Save
		modelsToSave := append(make([]ModelInterface, 0), model)

		// Add any child models (fire before hooks)
		if children := model.ChildModels(); len(children) > 0 {
			for _, child := range children {

				childBeforeHook := _beforeHook(child)
				if err = childBeforeHook(ctx); err != nil {
					return _closeTxWithError(tx, err)
				}
				// Set the record's timestamps
				child.SetRecordTime(child.IsNew())
			}

			// Add to list for saving
			modelsToSave = append(modelsToSave, children...)
		}

		// Logs for saving models
		model.Client().Logger().Debug().Msgf("saving %d models...", len(modelsToSave))

		// Save all models (or fail!)
		for index := range modelsToSave {
			modelsToSave[index].Client().Logger().Debug().
				Str("modelID", modelsToSave[index].GetID()).
				Msgf("starting to save model: %s", modelsToSave[index].Name())
			if err = modelsToSave[index].Client().Datastore().SaveModel(
				ctx, modelsToSave[index], tx, modelsToSave[index].IsNew(), false,
			); err != nil {
				return _closeTxWithError(tx, err)
			}
		}

		// Commit all the model(s) if needed
		if tx.CanCommit() {
			model.Client().Logger().Debug().Msg("committing db transaction...")
			if err = tx.Commit(); err != nil {
				return
			}
		}

		// Fire after hooks (only on commit success)
		var afterErr error
		for index := range modelsToSave {
			if modelsToSave[index].IsNew() {
				modelsToSave[index].NotNew() // NOTE: calling it before this method... after created assumes it's been saved already
				afterErr = modelsToSave[index].AfterCreated(ctx)
			} else {
				afterErr = modelsToSave[index].AfterUpdated(ctx)
			}
			if afterErr != nil {
				if err == nil { // First error - set the error
					err = afterErr
				} else { // Got more than one error, wrap it!
					err = errors.Wrap(err, afterErr.Error())
				}
			}
		}

		return
	})
}

// saveToCache will save the model to the cache using the given key(s)
//
// ttl of 0 will cache forever
func saveToCache(ctx context.Context, keys []string, model ModelInterface, ttl time.Duration) error { //nolint:nolintlint,unparam // this does not matter
	// NOTE: this check is in place in-case a model does not load its parent Client()
	if model.Client() != nil {
		for _, key := range keys {
			if err := model.Client().Cachestore().SetModel(ctx, key, model, ttl); err != nil {
				return err
			}
		}
	} else {
		model.Client().Logger().Debug().
			Str("modelID", model.GetID()).
			Msg("ignoring saveToCache: client or cachestore is missing")
	}
	return nil
}

// _closeTxWithError will close the transaction with the given error
// It's crucial to run this rollback to prevent hanging db connections.
func _closeTxWithError(tx *datastore.Transaction, baseError error) error {
	if tx == nil {
		if baseError != nil {
			return baseError
		}
		return errors.New("transaction is nil during rollback")
	}
	if err := tx.Rollback(); err != nil {
		if baseError != nil {
			return errors.Wrap(baseError, err.Error())
		}
		return err
	}
	if baseError != nil {
		return baseError
	}
	return errors.New("closing transaction with error")
}

func _beforeHook(model ModelInterface) func(context.Context) error {
	if model.IsNew() {
		return model.BeforeCreating
	}
	return model.BeforeUpdating
}
