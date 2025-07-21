package handlers

import (
    "fmt"
    "os"
    "net/http"
    "log"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/instances
func (h *ProjectHandler) CreateInstance(c *gin.Context) {
    var req struct {
        Name           string `json:"name" binding:"required"`
        Image          string `json:"image" binding:"required"`
        Flavor         string `json:"flavor" binding:"required"`
        NetworkID      string `json:"network_id" binding:"required"`
        KeyPair        string `json:"keypair" binding:"required"`
        SecurityGroup  string `json:"security_group" binding:"required"`
        Count          int    `json:"count"` // Optional: number of instances to create
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    instanceDir := fmt.Sprintf("/app/terraform/instances/%s", req.Name)
    if err := os.MkdirAll(instanceDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create instance directory: " + err.Error()})
        return
    }

    countStr := ""
    if req.Count > 1 {
        countStr = fmt.Sprintf("  count = %d\n", req.Count)
    }

    tfContent := fmt.Sprintf(`
resource "openstack_compute_instance_v2" "%s" {
  name            = "%s"
  image_name      = "%s"
  flavor_name     = "%s"
  key_pair        = "%s"
  security_groups = ["%s"]
  network {
    uuid = "%s"
  }
%s}
`, terraform_utilis.CleanResourceName(req.Name), req.Name, req.Image, req.Flavor, req.KeyPair, req.SecurityGroup, req.NetworkID, countStr)
    tfPath := fmt.Sprintf("%s/instance.tf", instanceDir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Terraform file: " + err.Error()})
        return
    }

    providerSrc := "/app/terraform/projects/provider.tf"
    providerDst := fmt.Sprintf("%s/provider.tf", instanceDir)
    if err := terraform_utilis.CopyFile(providerSrc, providerDst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy provider file: " + err.Error()})
        return
    }

    vars := terraform_utilis.GetTerraformConf()
    if err := terraform_utilis.ApplyTerraform(instanceDir, vars); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Terraform apply failed: " + err.Error()})
        return
    }

    log.Printf("Instance Terraform file created and applied at %s", tfPath)
    c.JSON(http.StatusCreated, gin.H{
        "terraform_file": tfPath,
        "message": "Instance(s) created in OpenStack.",
    })
}