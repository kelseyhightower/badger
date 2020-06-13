#!/bin/bash

set -x

declare project=$GOOGLE_CLOUD_PROJECT
declare service=$K_SERVICE
declare region=$GOOGLE_CLOUD_REGION

declare sa=badger@$project.iam.gserviceaccount.com

gcloud iam service-accounts describe $sa --project $project &> /dev/null

if [ $? -ne 0 ]; then
  echo "creating badger service account: $sa"
  gcloud iam service-accounts create badger \
    --display-name="badger" \
    --project=$project
    
  echo "waiting 30 seconds for service account consistency"
  sleep 30
fi

echo "allowing $sa to view cloud builds"
gcloud projects add-iam-policy-binding $project \
  --member=serviceAccount:$sa \
  ---role=roles/cloudbuild.builds.viewer &> /dev/null

echo "updating $service to use the service account $sa"
gcloud run services update $service \
  --platform=managed --project=$project --region=$region \
  --service-account=$sa &> /dev/null
