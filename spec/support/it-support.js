const env = require('dotenv').config();
const tmp = require('tmp');
const execSync = require('child_process').execSync;
const JasmineConsoleReporter = require('jasmine-console-reporter');

tmp.setGracefulCleanup();

const reporter = new JasmineConsoleReporter({
    colors: 1,           // (0|false)|(1|true)|2
    cleanStack: 1,       // (0|false)|(1|true)|2|3
    verbosity: 4,        // (0|false)|1|2|(3|true)|4|Object
    listStyle: 'indent', // "flat"|"indent"
    timeUnit: 'ms',      // "ms"|"ns"|"s"
    timeThreshold: { ok: 500, warn: 1000, ouch: 3000 }, // Object|Number
    activity: false,     // boolean or string ("dots"|"star"|"flip"|"bouncingBar"|...)
    emoji: true,
    beep: true
});

jasmine.getEnv().addReporter(reporter);

const testConfig = {
    billingId: env.parsed['STRM_TEST_BILLING_ID'],
    password: env.parsed['STRM_TEST_PASSWORD'],
    email: env.parsed['STRM_TEST_EMAIL'],
    s3AccessKey: env.parsed['STRM_TEST_S3_ACCESS_KEY_ID'],
    s3SecretKey: env.parsed['STRM_TEST_S3_SECRET_ACCESS_KEY']
};

// We ignore thrown errors, since these should already be exposed through stderr. stdin/out/err is the interface to the user, so these are
// the only things we test.
const ignoreErrors = exec => exec.catch(obj => obj);

const defaultTokenFile = 'spec/support/.default-login-token.json';

const execStrm = (subCommands, tokenFile) => {
    if (!tokenFile) {
        tokenFile = defaultTokenFile;
    }
    let configPath = process.cwd();
    let procString = `strm2 --token-file=${tokenFile} --config-path=${configPath} ${subCommands}`;
    let output = execSync(procString);

    return {
        output: '' + output,
        tokenFile: tokenFile.name,
        configPath: configPath
    }
};

const setUp = () => {
    // Make sure there is always a valid login token present for testing
    execStrm(`auth login ${testConfig.email} --password=${testConfig.password}`, defaultTokenFile)

    // Make sure that everything has been cleaned up before the test
    const streamsOutput = execStrm('list streams').output;
    let streams;
    try { streams = JSON.parse(streamsOutput).streams } catch (e) {}
    (streams || []).forEach(stream => {
        try {
            execStrm(`delete stream ${stream.stream.ref.name} --recursive`);
        } catch (e) {
            // Don't fail on missing streams, since list streams outputs both input and derived streams
        }
    });

    const sinksOutput = execStrm('list sinks').output;
    let sinks;
    try { sinks = JSON.parse(sinksOutput).sinks } catch (e) {}
    (sinks || []).forEach(sink => {
        execStrm(`delete sink ${sink.sink.ref.name} --recursive`);
    });
};

setUp();

const matchers = {
    toEqualExcept: function () {
        return {
            compare: function (actualResult, expectedResult, propertyPaths) {
                propertyPaths = propertyPaths.map(p => '.' + p);

                const checkObject = (actual, expected, path) => {
                    const hasToCheck = (key, value) => {
                        return typeof (value) !== 'object' && propertyPaths.indexOf(path + '.' + key) === -1;
                    };

                    if (JSON.stringify((Object.keys(actual)).sort()) !== JSON.stringify((Object.keys(expected)).sort())) {
                        throw {
                            pass: false,
                            message: `keys are different in: '${path}'. Actual: ${Object.keys(actual).sort()}, expected: ${Object.keys(expected).sort()}`
                        };
                    } else {
                        Object.keys(expected).forEach(expectedKey => {
                            let actualValue = actual[expectedKey];
                            let expectedValue = expected[expectedKey];

                            if (hasToCheck(expectedKey, expectedValue) && actualValue !== expectedValue) {
                                throw {
                                    pass: false,
                                    message: `Expected value: ${expectedValue} for key: ${path}.${expectedKey} is different than actual: ${actualValue}`
                                };
                            }

                            if (typeof (expectedValue) === 'object' && typeof (expectedValue) !== 'string') {
                                if (Array.isArray(expectedValue)) {
                                    if (!Array.isArray(actualValue)) {
                                        throw {
                                            pass: false,
                                            message: `Expected value for key: ${path}.${expectedKey} is an array, but actual is not: ${actualValue}`
                                        };
                                    }
                                    for (let i = 0; i < expectedValue.length; i++) {
                                        const actualItem = actualValue[i];
                                        const expectedItem = expectedValue[i];
                                        if (typeof (expectedItem) === 'object') {
                                            checkObject(actualItem, expectedItem, path + '.' + expectedKey);
                                        } else {
                                            if (actualItem !== expectedItem) {
                                                throw {
                                                    pass: false,
                                                    message: `Expected value for key: ${path}.${expectedKey} is invalid. Expected: ${expectedItem} is not: ${actualItem}`
                                                };
                                            }
                                        }
                                    }
                                } else {
                                    checkObject(actualValue, expectedValue, path + '.' + expectedKey)
                                }
                            } else {
                                if (hasToCheck(expectedKey, expectedValue) && actualValue !== expectedValue) {
                                    throw {
                                        pass: false,
                                        message: `Expected value for key: ${path}.${expectedKey} is invalid. Expected: ${expectedValue} is not: ${actualValue}`
                                    };
                                }
                            }
                        });
                    }
                };

                try {
                    checkObject(actualResult, expectedResult, '');
                    return {
                        pass: true
                    };
                } catch (e) {
                    return e;
                }
            }
        }
    }
};

module.exports = {
    execStrm,
    env,
    testConfig,
    matchers,
    setUp
};
