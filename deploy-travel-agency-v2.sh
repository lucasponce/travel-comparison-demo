#!/bin/bash

oc apply -f <(istioctl kube-inject -f travel_agency_v2.yaml) -n travel-agency
