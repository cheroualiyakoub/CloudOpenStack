
resource "openstack_images_image_v2" "debian12genericamd64_qcow2" {
  name             = "debian-12-generic-amd64.qcow2"
  image_source_url       = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-generic-amd64.qcow2"
  disk_format      = "qcow2"
  container_format = "bare"
}
