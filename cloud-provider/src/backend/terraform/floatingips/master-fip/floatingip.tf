
resource "openstack_networking_floatingip_v2" "master-fip" {
  pool = "public"
}

resource "openstack_networking_floatingip_associate_v2" "master-fip_assoc" {
  floating_ip = openstack_networking_floatingip_v2.master-fip.address
  port_id = "e005e8a7-a28f-43e3-9145-9497d1ee23ef"
}
