package datastore

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

// checkForMethod is an interface to check for existing methods
type checkForMethod interface {
	GetModelName() string
	GetModelTableName() string
}

// IsModelSlice returns true if the given interface is a slice of models
func IsModelSlice(model interface{}) bool {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		value = reflect.New(value.Type().Elem())
	}
	modelType := reflect.Indirect(value).Type()

	if modelType.Kind() == reflect.Interface {
		modelType = reflect.Indirect(reflect.ValueOf(model)).Elem().Type()
	}

	return modelType.Kind() == reflect.Slice
}

// GetModelName get the name of the model via reflection
func GetModelName(model interface{}) *string {
	// Model is nil
	if model == nil {
		return nil
	}

	// Model is a pointer
	k := GetModelType(model).Kind()
	if reflect.ValueOf(model).Type().Kind() == reflect.Ptr && k != reflect.Struct {
		if m, ok := model.(checkForMethod); ok {
			name := m.GetModelName()
			return &name
		}
		return nil
	}

	// Model is a struct
	val := reflect.New(GetModelType(model)).MethodByName("GetModelName")
	if val.Kind() == reflect.Invalid { // Struct does not contain the method
		return nil
	}
	name := val.Call([]reflect.Value{})
	modelName := name[0].String()
	return &modelName
}

// GetModelTableName get the db table name of the model via reflection
func GetModelTableName(model interface{}) *string {
	// Model is nil
	if model == nil {
		return nil
	}

	// Model is a pointer
	k := GetModelType(model).Kind()
	if reflect.ValueOf(model).Type().Kind() == reflect.Ptr && k != reflect.Struct {
		if m, ok := model.(checkForMethod); ok {
			name := m.GetModelTableName()
			return &name
		}
		return nil
	}

	// Model is a struct
	val := reflect.New(GetModelType(model)).MethodByName("GetModelTableName")
	if val.Kind() == reflect.Invalid { // Struct does not contain the method
		return nil
	}
	name := val.Call([]reflect.Value{})
	modelName := name[0].String()
	return &modelName
}

// GetModelType get the model type of the model interface via reflection
func GetModelType(model interface{}) reflect.Type {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr && value.IsNil() {
		value = reflect.New(value.Type().Elem())
	}
	modelType := reflect.Indirect(value).Type()

	if modelType.Kind() == reflect.Interface {
		modelType = reflect.Indirect(reflect.ValueOf(model)).Elem().Type()
	}

	// Using "for" here to traverse to an actual element
	// this will find the element even if something is for instance a Ptr to a Slice
	for modelType.Kind() == reflect.Slice ||
		modelType.Kind() == reflect.Array ||
		modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	return modelType
}

// GetModelStringAttribute the attribute from the model as a string
func GetModelStringAttribute(model interface{}, attribute string) *string {
	valueOf := reflect.ValueOf(model)
	if model == nil || (valueOf.Kind() == reflect.Ptr &&
		valueOf.IsNil()) {
		return nil
	}
	modelReflect := reflect.Indirect(valueOf)
	if modelReflect.IsValid() &&
		modelReflect.Kind() == reflect.Struct {
		modelID := modelReflect.FieldByName(attribute)
		if modelID.IsValid() {
			attr := modelID.String()
			return &attr
		}
	}

	return nil
}

// GetModelBoolAttribute the attribute from the model as a bool
func GetModelBoolAttribute(model interface{}, attribute string) *bool {
	modelReflect := reflect.Indirect(reflect.ValueOf(model))
	if modelReflect.IsValid() {
		modelID := modelReflect.FieldByName(attribute)
		if modelID.IsValid() {
			value := modelID.Bool()
			return &value
		}
	}

	return nil
}

// GetModelUnset gets any empty values on the model and makes sure
// the update actually unsets those values in the database, otherwise
// this never happens, and we cannot unset
func GetModelUnset(model interface{}) map[string]bool {
	unset := make(map[string]bool)
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = reflect.ValueOf(model).Elem()
	}
	if t.Kind() == reflect.Struct {
		fields := reflect.VisibleFields(t)
		for _, field := range fields {
			// only get base fields / not fields of the BaseModel
			if len(field.Index) == 1 {
				if field.Type.Name() == nullStringFieldType ||
					field.Type.Name() == nullTimeFieldType {
					vv := v.Field(field.Index[0])
					if vv.Kind() == reflect.Ptr {
						vv = v.Elem()
					}
					value := vv.Interface()
					valid := reflect.ValueOf(value).FieldByName("Valid").Interface().(bool)
					if !valid {
						fieldName := strcase.ToSnake(field.Name)
						bsonTag := field.Tag.Get(bsonTagName)
						if bsonTag == "-" {
							// do not process private fields
							continue
						}
						if bsonTag != "" {
							args := strings.Split(bsonTag, ",")
							fieldName = args[0]
						}
						unset[fieldName] = true
					}
				}
			}
		}
	}

	return unset
}

// StringInSlice check whether the string already is in the slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
