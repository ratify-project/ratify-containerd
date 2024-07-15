<div align="center">
<img src="logo.svg" width="200">
</div>

# Ratify Containerd Prototype

> [!CAUTION]
> This repository is marked EXPERIMENTAL. It demonstrates only a Proof of Concept. Contents may be altered at any time.

## Getting Started

Please refer to exploration [document](docs/overview.md) for more details.

### Prerequisites

* kubectl
* minikube

### Walkthrough

1. Create a `minikube` cluster with containerd container runtime

    ```bash
    minikube start -n 2 --container-runtime containerd
    ```

2. Configure node RBAC to get namespaced ConfigMap resources

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/akashsinghal/ratify-containerd/main/k8s-templates/clusterrolebinding.yaml
    ```

3. Configure nodes. Wait for 30-40 seconds for daemonset to complete (Note: daemonset pods will not terminate. check logs for completion)

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/akashsinghal/ratify-containerd/main/k8s-templates/configure-nodes.yaml
    ```

4. Apply Ratify ConfigMap

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/akashsinghal/ratify-containerd/main/k8s-templates/ratify-config.yaml
    ```

5. Test with signed image

    ```bash
    kubectl run demo-signed --image=ghcr.io/ratify-project/ratify/notary-image:signed
    kubectl describe pod demo-signed
    ```

6. Test with unsigned image. Pod should fail to pull image and start.

    ```bash
    kubectl run demo-unsigned --image=ghcr.io/ratify-project/ratify/notary-image:unsigned
    ```

7. Check Pod state and verify kublet is failing to pull due to verification plugin rejecting pull

    ```bash
    kubectl describe pod demo-unsigned
    ```


## Code of Conduct

ratify-containerd follows the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).

## Licensing

This project is released under the [Apache-2.0 License](./LICENSE).

## Trademark

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft trademarks or logos is subject to and must follow [Microsoft's Trademark & Brand Guidelines][microsoft-trademark]. Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship. Any use of third-party trademarks or logos are subject to those third-party's policies.

[microsoft-trademark]: https://www.microsoft.com/legal/intellectualproperty/trademarks
