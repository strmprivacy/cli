#!/bin/bash

set -e

# This script makes sure the docs can be built successfully after auto-generation.
./scripts/copy_generated_to_docs.sh

cd docs
npm i
npm run build

# If running locally, clean up the docs dir
#if [[ $APIS_EMAIL == "" ]]
#then
#  cd ..
#  rm -rf docs
#fi
