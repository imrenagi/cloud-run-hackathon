terraform {
  backend "gcs" {
    bucket  = "cloudrun-hackathon-terraform"
    prefix  = "waterfight"
  }
}
