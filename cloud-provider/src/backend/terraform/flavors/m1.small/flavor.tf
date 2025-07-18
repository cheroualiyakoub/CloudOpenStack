
resource "openstack_compute_flavor_v2" "m1.small" {
  name  = "m1.small"
  ram   = 4096
  vcpus = 2
  disk  = 20
}
