package randomsim

import (
	"fmt"
	"math/rand"
	"streammachine.io/strm/pkg/clickstream"
	"streammachine.io/strm/pkg/demoschema"
	"streammachine.io/strm/pkg/sim"
	"streammachine.io/strm/pkg/util"
)

func createRandomDemo102Event(consentLevels []int32, sessionId string) sim.StreammachineEvent {
	event := demoschema.NewDemoEvent()
	const eventContractRef = "streammachine/example/1.3.0"
	event.StrmMeta = &demoschema.StrmMeta{
		ConsentLevels:    consentLevels,
		EventContractRef: eventContractRef,
	}
	event.ConsistentValue = sessionId
	event.UniqueIdentifier = util.CreateUnionString(fmt.Sprintf("unique-%d", rand.Intn(100)))
	event.SomeSensitiveValue = util.CreateUnionString(fmt.Sprintf("sensitive-%d", rand.Intn(100)))
	event.NotSensitiveValue = util.CreateUnionString(fmt.Sprintf("not-sensitive-%d", rand.Intn(100)))
	return event
}

func createRandomClickstreamEvent(consentLevels []int32, sessionId string) sim.StreammachineEvent {
	event := clickstream.NewClickstreamEvent()
	event.StrmMeta = &clickstream.StrmMeta{ConsentLevels: consentLevels}
	event.ProducerSessionId = sessionId
	event.Customer = &clickstream.Customer{Id: "customer-" + event.ProducerSessionId}
	event.Url = "https://www.streammachine.io/rules"
	return event
}

var EventGenerators = map[string]func([]int32, string) sim.StreammachineEvent{
	"clickstream":              createRandomClickstreamEvent,
	"streammachine/demo/1.0.2": createRandomDemo102Event,
}
