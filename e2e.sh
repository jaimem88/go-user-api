#!/bin/bash
set -e;

set -ox pipefail; nohup bash -c './bin/$COMPONENT_NAME &'
set +x
echo "** sleeping to let some time for user service to start."
sleep 2
trap 'killall $COMPONENT_NAME' EXIT; \
set -x; SERVICE_ADDR=:8899 go test -tags=endtoend ./test/endtoend -count=1 -timeout 60s

