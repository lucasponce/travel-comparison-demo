#!/usr/bin/env bash

export INGRESS_HOST=$(oc get po -l istio=ingressgateway -n istio-system -o 'jsonpath={.items[0].status.hostIP}')
export INGRESS_PORT=$(oc -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')

TRAVEL_AGENCY_HOST="Host:www.travel-agency.com"
TRAVELS_SERVICE="http://${INGRESS_HOST}:${INGRESS_PORT}"

while true
do
    for city in paris rome london
    do
        for user in vip normal
        do
            RESULT=$(curl -H "${TRAVEL_AGENCY_HOST}" -H "user:${user}" -o /dev/null -s -w %{http_code} ${TRAVELS_SERVICE}/travels/${city})
            echo "Travel quota for ${city} from ${user} => $RESULT"
            sleep 1
        done
        RESULT=$(curl -H "${TRAVEL_AGENCY_HOST}" -H "user:${user}" -o /dev/null -s -w %{http_code} ${TRAVELS_SERVICE}/travels/${city})
        echo "Travel quota for ${city} from anonymous => $RESULT"
        sleep 1
    done
done