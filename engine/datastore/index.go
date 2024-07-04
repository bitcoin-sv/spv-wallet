package datastore

// IndexExists check whether the given index exists in the datastore
func (c *Client) IndexExists(_, _ string) (bool, error) {
	return false, ErrUnknownSQL
}

// IndexMetadata check and creates the metadata json index
func (c *Client) IndexMetadata(tableName, field string) error {
	indexName := "idx_" + tableName + "_" + field
	if c.Engine() == PostgreSQL {
		tx := c.Execute(`CREATE INDEX IF NOT EXISTS ` + indexName + ` ON ` + tableName + ` USING gin (` + field + ` jsonb_path_ops)`)
		return tx.Error
	}

	return nil
}
