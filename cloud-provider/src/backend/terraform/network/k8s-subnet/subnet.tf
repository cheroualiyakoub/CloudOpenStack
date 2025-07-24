
resource "openstack_networking_subnet_v2" "k8s-subnet" {
  name            = "k8s-subnet"
  network_id      = "dd9135fe-1ed3-4e0c-88a3-a14bb8fd9a42"
  cidr            = "192.168.100.0/24"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}
