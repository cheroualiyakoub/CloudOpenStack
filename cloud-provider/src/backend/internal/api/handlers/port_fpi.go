package handlers

import (
    "fmt"
    "os"
    "net/http"
    "log"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/ports-fips
func (h *ProjectHandler) CreatePortsAndFIPs(c *gin.Context) {
    var req struct {
        ClusterName      string `json:"cluster_name" binding:"required"`
        MasterCount      int    `json:"master_count" binding:"required"`
        WorkerCount      int    `json:"worker_count" binding:"required"`
        NetworkID        string `json:"network_id" binding:"required"`
        SubnetID         string `json:"subnet_id" binding:"required"`
        SecGroupID       string `json:"secgroup_id" binding:"required"`
        FloatingPoolName string `json:"floating_pool_name" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    dir := fmt.Sprintf("/app/terraform/cluster/%s", req.ClusterName)
    os.MkdirAll(dir, 0755)

    tfContent := fmt.Sprintf(`
resource "openstack_networking_port_v2" "cluster_ports" {
  count          = %d
  name           = "%s-port-${count.index + 1}"
  network_id     = "%s"
  admin_state_up = true
  security_group_ids = ["%s"]
  fixed_ip {
    subnet_id = "%s"
  }
}

resource "openstack_networking_floatingip_v2" "cluster_floatingips" {
  count = %d
  pool  = "%s"
}

resource "openstack_networking_floatingip_associate_v2" "cluster_fip_assoc" {
  count       = %d
  floating_ip = openstack_networking_floatingip_v2.cluster_floatingips[count.index].address
  port_id     = openstack_networking_port_v2.cluster_ports[count.index].id
}
`, req.MasterCount+req.WorkerCount, req.ClusterName, req.NetworkID, req.SecGroupID, req.SubnetID,
   req.MasterCount+req.WorkerCount, req.FloatingPoolName,
   req.MasterCount+req.WorkerCount)

    tfPath := fmt.Sprintf("%s/ports_fips.tf", dir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Terraform file: " + err.Error()})
        return
    }

    providerSrc := "/app/terraform/projects/provider.tf"
    providerDst := fmt.Sprintf("%s/provider.tf", dir)
    if err := terraform_utilis.CopyFile(providerSrc, providerDst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy provider file: " + err.Error()})
        return
    }

    vars := terraform_utilis.GetTerraformConf()
    if err := terraform_utilis.ApplyTerraform(dir, vars); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Terraform apply failed: " + err.Error()})
        return
    }

    log.Printf("Ports and Floating IPs Terraform file created and applied at %s", tfPath)
    c.JSON(http.StatusCreated, gin.H{
        "terraform_file": tfPath,
        "message": "Ports and Floating IPs created and applied in OpenStack.",
    })
}