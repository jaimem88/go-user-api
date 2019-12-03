#!/bin/bash
set -e;

# go test -c -o ./bin/$COMPONENT_NAME-with-cov  -cover -covermode=atomic -coverpkg=./... -count=1 ./cmd/$COMPONENT_NAME
set -ox pipefail; nohup bash -c './bin/$COMPONENT_NAME-with-cov -test.coverprofile=e2e-cov.out DEVEL --debug &' \
set +x
echo "** sleeping to let some time for user service to start."
sleep 1
# trap 'killall $COMPONENT_NAME-with-cov' EXIT
SERVICE_ADDR=:8899 go test -tags=endtoend ./test/endtoend -count=1 -timeout 60s
killall $COMPONENT_NAME-with-cov
echo "** sleeping to let some time for cover profile to flush"
sleep 1
