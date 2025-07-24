
resource "openstack_compute_instance_v2" "controlplane" {
  name            = "control-plane"
  image_name      = "debian-12"
  flavor_name     = "m2.small"
  key_pair        = "k8s-keypair"
  security_groups = ["k8s-secgroup"]
  network {
    uuid = "dd9135fe-1ed3-4e0c-88a3-a14bb8fd9a42"
  }
}
