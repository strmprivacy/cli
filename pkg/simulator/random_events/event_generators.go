package random_events

import (
	"fmt"
	"math/rand"
	"streammachine.io/strm/pkg/schemas/clickstream"
	"streammachine.io/strm/pkg/schemas/demoschema"
	"streammachine.io/strm/pkg/simulator"
)

func createRandomDemo102Event(consentLevels []int32, sessionId string) sim.StreamMachineEvent {
	event := demoschema.NewDemoEvent()
	const eventContractRef = "streammachine/example/1.3.0"
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

func createRandomClickstreamEvent(consentLevels []int32, sessionId string) sim.StreamMachineEvent {
	event := clickstream.NewClickstreamEvent()
	event.StrmMeta = &clickstream.StrmMeta{ConsentLevels: consentLevels}
	event.ProducerSessionId = sessionId
	event.Customer = &clickstream.Customer{Id: "customer-" + event.ProducerSessionId}
	event.Url = "https://www.streammachine.io/rules"
	return event
}

var EventGenerators = map[string]func([]int32, string) sim.StreamMachineEvent{
	"clickstream":              createRandomClickstreamEvent,
	"streammachine/demo/1.0.2": createRandomDemo102Event,
}

func createUnionString(s string) *demoschema.UnionNullString {
	v := demoschema.NewUnionNullString()
	v.UnionType = demoschema.UnionNullStringTypeEnumString
	v.String = s
	return v
}
