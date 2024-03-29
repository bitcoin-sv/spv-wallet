package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// IndexExists check whether the given index exists in the datastore
func (c *Client) IndexExists(tableName, indexName string) (bool, error) {
	if c.Engine() == MySQL {
		return c.indexExistsMySQL(tableName, indexName)
	}

	return false, ErrUnknownSQL
}

// indexExistsMySQL is unique for MySQL
func (c *Client) indexExistsMySQL(tableName, indexName string) (bool, error) {
	indexQuery := `SELECT 1
			FROM INFORMATION_SCHEMA.STATISTICS
			WHERE TABLE_SCHEMA = '` + c.GetDatabaseName() + `'
			  AND TABLE_NAME = '` + tableName + `'
			  AND INDEX_NAME = '` + indexName + `'`

	tx := c.Raw(indexQuery)
	if tx.Error != nil {
		return false, tx.Error
	}

	var count int
	if tx = tx.Scan(&count); tx.Error != nil {
		return false, tx.Error
	}
	return count > 0, nil
}

// IndexMetadata check and creates the metadata json index
func (c *Client) IndexMetadata(tableName, field string) error {
	indexName := "idx_" + tableName + "_" + field
	if c.Engine() == MySQL { //nolint:revive // leave for comment
		/*
			//No way to index JSON in a generic way in MySQL?

			// workaround is to index the keys you are using yourself in the database
			ALTER TABLE tableName
				ADD new_column_name VARCHAR(255)
				  AS (metadata->>'$.columnName') STORED
			ADD INDEX (new_column_name)

			exists, err := c.indexExistsMySQL(tableName, indexName)
			if err != nil {
				return err
			}
			if !exists {
				query := ""
				tx := c.Execute(query)
				return tx.Error
			}
		*/
	} else if c.Engine() == PostgreSQL {
		tx := c.Execute(`CREATE INDEX IF NOT EXISTS ` + indexName + ` ON ` + tableName + ` USING gin (` + field + ` jsonb_path_ops)`)
		return tx.Error
	} else if c.Engine() == MongoDB {
		ctx := context.Background()

		// todo: this changed in the new version of mongo (needs to be tested)
		/*return createMongoIndex(ctx, c.options, tableName, true, mongo.IndexModel{Keys: bsonx.Doc{{
			Key:   metadataField + ".k",
			Value: bsonx.Int32(1),
		}, {
			Key:   metadataField + ".v",
			Value: bsonx.Int32(1),
		}}})*/

		return createMongoIndex(ctx, c.options, tableName, true, mongo.IndexModel{
			Keys: bson.D{
				{Key: metadataField + ".k", Value: 1},
				{Key: metadataField + ".v", Value: 1},
			},
		})
	}

	return nil
}
