package field

import "reflect"

// Field is for storing field's metadata
type Field struct {
	Name       string
	DriverName string
	Primary    bool
	Ty         reflect.Type
	Value      interface{}
}
