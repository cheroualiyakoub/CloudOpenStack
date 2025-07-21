
resource "openstack_networking_network_v2" "first_network" {
  name           = "first_network"
  admin_state_up = true
}
