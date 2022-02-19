find . -name "*strm_*" -exec sh -c 'mv "$1" "${1#*_}"' _ {} \;
for i in `find . -name '*.md'`; do
   DIRNAME="${i%_*}"
   DIRNAME2="${DIRNAME%%/*}"
   FILENAME="${i##*_}"

  if [[ $DIRNAME != *md ]]
  then
    echo ${DIRNAME%%.md} | sed -r 's/_/\//g' | xargs -I {} mkdir -p {}
    echo ${DIRNAME%%.md} | sed -r 's/_/\//g' | xargs -I {} mv $i {}/$FILENAME
  else
    rm -rf ${DIRNAME}
  fi
done