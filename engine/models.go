package engine

import (
	"context"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	customTypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
)

var defaultPageSize = 25

// Model is the generic model field(s) and interface(s)
//
// gorm: https://gorm.io/docs/models.html
type Model struct {
	// ModelInterface `json:"-" toml:"-" yaml:"-" gorm:"-"` (@mrz: not needed, all models implement all methods)
	// ID string  `json:"id" toml:"id" yaml:"id" gorm:"primaryKey"`  (@mrz: custom per table)

	CreatedAt time.Time `json:"created_at" toml:"created_at" yaml:"created_at" gorm:"comment:The time that the record was originally created"`
	UpdatedAt time.Time `json:"updated_at" toml:"updated_at" yaml:"updated_at" gorm:"comment:The time that the record was last updated"`
	Metadata  Metadata  `gorm:"type:json;comment:The JSON metadata for the record" json:"metadata,omitempty"`

	// https://gorm.io/docs/indexes.html
	// DeletedAt gorm.DeletedAt `json:"deleted_at" toml:"deleted_at" yaml:"deleted_at" (@mrz: this was the original type)
	DeletedAt customTypes.NullTime `json:"deleted_at" toml:"deleted_at" yaml:"deleted_at" gorm:"index;comment:The time the record was marked as deleted"`

	// Private fields
	client        ClientInterface // Interface of the parent Client that loaded this SPV Wallet Engine model
	encryptionKey string          // Use for sensitive values that required encryption (IE: paymail public xpub)
	name          ModelName       // Name of model (table name)
	newRecord     bool            // Determine if the record is new (create vs update)
	pageSize      int             // Number of items per page to get if being used in for method getModels
	rawXpubKey    string          // Used on "CREATE" on some models
}

// ModelInterface is the interface that all models share
type ModelInterface interface {
	AfterCreated(ctx context.Context) (err error)
	AfterDeleted(ctx context.Context) (err error)
	AfterUpdated(ctx context.Context) (err error)
	BeforeCreating(ctx context.Context) (err error)
	BeforeUpdating(ctx context.Context) (err error)
	ChildModels() []ModelInterface
	Client() ClientInterface
	Display() interface{}
	GetID() string
	GetModelName() string
	GetModelTableName() string
	GetOptions(isNewRecord bool) (opts []ModelOps)
	IsNew() bool
	Migrate(client datastore.ClientInterface) error
	Name() string
	New()
	NotNew()
	RawXpub() string
	Save(ctx context.Context) (err error)
	SetOptions(opts ...ModelOps)
	SetRecordTime(bool)
	UpdateMetadata(metadata Metadata)
}

// ModelName is the model name type
type ModelName string

// NewBaseModel create an empty base model
func NewBaseModel(name ModelName, opts ...ModelOps) (m *Model) {
	m = &Model{name: name}
	m.SetOptions(opts...)
	return
}

// String is the string version of the name
func (n ModelName) String() string {
	return string(n)
}

// IsEmpty tests if the model name is empty
func (n ModelName) IsEmpty() bool {
	return n == ModelNameEmpty
}
