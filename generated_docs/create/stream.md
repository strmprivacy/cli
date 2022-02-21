## dstrm create stream

Create a stream

### Synopsis

A stream is the central resource in STRM Privacy. Clients can connect to
a stream to send and to receive events. A stream can be either an "input
stream" or a "derived stream".

Events are always sent to an input stream. Sending events to a derived
stream is not possible. After validation and encryption of all PII
fields, STRM Privacy sends all events to the input stream. Clients
consuming from the input stream will see all events, but with all PII
fields encrypted.

Derived streams can be made on top of a input stream. A derived stream
is configured with one or more consent levels and it only receives
events with matching consent levels (see details about this matching
process here). The PII fields with matching consent levels are decrypted
and sent to the derived stream. Clients connecting to the derived stream
will only receive the events on this stream.

Every stream has its own set of access tokens. Connecting to an input
stream requires different credentials than when connecting to a derived
stream.

### Usage

```
dstrm create stream [name] [flags]
```

### Options

```
      --consent-type string         CUMULATIVE or GRANULAR (default "CUMULATIVE")
  -D, --derived-from string         name of stream that this stream is derived from
      --description string          description
  -h, --help                        help for stream
  -L, --levels int32Slice           comma separated list of integers for derived streams (default [])
      --mask-seed string            A seed used for masking
  -M, --masked-fields stringArray   -M strmprivacy/example/1.3.0:sensitiveValue,consistentValue \
                                    -M strmprivacy/clickstream/1.0.0:sessionId
                                    
                                    Masks fields values in the output stream via hashing.
                                    	
      --save                        if true, save the result in the config directory (~/.config/strmprivacy/saved-entities). (default is true) (default true)
      --tags strings                tags
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

