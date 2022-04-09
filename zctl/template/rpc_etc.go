package template

const RPCETCTemplate = `rpc:
  name: {{.PackageName}}
  port: 8088
  is_online: false
  weight: 1
  group_name: format
  cluster_name: default_test
  enable_metrics: true
  metrics_port: 8843
  metrics_path: /metrics
  metadata:
    active: format
    metrics: "/metrics"
    metrics_port: "8843"
  center:
    enabled: false
    name: nacos
    log_level: warn
    namespace_id: 7d064d60-1cd0-4fc4-8637-a5fb6cfca57a
    server:
      - scheme: http
        context_path: /nacos
        host: 127.0.0.1
        port: 8849`

type EtcConfig struct {
	PackageName string
}
