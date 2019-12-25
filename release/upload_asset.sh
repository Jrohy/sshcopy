#!/bin/bash

GITHUB_TOKEN=""

PROJECT="Jrohy/sshcopy"

SHELL_PATH=$(cd `dirname $0`; pwd)

RELEASE_ID=`curl -H 'Cache-Control: no-cache' -s https://api.github.com/repos/$PROJECT/releases/latest|grep id|awk 'NR==1{print $2}'|sed 's/,//'`

function uploadfile() {
  FILE=$1
  CTYPE=$(file -b --mime-type $FILE)

  curl -H "Authorization: token ${GITHUB_TOKEN}" -H "Content-Type: ${CTYPE}" --data-binary @$FILE "https://uploads.github.com/repos/$PROJECT/releases/${RELEASE_ID}/assets?name=$(basename $FILE)"

  echo ""
}

function upload() {
  FILE=$1
  DGST=$1.dgst
  openssl dgst -md5 $FILE | sed 's/([^)]*)//g' >> $DGST
  openssl dgst -sha1 $FILE | sed 's/([^)]*)//g' >> $DGST
  openssl dgst -sha256 $FILE | sed 's/([^)]*)//g' >> $DGST
  openssl dgst -sha512 $FILE | sed 's/([^)]*)//g' >> $DGST
  uploadfile $FILE
  uploadfile $DGST
}

pushd `pwd` &>/dev/null

go get github.com/mitchellh/gox

gox -os="linux" -output="result/`basename $SHELL_PATH`_{{.OS}}_{{.Arch}}" -ldflags="-s -w" ..

cd result

UPLOAD_ITEM=($(ls -l|awk '{print $9}'|xargs -r))

for ITEM in ${UPLOAD_ITEM[@]}
do
    upload $ITEM
done

echo ""
echo "upload completed!"

popd &>/dev/null

rm -rf result