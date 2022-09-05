output "staging_ip_addr" {
  value = "${module.waterfight_staging_01.instance_ip_addr}"
}

output "staging_url" {
  value = "${module.waterfight_staging_01.instance_url}"
}

