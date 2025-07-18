package handlers

import (
    "fmt"
    "os"
    "net/http"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/routers
func (h *ProjectHandler) CreateRouter(c *gin.Context) {
    var req struct {
        Name             string `json:"name" binding:"required"`
        ExternalNetworkID string `json:"external_network_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    projectDir := fmt.Sprintf("/app/terraform/network/%s", req.Name)
    os.MkdirAll(projectDir, 0755)
    tfContent := fmt.Sprintf(`
resource "openstack_networking_router_v2" "%s" {
  name                = "%s"
  admin_state_up      = true
  external_network_id = "%s"
}
`, req.Name, req.Name, req.ExternalNetworkID)
    tfPath := fmt.Sprintf("%s/router.tf", projectDir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }


	providerSrc := "/app/terraform/projects/provider.tf"
    providerDst := fmt.Sprintf("%s/provider.tf", projectDir)
    if err := terraform_utilis.CopyFile(providerSrc, providerDst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy provider file: " + err.Error()})
        return
    }

    // 4. Prepare variables
    vars := terraform_utilis.GetTerraformConf()

    // 5. Apply Terraform using your utility function
    if err := terraform_utilis.ApplyTerraform(projectDir, vars); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Terraform apply failed: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Router created", "terraform_file": tfPath})
}