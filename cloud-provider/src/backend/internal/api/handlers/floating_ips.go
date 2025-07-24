package handlers

import (
    "fmt"
    "os"
    "net/http"
    "log"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/floatingips
func (h *ProjectHandler) CreateFloatingIP(c *gin.Context) {
    var req struct {
        Name              string `json:"name" binding:"required"`
        FloatingNetworkID string `json:"floating_network_id" binding:"required"` // External network ID
        PortID            string `json:"port_id"`                                // Optional: If you want to associate with a specific port
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fipDir := fmt.Sprintf("/app/terraform/floatingips/%s", req.Name)
    if err := os.MkdirAll(fipDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create floating IP directory: " + err.Error()})
        return
    }

    // If PortID is provided, associate with port; otherwise, associate with instance
    associateBlock := ""
    if req.PortID != "" {
        associateBlock = fmt.Sprintf("  port_id = \"%s\"\n", req.PortID)
    }

    tfContent := fmt.Sprintf(`
resource "openstack_networking_floatingip_v2" "%s" {
  pool = "%s"
}

resource "openstack_networking_floatingip_associate_v2" "%s_assoc" {
  floating_ip = openstack_networking_floatingip_v2.%s.address
%s}
`, req.Name, req.FloatingNetworkID, req.Name, req.Name, associateBlock)

    tfPath := fmt.Sprintf("%s/floatingip.tf", fipDir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Terraform file: " + err.Error()})
        return
    }

    providerSrc := "/app/terraform/projects/provider.tf"
    providerDst := fmt.Sprintf("%s/provider.tf", fipDir)
    if err := terraform_utilis.CopyFile(providerSrc, providerDst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy provider file: " + err.Error()})
        return
    }

    vars := terraform_utilis.GetTerraformConf()
    if err := terraform_utilis.ApplyTerraform(fipDir, vars); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Terraform apply failed: " + err.Error()})
        return
    }

    log.Printf("Floating IP Terraform file created and applied at %s", tfPath)
    c.JSON(http.StatusCreated, gin.H{
        "terraform_file": tfPath,
        "message": "Floating IP allocated and associated in OpenStack.",
    })
}