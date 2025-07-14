package terraform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"cloud-provider/internal/config"
)

type Executor struct {
	config *config.Config
}

func NewExecutor(cfg *config.Config) *Executor {
	return &Executor{
		config: cfg,
	}
}

func (e *Executor) getTerraformPath(environment string) string {
	return filepath.Join(e.config.Terraform.WorkingDir, "environments", environment)
}

func (e *Executor) initTerraform(ctx context.Context, environment string) (*tfexec.Terraform, error) {
	workingDir := e.getTerraformPath(environment)
	
	// Ensure directory exists
	if err := os.MkdirAll(workingDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create working directory: %w", err)
	}

	// Initialize terraform
	tf, err := tfexec.NewTerraform(workingDir, "terraform")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize terraform: %w", err)
	}

	// Set environment variables for OpenStack
	tf.SetEnv(e.config.GetOpenStackEnv())

	// Initialize terraform working directory
	if err := tf.Init(ctx); err != nil {
		return nil, fmt.Errorf("failed to run terraform init: %w", err)
	}

	return tf, nil
}

func (e *Executor) Plan(ctx context.Context, environment string, variables map[string]string) (string, error) {
	tf, err := e.initTerraform(ctx, environment)
	if err != nil {
		return "", err
	}

	// Convert variables to terraform format
	var tfVars []tfexec.PlanOption
	for key, value := range variables {
		tfVars = append(tfVars, tfexec.Var(fmt.Sprintf("%s=%s", key, value)))
	}

	// Run terraform plan
	hasChanges, err := tf.Plan(ctx, tfVars...)
	if err != nil {
		return "", fmt.Errorf("failed to run terraform plan: %w", err)
	}

	if hasChanges {
		return "Changes detected in terraform plan", nil
	}
	
	return "No changes detected in terraform plan", nil
}

func (e *Executor) Apply(ctx context.Context, environment string, variables map[string]string) (string, error) {
	tf, err := e.initTerraform(ctx, environment)
	if err != nil {
		return "", err
	}

	// Convert variables to terraform format
	var tfVars []tfexec.ApplyOption
	for key, value := range variables {
		tfVars = append(tfVars, tfexec.Var(fmt.Sprintf("%s=%s", key, value)))
	}

	// Run terraform apply
	if err := tf.Apply(ctx, tfVars...); err != nil {
		return "", fmt.Errorf("failed to run terraform apply: %w", err)
	}

	return "Terraform apply completed successfully", nil
}

func (e *Executor) Destroy(ctx context.Context, environment string, variables map[string]string) (string, error) {
	tf, err := e.initTerraform(ctx, environment)
	if err != nil {
		return "", err
	}

	// Convert variables to terraform format
	var tfVars []tfexec.DestroyOption
	for key, value := range variables {
		tfVars = append(tfVars, tfexec.Var(fmt.Sprintf("%s=%s", key, value)))
	}

	// Run terraform destroy
	if err := tf.Destroy(ctx, tfVars...); err != nil {
		return "", fmt.Errorf("failed to run terraform destroy: %w", err)
	}

	return "Terraform destroy completed successfully", nil
}

func (e *Executor) GetState(ctx context.Context, environment string) (interface{}, error) {
	tf, err := e.initTerraform(ctx, environment)
	if err != nil {
		return nil, err
	}

	// Get terraform state
	state, err := tf.Show(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get terraform state: %w", err)
	}

	return state, nil
}
