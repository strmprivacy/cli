#!/bin/bash

# This script makes sure the docs can be built successfully after auto-generation.
# We call publish docs, but since the APIS_EMAIL locally is not set, we don't publish
./scripts/publish_docs.sh

cd docs
npm i
npm run build

# If running locally, clean up the docs dir
if [[ $APIS_EMAIL == "" ]]
then
  cd ..
  rm -rf docs
fi
