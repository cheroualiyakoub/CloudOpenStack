
resource "openstack_networking_secgroup_v2" "k8s-secgroup" {
  name        = "k8s-secgroup"
  description = "Security group for Kubernetes cluster"
}

resource "openstack_networking_secgroup_rule_v2" "k8s-secgroup_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.k8s-secgroup.id
}

resource "openstack_networking_secgroup_rule_v2" "k8s-secgroup_ssh_ipv6" {
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "::/0"
  security_group_id = openstack_networking_secgroup_v2.k8s-secgroup.id
}

resource "openstack_networking_secgroup_rule_v2" "k8s-secgroup_nodeport" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.k8s-secgroup.id
}

resource "openstack_networking_secgroup_rule_v2" "k8s-secgroup_k8s_api" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6443
  port_range_max    = 6443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.k8s-secgroup.id
}

resource "openstack_networking_secgroup_rule_v2" "k8s-secgroup_internal" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "192.168.100.0/24" // Or use a variable if needed
  security_group_id = openstack_networking_secgroup_v2.k8s-secgroup.id
}
