package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

func main() {
	client, err := engine.NewClient(
		context.Background(), // Set context
		engine.WithDebugging(),
		engine.WithAutoMigrate(engine.BaseModels...),
		engine.WithModels(NewExample("example-field")), // Add additional custom models to SPV Wallet Engine
	)
	if err != nil {
		log.Fatalln("error: " + err.Error())
	}

	defer func() {
		_ = client.Close(context.Background())
	}()

	log.Println("client loaded!", client.UserAgent())
}

// Example is an example model
type Example struct {
	engine.Model `bson:",inline"` // Base SPV Wallet Engine model
	ID           string           `json:"id" toml:"id" yaml:"id" gorm:"<-:create;type:char(64);primaryKey;comment:This is the unique record id" bson:"_id"`                                       // Unique identifier
	ExampleField string           `json:"example_field" toml:"example_field" yaml:"example_field" gorm:"<-:create;type:varchar(64);comment:This is an example string field" bson:"example_field"` // Example string field
}

// ModelExample is an example model
const (
	ModelExample  = "example"
	tableExamples = "examples"
)

// NewExample create new example model
func NewExample(exampleString string, opts ...engine.ModelOps) *Example {
	id, _ := utils.RandomHex(32)

	// Standardize and sanitize!
	return &Example{
		Model:        *engine.NewBaseModel(ModelExample, opts...),
		ExampleField: exampleString,
		ID:           id,
	}
}

// GetModelName returns the model name
func (e *Example) GetModelName() string {
	return ModelExample
}

// GetModelTableName returns the model db table name
func (e *Example) GetModelTableName() string {
	return tableExamples
}

// Save the model
func (e *Example) Save(ctx context.Context) (err error) {
	return engine.Save(ctx, e)
}

// GetID will get the ID
func (e *Example) GetID() string {
	return e.ID
}

// BeforeCreating is called before the model is saved to the DB
func (e *Example) BeforeCreating(_ context.Context) (err error) {
	e.Client().Logger().Debug().Msgf("starting: %s BeforeCreating hook...", e.Name())

	// Do something here!

	e.Client().Logger().Debug().Msgf("end: %s BeforeCreating hook", e.Name())
	return
}

// Migrate model specific migration
func (e *Example) Migrate(client datastore.ClientInterface) error {
	return client.IndexMetadata(client.GetTableName(tableExamples), engine.ModelMetadata.String())
}
