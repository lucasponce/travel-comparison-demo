# Travels Demo - Version 2
A Microservices demo based on Istio Service Mesh tool. 

This goal of this demo is to demostrate [Istio](https:/istio.io/) cappabilities observed and managed by [Kiali](https://kiali.io) tool.

## Platform Install

This demo has been tested using [Minikube](https://istio.io/latest/docs/setup/platform-setup/minikube/) and [Istio 1.6 Demo Profile](https://istio.io/latest/docs/setup/install/istioctl/#install-a-different-profile)

## Quick Start

Create `travel-agency`, `travel-portal` and `travel-control` namespaces. Add `istio-injection` label and deploy demo app. 

```yaml
kubectl create namespace travel-agency
kubectl create namespace travel-portal
kubectl create namespace travel-control

kubectl label namespace travel-agency istio-injection=enabled
kubectl label namespace travel-portal istio-injection=enabled
kubectl label namespace travel-control istio-injection=enabled

kubectl apply -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_agency.yaml) -n travel-agency
kubectl apply -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_portal.yaml) -n travel-portal
kubectl apply -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_control.yaml) -n travel-control
```

Open Kiali dashboard:

```bash
istioctl dashboard kiali
```

Expose `travel-control` service to your localhost machine:

```bash
kubectl port-forward svc/control 8080:8080 -n travel-control
```

Open [Travels Dashboard](http://localhost:8080).

Undeploy the example:
```yaml
kubectl delete -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_agency.yaml) -n travel-agency
kubectl delete -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_portal.yaml) -n travel-portal
kubectl delete -f <(curl -L https://raw.githubusercontent.com/lucasponce/travel-comparison-demo/v2/travel_control.yaml) -n travel-control

kubectl delete namespace travel-agency
kubectl delete namespace travel-portal
kubectl delete namespace travel-control
```

