#!/bin/bash
./scripts/generate_docs.sh

if [[ $APIS_EMAIL != "" ]]
then
  git config --global user.email "${APIS_EMAIL}"
  git config --global user.name "${APIS_USERNAME}"
  tag_name="${GITHUB_REF##*/}"
else
  tag_name="local"
fi

if [[ $GITHUB_TOKEN == "" ]]
then
  git add -A
  git commit -m "update generated CLI reference docs (CLI version: ${tag_name})"
  git push -f origin $tag_name
else
  git add -A
  git commit -m "update generated CLI reference docs (CLI version: ${tag_name})"
  git push
fi

