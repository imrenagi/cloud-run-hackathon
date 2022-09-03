module "merchant_vm_01" {
  source                = "../modules/waterfight"
  name                  = "merchant-service-01"
  zone                  = "us-central1-a"
  network               = "default"
  subnet                = "default"
  app_name              = "staging"
  route53_hosted_zone   = "imrenagi.com"
}