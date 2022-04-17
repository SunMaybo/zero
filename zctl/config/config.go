package config

type Config struct {
	Proxy                     string `yaml:"proxy"`
	Maven                     string `yaml:"maven_exec"`
	MavenDeploymentRepository string `yaml:"maven_deployment_repository"`
	MavenSettings             string `yaml:"maven_settings"`
}
