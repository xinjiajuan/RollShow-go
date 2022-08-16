#!/usr/bin/env sh
mkdir public
#交叉编译
# amd64 darwin
  CGO_ENABLED=1
  GOOS=darwin
  GOARCH=amd64
  go build
  echo "darwin amd64 ENV Build Finish"
  chmod +x rollshow
  version=`./rollshow -v`
  ver=`echo ${version:0-5}`
  mv rollshow public/rollshow_"$ver"_darwin_amd64
# amd64 linux
  CGO_ENABLED=0
  GOOS=linux
  GOARCH=amd64
  go build
  echo "linux amd64 ENV Build Finish"
  chmod +x rollshow
  version=`./rollshow -v`
  ver=`echo ${version:0-5}`
  mv rollshow public/rollshow_"$ver"_linux_amd64
# amd64 windows
  CGO_ENABLED=0
  GOOS=windows
  GOARCH=amd64
  go build
  echo "windows amd64 ENV Build Finish"
  chmod +x rollshow
  version=`./rollshow -v`
  ver=`echo ${version:0-5}`
  mv rollshow public/rollshow_"$ver"_windows_amd64.exe