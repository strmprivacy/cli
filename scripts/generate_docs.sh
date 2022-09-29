#!/bin/bash

rm -rf generated_docs docs
mkdir generated_docs

if [[ $APIS_EMAIL == "" ]]
then
  # Local
  # We need to generate a dist/strm instead of dist/dstrm
  make target=strm
  echo "Generating docs locally"
  ./dist/strm --generate-docs > /dev/null 2>&1
else
  echo "Generating docs in CI"
  ./dist/strm_linux_amd64_v1/strm --generate-docs > /dev/null 2>&1
fi

echo "Rearranging docs"
./scripts/rearrange_docs.py
