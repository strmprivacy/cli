package stream

import "strmprivacy/strm/pkg/util"

var longDocCreate = util.LongDocsUsage(`
A stream is a pipeline implementation in STRM Privacy using streaming technology (Kafka).

Events are sent through the "Event Gateway" to a "Source Stream", where *all* sensitive (personal) data event attributes are encrypted
after validating that the event conforms to a data contract. "Privacy Streams" are derived from source streams.

A derived "Privacy Stream" is configured for one or more specific purposes. It only receives events that are allowed
to be processed for the configured purposes. The sensitive attributes matching those purposes are decrypted, while 
non-matching attributes remain encrypted.

Every stream has its own set of access tokens. Connecting to a source stream (i.e. sending events) requires different 
credentials than connecting to a derived Privacy Stream.

`)

var createExample = `
strm create stream test

A name is not required for a derived stream; when absent a name will be created based on the source stream
and the provided purposes.

strm create stream --derived-from test --purposes 1,3,8 test-marketing
`

var getExample = util.DedentTrim(`
strm get stream demo -o json
{
    "streamTree": {
        "stream": {
            "ref": { "name": "demo", },
            "enabled": true,
            "limits": { "eventRate": "10000", "eventCount": "10000000" },
            "credentials": [
                {
                    "clientId": "stream-2459..",
                    "clientSecret": "3rvIUfNi.."
                }
            ],
            "maskedFields": { "seed": "****" }
        }
    }
}
`)
var listExample = util.DedentTrim(`
strm list streams
 STREAM              DERIVED   PURPOSES   ENABLED   POLICY NAME

 e-commerce-masked   true      [1]        true
 ecommerce-1         true      [1]        true
 ecommerce           false     []         true
 ecommerce-2         true      [2]        true
 demo                false     []         true
 demo-0-1            true      [0 1]      true
`)
