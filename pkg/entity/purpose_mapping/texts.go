package purpose_mapping

import "strmprivacy/strm/pkg/util"

var longDocCreate = util.LongDocsUsage(`
A purpose mapping is a mapping between a name and an integer value. This integer value is automatically assigned to
the name when creating a new purpose mapping. The mapping is intended to make it easier to work with purposes when
creating data contracts and streams. The integer value that the purpose mapping is assigned to is used in the
actual transport / creation of events, as an integer value is more efficient to transport than a string.

Please note that purpose mappings cannot be deleted once created. This is to ensure that the integer values

`)

var createExample = `
strm create purpose-mapping "Legitimate Interest"
`

var getExample = util.DedentTrim(`
strm get purpose-mapping 0

 PURPOSE MAPPING              VALUE

 Legitimate Interest          0
`)
var listExample = util.DedentTrim(`
strm list purpose-mappings

 PURPOSE MAPPING              VALUE

 Legitimate Interest          0
 Marketing                    1
`)
