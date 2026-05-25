#!/usr/bin/env bash

set -euo pipefail

# ENV can be local, prod
argument_env="${1:-}"

if [ "$argument_env" == "prod" ]; then
  # manual steps to setup prod environments
  # ================================================================
  ## create cluster in prod
  # gcloud container clusters create notifier-gke --region asia-south1 --num-nodes=1 --machine-type=e2-small --disk-size=20

  ## get credentials for the cluster
  # gcloud container clusters get-credentials notifier-gke --region=asia-south1

  ## set secrets for the cluster in secrets manager
  # gcloud secrets create notifier-secrets --replication-policy="automatic"
  # ================================================================
  kubectl apply -f deploy/base/namespace.yaml

  kubectl create secret generic notifier-secrets --from-env-file=.env -n notifier

  # deploy the application
  kubectl apply -k deploy/overlays/gcp
  kubectl get svc api -n notifier -w
  kubectl get svc workers -n notifier -w

elif [ "$argument_env" == "local" ]; then
  # create cluster in local
  if ! kind get clusters | grep -q notifier; then
    kind create cluster --name notifier --config deploy/kind-config.yaml
  fi

  # set namespace
  kubectl apply -f ./deploy/base/namespace.yaml

  # set secrets for the cluster in secrets manager
  kubectl apply -f ./deploy/base/secrets.yaml

  # deploy the application
  kubectl apply -k deploy/overlays/local
else
  echo "ERROR: No environment provided or invalid environment: $argument_env"
  exit 1
fi


# get pods
# kubectl get pods -n notifier

# get services
# kubectl get services -n notifier

# get deployments
# kubectl get deployments -n notifier

# get secrets
# kubectl get secrets -n notifier

# get pods live
# kubectl get pods -n notifier -w

# get logs of a pod
# kubectl logs <pod-name> -n notifier

# get logs of a pod in real time
# kubectl logs -f <pod-name> -n notifier

# get migration job
# kubectl get job migrate -n notifier

# get migration job logs
# kubectl logs job/migrate -n notifier

# get migration job logs in real time
# kubectl logs -f job/migrate -n notifier

# delete migration job
# kubectl delete job migrate -n notifier

# restart a pod
# kubectl restart <pod-name> -n notifier

# restart a deployment
# kubectl rollout restart deployment <deployment-name> -n notifier

# rollout undo a deployment
# kubectl rollout undo deployment <deployment-name> -n notifier

# rollout undo a deployment to a specific revision
# kubectl rollout undo deployment <deployment-name> --to-revision=<revision> -n notifier

# delete secrets
#  kubectl delete secret notifier-secrets -n notifier
