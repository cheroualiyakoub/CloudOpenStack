
resource "openstack_compute_instance_v2" "controlplane" {
  name            = "control-plane"
  image_name      = "debian-12-generic-amd64.qcow2"
  flavor_name     = "m2.small"
  key_pair        = "first_keypair"
  security_groups = ["k8s_secgroup"]
  network {
    uuid = "d70a2f40-55df-40f8-9861-49a19e7372ee"
  }
}
