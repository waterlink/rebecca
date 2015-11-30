package rebecca

import (
	"reflect"
	"testing"

	"github.com/waterlink/rebecca/driver/fake"
)

func TestExec(t *testing.T) {
	d := fake.NewDriver()
	SetupDriver(d)

	if err := Exec("SOME QUERY $1, $2", 42, "hello"); err != nil {
		t.Fatal()
	}

	actual := d.ReceivedExec()
	expected := fake.ReceivedExec{
		Tx:    nil,
		Query: "SOME QUERY $1, $2",
		Args:  []interface{}{42, "hello"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected driver to receive exec %#v, but got %#v", expected, actual)
	}
}

func TestTransactionExec(t *testing.T) {
	d := fake.NewDriver()
	SetupDriver(d)

	tx, err := Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	if err := tx.Exec("SOME QUERY $1, $2", 42, "hello"); err != nil {
		t.Fatal()
	}

	actual := d.ReceivedExec()
	expected := fake.ReceivedExec{
		Tx:    tx.tx,
		Query: "SOME QUERY $1, $2",
		Args:  []interface{}{42, "hello"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected driver to receive exec %#v, but got %#v", expected, actual)
	}
}
