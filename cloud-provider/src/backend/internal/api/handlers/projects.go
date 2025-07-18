package handlers

import (
    "net/http"
	"fmt"
	"os"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/gophercloud/gophercloud"
    "github.com/gophercloud/gophercloud/openstack"
    "github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
    "cloud-provider/src/backend/internal/config"
    "cloud-provider/src/backend/terraform/terraform_utilis" 
)

type ProjectHandler struct {
    config         *config.Config
    identityClient *gophercloud.ServiceClient
}

func NewProjectHandler(cfg *config.Config) *ProjectHandler {
    // Initialize OpenStack client
    opts := gophercloud.AuthOptions{
        IdentityEndpoint: cfg.OpenStack.AuthURL,
        Username:         cfg.OpenStack.Username,
        Password:         cfg.OpenStack.Password,
        TenantName:       cfg.OpenStack.TenantName,
        DomainName:       cfg.OpenStack.DomainName,
    }

    provider, err := openstack.AuthenticatedClient(opts)
    if err != nil {
        panic("Failed to create OpenStack client: " + err.Error())
    }

    // Create Identity client for projects
    identityClient, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{
        Region: cfg.OpenStack.Region,
    })
    if err != nil {
        panic("Failed to create identity client: " + err.Error())
    }

    return &ProjectHandler{
        config:         cfg,
        identityClient: identityClient,
    }
}

// ListProjects retrieves all projects from OpenStack
func (h *ProjectHandler) ListProjects(c *gin.Context) {
    opts := projects.ListOpts{}
    allPages, err := projects.List(h.identityClient, opts).AllPages()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to retrieve projects",
            "details": err.Error(),
        })
        return
    }

    allProjects, err := projects.ExtractProjects(allPages)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to extract projects",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "total":    len(allProjects),
        "projects": allProjects,
    })
}

// GetProject retrieves a specific project by ID
func (h *ProjectHandler) GetProject(c *gin.Context) {
    projectID := c.Param("id")
    
    project, err := projects.Get(h.identityClient, projectID).Extract()
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error":   "Project not found",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "project": project,
    })
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
    var req struct {
        Name        string `json:"name" binding:"required"`
        Description string `json:"description"`
        enabled      bool   `json:"enable" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 1. Prepare project directory
    projectDir := fmt.Sprintf("/app/terraform/projects/%s", req.Name)
    if err := os.MkdirAll(projectDir, 0755); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project directory: " + err.Error()})
        return
    }

    // 2. Write main.tf (resource definition)
    tfContent := fmt.Sprintf(`
resource "openstack_identity_project_v3" "%s" {
  name        = "%s"
  description = "%s"
  enabled     = %t
  domain_id   = "default"
}
`, req.Name, req.Name, req.Description, req.enabled)
    tfPath := fmt.Sprintf("%s/main.tf", projectDir)
    if err := os.WriteFile(tfPath, []byte(tfContent), 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write Terraform file: " + err.Error()})
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
    c.JSON(http.StatusCreated, gin.H{
        "terraform_file": tfPath,
        "message": "Terraform file generated and applied to create the project in OpenStack.",
    })
}

