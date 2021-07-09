const {execStrm} = require('../support/it-support');

describe(`auth - access-token`, () => {
    it(`outputs the access-token when logged in`, () => {
        const strm = execStrm('auth access-token', 'spec/support/simple-token.json');

        expect(strm.output).toBe('id.token.test\n');
    });

    it(`outputs an error message when not logged in`, () => {

        try {
            execStrm('auth access-token', 'nonexistent.json');
            fail('should not get here');
        } catch (e) {
            expect(e.message).toContain('Error: No login information found. Use: `strm2 auth login` first.\n');
        }
    });
});

