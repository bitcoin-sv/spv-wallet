package engine

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	conditionAnd = "$and"
)

// processCustomFields will process all custom fields
func processCustomFields(conditions map[string]interface{}) {
	// Process the xpub_output_value
	_, ok := conditions["xpub_output_value"]
	if ok {
		processXpubOutputValueConditions(conditions)
	}

	// Process the xpub_output_value
	_, ok = conditions["xpub_metadata"]
	if ok {
		processXpubMetadataConditions(conditions)
	}
}

// processXpubOutputValueConditions will process xpub_output_value
func processXpubOutputValueConditions(conditions map[string]interface{}) {
	m, _ := json.Marshal(conditions["xpub_output_value"]) //nolint:errchkjson // this check might break the current code
	var r map[string]interface{}
	_ = json.Unmarshal(m, &r)

	xPubOutputValue := make([]map[string]interface{}, 0)
	for xPub, value := range r {
		outputKey := "xpub_output_value." + xPub
		xPubOutputValue = append(xPubOutputValue, map[string]interface{}{
			outputKey: value,
		})
	}
	if len(xPubOutputValue) > 0 {
		_, ok := conditions[conditionAnd]
		if ok {
			and := conditions[conditionAnd].([]map[string]interface{})
			and = append(and, xPubOutputValue...)
			conditions[conditionAnd] = and
		} else {
			conditions[conditionAnd] = xPubOutputValue
		}
	}

	delete(conditions, "xpub_output_value")
}

// processXpubMetadataConditions will process xpub_metadata
func processXpubMetadataConditions(conditions map[string]interface{}) {
	// marshal / unmarshal into standard map[string]interface{}
	m, _ := json.Marshal(conditions["xpub_metadata"]) //nolint:errchkjson // this check might break the current code
	var r map[string]interface{}
	_ = json.Unmarshal(m, &r)

	for xPub, xr := range r {
		xPubMetadata := make([]map[string]interface{}, 0)
		for key, value := range xr.(map[string]interface{}) {
			xPubMetadata = append(xPubMetadata, map[string]interface{}{
				"xpub_metadata.x": xPub,
				"xpub_metadata.k": key,
				"xpub_metadata.v": value,
			})
		}
		if len(xPubMetadata) > 0 {
			_, ok := conditions[conditionAnd]
			if ok {
				and := conditions[conditionAnd].([]map[string]interface{})
				and = append(and, xPubMetadata...)
				conditions[conditionAnd] = and
			} else {
				conditions[conditionAnd] = xPubMetadata
			}
		}
	}
	delete(conditions, "xpub_metadata")
}

// getMongoIndexes will get indexes from mongo
func getMongoIndexes() map[string][]mongo.IndexModel {
	return map[string][]mongo.IndexModel{
		"block_headers": {
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "height",
				Value: bsonx.Int32(1),
			}}},
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "synced",
				Value: bsonx.Int32(1),
			}}},
		},
		"destinations": {
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "address",
				Value: bsonx.Int32(1),
			}}},
		},
		"draft_transactions": {
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "status",
				Value: bsonx.Int32(1),
			}}},
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "xpub_id",
				Value: bsonx.Int32(1),
			}, {
				Key:   "status",
				Value: bsonx.Int32(1),
			}}},
		},
		"transactions": {
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "xpub_metadata.x",
				Value: bsonx.Int32(1),
			}, {
				Key:   "xpub_metadata.k",
				Value: bsonx.Int32(1),
			}, {
				Key:   "xpub_metadata.v",
				Value: bsonx.Int32(1),
			}}},
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "xpub_in_ids",
				Value: bsonx.Int32(1),
			}}},
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "xpub_out_ids",
				Value: bsonx.Int32(1),
			}}},
		},
		"utxos": {
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "transaction_id",
				Value: bsonx.Int32(1),
			}, {
				Key:   "output_index",
				Value: bsonx.Int32(1),
			}}},
			mongo.IndexModel{Keys: bsonx.Doc{{
				Key:   "xpub_id",
				Value: bsonx.Int32(1),
			}, {
				Key:   "type",
				Value: bsonx.Int32(1),
			}, {
				Key:   "draft_id",
				Value: bsonx.Int32(1),
			}, {
				Key:   "spending_tx_id",
				Value: bsonx.Int32(1),
			}}},
		},
	}
}
