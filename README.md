## Update Docker Image
```bash
docker build -t maximussxgod/json-rpc-node-proxy:<tag> . 

docker push maximussxgod/json-rpc-node-proxy:<tag>
```

## Run Docker Image Locally
```bash
docker run -d -p 8080:8080 \
  -v $(pwd)/config.yaml:/build/config.yaml \
  maximussxgod/json-rpc-node-proxy:<tag>
```

## Deploy to Hetzner Cloud
Create a `terraform/terraform.tfvars` file with the following content:
```hcl
hcloud_token = "your-hetzner-cloud-api-token"
hcloud_ssh_key_fingerprint = "your-ssh-key-fingerprint"
image = "maximussxgod/json-rpc-node-proxy:<tag>"
```

Run the following commands:
```bash
cd terraform
terraform plan
terraform apply
```
