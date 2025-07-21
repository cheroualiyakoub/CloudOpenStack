package handlers

import (
    "fmt"
    "os"
    "net/http"
    "log"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/security-groups
func (h *ProjectHandler) CreateSecurityGroup(c *gin.Context) {
    var req struct {
        Name        string `json:"name" binding:"required"`
        Description string `json:"description"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    sgDir := fmt.Sprintf("/app/terraform/security_groups/%s", req.Name)
    os.MkdirAll(sgDir, 0755)
    tfContent := fmt.Sprintf(`
resource "openstack_networking_secgroup_v2" "%s" {
  name        = "%s"
  description = "%s"
}

resource "openstack_networking_secgroup_rule_v2" "%s_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.%s.id
}

resource "openstack_networking_secgroup_rule_v2" "%s_k8s_api" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6443
  port_range_max    = 6443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.%s.id
}
`, req.Name, req.Name, req.Description, req.Name, req.Name, req.Name, req.Name)

    tfPath := fmt.Sprintf("%s/security_group.tf", sgDir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    providerSrc := "/app/terraform/projects/provider.tf"
    providerDst := fmt.Sprintf("%s/provider.tf", sgDir)
    if err := terraform_utilis.CopyFile(providerSrc, providerDst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy provider file: " + err.Error()})
        return
    }

    vars := terraform_utilis.GetTerraformConf()
    if err := terraform_utilis.ApplyTerraform(sgDir, vars); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Terraform apply failed: " + err.Error()})
        return
    }

    log.Printf("Security group Terraform file created and applied at %s", tfPath)
    c.JSON(http.StatusCreated, gin.H{
        "terraform_file": tfPath,
        "message": "Security group created and applied in OpenStack.",
    })
}