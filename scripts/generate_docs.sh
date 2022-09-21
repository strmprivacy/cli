#!/bin/bash

if [[ $APIS_EMAIL != "" ]]
then
  git config --global user.email "${APIS_EMAIL}"
  git config --global user.name "${APIS_USERNAME}"
  tag_name="${GITHUB_REF##*/}"
else
  tag_name="local_test"
fi

rm -rf generated_docs docs
mkdir generated_docs

if [[ $APIS_EMAIL == "" ]]
then
  make
  ./dist/strm --generate-docs > /dev/null 2>&1
else
  ./dist/strm_linux_amd64_v1/strm --generate-docs > /dev/null 2>&1
fi
./scripts/generate_docs.py
