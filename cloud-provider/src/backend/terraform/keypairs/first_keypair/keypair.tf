
resource "openstack_compute_keypair_v2" "first_keypair" {
  name       = "first_keypair"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAlWL3kFsTJW7aNI48g7/E7Gb4x8+nX+zh9L9Ohy7PKi"
}
