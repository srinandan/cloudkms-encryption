#!/bin/bash
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


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