package usage

import "strmprivacy/strm/pkg/util"

var longGetDoc = util.LongDocsUsage(`
Usage allows you to see how many events were sent on a certain stream. This is currently only the events received on the
event-gateway. Extracting events via Kafka or Batch exporters is not included.

The values are interpolated from cumulative event counts, and sampled over intervals
(the --by option). The default output is csv, but json is also available.

The default range is over the last 24 hours, with a default interval of 15 minutes.
`)

var getExample = util.DedentTrim(`
strm get usage demo --by 15m --from 2021/7/27-10:00  --until 2021/7/27-12:00

from,count,duration,change,rate
2021-07-27T10:00:00.000000+0200,173478,900,NaN,NaN
2021-07-27T10:15:00.000000+0200,182422,900,8944,9.94
2021-07-27T10:30:00.000000+0200,191363,900,8941,9.93
2021-07-27T10:45:00.000000+0200,200305,900,8942,9.94
2021-07-27T11:00:00.000000+0200,209248,900,8943,9.94
2021-07-27T11:15:00.000000+0200,218192,900,8944,9.94
2021-07-27T11:30:00.000000+0200,227134,900,8942,9.94
2021-07-27T11:45:00.000000+0200,236078,900,8944,9.94
2021-07-27T12:00:00.000000+0200,245023,900,8945,9.94
`)
