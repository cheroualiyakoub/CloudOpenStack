

package terraform_utilis

func GetTerraformConf() map[string]string {
    vars := map[string]string{
        "os_auth_url":    "http://94.237.83.40/identity/v3",
        "os_user_name":   "admin",
        "os_password":    "secret",
        "os_tenant_name": "admin",
        "os_domain_name": "default",
        "os_region":      "RegionOne",
    }
    return vars
}
