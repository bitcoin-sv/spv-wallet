package models

// Metadata is a SPV Wallet ecosystem metadata model.
type Metadata map[string]interface{}

// XpubMetadata is a SPV Wallet ecosystem xpub metadata model.
type XpubMetadata map[string]Metadata
