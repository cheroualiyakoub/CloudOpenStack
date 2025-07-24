terraform {
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.54.0" # or your preferred version
    }
  }
}

provider "openstack" {
  auth_url    = var.os_auth_url
  user_name   = var.os_user_name
  password    = var.os_password
  tenant_name = var.os_tenant_name
  domain_name = var.os_domain_name
  region      = var.os_region
}

variable "os_auth_url" {
  description = "OpenStack authentication URL"
  type        = string
  default     = "http://91.99.215.184/identity/v3"
}

variable "os_user_name" {
  description = "OpenStack username"
  type        = string
  default     = "admin"
}

variable "os_password" {
  description = "OpenStack password"
  type        = string
  default     = "secret"
  sensitive   = true
}

variable "os_tenant_name" {
  description = "OpenStack tenant name"
  type        = string
  default     = "admin"
}

variable "os_domain_name" {
  description = "OpenStack domain name"
  type        = string
  default     = "default"
}

variable "os_region" {
  description = "OpenStack region"
  type        = string
  default     = "RegionOne"
}