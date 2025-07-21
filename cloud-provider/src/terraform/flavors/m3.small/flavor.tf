
resource "openstack_compute_flavor_v2" "m3_small" {
  name  = "m3.small"
  ram   = 4096
  vcpus = 2
  disk  = 20
}
