#!/bin/bash
#
# Copyright 2017 Istio Authors. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
################################################################################

SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

sudo umount `cat /proc/mounts | grep /var/lib/kubelet | awk '{print $2}'` >/dev/null 2>&1
type -a systemctl > /dev/null
if [ $? -ne 0 ]; then
    sleep 15
    sudo umount /var/lib/kubelet
fi
sudo rm -rf /var/lib/kubelet
docker rm -f $(docker ps -aq)
