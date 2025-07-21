package handlers

import (
    "fmt"
    "os"
    "net/http"
    "log"


    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)




// POST /api/v1/flavors
func (h *ProjectHandler) CreateFlavor(c *gin.Context) {
    var req struct {
        Name  string `json:"name" binding:"required"`
        VCPUs int    `json:"vcpus" binding:"required"`
        RAM   int    `json:"ram" binding:"required"`   // in MB
        Disk  int    `json:"disk" binding:"required"`  // in GB
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    flavorDir := fmt.Sprintf("/app/terraform/flavors/%s", req.Name)
    if err := os.MkdirAll(flavorDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create flavor directory: " + err.Error()})
        return
    }

    tfContent := fmt.Sprintf(`
resource "openstack_compute_flavor_v2" "%s" {
  name  = "%s"
  ram   = %d
  vcpus = %d
  disk  = %d
}
`, terraform_utilis.CleanResourceName(req.Name), req.Name, req.RAM, req.VCPUs, req.Disk)
    tfPath := fmt.Sprintf("%s/flavor.tf", flavorDir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Terraform file: " + err.Error()})
        return
    }

    providerSrc := "/app/terraform/projects/provider.tf"
    providerDst := fmt.Sprintf("%s/provider.tf", flavorDir)
    if err := terraform_utilis.CopyFile(providerSrc, providerDst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy provider file: " + err.Error()})
        return
    }

    vars := terraform_utilis.GetTerraformConf()
    if err := terraform_utilis.ApplyTerraform(flavorDir, vars); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Terraform apply failed: " + err.Error()})
        return
    }

    log.Printf("Flavor Terraform file created and applied at %s", tfPath)
    c.JSON(http.StatusCreated, gin.H{
        "terraform_file": tfPath,
        "message": "Flavor created in OpenStack.",
    })
}