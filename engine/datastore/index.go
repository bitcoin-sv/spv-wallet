package datastore

// IndexMetadata check and creates the metadata json index
func (c *Client) IndexMetadata(tableName, field string) error {
	indexName := "idx_" + tableName + "_" + field
	tx := c.Execute(`CREATE INDEX IF NOT EXISTS ` + indexName + ` ON ` + tableName + ` USING gin (` + field + ` jsonb_path_ops)`)
	return tx.Error
}
