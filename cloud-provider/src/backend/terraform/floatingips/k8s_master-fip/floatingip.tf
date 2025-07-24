
resource "openstack_networking_floatingip_v2" "k8s_master-fip" {
  pool = "public"
}

resource "openstack_networking_floatingip_associate_v2" "k8s_master-fip_assoc" {
  floating_ip = openstack_networking_floatingip_v2.k8s_master-fip.address
  port_id = "b473567e-0950-471b-b83e-d73a9237d600"
}
