#!/bin/bash
set -ex

# following derived from the output of:
# gomobile bind -target=ios/amd64,ios/arm64 -v -x

rm -f build/*

GOEXE="go1.15.5"

# CLANGPATH="$(xcodebuild -find-executable clang)"
# CLANGXXPATH="$(xcodebuild -find-executable clang++)"
# SIMULATORFLAGS="-isysroot $(xcodebuild -version -sdk iphonesimulator 13.2 Path) -mios-simulator-version-min=13.2 -arch x86_64"
# DEVICEFLAGS="-isysroot $(xcodebuild -version -sdk iphoneos13.2 Path) -miphoneos-version-min=13.2 -arch arm64"
# MACOSFLAGS="-isysroot $(xcodebuild -version -sdk macosx10.15 Path) -mmacosx-version-min=10.14.6 -arch x86_64"

# x86_64-apple-darwin
if [ "$(uname)" == "Darwin" ]; then
GODEBUG=cgocheck=2 GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 ${GOEXE} build -buildmode=c-archive -o build/librclone-x86_64-apple-darwin.a
fi

# x86_64-unknown-linux-gnu
if [ "$(uname)" == "Linux" ]; then
go version
GODEBUG=cgocheck=2 GOOS=linux GOARCH=amd64 CGO_ENABLED=1 ${GOEXE} build -buildmode=c-archive -o build/librclone-x86_64-unknown-linux-gnu.a
fi

# x86_64-apple-ios
# GODEBUG=cgocheck=2 GOOS=darwin GOARCH=amd64 CC=${CLANGPATH} CXX=${CLANGXXPATH} CGO_CFLAGS=${SIMULATORFLAGS} CGO_CXXFLAGS=${SIMULATORFLAGS} CGO_LDFLAGS=${SIMULATORFLAGS} CGO_ENABLED=1 ${GOEXE} build -tags ios -buildmode=c-archive -o build/rclone-x86_64-apple-ios.a

# aarch64-apple-ios
# GODEBUG=cgocheck=2 GOOS=darwin GOARCH=arm64 CC=${CLANGPATH} CXX=${CLANGXXPATH} CGO_CFLAGS=${DEVICEFLAGS} CGO_CXXFLAGS=${DEVICEFLAGS} CGO_LDFLAGS=${DEVICEFLAGS} CGO_ENABLED=1 ${GOEXE} build -tags ios -buildmode=c-archive -o build/rclone-aarch64-apple-ios.a

# fat library
# lipo -create build/rclone-x86_64-apple-ios.a build/rclone-aarch64-apple-ios.a -o build/rclone-fat-apple-ios.a
# diff build/rclone-x86_64-apple-ios.h build/rclone-aarch64-apple-ios.h
# cp build/rclone-aarch64-apple-ios.h build/rclone-fat-apple-ios.h
