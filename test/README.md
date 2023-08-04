# Running tests

In order to run tests, you need to have an env file located at `test/.env`.
The contents of this file can be found
in [1Password](https://start.1password.com/open/i?a=7KRQKHUXK5G2PHTFVZ6O7L6AHU&v=qtrj3nbpr6744m3yz3bpom3lda&i=bq3d4aebenm64zluuvmthyk2ne&h=streammachine.1password.com).
The file should look like this:

```
STRM_TEST_USER_BILLING_ID=
STRM_TEST_USER_EMAIL=
STRM_TEST_USER_PASSWORD=
STRM_TEST_S3_USER_NAME=
STRM_TEST_S3_ACCESS_KEY_ID=
STRM_TEST_S3_SECRET_ACCESS_KEY=
STRM_CONFIG_PATH=
STRM_TEST_PROJECT_ID=
```

Once you have the env file, you can run the tests from the **root** of the repository with:

```
make test
```
