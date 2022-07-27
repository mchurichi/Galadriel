#!/bin/bash

set -e

INSTANCE_NO=$1

spire/spire/bin/spire-server run -config spire/server${INSTANCE_NO}.conf