package test

import (
	"testing"

	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/streams/v1"
	"google.golang.org/protobuf/proto"
)

var creds = &entities.Credentials{ClientId: "clientId", ClientSecret: "clientSecret"}
var limits = &entities.Limits{EventRate: 10000, EventCount: 1000000}
var streamWithTags, streamWithTagsWithoutSecret *entities.Stream

func init() {
	streamWithTags = &entities.Stream{
		Ref:          &entities.StreamRef{Name: "clitest-with-tags", ProjectId: testConfig().projectId},
		Enabled:      true,
		Limits:       limits,
		Tags:         []string{"foo", "bar", "baz"},
		Credentials:  []*entities.Credentials{creds},
		MaskedFields: &entities.MaskedFields{Seed: "****"},
	}
	streamWithTagsWithoutSecret = (proto.Clone(streamWithTags)).(*entities.Stream)
	streamWithTagsWithoutSecret.Credentials = []*entities.Credentials{{ClientId: "clientId"}}
}

func TestStreams(t *testing.T) {
	_ = newConfigDir()
	t.Run("listStreams", listStreamsTest)
	t.Run("createStreamTest1", createStreamTest1)
	t.Run("createStreamTest2", createStreamTest2)
	t.Run("createDerivedStream1", createDerivedStream1)
	t.Run("getStream1", getStream1)
	t.Run("deleteStream1", deleteStream1)
}

func listStreamsTest(t *testing.T) {
	ExecuteAndVerify(t, &streams.ListStreamsResponse{}, "list", "streams")
}

func createStreamTest1(t *testing.T) {
	expected := &streams.CreateStreamResponse{
		Stream: &entities.Stream{
			Ref:          &entities.StreamRef{Name: "clitest", ProjectId: testConfig().projectId},
			Enabled:      true,
			Limits:       limits,
			Credentials:  []*entities.Credentials{creds},
			MaskedFields: &entities.MaskedFields{}}}
	ExecuteAndVerify(t, expected, "create", "stream", "clitest")
}

func createStreamTest2(t *testing.T) {
	s := &entities.Stream{}
	proto.Merge(s, streamWithTags)
	s.MaskedFields = &entities.MaskedFields{}
	ExecuteAndVerify(t, &streams.CreateStreamResponse{Stream: s},
		"create", "stream", "clitest-with-tags", "--tags=foo,bar,baz")
}

func createDerivedStream1(t *testing.T) {
	expected := &streams.CreateStreamResponse{
		Stream: &entities.Stream{
			Ref:              &entities.StreamRef{Name: "clitest-with-tags-2", ProjectId: testConfig().projectId},
			ConsentLevels:    []int32{2},
			ConsentLevelType: entities.ConsentLevelType_CUMULATIVE,
			Enabled:          true,
			LinkedStream:     "clitest-with-tags",
			Credentials:      []*entities.Credentials{creds},
			MaskedFields:     &entities.MaskedFields{}}}
	ExecuteAndVerify(t, expected,
		"create", "stream", "--derived-from=clitest-with-tags", "--levels=2")
}

func getStream1(t *testing.T) {
	ExecuteAndVerify(t,
		&streams.GetStreamResponse{StreamTree: &entities.StreamTree{Stream: streamWithTagsWithoutSecret}},
		"get", "stream", "clitest-with-tags")
}

func deleteStream1(t *testing.T) {
	ExecuteAndVerify(t, &streams.DeleteStreamResponse{},
		"delete", "stream", "clitest-with-tags", "--recursive")
}
