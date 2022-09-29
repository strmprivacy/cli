#!/bin/bash
./scripts/generate_docs.sh

if [[ $APIS_EMAIL != "" ]]
then
  echo "Setting git config"
  git config --global user.email "${APIS_EMAIL}"
  git config --global user.name "${APIS_USERNAME}"
  tag_name="${GITHUB_REF##*/}"
else
  tag_name="local"
fi

echo "Tag = ${tag_name}"

function copy_cli_reference() {
  CLI_REF_PATH="./docs/04-reference/01-cli-reference"

  rm -rf "$CLI_REF_PATH/strm"
  cp -rf ../generated_docs/strm "$CLI_REF_PATH/strm"
}

if [[ $GITHUB_TOKEN == "" ]]
then
  git clone git@github.com:strmprivacy/docs.git
  cd docs
  git checkout -b $tag_name
  copy_cli_reference
else
  git clone "https://git:${GITHUB_TOKEN}@github.com/strmprivacy/docs.git"
  ls -al
  ls -al generated_docs
  cd docs
  ls -al
  copy_cli_reference
fi

