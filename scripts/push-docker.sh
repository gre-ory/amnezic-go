#!/bin/bash
set -e
# set -x

REGION="europe-west1"
APP_NAME="amnezic-app"
DOCKER_REPO="docker-repo"

print-cmd "gcloud config configurations activate amnezic"
gcloud config configurations activate amnezic

PROJECT_ID=$( gcloud config get-value project )
print-info "project id: ${PROJECT_ID}"

print-cmd "gcloud builds submit --region=${REGION} --tag ${REGION}-docker.pkg.dev/${PROJECT_ID}/${DOCKER_REPO}/${APP_NAME}:tag1"
gcloud builds submit --region=${REGION} --tag ${REGION}-docker.pkg.dev/${PROJECT_ID}/${DOCKER_REPO}/${APP_NAME}:tag1

print-cmd "config configurations activate default"
gcloud config configurations activate default