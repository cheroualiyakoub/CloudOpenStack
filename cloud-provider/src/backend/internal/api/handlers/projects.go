package handlers

import (
    "net/http"
	"fmt"
	"os"

    "github.com/gin-gonic/gin"
    "github.com/gophercloud/gophercloud"
    "github.com/gophercloud/gophercloud/openstack"
    "github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
    "cloud-provider/internal/config"
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
        DomainID    string `json:"domain_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Generate Terraform file only (do NOT create project in OpenStack yet)
    tfContent := fmt.Sprintf(`
resource "openstack_identity_project_v3" "%s" {
  name        = "%s"
  description = "%s"
  domain_id   = "%s"
}
`, req.Name, req.Name, req.Description, req.DomainID)

    tfPath := fmt.Sprintf("/app/terraform/clients/%s.tf", req.Name)
    err := os.WriteFile(tfPath, []byte(tfContent), 0644)
    if err != nil {
        pwd, _ := os.Getwd()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to write Terraform file: " + err.Error(),
            "pwd": pwd,
			"path": tfPath,
			"request": req,
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "terraform_file": tfPath,
        "message": "Terraform file generated. Apply it to create the project in OpenStack.",
    })
}