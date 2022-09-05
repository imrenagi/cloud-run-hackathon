output "instance_ip_addr" {
  value = google_compute_address.vm_ip_address.address
}

output "instance_url" {
  value = aws_route53_record.eatn_route53_record.name
}