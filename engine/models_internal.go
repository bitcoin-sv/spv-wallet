package engine

import (
	"context"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// AfterDeleted will fire after a successful delete in the Datastore
func (m *Model) AfterDeleted(_ context.Context) error {
	m.Client().Logger().Debug().Msgf("starting: %s AfterDeleted hook...", m.Name())
	m.Client().Logger().Debug().Msgf("end: %s AfterDeleted hook", m.Name())
	return nil
}

// BeforeUpdating will fire before updating a model in the Datastore
func (m *Model) BeforeUpdating(_ context.Context) error {
	m.Client().Logger().Debug().Msgf("starting: %s BeforeUpdate hook...", m.Name())
	m.Client().Logger().Debug().Msgf("end: %s BeforeUpdate hook", m.Name())
	return nil
}

// Client will return the current client
func (m *Model) Client() ClientInterface {
	return m.client
}

// ChildModels will return any child models
func (m *Model) ChildModels() []ModelInterface {
	return nil
}

// enrich is run after getting a record from the database
func (m *Model) enrich(name ModelName, opts ...ModelOps) {
	// Set the name
	m.name = name

	// Overwrite defaults
	m.SetOptions(opts...)
}

// GetOptions will get the options that are set on that model
func (m *Model) GetOptions(isNewRecord bool) (opts []ModelOps) {
	// Client was set on the model
	if m.client != nil {
		opts = append(opts, WithClient(m.client))
	}

	// New record flag
	if isNewRecord {
		opts = append(opts, New())
	}

	return
}

// IsNew returns true if the model is (or was) a new record
func (m *Model) IsNew() bool {
	return m.newRecord
}

// GetID will get the model id, if overwritten in the actual model
func (m *Model) GetID() string {
	return ""
}

// Name will get the collection name (model)
func (m *Model) Name() string {
	return m.name.String()
}

// New will set the record to new
func (m *Model) New() {
	m.newRecord = true
}

// NotNew sets newRecord to false
func (m *Model) NotNew() {
	m.newRecord = false
}

// RawXpub returns the rawXpubKey
func (m *Model) RawXpub() string {
	return m.rawXpubKey
}

// SetRecordTime will set the record timestamps (created is true for a new record)
func (m *Model) SetRecordTime(created bool) {
	if created {
		m.CreatedAt = time.Now().UTC()
	} else {
		m.UpdatedAt = time.Now().UTC()
	}
}

// UpdateMetadata will update the metadata on the model
// any key set to nil will be removed, other keys updated or added
func (m *Model) UpdateMetadata(metadata Metadata) {
	if m.Metadata == nil {
		m.Metadata = make(Metadata)
	}

	for key, value := range metadata {
		if value == nil {
			delete(m.Metadata, key)
		} else {
			m.Metadata[key] = value
		}
	}
}

// SetOptions will set the options on the model
func (m *Model) SetOptions(opts ...ModelOps) {
	for _, opt := range opts {
		opt(m)
	}
}

// Display filter the model for display
func (m *Model) Display() interface{} {
	return m
}

// AfterUpdated will fire after a successful update into the Datastore
func (m *Model) AfterUpdated(_ context.Context) error {
	m.Client().Logger().Debug().Msgf("starting: %s AfterUpdated hook...", m.Name())
	m.Client().Logger().Debug().Msgf("end: %s AfterUpdated hook", m.Name())
	return nil
}

// AfterCreated will fire after the model is created in the Datastore
func (m *Model) AfterCreated(_ context.Context) error {
	m.Client().Logger().Debug().Msgf("starting: %s AfterCreated hook...", m.Name())
	m.Client().Logger().Debug().Msgf("end: %s AfterCreated hook", m.Name())
	return nil
}

// incrementField will increment the given field atomically in the datastore
func incrementField(ctx context.Context, model ModelInterface, fieldName string,
	increment int64,
) (int64, error) {
	// Check for client
	c := model.Client()
	if c == nil {
		return 0, spverrors.ErrMissingClient
	}

	// Increment
	newValue, err := c.Datastore().IncrementModel(ctx, model, fieldName, increment)
	if err != nil {
		return 0, spverrors.Wrapf(err, "failed incrementing field %s", fieldName)
	}

	// AfterUpdate event should be called by parent function

	return newValue, nil
}

/*
// setFieldValueByJSONTag will parse the struct looking for the field (json tag) and updating the value if found
//
// todo: this was created because the increment field was not updating the model's value
func setFieldValueByJSONTag(item interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return spverrors.Newf("cannot assign to the item passed, item must be a pointer in order to assign")
	}
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJSONName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("json"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", spverrors.Newf("tag provided does not define a json tag: %s", fieldName)
	}
	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		structTypeField := v.Type().Field(i)
		jName, _ := findJSONName(structTypeField.Tag)
		if jName != "" && jName != "-" {
			fieldNames[jName] = i
		}
	}

	fieldNum, ok := fieldNames[fieldName]
	if !ok {
		return spverrors.Newf("field %s does not exist within the provided item", fieldName)
	}
	fieldVal := v.Field(fieldNum)
	switch fieldVal.Interface().(type) {
	case uint8:
		fieldVal.Set(reflect.ValueOf(uint8(value.(int64))))
	case uint16:
		fieldVal.Set(reflect.ValueOf(uint16(value.(int64))))
	case uint32:
		fieldVal.Set(reflect.ValueOf(uint32(value.(int64))))
	case uint64:
		fieldVal.Set(reflect.ValueOf(uint64(value.(int64))))
	}
	return nil
}
*/
