package stream

import "strmprivacy/strm/pkg/util"

var longDocCreate = util.LongDocsUsage(`
A stream is a pipeline implementation in STRM Privacy using streaming technology (Kafka).

Events are sent to the "Event Gateway" into an "Input Stream", where *all* personal data event attributes are encrypted
after validating that the event conforms to a data contract. "Privacy Streams" are derived from input streams.

A derived "Privacy Stream" is configured with one or more consent levels and it only receives events matching those
consent levels. The PII attributes with matching consent are decrypted, while non-matching attributes remain encrypted.

Every stream has its own set of access tokens. Connecting to an input stream requires different credentials than when
connecting to a derived Privacy Stream.

`)

var createExample = `
strm create stream test

A name is not required for a derived stream; when absent it will be created from the derived stream
and the consent-levels.

strm create stream --derived-from test --levels 1,3,8 --consent-type GRANULAR test-marketing
`

var getExample = util.DedentTrim(`
strm get stream demo -o json
{
    "streamTree": {
        "stream": {
            "ref": { "name": "demo", },
            "enabled": true,
            "limits": { "eventRate": "10000", "eventCount": "1000000" },
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
 STREAM              DERIVED   CONSENT LEVEL TYPE   CONSENT LEVELS   ENABLED   POLICY NAME

 e-commerce-masked   true      CUMULATIVE           [1]              true
 ecommerce-1         true      CUMULATIVE           [1]              true
 ecommerce           false                          []               true
 ecommerce-2         true      CUMULATIVE           [2]              true
 demo                false                          []               true
 demo-0-1            true      GRANULAR             [0 1]            true
`)
