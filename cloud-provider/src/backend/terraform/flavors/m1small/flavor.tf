
resource "openstack_compute_flavor_v2" "m1small" {
  name  = "m1small"
  ram   = 4096
  vcpus = 2
  disk  = 20
}
