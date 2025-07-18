package handlers


import (
    "fmt"
    "os"
    "net/http"
	"log"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/subnets
func (h *ProjectHandler) CreateSubnet(c *gin.Context) {
    var req struct {
        Name      string `json:"name" binding:"required"`
        NetworkID string `json:"network_id" binding:"required"`
        CIDR      string `json:"cidr" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    projectDir := fmt.Sprintf("/app/terraform/network/%s", req.Name)
    os.MkdirAll(projectDir, 0755)
    tfContent := fmt.Sprintf(`
resource "openstack_networking_subnet_v2" "%s" {
  name            = "%s"
  network_id      = "%s"
  cidr            = "%s"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}
`, req.Name, req.Name, req.NetworkID, req.CIDR)
    tfPath := fmt.Sprintf("%s/subnet.tf", projectDir)
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
	
    c.JSON(http.StatusCreated, gin.H{"message": "Subnet created", "terraform_file": tfPath})
}