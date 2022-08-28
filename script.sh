
gcloud run deploy go-bot-tanding-01 \
  --project=cloudrun-hackathon-359002 \
  --region=us-central1 \
  --set-env-vars=WHITELISTED_URLS=https://waterfight-staging-02-x2wnjf2lxq-uc.a.run.app \
  --set-env-vars=GOOGLE_CLOUD_PROJECT=cloudrun-hackathon-359002 \
  --allow-unauthenticated --source=.
#
#gcloud run deploy go-bot-tanding-06 \
#  --project=cloudrun-hackathon-359002 \
#  --region=us-central1 \
#  --allow-unauthenticated --source=.
#
#gcloud run deploy go-bot-tanding-07 \
#  --project=cloudrun-hackathon-359002 \
#  --region=us-central1 \
#  --allow-unauthenticated --source=.


