package datastore

import (
	"context"

	"gorm.io/gorm"
)

// NewTx will start a new datastore transaction
func (c *Client) NewTx(_ context.Context, fn func(*Transaction) error) error {
	// All GORM databases
	if c.options.db != nil {
		sessionDb := c.options.db.Session(getGormSessionConfig(c.options.db.PrepareStmt, c.IsDebug(), c.options.loggerDB))
		return fn(&Transaction{
			sqlTx: sessionDb.Begin(),
		})
	}

	// Empty transaction
	return fn(&Transaction{})
}

// NewRawTx will start a new datastore transaction
func (c *Client) NewRawTx() (*Transaction, error) {
	// All GORM databases
	if c.options.db != nil {
		sessionDb := c.options.db.Session(getGormSessionConfig(c.options.db.PrepareStmt, c.IsDebug(), c.options.loggerDB))
		return &Transaction{
			sqlTx: sessionDb.Begin(),
		}, nil
	}

	// Empty transaction
	return &Transaction{}, nil
}

// Transaction is the internal datastore transaction
type Transaction struct {
	committed    bool
	rowsAffected int64
	sqlTx        *gorm.DB
}

// CanCommit will return true if it can commit
func (tx *Transaction) CanCommit() bool {
	return !tx.committed && tx.sqlTx != nil
}

// Rollback the transaction
func (tx *Transaction) Rollback() error {
	if tx.sqlTx != nil {
		tx.sqlTx.Rollback()
	}

	return nil
}

// Commit will commit the transaction
func (tx *Transaction) Commit() error {
	// Have we already committed?
	if tx.committed {
		return nil
	} else if tx.sqlTx == nil {
		return nil
	}

	// Finally commit
	if tx.sqlTx != nil {
		result := tx.sqlTx.Commit()
		if result.Error != nil {
			_ = result.Rollback()
			return result.Error
		}
		tx.committed = true
		tx.rowsAffected = result.RowsAffected
	}

	return nil
}
