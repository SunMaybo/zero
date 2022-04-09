package template

const RPCConfigTemplate = "package config\n\nimport \"github.com/SunMaybo/zero/common/zrpc\"\n\ntype Config struct {\n\tRpc zrpc.RpcCfg `yaml:\"rpc\"`\n}"
