
resource "openstack_networking_subnet_v2" "first_sub_net" {
  name            = "first_sub_net"
  network_id      = "049d79ab-5fb1-4ebf-839c-a6c0f09ad99f"
  cidr            = "192.168.100.0/24"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}
