#!/bin/bash
set -e

DOCKER_VERSION=dev
DOCKER_TRAVEL_AGENCY=lucasponce/travel_agency
DOCKER_TRAVEL_AGENCY_TAG=${DOCKER_TRAVEL_AGENCY}:${DOCKER_VERSION}

go build -o docker/travel_agency/travel_agency travel_agency/travel_agency.go

docker build -t ${DOCKER_TRAVEL_AGENCY_TAG} docker/travel_agency

DOCKER_TRAVEL_PORTAL=lucasponce/travel_portal
DOCKER_TRAVEL_PORTAL_TAG=${DOCKER_TRAVEL_PORTAL}:${DOCKER_VERSION}

go build -o docker/travel_portal/travel_portal travel_portal/travel_portal.go

docker build -t ${DOCKER_TRAVEL_PORTAL_TAG} docker/travel_portal

