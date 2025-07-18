To deploy a Kubernetes cluster (control plane + 2 workers) on OpenStack, you need to provision several OpenStack resources in the correct order. Here’s a **step-by-step checklist** for what you should automate or create:

---

### 1. **Network Resources**
- **Network**: Create a private network for your cluster.
- **Subnet**: Create a subnet within the network.
- **Router**: Create a router and attach it to the subnet (and to the external/public network for internet access).
- **Security Groups**: Create security groups to allow necessary ports (SSH, Kubernetes API, etc.).

### 2. **Compute Resources**
- **Flavors**: Define or select flavors (VM sizes) for control plane and worker nodes.
- **Key Pair**: Create or upload an SSH key pair for VM access.
- **Instances**: Create instances (VMs) for:
  - 1 control plane node
  - 2 worker nodes
- **Floating IPs**: Allocate and associate floating IPs if you need public access to nodes.

### 3. **Storage (Optional)**
- **Volumes**: Create and attach volumes if you want persistent storage for your nodes.

### 4. **Other Resources**
- **Images**: Make sure you have a suitable OS image (Ubuntu, CentOS, etc.) available in OpenStack.

---

## **Recommended Order of Operations**

1. **Create network, subnet, and router.**
2. **Create security group(s) and rules.**
3. **Create or select flavors.**
4. **Create or upload key pair.**
5. **Create instances for control plane and workers, attaching them to the network and security group.**
6. **Assign floating IPs if needed.**
7. **(Optional) Attach volumes.**
8. **Install Kubernetes on the VMs (using kubeadm, Ansible, or a tool like k3s, etc.).**

---

## **Summary Table**

| Step | Resource         | Purpose                                 |
|------|------------------|-----------------------------------------|
| 1    | Network/Subnet   | Cluster communication                   |
| 2    | Router           | Internet access                         |
| 3    | Security Group   | Open required ports                     |
| 4    | Flavor           | Define VM specs                         |
| 5    | Key Pair         | SSH access                              |
| 6    | Instances        | Control plane + 2 workers               |
| 7    | Floating IP      | Public access (optional)                |
| 8    | Volume           | Persistent storage (optional)           |
| 9    | Image            | OS for VMs                              |

---

**After all resources are created, you can proceed to install and configure Kubernetes on your instances.**

Let me know if you want Terraform examples for any of these steps!