## dstrm create sink

Create sink

### Synopsis

A Sink is a STRM Privacy configuration object for a remote file storage.
For now, AWS S3 and Google Cloud Storage Buckets are supported. By
itself, a Sink does nothing. A Batch Exporter needs to be connected to a
Sink and a Stream to start outputting events.

Upon creation, STRM Privacy validates whether or not the Bucket exists
and if it is accessible with the given credentials.

### Usage

```
dstrm create sink [sink-name] [bucket-name] [flags]
```

### Options

```
      --assume-role-arn string    ARN of the role to assume
      --credentials-file string   file with credentials
  -h, --help                      help for sink
      --sink-type string          S3 or GCLOUD
```

### Options inherited from parent commands

```
      --api-auth-url string            User authentication host (default "https://accounts.strmprivacy.io")
      --api-host string                API host and port (default "api.strmprivacy.io:443")
      --events-auth-url string         Event authentication host (default "https://sts.strmprivacy.io")
      --kafka-bootstrap-hosts string   Kafka bootstrap brokers, separated by comma (default "export-bootstrap.kafka.strmprivacy.io:9092")
  -o, --output string                  Output format [json, json-raw, table, plain] (default "table")
      --token-file string              Token file that contains an access token (default is $HOME/.config/strmprivacy/credentials-<api-auth-url>.json)
      --web-socket-url string          Websocket to receive events from (default "wss://websocket.strmprivacy.io/ws")
```

### SEE ALSO

* [dstrm create](dstrm_create.md)	 - Create an entity

