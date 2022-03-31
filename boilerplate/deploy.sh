###############################################
###  WebChest - Deployment: Single Service  ###
###############################################

# Log service deployment
echo 'Running Main Deployment...'

# Configure GCloud Project
echo 'Configing gcloud project...'
gcloud config set project webchest

# Build GCloud Docker Image
echo 'Build new gcloud image...'
gcloud builds submit --tag gcr.io/webchest/webchest-main

# Deploy new Docker Image to Cloud Run
echo 'Deploying to gcloud run...'
gcloud run deploy webchest-main --image gcr.io/webchest/webchest-main --platform managed --region us-east1

