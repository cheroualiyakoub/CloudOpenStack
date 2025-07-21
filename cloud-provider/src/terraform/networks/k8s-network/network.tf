
resource "openstack_networking_network_v2" "k8s-network" {
  name           = "k8s-network"
  admin_state_up = true
}
