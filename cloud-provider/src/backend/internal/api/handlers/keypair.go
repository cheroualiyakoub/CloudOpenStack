package handlers

import (
    "fmt"
    "os"
    "net/http"
    "log"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/keypairs
func (h *ProjectHandler) CreateKeyPair(c *gin.Context) {
    var req struct {
        Name      string `json:"name" binding:"required"`
        PublicKey string `json:"public_key" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    keypairDir := fmt.Sprintf("/app/terraform/keypairs/%s", req.Name)
    if err := os.MkdirAll(keypairDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create keypair directory: " + err.Error()})
        return
    }

    tfContent := fmt.Sprintf(`
resource "openstack_compute_keypair_v2" "%s" {
  name       = "%s"
  public_key = "%s"
}
`, req.Name, req.Name, req.PublicKey)
    tfPath := fmt.Sprintf("%s/keypair.tf", keypairDir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Terraform file: " + err.Error()})
        return
    }

    providerSrc := "/app/terraform/projects/provider.tf"
    providerDst := fmt.Sprintf("%s/provider.tf", keypairDir)
    if err := terraform_utilis.CopyFile(providerSrc, providerDst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy provider file: " + err.Error()})
        return
    }

    vars := terraform_utilis.GetTerraformConf()
    if err := terraform_utilis.ApplyTerraform(keypairDir, vars); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Terraform apply failed: " + err.Error()})
        return
    }

    log.Printf("Keypair Terraform file created and applied at %s", tfPath)
    c.JSON(http.StatusCreated, gin.H{
        "terraform_file": tfPath,
        "message": "Key pair created and applied in OpenStack.",
    })
}