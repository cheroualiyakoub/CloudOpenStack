
resource "openstack_networking_router_v2" "my_router" {
  name                = "my_router"
  admin_state_up      = true
  external_network_id = "3db4d4b4-ce8d-437e-8db1-3fc1cde7d0aa"
}
