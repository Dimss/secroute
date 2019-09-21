## SECRouter - Admission Controller Webhook Server for Securing OCP Routes and Services

### Use case 1 - Simple Validation Webhook
```bash
# Deploy Webhook Configuration 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/validate-route-webhook.yaml
# Create insecure route 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/routes/insecure-route.yaml
# Create secure route
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/routes/secure-route.yaml
 
# Cleanup
# Delete route
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/routes/secure-route.yaml
# Delete webhook configuration
oc delete -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/validate-route-webhook.yaml
```

### Use case 2 - Mutation Webhook

```bash
# Deploy Webhook Configuration 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/webhooks/mutate-route-webhook.yaml
# Create insecure route 
oc create -f https://raw.githubusercontent.com/Dimss/secroute/master/ocp/routes/insecure-route.yaml
```
