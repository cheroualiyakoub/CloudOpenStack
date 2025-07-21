
resource "openstack_networking_network_v2" "first network" {
  name           = "first network"
  admin_state_up = true
}
