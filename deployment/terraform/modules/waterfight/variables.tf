variable "name" {
  type = string
}

variable "zone" {
  type = string
}

variable "network" {
  type = string
}

variable "subnet" {
  type = string
}

variable "route53_hosted_zone" {
  type = string
}

variable "app_name" {
  type = string
}

variable "spot" {
  type    = bool
  default = false
}