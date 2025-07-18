package terraform_utilis

import (
    "context"
    "log"
    "os"
    "github.com/hashicorp/terraform-exec/tfexec"
)

func ApplyTerraform(dir string, vars map[string]string) error {
    tf, err := tfexec.NewTerraform(dir, "terraform")
    if err != nil {
        return err
    }

    // Prepare environment variables for Terraform variables
    envMap := map[string]string{}
    for _, e := range os.Environ() {
        // Split "key=value"
        pair := []rune(e)
        for i, c := range pair {
            if c == '=' {
                envMap[string(pair[:i])] = string(pair[i+1:])
                break
            }
        }
    }
    for k, v := range vars {
        envMap["TF_VAR_"+k] = v
    }
    tf.SetEnv(envMap)

    // Initialize Terraform
    if err := tf.Init(context.Background(), tfexec.Upgrade(true)); err != nil {
        log.Printf("Terraform init failed: %v", err)
        return err
    }
    log.Println("Terraform initialized successfully")

    // Apply Terraform
    if err := tf.Apply(context.Background()); err != nil {
        log.Printf("Terraform apply failed: %v", err)
        return err
    }
    log.Println("Terraform apply completed successfully")

    return nil
}