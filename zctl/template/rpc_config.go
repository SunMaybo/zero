package template

const RPCConfigTemplate = "package config\n\nimport \"zero/common/zrpc\"\n\ntype Config struct {\n\tRpc zrpc.RpcCfg `yaml:\"rpc\"`\n}"
