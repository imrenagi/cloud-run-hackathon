
#gcloud run deploy go-bot-tanding-01 \
#  --project=cloudrun-hackathon-359002 \
#  --region=us-central1 \
#  --set-env-vars=WHITELISTED_URLS=https://waterfight-staging-01-x2wnjf2lxq-uc.a.run.app \
#  --set-env-vars=GOOGLE_CLOUD_PROJECT=cloudrun-hackathon-359002 \
#  --allow-unauthenticated --source=.
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

#gcloud compute instances create-with-container waterfight-bot \
#  --container-image us-central1-docker.pkg.dev/cloudrun-hackathon-359002/cloud-run-source-deploy/go-bot-tanding-01:latest \
#  --container-privileged \
#  --container-env=GOOGLE_CLOUD_PROJECT=cloudrun-hackathon-359002 \
#  --container-env=K_SERVICE=waterfight-server \
#  --preemptible \
#  --tags http-server,waterfight \
#  --zone us-central1-a

gcloud run deploy waterfight-staging-04 \
  --project=cloudrun-hackathon-359002 \
  --region=us-central1 \
  --set-env-vars=^##^WHITELISTED_URLS=https://staging-01.waterfight.imrenagi.com,https://production-qualification.waterfight.imrenagi.com,https://production-final.waterfight.imrenagi.com \
  --set-env-vars=GOOGLE_CLOUD_PROJECT=cloudrun-hackathon-359002 \
  --set-env-vars=LOG_LEVEL=DISABLED \
  --set-env-vars=PLAYER_MODE=guard \
  --platform=managed --max-instances=1 --allow-unauthenticated

#
#
#  args: [
#        'run', 'deploy', '--image', 'gcr.io/$PROJECT_ID/waterfight:$COMMIT_SHA', '--region', 'us-central1',
#        '--project=$PROJECT_ID', '--region=us-central1',
#        '--set-env-vars=GOOGLE_CLOUD_PROJECT=$PROJECT_ID',
#        '--set-env-vars=LOG_LEVEL=${_LOG_LEVEL}',
#        '--set-env-vars=PLAYER_MODE=guard',
#        '--platform=managed', '--max-instances=1', '--allow-unauthenticated', 'waterfight-staging-04']
#    env:
#      - '_LOG_LEVEL=${_LOG_LEVEL}'
#images:
#  - gcr.io/$PROJECT_ID/waterfight:$COMMIT_SHA

#        '--set-env-vars=WHITELISTED_URLS="https://staging-01.waterfight.imrenagi.com,https://production-qualification.waterfight.imrenagi.com,https://production-final.waterfight.imrenagi.com"',


#gcr.io/cloudrun-hackathon-359002/waterfight:3e738da0e7cb21cec8e62041d692cb30c092a7bc:
