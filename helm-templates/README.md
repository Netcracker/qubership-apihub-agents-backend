# Qubership APIHUB Agents Backend Helm Chart

This folder contains `qubership-apihub-agents-backend` Helm chart for Qubership API Linter Service deployment to k8s
cluster.

It is ready for usage Helm chart.

## 3rd party dependencies

| Name       | Version | Mandatory/Optional | Comment |
|------------|---------|--------------------|---------|
| Kubernetes | 1.23+   | Mandatory          |         |

## HWE

|                | CPU request | CPU limit | RAM request | RAM limit |
|----------------|-------------|-----------|-------------|-----------|
| Default values | 30m         | 1         | 256Mi       | 256Mi     |

## Prerequisites

1. kubectl installed and configured for k8s cluster access.
2. Helm installed

## Set up values.yml

1. Download Qubership APIHUB Agents Backend helm chart
2. Fill `values.yaml` with corresponding deploy parameters. `values.yaml` is self-documented, so please refer to it

## Execute helm install

In order to deploy Qubership APIHUB Agents Backend to your k8s cluster execute the following command:

```
helm install qubership-apihub-agents-backend -n qubership-apihub-agents-backend --create-namespace  -f ./qubership-apihub-agents-backend/values.yaml ./qubership-apihub-agents-backend
```

In order to uninstall Qubership APIHUB Agents Backend from your k8s cluster execute the following command:

```
helm uninstall qubership-apihub-agents-backend -n qubership-apihub-agents-backend
```

## Dev cases

**Installation to local k8s cluster**

File `local-k8s-values.yaml` has predefined deploy parameters for deploy to local k8s cluster on your PC.

Execute the following command to deploy Qubership APIHUB Agents Backend:

```
helm install qubership-apihub-agents-backend -n qubership-apihub-agents-backend --create-namespace  -f ./qubership-apihub-agents-backend/local-k8s-values.yaml ./qubership-apihub-agents-backend
```
