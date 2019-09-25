## SECRouter - Admission Controller Webhook Server for Securing OCP Routes and Services

## Development environment setup

### Write your POST handlers 


### Use case 1 - Simple Validation 
```bash
# Deploy Webhook Configuration 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/block-insecure-route.yaml
# Create insecure route 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/routes/insecure-route.yaml
# Create secure route
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/routes/secure-route.yaml
 
# Cleanup route
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/routes/secure-route.yaml
# Cleanup webhook configuration
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/block-insecure-route.yaml
```

### Use case 2 - Simple Mutation 
```bash
# Deploy Webhook Configuration 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/create-secure-route.yaml
# Create insecure route 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/routes/insecure-route.yaml
# Check created route
oc get route -o yaml demo-insecure-route 
# Cleanup route 
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/routes/insecure-route.yaml
# Cleanup webhook configuration
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/create-secure-route.yaml

```


### Use case 3 - Mutation Side effects 
```bash
# Deploy Webhooks configuration
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/create-route-for-service.yaml
# Create service 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/services/service.yaml

# Cleanup service
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/services/service.yaml
# Cleanup webhook
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/create-route-for-service.yaml
```

### Use case 4 - Webhooks chain  
```bash
# Deploy block insecure route Webhook Configuration 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/block-insecure-route.yaml
# Deploy create-route-for-service Webhooks configuration
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/create-route-for-service.yaml
# Create service 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/services/service.yaml
# Create mutate admission webhook 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/create-secure-route.yaml
# Create service again
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/services/service.yaml

# Cleanup service
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/services/service.yaml
# Cleanup webhook for routes 
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/create-route-for-service.yaml
# Deploy block insecure route Webhook Configuration 
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/block-insecure-route.yaml
```