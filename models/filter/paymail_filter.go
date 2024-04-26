package filter

/*
ID         string `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the unique paymail record id" bson:"_id"`                                  // Unique identifier
	XpubID     string `json:"xpub_id" toml:"xpub_id" yaml:"xpub_id" gorm:"<-:create;type:char(64);index;comment:This is the related xPub" bson:"xpub_id"`                                // Related xPub ID
	Alias      string `json:"alias" toml:"alias" yaml:"alias" gorm:"<-;type:varchar(64);comment:This is alias@" bson:"alias"`                                                            // Alias part of the paymail
	Domain     string `json:"domain" toml:"domain" yaml:"domain" gorm:"<-;type:varchar(255);comment:This is @domain.com" bson:"domain"`                                                  // Domain of the paymail
	PublicName string `json:"public_name" toml:"public_name" yaml:"public_name" gorm:"<-;type:varchar(255);comment:This is public name for public profile" bson:"public_name,omitempty"` // Full username
	Avatar     string `json:"avatar" toml:"avatar" yaml:"avatar" gorm:"<-;type:text;comment:This is avatar url" bson:"avatar"`                                                           // This is the url of the user (public profile)

	ExternalXpubKey    string `json:"external_xpub_key" toml:"external_xpub_key" yaml:"external_xpub_key" gorm:"<-:create;type:varchar(512);index;comment:This is full xPub for external use, encryption optional" bson:"external_xpub_key"` // PublicKey hex encoded
	ExternalXpubKeyNum uint32 `json:"external_xpub_num" toml:"external_xpub_num" yaml:"external_xpub_num" gorm:"<-;type:int;default:0;comment:Derivation number used to generate ExternalXpubKey:external_xpub_num"`
	PubKeyNum          uint32 `json:"pubkey_num" toml:"pubkey_num" yaml:"pubkey_num" gorm:"<-;type:int;default:0;comment:Derivation number use to create PKI public key:pubkey_num"`
	XpubDerivationSeq  uint32 `json:"xpub_derivation_seq" toml:"xpub_derivation_seq" yaml:"xpub_derivation_seq" gorm:"<-;type:int;default:0;comment:The index derivation number use to generate new external xpub child keys and rotate PubKey:xpub_derivation_seq"`
*/

// PaymailFilter is a struct for handling request parameters for paymail_addresses search requests
type PaymailFilter struct {
	ModelFilter `json:",inline"`

	ID                 *string `json:"id,omitempty"`
	Alias              *string `json:"alias,omitempty"`
	Domain             *string `json:"domain,omitempty"`
	PublicName         *string `json:"publicName,omitempty"`
	Avatar             *string `json:"avatar,omitempty"`
	ExternalXpubKey    *string `json:"externalXpubKey,omitempty"`
	ExternalXpubKeyNum *uint32 `json:"externalXpubNum,omitempty"`
	PubKeyNum          *uint32 `json:"pubkeyNum,omitempty"`
	XpubDerivationSeq  *uint32 `json:"xpubDerivationSeq,omitempty"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *PaymailFilter) ToDbConditions() map[string]interface{} {
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_paymail_addresses.go
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "alias", d.Alias)
	applyIfNotNil(conditions, "domain", d.Domain)
	applyIfNotNil(conditions, "public_name", d.PublicName)
	applyIfNotNil(conditions, "avatar", d.Avatar)
	applyIfNotNil(conditions, "external_xpub_key", d.ExternalXpubKey)
	applyIfNotNil(conditions, "external_xpub_num", d.ExternalXpubKeyNum)
	applyIfNotNil(conditions, "pubkey_num", d.PubKeyNum)
	applyIfNotNil(conditions, "xpub_derivation_seq", d.XpubDerivationSeq)

	return conditions
}
