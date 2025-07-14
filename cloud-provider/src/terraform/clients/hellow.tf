
resource "openstack_identity_project_v3" "hellow" {
  name        = "hellow"
  description = "just testing"
  domain_id   = "default"
}
