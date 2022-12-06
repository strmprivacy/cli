package kafka_exporter

import "strmprivacy/strm/pkg/util"

var longDoc = util.LongDocsUsage(`
A Kafka Exporter, like a Batch Exporter, can be used to export events from Stream Machine to somewhere outside of STRM
Privacy. But in contrast to a Batch Exporter, a Kafka Exporter does not work in batches, but processes the events in
real time.

The Kafka exporter produces your events in JSON format, even when originally in Avro binary for easier
downstream processing.`)

var longDeleteDoc = util.LongDocsUsage(`
Deletes a Kafka Exporter. 

If it has dependent entities (like Kafka Users), you can use
the 'recursive' option to get rid of those also.

Returns everything that was deleted.
`)

var exampleList = util.DedentTrim(`
Somewhat shortened.

strm list kafka-exporters -o json
{
    "kafkaExporters": [
        {
            "ref": {
                "name": "shared-export-austindemo"
            },
            "streamRef": {
                "name": "austindemo"
            },
            "target": {
                "clusterRef": {
                    "name": "shared-export"
                },
                "topic": "export-c42dc1f5-43f9-4672-8ddc-8865df355ea9"
            }
        },
        ...
    ]
}`)
