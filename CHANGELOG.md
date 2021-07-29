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
