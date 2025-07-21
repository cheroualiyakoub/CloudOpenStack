package handlers

import (
    "fmt"
    "os"
    "net/http"
    "log"

    "github.com/gin-gonic/gin"
    "cloud-provider/src/backend/terraform/terraform_utilis"
)

// POST /api/v1/images
func (h *ProjectHandler) CreateImage(c *gin.Context) {
    var req struct {
        Name     string `json:"name" binding:"required"`
        ImageURL string `json:"image_url" binding:"required"` // URL or file path to the image
        DiskFormat string `json:"disk_format" binding:"required"` // e.g., "qcow2", "raw"
        ContainerFormat string `json:"container_format" binding:"required"` // e.g., "bare"
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    imageDir := fmt.Sprintf("/app/terraform/images/%s", req.Name)
    if err := os.MkdirAll(imageDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create image directory: " + err.Error()})
        return
    }

    tfContent := fmt.Sprintf(`
resource "openstack_images_image_v2" "%s" {
  name             = "%s"
  image_source_url       = "%s"
  disk_format      = "%s"
  container_format = "%s"
}
`, terraform_utilis.CleanResourceName(req.Name), req.Name, req.ImageURL, req.DiskFormat, req.ContainerFormat)
    tfPath := fmt.Sprintf("%s/image.tf", imageDir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Terraform file: " + err.Error()})
        return
    }

    providerSrc := "/app/terraform/projects/provider.tf"
    providerDst := fmt.Sprintf("%s/provider.tf", imageDir)
    if err := terraform_utilis.CopyFile(providerSrc, providerDst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy provider file: " + err.Error()})
        return
    }

    vars := terraform_utilis.GetTerraformConf()
    if err := terraform_utilis.ApplyTerraform(imageDir, vars); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Terraform apply failed: " + err.Error()})
        return
    }

    log.Printf("Image Terraform file created and applied at %s", tfPath)
    c.JSON(http.StatusCreated, gin.H{
        "terraform_file": tfPath,
        "message": "Image created and applied in OpenStack.",
    })
}