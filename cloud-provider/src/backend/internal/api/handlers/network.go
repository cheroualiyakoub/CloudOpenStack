package handlers


import (
    "fmt"
    "os"
    "net/http"
	"log"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/networks
func (h *ProjectHandler) CreateNetwork(c *gin.Context) {
    var req struct {
        Name string `json:"name" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    projectDir := fmt.Sprintf("/app/terraform/networks/%s", req.Name)
    os.MkdirAll(projectDir, 0755)
    tfContent := fmt.Sprintf(`
resource "openstack_networking_network_v2" "%s" {
  name           = "%s"
  admin_state_up = true
}
`, req.Name, req.Name)
    tfPath := fmt.Sprintf("%s/network.tf", projectDir)
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

    log.Printf("Terraform file created and applied at %s", tfPath)

    c.JSON(http.StatusCreated, gin.H{"message": "Network created", "terraform_file": tfPath})
}