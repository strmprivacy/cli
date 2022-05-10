#!/bin/bash
./scripts/generate_docs.sh

if [[ $APIS_EMAIL != "" ]]
then
  git config --global user.email "${APIS_EMAIL}"
  git config --global user.name "${APIS_USERNAME}"
  tag_name="${GITHUB_REF##*/}"
else
  tag_name="local_test"
fi

if [[ $GITHUB_TOKEN == "" ]]
then
  git clone git@github.com:strmprivacy/docs.git
  cd docs
  git checkout -b $tag_name
  rm -rf ./docs/cli-reference
  cp -rf ../generated_docs ./docs/cli-reference
  git add .
  git commit -m "add generated docs (cli branch: ${tag_name})"
  git push -f origin $tag_name
else
  git clone "https://git:${GITHUB_TOKEN}@github.com/strmprivacy/docs.git"
  cd docs
  rm -rf ./docs/cli-reference
  cp -rf ../generated_docs ./docs/cli-reference
  git add .
  git commit -m "add generated docs (cli branch: ${tag_name})"
  git push
fi
