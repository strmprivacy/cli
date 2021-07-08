const { execStrm, testConfig, matchers, setUp } = require('../support/it-support');
const tmp = require('tmp');
const fs = require('fs')

const awsAccessKey = `{"AccessKey":{"UserName":"cli-tester-writeonly","AccessKeyId":"${testConfig.s3AccessKey}","SecretAccessKey":"${testConfig.s3SecretKey}"}}`;

describe(`sinks`, () => {
    beforeEach(() => {
        jasmine.addMatchers(matchers);
        setUp();

    });
    it(`spec`,  () => {
        const accessKeyFile = tmp.fileSync({postfix: '.json'});
        const accessKeyData = JSON.stringify({
           AccessKey: {
               UserName: 'cli-tester-writeonly',
               AccessKeyId: testConfig.s3AccessKey,
               SecretAccessKey: testConfig.s3SecretKey,
           }
        });

        fs.writeFileSync(accessKeyFile.name, accessKeyData, {encoding: 'utf8'})

        execStrm('create stream teststream');

        let strm = execStrm('list sinks');

        expect(JSON.parse(strm.output)).toEqual({});
        strm = execStrm(`create sink s3sink strm-cli-tester --sink-type S3 --credentials-file ${accessKeyFile.name}`);
        expect(JSON.parse(strm.output)).toEqual({ref: {billingId: testConfig.billingId, name: 's3sink'}, sinkType: 'S3', bucket: {bucketName: 'strm-cli-tester', credentials: awsAccessKey}});

        strm = execStrm('list sinks');
        expect(JSON.parse(strm.output)).toEqual({sinks: [ {sink: {ref: {billingId: testConfig.billingId, name: 's3sink'}, sinkType: 'S3', bucket: {bucketName: 'strm-cli-tester'}}} ]});

        strm = execStrm(`create batch-exporter teststream`);
        expect(JSON.parse(strm.output)).toEqual({ref: { billingId: testConfig.billingId, name: 's3sink-teststream' }, streamRef: {billingId: testConfig.billingId, name: 'teststream'}, interval: '60s', sinkName: 's3sink' });

        strm = execStrm(`create sink another-sink strm-cli-tester --sink-type S3 --credentials-file ${accessKeyFile.name}`);
        expect(JSON.parse(strm.output)).toEqual({ref: {billingId: testConfig.billingId, name: 'another-sink'}, sinkType: 'S3', bucket: { bucketName: 'strm-cli-tester', credentials: awsAccessKey } });

        strm = execStrm(`create batch-exporter teststream --sink another-sink --interval 300 --name another-batch-exporter --path-prefix some-prefix`);
        expect(JSON.parse(strm.output)).toEqual({ref: {billingId: testConfig.billingId, name: 'another-batch-exporter'}, streamRef: {billingId: testConfig.billingId, name: 'teststream'}, interval: '300s', sinkName: 'another-sink', pathPrefix: 'some-prefix' });

        strm = execStrm(`list batch-exporters`);
        expect(JSON.parse(strm.output)).toEqual({ batchExporters: [ { ref: { billingId: testConfig.billingId, name: 's3sink-teststream' }, streamRef: {billingId: testConfig.billingId, name: 'teststream'}, interval: '60s', sinkName: 's3sink' }, { ref: {billingId: testConfig.billingId, name: 'another-batch-exporter'}, streamRef: {billingId: testConfig.billingId, name: 'teststream'}, interval: '300s', sinkName: 'another-sink', pathPrefix: 'some-prefix' } ]});

        strm = execStrm(`get batch-exporter another-batch-exporter`);
        expect(JSON.parse(strm.output)).toEqual({ batchExporter: {ref: {billingId: testConfig.billingId, name: 'another-batch-exporter'}, streamRef: {billingId: testConfig.billingId, name: 'teststream'}, interval: '300s', sinkName: 'another-sink', pathPrefix: 'some-prefix' }});

        strm = execStrm(`list sinks --recursive`);
        expect(JSON.parse(strm.output)).toEqual({sinks: [ {sink: {ref: {billingId: testConfig.billingId, name: 's3sink'}, sinkType: 'S3', bucket: {bucketName: 'strm-cli-tester'}}}, { sink: { ref: {billingId: testConfig.billingId, name: 'another-sink'}, sinkType: 'S3', bucket: { bucketName: 'strm-cli-tester' } } } ]});

        strm = execStrm(`get sink another-sink --recursive`);
        expect(JSON.parse(strm.output)).toEqual({sinkTree:{ sink: { ref: {billingId: testConfig.billingId, name: 'another-sink'}, sinkType: 'S3', bucket: { bucketName: 'strm-cli-tester' } } }});

        strm = execStrm('delete sink another-sink --recursive');
        expect(JSON.parse(strm.output)).toEqual({});

        try {
            execStrm('delete sink s3sink');
            fail('should not get here');
        } catch (e) {
            expect(e.message).toContain('Cannot delete sink with name s3sink, as it still has exporters linked to it. Delete those first before deleting this sink.');
        }
    });

});
