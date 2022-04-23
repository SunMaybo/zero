package template

const HttpConfigTemplate = "package config\n\nimport \"github.com/SunMaybo/zero/common/zcfg\"\n\ntype Config struct {\n\tZero zcfg.ZeroConfig `yaml:\"zero\"`\n}\n"
