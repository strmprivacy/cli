#!/bin/bash

# This script makes sure the docs can be built successfully after auto-generation.
# Only for local use.

# We need to generate a dist/strm instead of dist/dstrm
make target=strm
./scripts/generate_docs.sh

CLI_REF_PATH="./docs/04-reference/01-cli-reference"

rm -rf docs
git clone git@github.com:strmprivacy/docs.git
cd docs
rm -rf "$CLI_REF_PATH/strm"
cp -rf ../generated_docs/strm "$CLI_REF_PATH/strm"

npm i
npm run build

cd ..
rm -rf docs
