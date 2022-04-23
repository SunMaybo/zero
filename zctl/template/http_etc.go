package template

const HttpETCTemplate = `zero:
  server:
    port: 3000
    timeout: 30
  center:
    registry: false
    name: nacos
    log_level: warn
    namespace_id: 7d064d60-1cd0-4fc4-8637-a5fb6cfca57a
    server:
      - scheme: http
        context_path: /nacos
        host: 127.0.0.1
        port: 8849`
