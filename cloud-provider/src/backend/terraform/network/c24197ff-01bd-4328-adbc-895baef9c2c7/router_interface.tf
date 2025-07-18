
resource "openstack_networking_router_interface_v2" "router_interface" {
  router_id = "c24197ff-01bd-4328-adbc-895baef9c2c7"
  subnet_id = "805d80a7-ba9b-40d1-8250-5f2ea3fb23ce"
}
