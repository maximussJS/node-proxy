## Docker
```bash

docker build -t maximussxgod/json-rpc-node-proxy:v2 . 

docker run -p 8080:8080 \
  --env-file .env -v $(pwd)/.env:/build/.env \
  -v $(pwd)/config.yaml:/build/config.yaml \
  maximussxgod/json-rpc-node-proxy:v2

```