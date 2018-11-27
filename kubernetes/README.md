### Create Namespace

    kubectl create namespace genghis-khan


### Get Certification & Token
    
    kubectl -n genghis-khan get secrets
    

    
    kubectl get secret -n genghis-khan default-token-hjv2z -o yaml | egrep 'ca.crt:|token:'
    

    Decode token with base64:

    
    echo "<encoded token>" | base64 -d
    

### Create secret under namespace

    
    kubectl --namespace=genghis-khan create secret docker-registry regcred --docker-server=registry.hub.docker.com --docker-username={ACCOUNT} --docker-password={PASSWORD} --docker-email={EMAIL}
    

### Assign roles to service account for drone doing the deployment

    
    kubectl create -f kubernetes/drone-deployment.yml
    
