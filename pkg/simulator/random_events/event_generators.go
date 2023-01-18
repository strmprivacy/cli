package random_events

import (
	"fmt"
	"io"
	"math/rand"
	"strmprivacy/strm/pkg/schemas/demoschema"
)

// StrmPrivacyEvent is an interface that is implemented by all generated code with actgardner/gogen-avro
type StrmPrivacyEvent interface {
	Serialize(w io.Writer) error
}

func createRandomDemoEvent(consentLevels []int32, sessionId string) StrmPrivacyEvent {
	event := demoschema.NewDemoEvent()

	event.StrmMeta = &demoschema.StrmMeta{
		ConsentLevels:    consentLevels,
		EventContractRef: "strmprivacy/example/1.5.0",
	}
	event.ConsistentValue = sessionId
	event.UniqueIdentifier = createUnionString(fmt.Sprintf("unique-%d", rand.Intn(100)))
	event.SomeSensitiveValue = createUnionString(fmt.Sprintf("sensitive-%d", rand.Intn(100)))
	event.NotSensitiveValue = createUnionString(fmt.Sprintf("not-sensitive-%d", rand.Intn(100)))
	return event
}

var EventGenerators = map[string]func([]int32, string) StrmPrivacyEvent{
	"strmprivacy/example/1.5.0": createRandomDemoEvent,
}

func createUnionString(s string) *demoschema.UnionNullString {
	v := demoschema.NewUnionNullString()
	v.UnionType = demoschema.UnionNullStringTypeEnumString
	v.String = s
	return v
}
