package random_events

import (
	"fmt"
	"math/rand"
	"strmprivacy/strm/pkg/schemas/clickstream"
	"strmprivacy/strm/pkg/schemas/demoschema"
	"strmprivacy/strm/pkg/simulator"
)

func createRandomDemoEvent(consentLevels []int32, sessionId string) sim.StrmPrivacyEvent {
	event := demoschema.NewDemoEvent()
	const eventContractRef = "strmprivacy/example/1.3.0"
	event.StrmMeta = &demoschema.StrmMeta{
		ConsentLevels:    consentLevels,
		EventContractRef: eventContractRef,
	}
	event.ConsistentValue = sessionId
	event.UniqueIdentifier = createUnionString(fmt.Sprintf("unique-%d", rand.Intn(100)))
	event.SomeSensitiveValue = createUnionString(fmt.Sprintf("sensitive-%d", rand.Intn(100)))
	event.NotSensitiveValue = createUnionString(fmt.Sprintf("not-sensitive-%d", rand.Intn(100)))
	return event
}

var EventGenerators = map[string]func([]int32, string) sim.StrmPrivacyEvent{
	"strmprivacy/demo/1.0.2": createRandomDemoEvent,
}

func createUnionString(s string) *demoschema.UnionNullString {
	v := demoschema.NewUnionNullString()
	v.UnionType = demoschema.UnionNullStringTypeEnumString
	v.String = s
	return v
}
