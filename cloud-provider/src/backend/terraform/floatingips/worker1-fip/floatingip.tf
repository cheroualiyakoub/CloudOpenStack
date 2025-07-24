
resource "openstack_networking_floatingip_v2" "worker1-fip" {
  pool = "public"
}

resource "openstack_networking_floatingip_associate_v2" "worker1-fip_assoc" {
  floating_ip = openstack_networking_floatingip_v2.worker1-fip.address
  port_id = "dcedd4b9-e534-4099-871a-f9e5315506c8"
}
