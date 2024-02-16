package models

// Metadata is a SPV wallet ecosystem metadata model.
type Metadata map[string]interface{}

// XpubMetadata is a SPV wallet ecosystem xpub metadata model.
type XpubMetadata map[string]Metadata
