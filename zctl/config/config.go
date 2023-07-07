package config

type Config struct {
	Proxy                     string `yaml:"proxy" json:"proxy"`
	Maven                     string `yaml:"maven_exec" json:"maven"`
	MavenDeploymentRepository string `yaml:"maven_deployment_repository" json:"maven_deployment_repository"`
	MavenSettings             string `yaml:"maven_settings" json:"maven_settings"`
	FrontWebPk                string `yaml:"front_web_pk" json:"front_web_pk"`
	DingTalkSecret            string `yaml:"ding_talk_secret" json:"ding_talk_secret"`
	CdnKey                    string `yaml:"cdn_key"`
}
