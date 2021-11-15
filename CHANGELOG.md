## [1.12.6](https://github.com/streammachineio/cli/compare/v1.12.5...v1.12.6) (2021-11-15)


### Bug Fixes

* cross compile fixed + switched to sarama kafka consumer ([#68](https://github.com/streammachineio/cli/issues/68)) ([0e728e3](https://github.com/streammachineio/cli/commit/0e728e338b77c5fb8c760195067239e759fac149))

## [1.12.5](https://github.com/streammachineio/cli/compare/v1.12.4...v1.12.5) (2021-11-12)


### Bug Fixes

* cgo=1 windows=0 ([12e6811](https://github.com/streammachineio/cli/commit/12e6811ff2cadee7025bc6ebc22793f4f990bcc4))

## [1.12.4](https://github.com/streammachineio/cli/compare/v1.12.3...v1.12.4) (2021-11-12)


### Bug Fixes

* cgo disabled ([c42f1b9](https://github.com/streammachineio/cli/commit/c42f1b98b016be4f3619c40d06dc81166dcd977c))

## [1.12.3](https://github.com/streammachineio/cli/compare/v1.12.2...v1.12.3) (2021-11-12)


### Bug Fixes

* goarch amd64 only ([2457a88](https://github.com/streammachineio/cli/commit/2457a8882a8c3a5931069a6b7cf83d87190cf26d))

## [1.12.2](https://github.com/streammachineio/cli/compare/v1.12.1...v1.12.2) (2021-11-12)


### Bug Fixes

* cgo ([47f1013](https://github.com/streammachineio/cli/commit/47f10131956171c3f57134bbb7424dc3085918b6))

## [1.12.1](https://github.com/streammachineio/cli/compare/v1.12.0...v1.12.1) (2021-11-12)


### Bug Fixes

* temporarily disable failing windows build ([6eef102](https://github.com/streammachineio/cli/commit/6eef102aa3e288be9e79fd8efc1c21fab9aa2c58))

# [1.12.0](https://github.com/streammachineio/cli/compare/v1.11.0...v1.12.0) (2021-11-12)


### Features

* updated readme ([3a3770e](https://github.com/streammachineio/cli/commit/3a3770e878ba0670b9a4addadab3db26761c77f1))

# [1.11.0](https://github.com/streammachineio/cli/compare/v1.10.1...v1.11.0) (2021-11-01)


### Features

* masked fields ([#63](https://github.com/streammachineio/cli/issues/63)) ([fbecd35](https://github.com/streammachineio/cli/commit/fbecd35266177f71f12c4be01540b1bbc9b73a5d))

## [1.10.1](https://github.com/streammachineio/cli/compare/v1.10.0...v1.10.1) (2021-10-26)


### Bug Fixes

* error messages via util.CliExit ([fc80631](https://github.com/streammachineio/cli/commit/fc80631f87573b23f2a4599778c7ccdf1d990ad0))

# [1.10.0](https://github.com/streammachineio/cli/compare/v1.9.2...v1.10.0) (2021-10-15)


### Features

* add AssumeRole to S3 Sink ([#55](https://github.com/streammachineio/cli/issues/55)) ([625bc5a](https://github.com/streammachineio/cli/commit/625bc5acc989635d4e91200c0e8a95cbd7ea3251))

## [1.9.2](https://github.com/streammachineio/cli/compare/v1.9.1...v1.9.2) (2021-10-05)


### Bug Fixes

* **event-auth:** only run authenticate at the start of event auth and not for every sent event ([#53](https://github.com/streammachineio/cli/issues/53)) ([907097f](https://github.com/streammachineio/cli/commit/907097fad7abdae5381fcc5694c460a563257ae1))

## [1.9.1](https://github.com/streammachineio/cli/compare/v1.9.0...v1.9.1) (2021-10-05)


### Bug Fixes

* proper auth success and failure redirects after login ([#52](https://github.com/streammachineio/cli/issues/52)) ([adac9a6](https://github.com/streammachineio/cli/commit/adac9a60e47d46f3d59f0656f80923c7fe508c93))

# [1.9.0](https://github.com/streammachineio/cli/compare/v1.8.1...v1.9.0) (2021-10-04)


### Features

* replace login flow with OAuth2.0 Authorization Code flow ([#51](https://github.com/streammachineio/cli/issues/51)) ([9b6d871](https://github.com/streammachineio/cli/commit/9b6d871f4dac8dda6445a2373043cabe9fa8c679))

## [1.8.1](https://github.com/streammachineio/cli/compare/v1.8.0...v1.8.1) (2021-08-19)


### Bug Fixes

* **strm.yaml:** configuration directory wasn't set yet ([#46](https://github.com/streammachineio/cli/issues/46)) ([4023657](https://github.com/streammachineio/cli/commit/4023657fd629eab88dac2b8ca0d6c0ea049def68))

# [1.8.0](https://github.com/streammachineio/cli/compare/v1.7.0...v1.8.0) (2021-08-18)


### Features

* **config:** create default config if not exists; save all entities to a a subdir under configpath ([e8380d5](https://github.com/streammachineio/cli/commit/e8380d577d0e4f503c8181499800d4079a4a95ac))

# [1.7.0](https://github.com/streammachineio/cli/compare/v1.6.0...v1.7.0) (2021-08-18)


### Features

* **context:** add context command to print CLI context and show details of saved entities ([a14b583](https://github.com/streammachineio/cli/commit/a14b583a90892e4a8f2c809e0ddba06d2582c465))

# [1.6.0](https://github.com/streammachineio/cli/compare/v1.5.0...v1.6.0) (2021-08-18)


### Features

* **batch-exporters:** implement plain, table, json, json-raw printers for batch-exporters ([1a8ebc0](https://github.com/streammachineio/cli/commit/1a8ebc030da793c8b8e8c4a80b03d3159e1b836a))
* **event-contract:** create event contract now partially asks for cli flags and a definition ([dfa0d74](https://github.com/streammachineio/cli/commit/dfa0d74387936c730504eecd96c9a3c32b627806)), closes [#33](https://github.com/streammachineio/cli/issues/33)
* **event-contracts:** implement plain, table, json, json-raw printers for event-contracts ([8f5931f](https://github.com/streammachineio/cli/commit/8f5931f0cd0fe42f9c53355af857c56ca6f90295))
* **kafka-clusters:** implement plain, table, json, json-raw printers for kafka-clusters ([6dbe06e](https://github.com/streammachineio/cli/commit/6dbe06ed6c08e9d32749e5140315b7a1bcc44404))
* **kafka-exporters:** implement plain, table, json, json-raw printers for kafka-exporters ([04ee444](https://github.com/streammachineio/cli/commit/04ee44467c338527cbfa135c751e2d74f435ca3c))
* **kafka-users:** implement plain, table, json, json-raw printers for kafka-users ([fa6473e](https://github.com/streammachineio/cli/commit/fa6473efe68abadca4c98ea0b0b2a453ccf1c435))
* **key-streams:** implement plain, table, json, json-raw printers for key-streams ([e8efd41](https://github.com/streammachineio/cli/commit/e8efd4162b8e785b00db6bf6c3a62b061d2e8a6e))
* **printers:** added interface for printers; implemented list streams printer ([39b183a](https://github.com/streammachineio/cli/commit/39b183ae5b55959868b35d36026d32deddba7a9e)), closes [#14](https://github.com/streammachineio/cli/issues/14)
* **schemas:** implement plain, table, json, json-raw printers for schemas ([f205174](https://github.com/streammachineio/cli/commit/f20517479bcf09178502af13ba4372b69b563871))
* **sinks:** implement plain, table, json, json-raw printers for sinks ([6dc2c77](https://github.com/streammachineio/cli/commit/6dc2c772d00b27414da15856c8cff8115d824e48))
* **streams:** implement plain, table, json, json-raw printers for streams ([7d10c42](https://github.com/streammachineio/cli/commit/7d10c42cad8d52a6a4d3f3ca64420a46798b65a0))
* **usage:** added missing column in usage csv ([4af404d](https://github.com/streammachineio/cli/commit/4af404dd1ed88fc4b68e01781c7f14f7b74e9c71))
* **usage:** implement csv, json, json-raw printers for usage ([6e96b2e](https://github.com/streammachineio/cli/commit/6e96b2e00f50dcd7d3f5fc8a395eac2e8df46230))

# [1.5.0](https://github.com/streammachineio/cli/compare/v1.4.0...v1.5.0) (2021-08-18)


### Features

* **event-contract:** create event contract now partially asks for cli flags and a definition ([#40](https://github.com/streammachineio/cli/issues/40)) ([ee06c74](https://github.com/streammachineio/cli/commit/ee06c749be3b272e1796075c2453e1bca1723287)), closes [#33](https://github.com/streammachineio/cli/issues/33)

# [1.4.0](https://github.com/streammachineio/cli/compare/v1.3.0...v1.4.0) (2021-08-12)


### Features

* **sim:** added demo schema random simulator ([3311aec](https://github.com/streammachineio/cli/commit/3311aecbea3197fa91fb14041255e364fc37eb8f))

# [1.3.0](https://github.com/streammachineio/cli/compare/v1.2.1...v1.3.0) (2021-08-03)


### Features

* comments jk ([230d12d](https://github.com/streammachineio/cli/commit/230d12d05b653ba0906ca4d1f4cd2c00a5624cc8))
* **schema-code:** getting schema-code works ([d9cf7b7](https://github.com/streammachineio/cli/commit/d9cf7b7819759ae215b7e2f95c0a6f532510bd53))
* handles creation of schemas and event-contracts ([8c3d267](https://github.com/streammachineio/cli/commit/8c3d26798430a98d79ec425f85af289730e5e119))
* merged create-schemas ([10d2ccc](https://github.com/streammachineio/cli/commit/10d2ccc9a00eb3eff3afcfa20825e7fb7a521de5))
* merged usage branch ([bdb28f3](https://github.com/streammachineio/cli/commit/bdb28f322dbaa1ca221a373e0ce72ca35c6ddc10))
* **schema-code:** interacts ([7aaae9a](https://github.com/streammachineio/cli/commit/7aaae9a540f54c97cdb13b6568a3c86f93e57bc9))
* **schemas:** create added ([3dd41c9](https://github.com/streammachineio/cli/commit/3dd41c907d59066df90e1bc0c392f65ad059d343))
* **usage:** better cli argument parsing ([8773536](https://github.com/streammachineio/cli/commit/8773536d85f7e21e3fbf3cd76650e9b01cd1a5c5))
* wIP creation of schemas  and event-contracts ([afdce5d](https://github.com/streammachineio/cli/commit/afdce5d7e63e0971ad2fa5bf40ebccc1d7532f48))
* **usage:** get usage added ([9a3d6e5](https://github.com/streammachineio/cli/commit/9a3d6e507bf6ac5125936190a4bd8581421d43b9))

## [1.2.1](https://github.com/streammachineio/cli/compare/v1.2.0...v1.2.1) (2021-07-29)


### Bug Fixes

* **kafka-exporter:** fixed can't create bug ([85310b8](https://github.com/streammachineio/cli/commit/85310b8b564cb599697d49091fcb68a0d3c0b2c1))

# [1.2.0](https://github.com/streammachineio/cli/compare/v1.1.3...v1.2.0) (2021-07-19)


### Features

* **schema:** add --kafka-cluster to get schema ([b48c019](https://github.com/streammachineio/cli/commit/b48c01935802b872b3aa6ceb4048d6c719f29740))

## [1.1.3](https://github.com/streammachineio/cli/compare/v1.1.2...v1.1.3) (2021-07-15)


### Bug Fixes

* **egress:** fixed egress usage info ([#24](https://github.com/streammachineio/cli/issues/24)) ([7dbe339](https://github.com/streammachineio/cli/commit/7dbe33968eec094e648c61753ef3d348452b5c7f))

## [1.1.2](https://github.com/streammachineio/cli/compare/v1.1.1...v1.1.2) (2021-07-14)


### Bug Fixes

* set ldflags correctly in order to set values for the version command ([7071aaf](https://github.com/streammachineio/cli/commit/7071aaf686f2ea14fb152938c284a138504c0cc6))

## [1.1.1](https://github.com/streammachineio/cli/compare/v1.1.0...v1.1.1) (2021-07-14)


### Bug Fixes

* fix autocomplete; add logging and logrotate ([ec4a0f7](https://github.com/streammachineio/cli/commit/ec4a0f7872a52985ee5025fb47b745c87458f37d))

# [1.1.0](https://github.com/streammachineio/cli/compare/v1.0.0...v1.1.0) (2021-07-12)


### Bug Fixes

* stream create help text ([164dbb5](https://github.com/streammachineio/cli/commit/164dbb5d157334aef8fb96cd21534ea06c9898f0))
* updated Makefile default; aligned all help texts; misc ([175f23b](https://github.com/streammachineio/cli/commit/175f23be9ebdc4461d2b0589e80decc30de8b836))
* **event-auth-host:** default lacked https prefix ([574aaef](https://github.com/streammachineio/cli/commit/574aaef23a763072f2098fac932d144db821c844))


### Features

* **config-path:** config path used for saving ([9766435](https://github.com/streammachineio/cli/commit/976643566aa21b08283e391f5e16b00c5aec40e0))
* **configpath:** configpath used everywhere ([243af63](https://github.com/streammachineio/cli/commit/243af63e2e822c21a20a25be754f3b0edc09dfbf))
* **egress:** persistent flag ([c7dad70](https://github.com/streammachineio/cli/commit/c7dad70d50e3c95d569a3c8283208397f4334653))

# [1.0.0](https://github.com/streammachineio/cli/compare/v0.4.0...v1.0.0) (2021-07-09)


### Features

* rewrite entire CLI in go ([#15](https://github.com/streammachineio/cli/issues/15)) ([71e6b5a](https://github.com/streammachineio/cli/commit/71e6b5adf7e91391ed89c1bda6fa53750aafb695))


### BREAKING CHANGES

* changed to gRPC API definitions
