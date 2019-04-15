#!/usr/bin/env bash

oc apply -f <(istioctl kube-inject -f travels-chaosmonkey.yaml) -n travel-agency
