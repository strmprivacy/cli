#!/bin/bash
if [[ "$OSTYPE" == "darwin"* ]]; then
    SED="gsed"
else
    SED="sed"
fi

make
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
  dstrm --generate-docs > /dev/null 2>&1
else
  ./dist/dstrm --generate-docs > /dev/null 2>&1
fi

cd generated_docs
#find . -name "*strm_*" -exec sh -c 'mv "$1" "${1#*strm_}"' _ {} \;
for i in $(find . -name '*.md'); do
  DIRNAME_TO_PARSE="${i%_*}"
  FILENAME="${i##*_}"
  if [[ $DIRNAME_TO_PARSE != *md ]]
  then
    DIRNAME=$(echo "${DIRNAME_TO_PARSE%%.md}" | sed -r 's/_/\//g')
    mkdir -p "$DIRNAME"
    mv "$i" "$DIRNAME/$FILENAME"
  fi
done

for i in $(find . -name '*.md'); do
  FILENAME="${i##*_}"
  if [[ -d "${FILENAME%.md}" ]]
  then
    X1=${FILENAME:2}
    X1=${X1::${#X1}-3}
    X2=${FILENAME%.md}
    X2=${X2:2}
    "$SED" -i "s/cli-reference\/${X1//\//\\/}.md/cli-reference\/${X2//\//\\/}\/index.md/g" {,**/}*/*.md
    mv "$FILENAME" "${FILENAME%.md}/index.md"
  fi
done

