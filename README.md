# vCluster Prometheus operator plugin

[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/codefresh-inc/codefresh-contrib%2Fvcluster-prometheus-operator-plugin%2Fci?type=cf-1)]( https://g.codefresh.io/public/accounts/codefresh-inc/pipelines/new/667d09d7db9387f8b1294868)

This repository contains a [vCluster plugin](https://www.vcluster.com/docs/v0.19/advanced-topics/plugins-overview) that syncs [Prometheus operator](https://github.com/prometheus-operator/prometheus-operator) resources from virtual clusters to the host cluster.

Currently only the very basic functionality is implemented so the plugin only supports syncing of `PodMonitor` and `ServiceMonitor` resources. This is to allow scraping of metrics from workloads running on virtual clusters from a signle Prometheus or Open Telemetry collector on the host (with [target allocator](https://github.com/open-telemetry/opentelemetry-operator/blob/main/cmd/otel-allocator/README.md) that supports Prometheus operator CRDs).

This plugin is compatibe with vCluster version `0.21.0` and above, tested on `v0.23.0`, which was latest at the time of writing this plugin.

# Development

For more information how to develop plugins in vcluster, please refer to the [official vcluster docs](https://www.vcluster.com/docs/plugins/overview).

## Using the Plugin in vcluster

To use the plugin, create a new vcluster with the `plugin.yaml`:

```
# Deploy Prometheus operator on host cluster with Helm:
For more info see -https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack?modal=install

# Create vcluster with plugin
vcluster create my-vcluster -n my-vcluster -f https://raw.githubusercontent.com/codefresh-contrib/vcluster-prometheus-operator-plugin/main/v2/plugin.yaml
```

This will create a new vcluster with the plugin installed. Then test the plugin with:

```
# Apply example ServicerMonitor
vcluster connect my-vcluster -n my-vcluster -- kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/example/user-guides/getting-started/example-app-service-monitor.yaml

# Check if car got correctly synced
kubectl get servicemonitor -n my-vcluster
```

## Building the Plugin
To just build the plugin image and push it to the registry, run:
```
# Build
docker build . -t my-repo/my-plugin:0.0.1

# Push
docker push my-repo/my-plugin:0.0.1
```

Then exchange the image in the `plugin.yaml`

## Development Process

General vcluster plugin project structure:
```
.
├── go.mod              # Go module definition
├── go.sum
├── devspace.yaml       # Development environment definition
├── devspace_start.sh   # Development entrypoint script
├── Dockerfile          # Production Dockerfile
├── main.go             # Go Entrypoint
├── plugin.yaml         # Plugin Helm Values
├── syncers/            # Plugin Syncers
└── manifests/          # Additional plugin resources
```

Before starting to develop, make sure you have installed the following tools on your computer:
- [docker](https://docs.docker.com/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) with a valid kube context configured
- [helm](https://helm.sh/docs/intro/install/), which is used to deploy vcluster and the plugin
- [vcluster CLI](https://www.vcluster.com/docs/getting-started/setup) v0.20.0 or higher
- [DevSpace](https://devspace.sh/cli/docs/quickstart), which is used to spin up a development environment

After successfully setting up the tools, start the development environment with:
```
devspace dev -n vcluster
```

After a while a terminal should show up with additional instructions. Enter the following command to start the plugin:
```
go build -mod vendor -o plugin main.go && /vcluster/syncer start
```

You can now change a file locally in your IDE and then restart the command in the terminal to apply the changes to the plugin.

Delete the development environment with:
```
devspace purge -n vcluster
```
