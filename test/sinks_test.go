package test

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/strmprivacy/api-definitions-go/v2/api/batch_exporters/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/sinks/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/streams/v1"
	"google.golang.org/protobuf/proto"
)

type AccessKey struct {
	UserName        string `json:"UserName"`
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
}

var awsCredentialsFileName string
var streamRef *entities.StreamRef
var sinkRef, anotherSinkRef *entities.SinkRef
var sink, anotherSink *entities.Sink
var bucketConfig *entities.BucketConfig
var batchExporterRef *entities.BatchExporterRef
var batchExporter, anotherBatchExporter *entities.BatchExporter

const (
	bucketCredentials = `{"AccessKey":{"UserName":"UserName","AccessKeyId":"AccessKeyId","SecretAccessKey":"SecretAccessKey"}}`
	billingId         = "testBillingId"
)

func TestSinks(t *testing.T) {
	awsCredentialsFileName = createAwsCredentialsFile(t)
	t.Run("createStreamForSinkTest", createStreamForSinkTest)
	t.Run("listSinks", listSinks)
	t.Run("createSink", createSink)
	t.Run("listSinks2", listSinks2)
	t.Run("createBatchExporter", createBatchExporter)
	t.Run("createSink2", createSink2)
	t.Run("createAnotherBatchExporter", createAnotherBatchExporter)
	t.Run("listBatchExporters", listBatchExporters)
	t.Run("getBatchExporter", getBatchExporter)
	t.Run("listSinksRecursive", listSinksRecursive)
	t.Run("getAnotherSinkRecursive", getAnotherSinkRecursive)
	t.Run("deleteAnotherSinkRecursive", deleteAnotherSinkRecursive)
	t.Run("deleteSinkNonRecursive", deleteSinkNonRecursive)
}

/*
setup some constant values.
*/
func init() {
	streamRef = &entities.StreamRef{BillingId: billingId, Name: "teststream"}
	sinkRef = &entities.SinkRef{BillingId: billingId, Name: "s3sink"}
	anotherSinkRef = &entities.SinkRef{BillingId: billingId, Name: "another-sink"}
	bucketConfig = &entities.BucketConfig{BucketName: "strm-cli-tester"}

	sink = &entities.Sink{
		Ref: sinkRef, SinkType: entities.SinkType_S3,
		Config: &entities.Sink_Bucket{Bucket: bucketConfig},
	}
	anotherSink = &entities.Sink{
		Ref: anotherSinkRef, SinkType: entities.SinkType_S3,
		Config: &entities.Sink_Bucket{Bucket: bucketConfig},
	}
	batchExporterRef = &entities.BatchExporterRef{BillingId: billingId, Name: sinkRef.Name + "-" + streamRef.Name}
	batchExporter = &entities.BatchExporter{
		Ref:                  batchExporterRef,
		StreamOrKeyStreamRef: &entities.BatchExporter_StreamRef{StreamRef: streamRef},
		SinkName:             sinkRef.Name,
		Interval:             &duration.Duration{Seconds: 60},
	}

	anotherBatchExporter = &entities.BatchExporter{}
	proto.Merge(anotherBatchExporter, batchExporter)
	anotherBatchExporter.SinkName = anotherSink.Ref.Name
	anotherBatchExporter.Ref.Name = "another-batch-exporter"
	anotherBatchExporter.PathPrefix = "some-prefix"
	anotherBatchExporter.Interval = &duration.Duration{Seconds: 300}

	_ = newConfigDir()
}

type AwsCredentials struct {
	AccessKey AccessKey `json:"AccessKey"`
}

func createStreamForSinkTest(t *testing.T) {
	ExecuteAndVerify(t, &streams.CreateStreamResponse{
		Stream: &entities.Stream{
			Ref:          streamRef,
			Enabled:      true,
			Limits:       limits,
			Credentials:  []*entities.Credentials{creds},
			MaskedFields: &entities.MaskedFields{}}},
		"create", "stream", streamRef.Name)
}

func listSinks(t *testing.T) {
	ExecuteAndVerify(t, &sinks.ListSinksResponse{}, "list", "sinks")
}

func createSink(t *testing.T) {
	s := &entities.Sink{}
	proto.Merge(s, sink)
	ExecuteAndVerify(t, &sinks.CreateSinkResponse{Sink: s}, "create", "sink", sink.Ref.Name, "strm-cli-tester", "--sink-type=S3", "--credentials-file="+awsCredentialsFileName)
}

