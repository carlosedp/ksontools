# yaml2jsonnet

Given YAML to create a Kubernetes objects, convert it jsonnet.

```yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

## This could generate

```json
local k = import "k.libsonnet";
local deployment = "k.apps.v1beta2.deployment";

local deploymentInstance = deployment.new()
```