// Code generated by github.com/actgardner/gogen-avro/v7. DO NOT EDIT.
/*
 * SOURCE:
 *     clickstream.avsc
 */
package clickstream

import (
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
	"io"
)

type Customer struct {
	Id string `json:"id"`
}

const CustomerAvroCRC64Fingerprint = "\xaf;]\x04*\x11\xf7-"

func NewCustomer() *Customer {
	return &Customer{}
}

func DeserializeCustomer(r io.Reader) (*Customer, error) {
	t := NewCustomer()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func DeserializeCustomerFromSchema(r io.Reader, schema string) (*Customer, error) {
	t := NewCustomer()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return nil, err
	}

	err = vm.Eval(r, deser, t)
	if err != nil {
		return nil, err
	}
	return t, err
}

func writeCustomer(r *Customer, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Id, w)
	if err != nil {
		return err
	}
	return err
}

func (r *Customer) Serialize(w io.Writer) error {
	return writeCustomer(r, w)
}

func (r *Customer) Schema() string {
	return "{\"fields\":[{\"name\":\"id\",\"type\":\"string\"}],\"name\":\"io.streammachine.schemas.strmcatalog.clickstream.Customer\",\"type\":\"record\"}"
}

func (r *Customer) SchemaName() string {
	return "io.streammachine.schemas.strmcatalog.clickstream.Customer"
}

func (_ *Customer) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *Customer) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *Customer) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *Customer) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *Customer) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *Customer) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *Customer) SetString(v string)   { panic("Unsupported operation") }
func (_ *Customer) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Customer) Get(i int) types.Field {
	switch i {
	case 0:
		return &types.String{Target: &r.Id}
	}
	panic("Unknown field index")
}

func (r *Customer) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *Customer) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ *Customer) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *Customer) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *Customer) Finalize()                        {}

func (_ *Customer) AvroCRC64Fingerprint() []byte {
	return []byte(CustomerAvroCRC64Fingerprint)
}
