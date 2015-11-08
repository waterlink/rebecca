package rebecca

import "github.com/waterlink/rebecca/field"

// ModelMetadata is for storing any metadata for the whole model
type ModelMetadata struct{}

type metadata struct {
	tablename string
	fields    []field.Field
	primary   field.Field
}
