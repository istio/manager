# Copyright 2017 Istio Authors

#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

#!/bin/bash

# Bring up control plane
deploy_istio() {
    echo "Deploying ISTIO"
    $K8CLI create -f istio/controlplane.yaml
    for (( i=0; i<=4; i++ ))
    do
        ready=$($K8CLI get pods | awk 'NR>1 {print $1 "\t" $2}' | grep "istio" | grep "1/1" | wc -l)
        if [ $ready -eq 2 ]
        then
            echo "ISTIO control plane deployed"
            return 0
        fi
        sleep 10
    done
    echo "Unable to deploy ISTIO"
    return 1
}

# Clean up the controlplane when we are done with it
cleanup_istio(){
    echo "Cleaning up ISTIO"
    $K8CLI delete -f istio/controlplane.yaml
}

# Deploy the bookinfo microservices
deploy_bookinfo(){
    echo "Deploying BookInfo to kube"
    $K8CLI create -f apps/bookinfo/bookinfo.yaml
    for (( i=0; i<=4; i++ ))
    do
        ready=$($K8CLI get pods | awk 'NR>1 {print $1 "\t" $2}' | grep -v "istio" | grep "2/2" | wc -l)
        if [ $ready -eq 6 ]
        then
            echo "BookInfo deployed"
            return 0
        fi
        sleep 10
    done
    echo "Unable to deploy BookInfo"
    return 1
}

# Clean up bookinfo when we are done with it
cleanup_bookinfo(){
    echo "Cleaning up BookInfo"
    $K8CLI delete -f apps/bookinfo/bookinfo.yaml
}

# Debug dump for failures
dump_debug() {
    # kubectl get pods
    # kubectl get thirdpartyresources
    # kubectl get thirdpartyresources -o json
    # GATEWAY_PODNAME=$(kubectl get pods | grep istio-ingress | awk '{print $1}')
    # kubectl logs $GATEWAY_PODNAME
    # PRODUCTPAGE_PODNAME=$(kubectl get pods | grep productpage | awk '{print $1}')
    # kubectl logs $PRODUCTPAGE_PODNAME -c productpage
    # kubectl logs $PRODUCTPAGE_PODNAME -c proxy
    # RULES_PODNAME=$(kubectl get pods | grep rules | awk '{print $1}')
    # kubectl logs $RULES_PODNAME
    echo
}
