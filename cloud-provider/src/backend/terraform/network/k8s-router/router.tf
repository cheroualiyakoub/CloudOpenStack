
resource "openstack_networking_router_v2" "k8s-router" {
  name                = "k8s-router"
  admin_state_up      = true
  external_network_id = "8cfe7f9a-75cc-4b51-833d-77de734678e3"
}
