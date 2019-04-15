#!/bin/bash

oc new-project travel-portal
oc project travel-portal
oc adm policy add-scc-to-user privileged -z default

oc apply -f <(istioctl kube-inject -f travel_portal.yaml) -n travel-portal