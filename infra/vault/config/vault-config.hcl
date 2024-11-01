# vault-config.hcl

disable_mlock = true



api_addr = "http://localhost:8200"

storage "file" {
  path = "/vault/data"
}

ui = true
