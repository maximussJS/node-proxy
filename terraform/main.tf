terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = ">= 1.46.1"
    }
  }
}

provider "hcloud" {
  token = var.hcloud_token
}

resource "hcloud_server" "node_proxy" {
  name        = var.server_name
  server_type = "cpx11" # Adjust as needed (cx11, cx21, etc.)
  image       = "ubuntu-22.04" # Use the desired base image
  location    = "nbg1"         # Location (e.g., nbg1, fsn1, hel1)
  ssh_keys    = [data.hcloud_ssh_key.ssh_key.id]
}

data "hcloud_ssh_key" "ssh_key" {
  fingerprint = var.hcloud_ssh_key_fingerprint
}

resource "null_resource" "docker_setup" {
  depends_on = [hcloud_server.node_proxy]

  provisioner "file" {
    connection {
      type        = "ssh"
      host        = hcloud_server.node_proxy.ipv4_address
      user        = "root"
      private_key = var.private_key_content != "" ? var.private_key_content : file(var.private_key_path)
    }

    source      = "${path.module}/../config.yaml"
    destination = "/root/config.yaml"
  }

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      host        = hcloud_server.node_proxy.ipv4_address
      user        = "root"
      private_key = var.private_key_content != "" ? var.private_key_content : file(var.private_key_path)
    }

    inline = [
      "apt-get update && apt-get install -y docker.io",
      "docker pull ${var.image}",
      "docker run -d -p 8080:8080 -v /root/config.yaml:/build/config.yaml ${var.image}"
    ]
  }
}



