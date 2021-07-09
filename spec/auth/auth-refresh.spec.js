const {execStrm, testConfig} = require('../support/it-support');
const tmp = require('tmp');
const fs = require('fs')

describe(`auth - refresh`, () => {

    it(`refreshes the token`,  (done) => {
        const tokenFile = tmp.fileSync({postfix: '.json'});

        const strmLogin = execStrm(`auth login ${testConfig.email} --password=${testConfig.password}`, tokenFile.name);

        expect(strmLogin.output).toBe(`billingId ${testConfig.billingId}\nsaved login to ${tokenFile.name}\n`);
        const oldToken = JSON.parse(fs.readFileSync(tokenFile.name, 'utf8')).idToken;

        setTimeout(() => {
            const strm = execStrm(`auth refresh`, tokenFile.name);
            expect(strm.output).toBe(``);

            const newToken = JSON.parse(fs.readFileSync(tokenFile.name, 'utf8')).idToken;
            expect(newToken).not.toBe(oldToken);

            done();
        }, 1000)
    });
});
