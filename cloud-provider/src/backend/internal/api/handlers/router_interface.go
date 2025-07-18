
package handlers	

import (
    "fmt"
    "os"
    "net/http"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/router-interfaces
func (h *ProjectHandler) CreateRouterInterface(c *gin.Context) {
    var req struct {
        RouterID string `json:"router_id" binding:"required"`
        SubnetID string `json:"subnet_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	projectDir := fmt.Sprintf("/app/terraform/network/%s", req.RouterID)
	os.MkdirAll(projectDir, 0755)
	tfContent := fmt.Sprintf(`
resource "openstack_networking_router_interface_v2" "router_interface" {
  router_id = "%s"
  subnet_id = "%s"
}
`, req.RouterID, req.SubnetID)
    tfPath := fmt.Sprintf("%s/router_interface.tf", projectDir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	    // 3. Copy provider.tf using your utility function
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

    c.JSON(http.StatusCreated, gin.H{"message": "Router interface created", "terraform_file": tfPath})
}