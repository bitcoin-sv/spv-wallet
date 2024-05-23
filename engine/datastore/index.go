package datastore

// TODO: Remove this method from client interface and then from here, when the datastore will be refactored
// IndexExists check whether the given index exists in the datastore
func (c *Client) IndexExists(tableName, indexName string) (bool, error) {
	return false, ErrUnknownSQL
}

// IndexMetadata check and creates the metadata json index
func (c *Client) IndexMetadata(tableName, field string) error {
	indexName := "idx_" + tableName + "_" + field
	tx := c.Execute(`CREATE INDEX IF NOT EXISTS ` + indexName + ` ON ` + tableName + ` USING gin (` + field + ` jsonb_path_ops)`)
	return tx.Error
}
