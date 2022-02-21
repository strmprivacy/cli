A Batch Exporter listens to a stream and outputs all events to files in
a Sink. This happens with a regular interval.

Each file follows the JSON Lines format, which is one full JSON document
per line.

A [sink](sink.md) is a configuration item that defines location
(Gcloud, AWS, ..) bucket and associated credentials.

A sink needs to be created *before* you can create a batch exporter that
uses it.

### Usage