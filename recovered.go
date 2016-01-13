package rebecca

import "fmt"

// Recovered represents recovered error
type Recovered struct {
	Err error
}

// Error implements error interface
func (r *Recovered) Error() string {
	return fmt.Sprintf("%s (recovered)", r.Err.Error())
}
