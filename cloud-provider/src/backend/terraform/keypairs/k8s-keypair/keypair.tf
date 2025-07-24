
resource "openstack_compute_keypair_v2" "k8s-keypair" {
  name       = "k8s-keypair"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIVqGaHXTF9WUerS9XBsYGY/ouLbA/2Qbqbc7aa58NJJ cherouali@gmail.com"
}
