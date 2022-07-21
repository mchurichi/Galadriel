#!/bin/bash

function create_cluster() {
    cat <<EOF | kind create cluster --config=-
        kind: Cluster
        apiVersion: kind.x-k8s.io/v1alpha4
        nodes:
        - role: control-plane
          extraMounts:
          - hostPath: $(pwd)/sockets
            containerPath: /host/sockets
EOF
}

function deploy_spire() {
    echo "Deploying SPIRE Servers..."
    helm install spire1 ./spire -f ./spire/values-spire1.yaml
    helm install spire2 ./spire -f ./spire/values-spire2.yaml
}

function cleanup() {
    echo "Cleaning up SPIRE deployments..."
    helm uninstall spire1
    helm uninstall spire2
}

$1