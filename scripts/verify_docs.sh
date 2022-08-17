#!/bin/bash

# This script makes sure the docs can be built successfully after auto-generation.
# Only for local use.

# We need to generate a dist/strm instead of dist/dstrm
make target=strm
./scripts/generate_docs.sh

rm -rf docs
git clone git@github.com:strmprivacy/docs.git
cd docs
git checkout -b verification
rm -rf ./docs/cli-reference
cp -rf ../generated_docs ./docs/cli-reference

npm i
npm run build

cd ..
rm -rf docs
