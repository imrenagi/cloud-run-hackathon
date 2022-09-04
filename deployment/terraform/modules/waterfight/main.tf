data "google_compute_network" "default" {
  name = "default"
}

data "google_compute_subnetwork" "default_us_central_1" {
  name          = "default-us-central1"
  region        = "us-central1"  
}

resource "google_compute_address" "vm_ip_address" {  
  name         = "${var.app_name}-waterfight-address"
  address_type = "EXTERNAL"
  region        = "us-central1"
}

data "aws_route53_zone" "imrenagi_com_zone" {
  name         = var.route53_hosted_zone
}

resource "aws_route53_record" "eatn_route53_record" {
  zone_id = data.aws_route53_zone.imrenagi_com_zone.zone_id
  name    = "${var.app_name}.waterfight.${var.route53_hosted_zone}"
  type    = "A"
  ttl     = "300"
  records = [google_compute_address.vm_ip_address.address]
}

resource "google_compute_instance" "waterfight_vm" {
  name         = "${var.app_name}-waterfight"
  machine_type = "n1-standard-1"
  zone         = var.zone

  scheduling {
    provisioning_model        = var.spot ? "SPOT" : "STANDARD"
    preemptible               = var.spot ? true : false
    on_host_maintenance       = var.spot ? "TERMINATE" : "MIGRATE"
    automatic_restart         = var.spot ? false : true
  }
  
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2004-lts"
      size  = 10
      type  = "pd-standard"
    }
  }

  tags = ["waterfight", "http-server", "https-server"]

  network_interface {
    network     = data.google_compute_network.default.id
    subnetwork  = data.google_compute_subnetwork.default_us_central_1.id
    access_config {
      nat_ip = google_compute_address.vm_ip_address.address
    }
  }

  service_account {
    scopes = [
      "userinfo-email",
      "cloud-platform"
    ]    
  }
}


