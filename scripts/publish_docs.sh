make
if [[ $APIS_EMAIL != "" ]]
then
  git config --global user.email "${APIS_EMAIL}"
  git config --global user.name "${APIS_USERNAME}"
  tag_name="${GITHUB_REF##*/}"
  LOCAL=0
else
  tag_name="local_test"
  LOCAL=1
fi

rm -rf generated_docs docs
mkdir generated_docs

if [[ $LOCAL == 1 ]]
then
  dstrm --generate-docs
else
  ./dist/strm --generate-docs
fi

cd generated_docs
find . -name "*strm_*" -exec sh -c 'mv "$1" "${1#*strm_}"' _ {} \;
rm -rf dstrm.md
for i in $(find . -name '*.md'); do
  DIRNAME_TO_PARSE="${i%_*}"
  FILENAME="${i##*_}"
  if [[ $DIRNAME_TO_PARSE != *md ]]
  then
    DIRNAME=$(echo "${DIRNAME_TO_PARSE%%.md}" | sed -r 's/_/\//g')
    mkdir -p "$DIRNAME"
    mv "$i" "$DIRNAME/$FILENAME"
  else
    mkdir -p "${FILENAME%.md}"
    mv "$FILENAME" "${FILENAME%.md}/index.md"
    rm -rf "${DIRNAME_TO_PARSE}"
  fi
done

cd ..

if [[ $GITHUB_TOKEN == "" ]]
then
  git clone git@github.com:strmprivacy/docs.git
else
  git clone "https://git:${GITHUB_TOKEN}@github.com/strmprivacy/docs.git"
fi
cd docs
git checkout -b $tag_name
rm -rf ./docs/cli-reference/*
cp -rf ../generated_docs/* ./docs/cli-reference
git add .
git commit -m "add generated docs (${tag_name})"
git push -u origin $tag_name