
resource "openstack_networking_network_v2" "second_network" {
  name           = "second_network"
  admin_state_up = true
}
