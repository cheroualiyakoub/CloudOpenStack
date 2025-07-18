
resource "openstack_compute_flavor_v2" "m1_small" {
  name  = "m1_small"
  ram   = 4096
  vcpus = 2
  disk  = 20
}
