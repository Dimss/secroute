#!/usr/bin/env bash
oc delete -f 3-service.yaml
oc delete -f 2-create-secure-route.yaml
oc delete -f 1-create-route-for-service.yaml