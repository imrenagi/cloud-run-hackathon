output "instance_ip_addr" {
  value = google_compute_address.vm_ip_address.address
}