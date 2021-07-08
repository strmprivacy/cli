const { execStrm, testConfig, matchers, setUp } = require('../support/it-support');

describe(`streams`, () => {

    beforeEach(() => {
        jasmine.addMatchers(matchers);
        setUp();
    });

    it(`spec`,  () => {
        let strm = execStrm('list streams');

        // TODO: CLL should not return a literal value of "undefined"
        expect(strm.output).toEqual('{}\n');

        strm = execStrm('create stream clitest');
        expect(JSON.parse(strm.output)).toEqualExcept({"ref":{"billingId":testConfig.billingId, "name":"clitest"}, "enabled":true, "limits":{"eventRate":"999999", "eventCount":"999999999"}, "credentials":[{"clientId":"", "clientSecret":""}]}, ['credentials.clientId', 'credentials.clientSecret']);

        strm = execStrm('list streams');
        // TODO: The consentLevelType has to be removed here. This is invalid for a source stream
        expect(JSON.parse(strm.output)).toEqualExcept({"streams":[{"stream":{"ref":{"billingId":testConfig.billingId, "name":"clitest"}, "enabled":true, "limits":{"eventRate":"999999", "eventCount":"999999999"}, "credentials":[{"clientId":""}]}}]}, ['streams.stream.credentials.clientId', 'streams.stream.credentials.clientSecret']);

        strm = execStrm('create stream clitest-with-tags --tags=foo,bar,baz');
        expect(JSON.parse(strm.output)).toEqualExcept({"ref":{"billingId": testConfig.billingId, "name":"clitest-with-tags"}, "enabled":true, "limits":{"eventRate":"999999", "eventCount":"999999999"}, "tags":["foo", "bar", "baz"], "credentials":[{"clientId":".+", "clientSecret":".+"}]}, ['credentials.clientId', 'credentials.clientSecret']);

        strm = execStrm('create stream --derived-from=clitest-with-tags --levels=2');
        expect(JSON.parse(strm.output)).toEqualExcept({"ref":{"billingId":testConfig.billingId, "name":"clitest-with-tags-2"}, "consentLevels":[2], "consentLevelType":"CUMULATIVE", "enabled":true, "linkedStream":"clitest-with-tags", "credentials":[{"clientId":"", "clientSecret":""}]}, ['credentials.clientId', 'credentials.clientSecret']);

        strm = execStrm('get stream clitest-with-tags');
        expect(JSON.parse(strm.output)).toEqualExcept({"streamTree":{"stream":{"ref":{"billingId":testConfig.billingId, "name":"clitest-with-tags"}, "enabled":true, "limits":{"eventRate":"999999", "eventCount":"999999999"}, "tags":["foo", "bar", "baz"], "credentials":[{"clientId":""}]}}}, ['streamTree.stream.credentials.clientId']);

        strm = execStrm('delete stream clitest-with-tags --recursive');
        expect(JSON.parse(strm.output)).toEqualExcept({"streamTree":{"stream":{"ref":{"billingId":testConfig.billingId, "name":"clitest-with-tags"}, "enabled":true, "limits":{"eventRate":"999999", "eventCount":"999999999"}, "tags":["foo", "bar", "baz"], "credentials":[{"clientId":"d9al1sljl17bve1m1y5r432fcjx1ly"}]}, "keyStream":{"ref":{"billingId":"clitestdev7744696560", "name":"clitest-with-tags"}}, "derived":[{"ref":{"billingId":"clitestdev7744696560", "name":"clitest-with-tags-2"}, "consentLevels":[2], "consentLevelType":"CUMULATIVE", "enabled":true, "limits":{}, "linkedStream":"clitest-with-tags", "credentials":[{"clientId":""}]}]}}, ['streamTree.stream.credentials.clientId', 'streamTree.derived.credentials.clientId']);
    });

});
