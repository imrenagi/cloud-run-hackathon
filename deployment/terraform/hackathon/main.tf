module "waterfight_staging_01" {
  source                = "../modules/waterfight"
  name                  = "waterfight-staging-01"
  zone                  = "us-central1-a"
  network               = "default"
  subnet                = "default"
  app_name              = "staging-01"
  route53_hosted_zone   = "imrenagi.com"
}

module "waterfight_production_01" {
  source                = "../modules/waterfight"
  name                  = "waterfight-production-01"
  zone                  = "us-central1-a"
  network               = "default"
  subnet                = "default"
  app_name              = "production-qualification"
  route53_hosted_zone   = "imrenagi.com"
}

module "waterfight_production_02" {
  source                = "../modules/waterfight"
  name                  = "waterfight-production-02"
  zone                  = "us-central1-a"
  network               = "default"
  subnet                = "default"
  app_name              = "production-final"
  route53_hosted_zone   = "imrenagi.com"
}