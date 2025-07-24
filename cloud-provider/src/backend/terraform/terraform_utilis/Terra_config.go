

package terraform_utilis


import (
    "regexp"
    "strings"
)

func GetTerraformConf() map[string]string {
    vars := map[string]string{
        "os_auth_url":    "http://91.99.215.184/identity/v3",
        "os_user_name":   "admin",
        "os_password":    "secret",
        "os_tenant_name": "admin",
        "os_domain_name": "default",
        "os_region":      "RegionOne",
    }
    return vars
}


// CleanResourceName returns a Terraform-safe resource name
func CleanResourceName(name string) string {
    // Replace spaces and dots with underscores
    cleaned := strings.ReplaceAll(name, " ", "_")
    cleaned = strings.ReplaceAll(cleaned, ".", "_")
    // Remove all non-alphanumeric and non-underscore characters
    re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
    cleaned = re.ReplaceAllString(cleaned, "")
    // Ensure it doesn't start with a digit
    if len(cleaned) > 0 && cleaned[0] >= '0' && cleaned[0] <= '9' {
        cleaned = "_" + cleaned
    }
    return cleaned
}