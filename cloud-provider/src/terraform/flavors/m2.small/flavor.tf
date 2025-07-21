
resource "openstack_compute_flavor_v2" "m2_small" {
  name  = "m2.small"
  ram   = 4096
  vcpus = 2
  disk  = 20
}
