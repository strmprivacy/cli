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

cd docs

if [[ $GITHUB_TOKEN == "" ]]
then
  git add -A
  git commit -m "update generated CLI reference docs (CLI version: ${tag_name})" --allow-empty
  git push -f origin $tag_name
else
  git add -A
  git commit -m "update generated CLI reference docs (CLI version: ${tag_name})" --allow-empty
  git push
fi

