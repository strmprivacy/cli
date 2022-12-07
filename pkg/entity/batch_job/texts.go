package batch_job

import (
	"strmprivacy/strm/pkg/util"
)

var longDoc = util.LongDocsUsage(`
A Batch Job reads all events in a file via a Data Connector and writes them to one or more Data Connectors,
applying one of our privacy algorithms as defined by the job's configuration file. An encryption batch job
encrypts sensitive data, while a micro-aggregation batch job applies k-member clustering and replaces
the values of quasi identifier fields with an aggregated value (e.g. mean value of a cluster). 

A [Data Connector](docs/04-reference/01-cli-reference/!strm/create/data-connector.md) is a configuration
entity that comprises a location (GCS bucket, AWS S3 bucket, ...) and associated credentials.

A Data Connector must be created in the same project *before* you can create a batch job that uses it.

The policy in the Batch Job configuration file can be overridden with the policy flags.

Batch Jobs are [explained in the documentation](https://docs.strmprivacy.io/docs/latest/quickstart/batch/batch-jobs/).
`)

var example = util.DedentTrim(`
A simplified example Batch Job configuration file

{
  "policyId": "5c8e653a-8102-4444-ac15-a3d1aa0ff109",
	"source_data": {
	  "data_connector_ref": { "name": "s3-batch-demo"},
	  "file_name": "online_retail_II-small.csv"
	},
	"consent": { "default_consent_levels": [ 2 ] },
	"encryption": {
	  "batch_job_group_id": "35ced9a7-413f-49e8-9320-d17ebbc7e2d2",
	  "timestamp_config": {
		"field": "InvoiceDate",
		"format": "M/d/yyyy H:m",
		"default_time_zone": { "id": "Europe/Amsterdam" }
	  }
	},
	"event_contract_ref": { "handle": "strmprivacy", "name": "online-retail",
			"version": "1.0.0" },
	"encrypted_data": {
	  "target": {
		"data_connector_ref": { "name": "s3-batch-demo"},
		"file_name": "online_retail_II/encrypted-small.csv"
	  }
	},
	"encryption_keys_data": {
	  "target": {
		"data_connector_ref": { "name": "s3-batch-demo"},
		"file_name": "online_retail_II/keys-small.csv"
	  }
	},
	"derived_data": [      {
		"target": {
		  "data_connector_ref": { "name": "s3-batch-demo"},
		  "file_name": "online_retail_II/decrypted-0-small.csv"
		},
		"consent_levels": [ 2 ],
		"consent_level_type": "CUMULATIVE"
	  }
	]
  }
`)
