#!/bin/bash

NPM_TOKEN=${1}

if [ -z ${NPM_TOKEN} ] ; then
    echo "NPM_TOKEN must be passed as argument"
    exit 2
fi

cat > ~/.npmrc << EOF
//registry.npmjs.org/:_authToken=${NPM_TOKEN}
EOF