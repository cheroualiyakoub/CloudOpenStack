
resource "openstack_networking_port_v2" "cluster_ports" {
  count          = 3
  name           = "k8s-cluster-port-${count.index + 1}"
  network_id     = "dd9135fe-1ed3-4e0c-88a3-a14bb8fd9a42"
  admin_state_up = true
  security_group_ids = ["046d7466-6031-4ce7-84f5-0437c3cd3e08"]
  fixed_ip {
    subnet_id = "b4887a2e-c784-4668-b184-94c152125a55"
  }
}

resource "openstack_networking_floatingip_v2" "cluster_floatingips" {
  count = 3
  pool  = "public"
}

resource "openstack_networking_floatingip_associate_v2" "cluster_fip_assoc" {
  count       = 3
  floating_ip = openstack_networking_floatingip_v2.cluster_floatingips[count.index].address
  port_id     = openstack_networking_port_v2.cluster_ports[count.index].id
}
