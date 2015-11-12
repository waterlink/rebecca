package driver

import (
	"sync"

	"github.com/waterlink/rebecca/context"
	"github.com/waterlink/rebecca/field"
)

var (
	driver    Driver
	driverMux = &sync.RWMutex{}
)

// Driver is for abstracting interaction with specific database
type Driver interface {
	Get(tx interface{}, tablename string, fields []field.Field, ID field.Field) ([]field.Field, error)
	Create(tx interface{}, tablename string, fields []field.Field, ID *field.Field) error
	Update(tx interface{}, tablename string, fields []field.Field, ID field.Field) error
	All(tablename string, fields []field.Field, ctx context.Context) ([][]field.Field, error)
	Where(tablename string, fields []field.Field, ctx context.Context, where string, args ...interface{}) ([][]field.Field, error)
	First(tablename string, fields []field.Field, ctx context.Context, where string, args ...interface{}) ([]field.Field, error)
	Remove(tx interface{}, tablename string, ID field.Field) error
	HasTransactions() bool
	Begin() (interface{}, error)
	Rollback(tx interface{})
	Commit(tx interface{}) error
}

// SetupDriver is for setting up driver manually
func SetupDriver(d Driver) {
	driverMux.Lock()
	defer driverMux.Unlock()

	driver = d
}

// Get is for fetching the driver instance, it returns reader lock
func Get() (Driver, sync.Locker) {
	lock := driverMux.RLocker()
	lock.Lock()
	return driver, lock
}
