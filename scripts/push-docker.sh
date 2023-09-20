#!/bin/bash
set -e
# set -x

APP_NAME="amnezic-app"

#
# build app
#

#rm -rf bin
#make build

#
# bump
#

TAG="v0.0.2"

#
# build docker
#

IMAGE="${APP_NAME}:${TAG}"
#docker build -t "${IMAGE}" .

#
# push docker
#

# gcloud auth configure-docker europe-west1-docker.pkg.dev

REGION="europe-west1"
LOCATION="europe-west1"
DOCKER_REPO="docker-repo"

print-cmd "gcloud config configurations activate amnezic"
gcloud config configurations activate amnezic

PROJECT_ID=$( gcloud config get-value project )
print-info "project id: ${PROJECT_ID}"

print-cmd "docker tag ${IMAGE} ${LOCATION}-docker.pkg.dev/${PROJECT_ID}/${DOCKER_REPO}/${IMAGE}"
docker tag ${IMAGE} ${LOCATION}-docker.pkg.dev/${PROJECT_ID}/${DOCKER_REPO}/${IMAGE}

print-cmd "docker push ${LOCATION}-docker.pkg.dev/${PROJECT_ID}/${DOCKER_REPO}/${IMAGE}"
docker push ${LOCATION}-docker.pkg.dev/${PROJECT_ID}/${DOCKER_REPO}/${IMAGE}

#print-cmd "gcloud builds submit --region=${REGION} --tag ${REGION}-docker.pkg.dev/${PROJECT_ID}/${DOCKER_REPO}/${APP_NAME}:${TAG}"
#gcloud builds submit --region=${REGION} --tag ${LOCATION}-docker.pkg.dev/${PROJECT_ID}/${DOCKER_REPO}/${APP_NAME}:${TAG}

print-cmd "config configurations activate default"
gcloud config configurations activate default

# gcloud container clusters create-auto amnezic-cluster --region=europe-west1
# gcloud container clusters get-credentials amnezic-cluster --region europe-west1
# kubectl create deployment amnezic-server --image=${LOCATION}-docker.pkg.dev/${PROJECT_ID}/${DOCKER_REPO}/${IMAGE}
# kubectl create deployment amnezic-server --image=europe-west1-docker.pkg.dev/amnezic-app/docker-repo/amnezic-app:v0.0.2
# kubectl expose deployment amnezic-server --type LoadBalancer --port 80 --target-port 8080

