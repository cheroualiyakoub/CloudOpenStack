
resource "openstack_networking_floatingip_v2" "worker2-fip" {
  pool = "public"
}

resource "openstack_networking_floatingip_associate_v2" "worker2-fip_assoc" {
  floating_ip = openstack_networking_floatingip_v2.worker2-fip.address
  port_id = "f22fbfbe-10e3-47ad-b1ad-cd2b2937467f"
}
