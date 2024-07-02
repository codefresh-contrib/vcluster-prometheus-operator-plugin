# vCluster Prometheus operator plugin

[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/codefresh-inc/codefresh-contrib%2Fvcluster-prometheus-operator-plugin%2Fci?type=cf-1)]( https://g.codefresh.io/public/accounts/codefresh-inc/pipelines/new/667d09d7db9387f8b1294868)

This repository contains a [vCluster plugin](https://www.vcluster.com/docs/v0.19/advanced-topics/plugins-overview) that syncs [Prometheus operator](https://github.com/prometheus-operator/prometheus-operator) resources from virtual clusters to the host cluster.

Currently only the very basic functionality is implemented so the plugin only supports syncing of `PodMonitor` and `ServiceMonitor` resources. This is to allow scraping of metrics from workloads running on virtual clusters from a signle Prometheus or Open Telemetry collector on the host (with [target allocator](https://github.com/open-telemetry/opentelemetry-operator/blob/main/cmd/otel-allocator/README.md) that supports Prometheus operator CRDs).

The repository contains 2 versions of the plugin, each version is compatible with different versions of vCluster, but both versions provide the same functionality.
This compatibility is required as Syncer arhitecture was overhauled by loft.sh in version 0.20.0 and plugin-sdk changed accordingly.

`v1` - Compatible with older versions of vCluster - latest confirmed & tested version is `0.16.4`

`v2` - Compatibe with vCluster version `0.20.0-beta.9` which was the latest version at the time of writing.
