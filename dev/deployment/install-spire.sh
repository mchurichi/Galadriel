#!/bin/bash

set -e

SPIRE_VERSION="1.3.3"

wget https://github.com/spiffe/spire/releases/download/v1.3.3/spire-${SPIRE_VERSION}-linux-x86_64-glibc.tar.gz
tar xvf spire-${SPIRE_VERSION}-linux-x86_64-glibc.tar.gz
rm spire-${SPIRE_VERSION}-linux-x86_64-glibc.tar.gz
mv spire-${SPIRE_VERSION} spire