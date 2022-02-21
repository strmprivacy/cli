## dstrm create event-contract

Create an event-contract with reference 'handle/name/version'

### Synopsis

A Batch Exporter listens to a stream and outputs all events to files in
a Sink. This happens with a regular interval.

Each file follows the JSON Lines format, which is one full JSON document
per line.

A [sink](sink.md) is a configuration item that defines location
(Gcloud, AWS, ..) bucket and associated credentials.

A sink needs to be created *before* you can create a batch exporter that
uses it.

### Usage

```
dstrm create event-contract (handle/name/version) [flags]
```

### Options

```
  -F, --definition-file string   The path to the file with the keyField, and possibly piiFields and validations. Example JSON definition file:
                                 {
                                     "keyField": "sessionId",
                                     "piiFields": {
                                         "sessionId": 2,
                                         "referrerUrl": 1
                                     },
                                     "validations": [
                                         {
                                             "field": "referrerUrl",
                                             "type": "regex",
                                             "value": "^.*strmprivacy.*$"
                                         }
                                     ]
                                 }
  -h, --help                     help for event-contract
  -P, --public                   Public visibility of the Event Contract (allow others to use this contract)
  -S, --schema-ref string        The Serialization Schema to which this Event Contract is linked
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

