package test

import (
	"encoding/json"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"testing"
)

type AccessKey struct {
	UserName        string `json:"UserName"`
	AccessKeyId     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
}

type AwsCredentials struct {
	AccessKey AccessKey `json:"AccessKey"`
}

func TestSinks(t *testing.T) {
	_ = newConfigDir()

	awsCredentialsFileName := createAwsCredentialsFile(t)

	out := ExecuteCliAndGetOutput(t, "", "create", "stream", "teststream")
	assert.Equal(t, out, `{"ref":{"billingId":"testBillingId", "name":"teststream"}, "enabled":true, "limits":{"eventRate":"999999", "eventCount":"999999999"}, "credentials":[{"clientId":"clientId", "clientSecret":"clientSecret"}]}
`)
	out = ExecuteCliAndGetOutput(t, "", "list", "sinks")
	assert.Equal(t, out, `{}
`)
	out = ExecuteCliAndGetOutput(t, "", "create", "sink", "s3sink", "strm-cli-tester", "--sink-type=S3", "--credentials-file="+awsCredentialsFileName)
	assert.Equal(t, out, `{"ref":{"billingId":"testBillingId", "name":"s3sink"}, "sinkType":"S3", "bucket":{"bucketName":"strm-cli-tester", "credentials":"{\"AccessKey\":{\"UserName\":\"UserName\",\"AccessKeyId\":\"AccessKeyId\",\"SecretAccessKey\":\"SecretAccessKey\"}}"}}
`)
	out = ExecuteCliAndGetOutput(t, "", "list", "sinks")
	assert.Equal(t, out, `{"sinks":[{"sink":{"ref":{"billingId":"testBillingId", "name":"s3sink"}, "sinkType":"S3", "bucket":{"bucketName":"strm-cli-tester"}}}]}
`)
	out = ExecuteCliAndGetOutput(t, "", "create", "batch-exporter", "teststream")
	assert.Equal(t, out, `{"ref":{"billingId":"testBillingId", "name":"s3sink-teststream"}, "streamRef":{"billingId":"testBillingId", "name":"teststream"}, "interval":"60s", "sinkName":"s3sink"}
`)
	out = ExecuteCliAndGetOutput(t, "", "create", "sink", "another-sink", "strm-cli-tester", "--sink-type=S3", "--credentials-file="+awsCredentialsFileName)
	assert.Equal(t, out, `{"ref":{"billingId":"testBillingId", "name":"another-sink"}, "sinkType":"S3", "bucket":{"bucketName":"strm-cli-tester", "credentials":"{\"AccessKey\":{\"UserName\":\"UserName\",\"AccessKeyId\":\"AccessKeyId\",\"SecretAccessKey\":\"SecretAccessKey\"}}"}}
`)
	out = ExecuteCliAndGetOutput(t, "", "create", "batch-exporter", "teststream", "--sink=another-sink", "--interval=300", "--name=another-batch-exporter", "--path-prefix=some-prefix")
	assert.Equal(t, out, `{"ref":{"billingId":"testBillingId", "name":"another-batch-exporter"}, "streamRef":{"billingId":"testBillingId", "name":"teststream"}, "interval":"300s", "sinkName":"another-sink", "pathPrefix":"some-prefix"}
`)
	out = ExecuteCliAndGetOutput(t, "", "list", "batch-exporters")
	assert.Equal(t, out, `{"batchExporters":[{"ref":{"billingId":"testBillingId", "name":"s3sink-teststream"}, "streamRef":{"billingId":"testBillingId", "name":"teststream"}, "interval":"60s", "sinkName":"s3sink"}, {"ref":{"billingId":"testBillingId", "name":"another-batch-exporter"}, "streamRef":{"billingId":"testBillingId", "name":"teststream"}, "interval":"300s", "sinkName":"another-sink", "pathPrefix":"some-prefix"}]}
`)
	out = ExecuteCliAndGetOutput(t, "", "get", "batch-exporter", "another-batch-exporter")
	assert.Equal(t, out, `{"batchExporter":{"ref":{"billingId":"testBillingId", "name":"another-batch-exporter"}, "streamRef":{"billingId":"testBillingId", "name":"teststream"}, "interval":"300s", "sinkName":"another-sink", "pathPrefix":"some-prefix"}}
`)
	out = ExecuteCliAndGetOutput(t, "", "list", "sinks", "--recursive")
	assert.Equal(t, out, `{"sinks":[{"sink":{"ref":{"billingId":"testBillingId", "name":"s3sink"}, "sinkType":"S3", "bucket":{"bucketName":"strm-cli-tester"}}, "batchExporters":[{"ref":{"billingId":"testBillingId", "name":"s3sink-teststream"}, "streamRef":{"billingId":"testBillingId", "name":"teststream"}, "interval":"60s", "sinkName":"s3sink"}]}, {"sink":{"ref":{"billingId":"testBillingId", "name":"another-sink"}, "sinkType":"S3", "bucket":{"bucketName":"strm-cli-tester"}}, "batchExporters":[{"ref":{"billingId":"testBillingId", "name":"another-batch-exporter"}, "streamRef":{"billingId":"testBillingId", "name":"teststream"}, "interval":"300s", "sinkName":"another-sink", "pathPrefix":"some-prefix"}]}]}
`)
	out = ExecuteCliAndGetOutput(t, "", "get", "sink", "another-sink", "--recursive")
	assert.Equal(t, out, `{"sinkTree":{"sink":{"ref":{"billingId":"testBillingId", "name":"another-sink"}, "sinkType":"S3", "bucket":{"bucketName":"strm-cli-tester"}}, "batchExporters":[{"ref":{"billingId":"testBillingId", "name":"another-batch-exporter"}, "streamRef":{"billingId":"testBillingId", "name":"teststream"}, "interval":"300s", "sinkName":"another-sink", "pathPrefix":"some-prefix"}]}}
`)
	out = ExecuteCliAndGetOutput(t, "", "delete", "sink", "another-sink", "--recursive")
	assert.Equal(t, out, `{}
`)
	out = ExecuteCliAndGetOutput(t, "", "delete", "sink", "s3sink")
	assert.Equal(t, out, `Error: rpc error: code = FailedPrecondition desc = Cannot delete sink with name s3sink, as it still has exporters linked to it. Delete those first before deleting this sink.
`)
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
