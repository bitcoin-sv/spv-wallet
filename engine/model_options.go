package engine

import "encoding/json"

// ModelOps allow functional options to be supplied
// that overwrite default model options
type ModelOps func(m *Model)

// New set this model to a new record
func New() ModelOps {
	return func(m *Model) {
		m.New()
	}
}

// WithClient will set the Client on the model
func WithClient(client ClientInterface) ModelOps {
	return func(m *Model) {
		if client != nil {
			m.client = client
		}
	}
}

// WithXPub will set the xPub key on the model
func WithXPub(rawXpubKey string) ModelOps {
	return func(m *Model) {
		if len(rawXpubKey) > 0 {
			m.rawXpubKey = rawXpubKey
		}
	}
}

// WithEncryptionKey will set the encryption key on the model (if needed)
func WithEncryptionKey(encryptionKey string) ModelOps {
	return func(m *Model) {
		if len(encryptionKey) > 0 {
			m.encryptionKey = encryptionKey
		}
	}
}

// WithMetadata will add the metadata record to the model
func WithMetadata(key string, value interface{}) ModelOps {
	return func(m *Model) {
		if m.Metadata == nil {
			m.Metadata = make(Metadata)
		}
		m.Metadata[key] = value
	}
}

// WithMetadatas will add multiple metadata records to the model
func WithMetadatas(metadata map[string]interface{}) ModelOps {
	return func(m *Model) {
		if len(metadata) > 0 {
			if m.Metadata == nil {
				m.Metadata = make(Metadata)
			}
			for key, value := range metadata {
				m.Metadata[key] = value
			}
		}
	}
}

// WithMetadataFromJSON will add the metadata record to the model
func WithMetadataFromJSON(jsonData []byte) ModelOps {
	return func(m *Model) {
		if len(jsonData) > 0 {
			if m.Metadata == nil {
				m.Metadata = make(Metadata)
			}
			_ = json.Unmarshal(jsonData, &m.Metadata)
		}
	}
}

// WithPageSize will set the pageSize to use on the model in queries
func WithPageSize(pageSize int) ModelOps {
	return func(m *Model) {
		m.pageSize = pageSize
	}
}
