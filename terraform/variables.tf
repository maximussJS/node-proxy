variable "hcloud_token" {
  description = "Hetzner Cloud API Token"
  type        = string
  sensitive   = true
}

variable "hcloud_ssh_key_fingerprint" {
    description = "SSH key fingerprint to use for the server"
    type        = string
    sensitive   = true
}

variable "server_name" {
  description = "Kava Json RPC Node Proxy"
  type        = string
  default     = "kava-json-rpc-node-proxy"
}

variable "image" {
  description = "Golang Docker image to run"
  type        = string
  default     = "maximussxgod/json-rpc-node-proxy:v3"
}

variable "private_key_path" {
  description = "Path to the private SSH key for connecting to the server"
  type        = string
  default     = "~/.ssh/hetzner_id_rsa"
}

variable "private_key_content" {
  description = "Content of the private SSH key for connecting to the server"
  type        = string
  sensitive   = true
  default     = ""
}
