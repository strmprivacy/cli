const {execStrm, testConfig} = require('../support/it-support');
const tmp = require('tmp');

describe(`auth - login`, () => {
    it(`with a password specified logs the user in`,  () => {
        const tokenFile = tmp.fileSync({postfix: '.json'});

        const strm = execStrm(`auth login ${testConfig.email} --password=${testConfig.password}`, tokenFile.name);

        expect(strm.output).toBe(`billingId ${testConfig.billingId}\nsaved login to ${tokenFile.name}\n`);
    });

    // I can't get interactive input to work
    // xit(`with no password specified logs the user in after asking interactively`, async () => {
    //     const tokenFile = tmp.fileSync({postfix: '.json'});
    //     await ignoreErrors(exec(
    //         `strm2 auth login clitest@streammachine.io --token-file=${tokenFile.name}`));
    // });

});
