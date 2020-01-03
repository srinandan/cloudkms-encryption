#!/bin/bash

./check-prereqs.sh $1 $2 $3
RESULT=$?
if [ $RESULT -ne 0 ]; then
  exit 1
fi

# zip the bundles
zip -r cloudkms-sample.zip apiroxy

# import the sample
apigeecli apis create -o $1 -n CloudKMS_Demo -p cloudkms-sample.zip -a $3

# deploy the sample
apigeecli apis deploy -o $1 -e $2 -n CloudKMS_Demo -v 1