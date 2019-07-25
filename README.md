# Grackle Operator

A Kubernetes [Operator](https://coreos.com/operators/) for managing [Grackle](https://github.com/jmckind/grackle)
clusters.

## Usage

Create the RBAC resources for the operator.

```bash
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
```

Deploy the operator.

```bash
kubectl create -f deploy/crds/k8s_v1alpha1_grackle_crd.yaml
kubectl create -f deploy/operator.yaml
```

Create a new Grackle resource.

```bash
kubectl create -f deploy/crds/k8s_v1alpha1_grackle_cr.yaml
```

## License

Grackle Operator is released under the Apache 2.0 license. See the [LICENSE][license_file] file for details.

[license_file]:./LICENSE
