// Code generated by github.com/actgardner/gogen-avro/v7. DO NOT EDIT.
/*
 * SOURCE:
 *     demo.avsc
 */
package demoschema

import (
	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/vm"
	"github.com/actgardner/gogen-avro/v7/vm/types"
	"io"
)

type DemoEvent struct {
	StrmMeta *StrmMeta `json:"strmMeta"`
	// any value. For illustration purposes: use a value that is consistent over time like a customer or device ID.
	UniqueIdentifier *UnionNullString `json:"uniqueIdentifier"`
	// any value. For illustration purposes: use a value that is consistent over a limited period like a session.
	ConsistentValue string `json:"consistentValue"`
	// any value. For illustration purposes: use a value that could identify a user over time based on behavior, like browsing behavior (e.g. urls).
	SomeSensitiveValue *UnionNullString `json:"someSensitiveValue"`
	// any value. For illustration purposes: use a value that is not sensitive at all, like the rank of an item in a set.
	NotSensitiveValue *UnionNullString `json:"notSensitiveValue"`
}

const DemoEventAvroCRC64Fingerprint = "\xf7\x84\xf6薠\x8fT"

func NewDemoEvent() *DemoEvent {
	return &DemoEvent{}
}

func DeserializeDemoEvent(r io.Reader) (*DemoEvent, error) {
	t := NewDemoEvent()
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

func DeserializeDemoEventFromSchema(r io.Reader, schema string) (*DemoEvent, error) {
	t := NewDemoEvent()

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

func writeDemoEvent(r *DemoEvent, w io.Writer) error {
	var err error
	err = writeStrmMeta(r.StrmMeta, w)
	if err != nil {
		return err
	}
	err = writeUnionNullString(r.UniqueIdentifier, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.ConsistentValue, w)
	if err != nil {
		return err
	}
	err = writeUnionNullString(r.SomeSensitiveValue, w)
	if err != nil {
		return err
	}
	err = writeUnionNullString(r.NotSensitiveValue, w)
	if err != nil {
		return err
	}
	return err
}

func (r *DemoEvent) Serialize(w io.Writer) error {
	return writeDemoEvent(r, w)
}

func (r *DemoEvent) Schema() string {
	return "{\"fields\":[{\"name\":\"strmMeta\",\"type\":{\"fields\":[{\"name\":\"eventContractRef\",\"type\":\"string\"},{\"default\":null,\"name\":\"nonce\",\"type\":[\"null\",\"int\"]},{\"default\":null,\"logicalType\":\"date\",\"name\":\"timestamp\",\"type\":[\"null\",\"long\"]},{\"default\":null,\"name\":\"keyLink\",\"type\":[\"null\",\"string\"]},{\"default\":null,\"name\":\"billingId\",\"type\":[\"null\",\"string\"]},{\"name\":\"consentLevels\",\"type\":{\"items\":\"int\",\"type\":\"array\"}}],\"name\":\"StrmMeta\",\"type\":\"record\"}},{\"default\":null,\"doc\":\"any value. For illustration purposes: use a value that is consistent over time like a customer or device ID.\",\"name\":\"uniqueIdentifier\",\"type\":[\"null\",\"string\"]},{\"doc\":\"any value. For illustration purposes: use a value that is consistent over a limited period like a session.\",\"name\":\"consistentValue\",\"type\":\"string\"},{\"default\":null,\"doc\":\"any value. For illustration purposes: use a value that could identify a user over time based on behavior, like browsing behavior (e.g. urls).\",\"name\":\"someSensitiveValue\",\"type\":[\"null\",\"string\"]},{\"default\":null,\"doc\":\"any value. For illustration purposes: use a value that is not sensitive at all, like the rank of an item in a set.\",\"name\":\"notSensitiveValue\",\"type\":[\"null\",\"string\"]}],\"name\":\"io.strmprivacy.schemas.demo.v1.DemoEvent\",\"type\":\"record\"}"
}

func (r *DemoEvent) SchemaName() string {
	return "io.strmprivacy.schemas.demo.v1.DemoEvent"
}

func (_ *DemoEvent) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ *DemoEvent) SetInt(v int32)       { panic("Unsupported operation") }
func (_ *DemoEvent) SetLong(v int64)      { panic("Unsupported operation") }
func (_ *DemoEvent) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ *DemoEvent) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ *DemoEvent) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ *DemoEvent) SetString(v string)   { panic("Unsupported operation") }
func (_ *DemoEvent) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *DemoEvent) Get(i int) types.Field {
	switch i {
	case 0:
		r.StrmMeta = NewStrmMeta()

		return r.StrmMeta
	case 1:
		r.UniqueIdentifier = NewUnionNullString()

		return r.UniqueIdentifier
	case 2:
		return &types.String{Target: &r.ConsistentValue}
	case 3:
		r.SomeSensitiveValue = NewUnionNullString()

		return r.SomeSensitiveValue
	case 4:
		r.NotSensitiveValue = NewUnionNullString()

		return r.NotSensitiveValue
	}
	panic("Unknown field index")
}

func (r *DemoEvent) SetDefault(i int) {
	switch i {
	case 1:
		r.UniqueIdentifier = nil
		return
	case 3:
		r.SomeSensitiveValue = nil
		return
	case 4:
		r.NotSensitiveValue = nil
		return
	}
	panic("Unknown field index")
}

func (r *DemoEvent) NullField(i int) {
	switch i {
	case 1:
		r.UniqueIdentifier = nil
		return
	case 3:
		r.SomeSensitiveValue = nil
		return
	case 4:
		r.NotSensitiveValue = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ *DemoEvent) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *DemoEvent) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *DemoEvent) Finalize()                        {}

func (_ *DemoEvent) AvroCRC64Fingerprint() []byte {
	return []byte(DemoEventAvroCRC64Fingerprint)
}
