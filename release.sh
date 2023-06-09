# !/usr/bin/env sh

ESCAPE=$'\e'
export VERSION="1.23"

./build.sh && \

echo "${ESCAPE}[0;32mReleasing v${VERSION}...${ESCAPE}[0m" && \

docker build -t antfie/scan_compare:$VERSION . && \
docker build -t antfie/scan_compare . && \
docker push antfie/scan_compare:$VERSION && \
docker push antfie/scan_compare && \

echo "${ESCAPE}[0;32mRelease Success${ESCAPE}[0m"