#!/bin/bash

# runs the integration tests

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"
ginkgo -tags integration -p