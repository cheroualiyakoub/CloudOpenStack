
resource "openstack_networking_subnet_v2" "k8s-subnet" {
  name            = "k8s-subnet"
  network_id      = "d70a2f40-55df-40f8-9861-49a19e7372ee"
  cidr            = "192.168.100.0/24"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}