func listSinks2(t *testing.T) {
	sinkTree := &entities.SinkTree{Sink: sink}
	ExecuteAndVerify(t, &sinks.ListSinksResponse{Sinks: []*entities.SinkTree{sinkTree}},
		"list", "sinks")
}

func createBatchExporter(t *testing.T) {
	b := &entities.BatchExporter{}
	proto.Merge(b, batchExporter)
	b.Interval = &duration.Duration{Seconds: 60}
	b.SinkName = sink.Ref.Name

	ExecuteAndVerify(t, &batch_exporters.CreateBatchExporterResponse{
		BatchExporter: b}, "create", "batch-exporter", streamRef.Name)
}

func createSink2(t *testing.T) {
	s := &entities.Sink{}
	proto.Merge(s, anotherSink)
	ExecuteAndVerify(t, &sinks.CreateSinkResponse{Sink: s},
		"create", "sink", s.Ref.Name, "strm-cli-tester", "--sink-type=S3", "--credentials-file="+awsCredentialsFileName)
}

func createAnotherBatchExporter(t *testing.T) {
	ExecuteAndVerify(t, &batch_exporters.CreateBatchExporterResponse{BatchExporter: anotherBatchExporter},
		"create", "batch-exporter", streamRef.Name, "--sink="+anotherSink.Ref.Name, "--interval=300",
		"--name=another-batch-exporter", "--path-prefix=some-prefix")
}

func listBatchExporters(t *testing.T) {
	ExecuteAndVerify(t, &batch_exporters.ListBatchExportersResponse{
		BatchExporters: []*entities.BatchExporter{batchExporter, anotherBatchExporter}},
		"list", "batch-exporters")
}

func getBatchExporter(t *testing.T) {
	ExecuteAndVerify(t, &batch_exporters.GetBatchExporterResponse{BatchExporter: anotherBatchExporter},
		"get", "batch-exporter", "another-batch-exporter")

}

func listSinksRecursive(t *testing.T) {
	sinkTree := &entities.SinkTree{
		Sink:           sink,
		BatchExporters: []*entities.BatchExporter{batchExporter},
	}
	anotherSinkTree := &entities.SinkTree{
		Sink:           anotherSink,
		BatchExporters: []*entities.BatchExporter{anotherBatchExporter},
	}
	ExecuteAndVerify(t, &sinks.ListSinksResponse{Sinks: []*entities.SinkTree{sinkTree, anotherSinkTree}},
		"list", "sinks", "--recursive")
}

func getAnotherSinkRecursive(t *testing.T) {
	ExecuteAndVerify(t, &sinks.GetSinkResponse{
		SinkTree: &entities.SinkTree{
			Sink:           anotherSink,
			BatchExporters: []*entities.BatchExporter{anotherBatchExporter},
		},
	}, "get", "sink", anotherSink.Ref.Name, "--recursive")
}

func deleteAnotherSinkRecursive(t *testing.T) {
	ExecuteAndVerify(t, &sinks.DeleteSinkResponse{}, "delete", "sink", anotherSink.Ref.Name, "--recursive")
}

func deleteSinkNonRecursive(t *testing.T) {
	output := ExecuteCliAndGetOutput(t, "", "delete", "sink", sink.Ref.Name)
	if strings.Index(output, "FailedPrecondition") == -1 {
		t.Fail()
	}
}

func createAwsCredentialsFile(t *testing.T) string {
	awsCredentials := AwsCredentials{
		AccessKey: AccessKey{
			UserName:        testConfig().s3UserName,
			AccessKeyId:     testConfig().s3AccessKeyId,
			SecretAccessKey: testConfig().s3SecretAccessKey,
		},
	}
	awsCredentialsFile, awsCredentialsFileErr := ioutil.TempFile("", "aws_*.json")
	if awsCredentialsFileErr != nil {
		t.Error(awsCredentialsFileErr)
	}
	awsCredentialsJson, awsCredentialsJsonErr := json.Marshal(awsCredentials)
	if awsCredentialsJsonErr != nil {
		t.Error(awsCredentialsJsonErr)
	}
	_, err := awsCredentialsFile.WriteString(string(awsCredentialsJson))
	if err != nil {
		t.Error(err)
	}

	return awsCredentialsFile.Name()
}
