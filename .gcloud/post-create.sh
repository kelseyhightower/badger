#!/bin/bash

set -x

declare project=$GOOGLE_CLOUD_PROJECT
declare service=$K_SERVICE
declare region=$GOOGLE_CLOUD_REGION

# service accounts once deleted should not be reused so we make it unique(ish)
declare name=badger-$(cat /dev/urandom|tr -dc '0-9'|fold -w 4|head -n 1)
declare sa=$name@$project.iam.gserviceaccount.com

echo "creating badger service account: $sa"
gcloud iam service-accounts create $name \
  --display-name="badger" \
  --project=$project

echo "allowing $sa to view cloud builds"
gcloud projects add-iam-policy-binding $project \
  --member=serviceAccount:$sa \
  --role=roles/cloudbuild.builds.viewer &> /dev/null

echo "updating $service to use the service account $sa"
gcloud run services update $service \
  --platform=managed --project=$project --region=$region \
  --service-account=$sa &> /dev/null
