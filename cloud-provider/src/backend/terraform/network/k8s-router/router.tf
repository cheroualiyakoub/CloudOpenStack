
resource "openstack_networking_router_v2" "k8s-router" {
  name                = "k8s-router"
  admin_state_up      = true
  external_network_id = "5e29e491-c02f-43b7-a69d-ac7540c80377"
}
