#!/bin/bash

oc delete services --all -n travel-agency
oc delete deployments --all -n travel-agency
oc delete virtualservices --all -n travel-agency
oc delete gateways --all -n travel-agency
