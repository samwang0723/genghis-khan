#!/usr/bin/env bash
kernel_name=$(uname -s)
kubectl cluster-info > /dev/null 2>&1
if [ $? -eq 1 ]
then
  echo "kubectl was unable to reach your Kubernetes cluster. Make sure that" \
       "you have selected one using the 'gcloud container' commands."
  exit 1
fi

# Clear out any existing configmap. Fail silently if there are none to delete.
kubectl delete namspace genghis-khan 2> /dev/null

echo "Create genghis-khan namespace..."
kubectl create -f genghis-khan-namespace.yaml 2> /dev/null

kubectl create -f redis-deployment.yaml
kubectl create -f redis-service.yaml
echo
echo "===== Redis installed ============================================"
echo "Your cluster is now downloading the Docker image for Redis."
echo "You can check the progress of this by typing 'kubectl get pods' in another"
echo "tab. Please check if you see 1/1 READY for your redis-* pod."
echo
read -p "<Press enter once you've verified that your Redis is up>"
echo
echo "===== genghis-khan Server installation =========================================="

kubectl create -f genghis-khan-claim0-persistentvolumeclaim.yaml
kubectl create -f genghis-khan-deployment.yaml
kubectl create -f genghis-khan-service.yaml

echo "Your cluster is now downloading the Docker image for genghis-khan Server."
echo "You can check the progress of this by typing 'kubectl get pods'"
echo "Once you see 1/1 READY for your genghis-khan-* pod, your Agent is ready"
echo "to start pulling and running builds."
echo
read -p "<Press enter once you've verified that your genghis-khan Server is up>"
echo
echo "Installation Completed"

#kubectl delete -f genghis-khan-namespace.yaml