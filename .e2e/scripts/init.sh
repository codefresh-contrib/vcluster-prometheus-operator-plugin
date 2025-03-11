#!/bin/bash
# Deploy CRDs and vcluster with the plugin, connect to vcluster and apply resources in resources folder for testing

PLUGIN_IMAGE=$1
MYDIR=$(dirname $0)
ROOT_DIR=$MYDIR/../../
RESOURCES_FILE=$MYDIR/../../.e2e/vcluster-resources.yaml

kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/v0.70.0/example/prometheus-operator-crd-full/monitoring.coreos.com_podmonitors.yaml
kubectl apply -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/v0.70.0/example/prometheus-operator-crd-full/monitoring.coreos.com_servicemonitors.yaml
helm upgrade --install --repo https://charts.loft.sh vcluster vcluster --version v0.20.0-beta.5 --values $MYDIR/vcluster-values.yaml --values $ROOT_DIR/plugin.yaml --set plugins.prometheus-operator-resources.image=$PLUGIN_IMAGE --wait

vcluster connect vcluster -n default -- kubectl get servicemonitor && vcluster connect vcluster -n default -- kubectl get podmonitor

res=$?
secondsWaited=0
timeout=300

while [ $res -ne 0 ] && [ $secondsWaited -lt $timeout ]; do
    echo "Waiting for CRDs to get created, sleep for 30 seconds..."
    sleep 30
    secondsWaited=$((secondsWaited + 30))
    vcluster connect vcluster -n default -- kubectl get servicemonitor && vcluster connect vcluster -n default -- kubectl get podmonitor
    res=$?
done

if [ $res -ne 0 ]; then
  echo "Timed out waiting for CRDs to get created in vcluster"
  exit 1
fi

cat $RESOURCES_FILE | vcluster connect vcluster -n default -- kubectl -n default apply -f -
