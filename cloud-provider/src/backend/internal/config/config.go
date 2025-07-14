package config

import (
	"fmt"
    "github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"environment"`
	Port        string `mapstructure:"port"`
	
	OpenStack struct {
		AuthURL    string `mapstructure:"auth_url"`
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		TenantName string `mapstructure:"tenant_name"`
		Region     string `mapstructure:"region"`
		DomainName string `mapstructure:"domain_name"`
	} `mapstructure:"openstack"`
	
	Terraform struct {
		WorkingDir string `mapstructure:"working_dir"`
		StateDir   string `mapstructure:"state_dir"`
	} `mapstructure:"terraform"`
}

func Load() (*Config, error) {
	var cfg Config
	
	// Set default values
	viper.SetDefault("environment", "development")
	viper.SetDefault("port", "8080")
	viper.SetDefault("terraform.working_dir", "/app/terraform")
	viper.SetDefault("terraform.state_dir", "/app/terraform/state")
	
	// Read from config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	
	// Read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CLOUD_PROVIDER")
	
	// Environment variable mappings
	viper.BindEnv("openstack.auth_url", "OS_AUTH_URL")
	viper.BindEnv("openstack.username", "OS_USERNAME")
	viper.BindEnv("openstack.password", "OS_PASSWORD")
	viper.BindEnv("openstack.tenant_name", "OS_TENANT_NAME")
	viper.BindEnv("openstack.region", "OS_REGION_NAME")
	viper.BindEnv("openstack.domain_name", "OS_DOMAIN_NAME")
	
	// Try to read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}
	
	// Unmarshal config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	
	// Validate required OpenStack credentials
	if cfg.OpenStack.AuthURL == "" {
		return nil, fmt.Errorf("OpenStack auth URL is required")
	}
	if cfg.OpenStack.Username == "" {
		return nil, fmt.Errorf("OpenStack username is required")
	}
	if cfg.OpenStack.Password == "" {
		return nil, fmt.Errorf("OpenStack password is required")
	}
	
	return &cfg, nil
}

func (c *Config) GetOpenStackEnv() []string {
	return []string{
		fmt.Sprintf("OS_AUTH_URL=%s", c.OpenStack.AuthURL),
		fmt.Sprintf("OS_USERNAME=%s", c.OpenStack.Username),
		fmt.Sprintf("OS_PASSWORD=%s", c.OpenStack.Password),
		fmt.Sprintf("OS_TENANT_NAME=%s", c.OpenStack.TenantName),
		fmt.Sprintf("OS_REGION_NAME=%s", c.OpenStack.Region),
		fmt.Sprintf("OS_DOMAIN_NAME=%s", c.OpenStack.DomainName),
	}
}
