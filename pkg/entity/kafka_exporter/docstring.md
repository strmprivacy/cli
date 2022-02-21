A Kafka Exporter, like a Batch Exporter, can be used to export events
from Stream Machine to somewhere outside of STRM Privacy. But in
contrast to a Batch Exporter, a Kafka Exporter does not work in batches,
but processes the events in real time.

After creation, the CLI exposes the authentication information that is
needed to connect to it with your own Kafka Consumer.

In case your data are Avro encoded, the Kafka exporter provides a *json
format* conversion of your data for easier downstream processing. See
the [exporting Kafka](quickstart/exporting-kafka.md) page for how to
consume from the exporter.

### Usage