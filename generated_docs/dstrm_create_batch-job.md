---
date: 2022-02-18T14:37:11+01:00
title: "dstrm create batch-job"
slug: dstrm_create_batch-job
url: /commands/dstrm_create_batch-job/
---
## dstrm create batch-job

Create a Batch Job

```
dstrm create batch-job [flags]
```

### Options

```
  -F, --file string   The path to the JSON file containing the batch job configuration
  -h, --help          help for batch-job
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

###### Auto generated by spf13/cobra on 18-Feb-2022