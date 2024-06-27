#!/bin/bash
# Deploy CRDs and vcluster with the plugin, connect to vcluster and apply resources in resources folder for testing

PLUGIN_IMAGE=$1
MYDIR=$(dirname $0)
ROOT_DIR=$MYDIR/../../
RESOURCES_FILE=$MYDIR/../vcluster-resources.yaml

kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/v0.70.0/example/prometheus-operator-crd-full/monitoring.coreos.com_podmonitors.yaml
kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/v0.70.0/example/prometheus-operator-crd-full/monitoring.coreos.com_servicemonitors.yaml
helm upgrade --install --repo https://charts.loft.sh vcluster vcluster --version 0.16.4 --values $MYDIR/vcluster-values.yaml --values $ROOT_DIR/plugin.yaml --set plugin.prometheus-operator-resources.image=$PLUGIN_IMAGE --wait

sleep 120 # Give time for for CRDs to sync

cat $RESOURCES_FILE | vcluster connect vcluster -n default -- kubectl -n default apply -f -
