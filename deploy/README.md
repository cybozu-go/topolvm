Deploying TopoLVM
=================

Each of these steps are shown in depth in the following sections:

1. Deploy [lvmd][] as a `systemd` service on a worker node with LVM installed.
1. Prepare [cert-manager][] for [topolvm-controller][]. You may supplement an existing instance.
1. Create the `topolvm-system` namespace using `deploy/manifests/base/namespace.yaml`.
1. Determine how [topolvm-scheduler][] to be run:
   - If you run with a managed control plane (such as GKE, AKS, etc), `topolvm-scheduler` should be deployed as Deployment and Service
   - `topolvm-scheduler` should otherwise be deployed as DaemonSet in unmanaged (i.e. bare metal) deployments
1. Add `topolvm.cybozu.com/webhook: ignore` label to system namespaces such as `kube-system`.
1. Apply remaining manifests for TopoLVM from `deploy/manifests/base` plus overlays as appropriate to your installation.
1. Configure `kube-scheduler` to use `topolvm-scheduler`. 
1. Prepare StorageClasses for TopoLVM.

Example configuration files are included in the following sub directories:

- `lvmd-config/`: Configuration file for `lvmd`
- `manifests/`: Manifests for Kubernetes.
- `scheduler-config/`: Configurations to extend `kube-scheduler` with `topolvm-scheduler`.
- `systemd/`: A systemd unit file for `lvmd`.

These configuration files may need to be modified for your environment.
Read carefully the following descriptions.

lvmd
----

[lvmd][] is a gRPC service to manage an LVM volume group.  The pre-built binary can be downloaded from [releases page](https://github.com/topolvm/topolvm/releases).
It can be built from source code by `mkdir build; go build -o build/lvmd ./pkg/lvmd`.

`lvmd` can setup as a daemon or a Kubernetes Daemonset.

### Setup `lvmd` as daemon

1. Prepare LVM volume groups.  A non-empty volume group can be used because LV names wouldn't conflict.
2. Edit [lvmd.yaml](./lvmd-config/lvmd.yaml) if you want to specify the device-class settings to use multiple volume groups. See [lvmd.md](../docs/lvmd.md) for details.

    ```yaml
    device-classes:
      - name: ssd
        volume-group: myvg1
        default: true
        spare-gb: 10
    ```

3. Install `lvmd` and `lvmd.service`, then start the service.

### Setup `lvmd` using Kubernetes DaemonSet

Also, you can setup [lvmd][] using Kubernetes DaemonSet.

Notice: The lvmd container uses `nsenter` to run some lvm commands(like `lvcreate`) as a host process, so you can't launch lvmd with DaemonSet when you're using [kind](https://kind.sigs.k8s.io/).

To setup `lvmd` with Daemonset:

1. Prepare LVM volume groups in the host. A non-empty volume group can be used because LV names wouldn't conflict.
2. Edit `ConfigMap` in [lvmd-configmap.yaml](./manifests/lvmd/lvmd-configmap.yaml) as follows:

    ```yaml
    ---
    apiVersion: v1
    kind: ConfigMap
    metadata:
      namespace: topolvm-system
      name: lvmd
    data:
      lvmd.yaml: |
        socket-name: /run/topolvm/lvmd.sock
        device-classes:
          - name: ssd
            volume-group: myvg1 # Change this value to your VG name.
            default: true
            spare-gb: 10
    ```

3. Apply `lvmd` manifests.

    ```console
    kustomize build ./manifests/overlays/lvmd/ | kubectl apply -f -
    ```

4. Check DaemonSet `lvmd` is ready

    ```console
    # If DaemonSet lvmd is ready, launching lvmd is a success
    $ kubectl -n topolvm-system get daemonset topolvm-lvmd
    NAME           DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
    topolvm-lvmd   1         1         1       1            1           <none>          45m
    ```

### Migrate lvmd, which is running as daemon, to DaemonSet

This section describes how to switch to DaemonSet lvmd from lvmd running as daemons.

1. Apply the manifest in `deploy/lvmd` folder and start lvmd with DaemonSet.
   **You need to set the temporal `socket-name` which is not the same as the value in lvmd running as daemon.**
   After applying the manifest, DaemonSet lvmd and lvmd running as daemon exist at the same time using different sockets.

    ```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      namespace: topolvm-system
      name: lvmd
    data:
      lvmd.yaml: |
        socket-name: /run/topolvm/lvmd.sock # Change this value to something like `/run/topolvm/lvmd-work.sock`.
        device-classes:
          - name: ssd
            volume-group: myvg1
            default: true
            spare-gb: 10
    ```

2. Change the options of topolvm-node to communicate with the DaemonSet lvmd instead of lvmd running as daemon.
   **You should set the temporal socket name which is not the same as in lvmd running as daemon.**

    ```yaml
    $ kubect edit daemonset -n topolvm-system node
    <snip>
        spec:
          serviceAccountName: node
          containers:
            - name: topolvm-node
              image: quay.io/topolvm/topolvm-with-sidecar:latest
              securityContext:
                privileged: true
              command:
                - /topolvm-node
                - --lvmd-socket=/run/lvmd/lvmd.sock # Change this value to to something linke `/run/lvmd/lvmd-work.sock`.
    <snip>
    ```
3. Check if you can create Pod/PVC and can access to existing PV.

4. Stop and remove lvmd running as daemon.

5. Change the `socket-name` and `--lvmd-socket` options to the original one.
   To reflect the changes of ConfigMap, restart DamonSet lvmd manually.

    ```yaml
    $ kubect edit cm -n topolvm-system lvmd
    <snip>
    data:
      lvmd.yaml: |
        socket-name: /run/topolvm/lvmd-work.sock # Change this value to something like `/run/topolvm/lvmd.sock`.
    <snip>

    $ kubect edit daemonset -n topolvm-system node
    <snip>
        spec:
          serviceAccountName: node
          containers:
            - name: topolvm-node
              image: quay.io/topolvm/topolvm-with-sidecar:latest
              securityContext:
                privileged: true
              command:
                - /topolvm-node
                - --lvmd-socket=/run/lvmd/lvmd-work.sock # Change this value to something linke `/run/lvmd/lvmd.sock`.
    <snip>
    ```
cert-manager
------------

[cert-manager][] is used to issue self-signed TLS certificate for [topolvm-controller][].
Follow the [documentation](https://docs.cert-manager.io/en/latest/getting-started/install/kubernetes.html) to install it into your Kubernetes cluster.

### OPTIONAL: Prepare the certificate without cert-manager

You can prepare the certificate manually without `cert-manager`.
When doing so, do not apply [certificates.yaml](./manifests/base/certificates.yaml).

1. Prepare PEM encoded self-signed certificate and key files.  
    The certificate must be valid for hostname `controller.topolvm-system.svc`.
2. Base64-encode the CA cert (in its PEM format
3. Create Secret in `topolvm-system` namespace as follows:

    ```console
    kubectl -n topolvm-system create secret tls mutatingwebhook \
        --cert=<CERTIFICATE FILE> --key=<KEY FILE>
    ```

4. Edit `MutatingWebhookConfiguration` in [webhooks.yaml](./manifests/base/mutating/webhooks.yaml) as follows:

    ```yaml
    apiVersion: admissionregistration.k8s.io/v1
    kind: MutatingWebhookConfiguration
    metadata:
      name: topolvm-hook
    # snip
    webhooks:
      - name: pvc-hook.topolvm.cybozu.com
        # snip
        clientConfig:
          caBundle: |  # Base64-encoded, PEM-encoded CA certificate that signs the server certificate
            ...
      - name: pod-hook.topolvm.cybozu.com
        # snip
        clientConfig:
          caBundle: |  # The same CA certificate as above
            ...
    ```

Create namespace
----------------

Create the `topolvm-system` namespace using `deploy/manifests/base/namespace.yaml`. This manifest has labels that are required for proper operation. Do not apply the other manifests yet, these will be applied after editing the values for Kustomize in the following steps.

Scheduing
-----------------

### topolvm-scheduler

[topolvm-scheduler][] is a [scheduler extender](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/scheduling/scheduler_extender.md) for `kube-scheduler`.
It must be deployed to where `kube-scheduler` can connect.

If your Kubernetes cluster runs the control plane on Nodes, `topolvm-scheduler` should be run as DaemonSet
limited to the control plane nodes.  `kube-scheduler` then connects to the extender via loopback network device.

Otherwise, `topolvm-scheduler` should be run as Deployment and Service.
`kube-scheduler` then connects to the Service address.

#### Running topolvm-scheduler using DaemonSet

The [example manifest](./manifests/overlays/daemonset-scheduler/scheduler.yaml) can be used almost as is.
You may need to change the taint key or label name of the DaemonSet.

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  namespace: topolvm-system
  name: topolvm-scheduler
spec:
  # snip
      hostNetwork: true               # If kube-scheduler does not use host network, change this false.
      tolerations:                    # Add tolerations needed to run pods on control plane nodes.
        - key: CriticalAddonsOnly
          operator: Exists
        - effect: NoSchedule
          key: node-role.kubernetes.io/control-plane
        - effect: NoSchedule
          key: node-role.kubernetes.io/master
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: node-role.kubernetes.io/control-plane # match the control plane node specific labels
                    operator: Exists
              - matchExpressions:
                  - key: node-role.kubernetes.io/master # match the control plane node specific labels
                    operator: Exists
```

#### Running topolvm-scheduler using Deployment and Service

In this case, you can use [deployment-scheduler/scheduler.yaml](./manifests/overlays/deployment-scheduler/scheduler.yaml) instead of [daemonset-scheduler/scheduler.yaml](./manifests/overlays/daemonset-scheduler/scheduler.yaml).

This way, `topolvm-scheduler` is exposed by LoadBalancer service.

Then edit `urlPrefix` in [scheduler-config-v1beta1.yaml](./scheduler-config/scheduler-config-v1beta1.yaml) for K8s 1.19 or later, to specify the LoadBalancer address.

#### OPTIONAL: tune the node scoring

The node scoring for Pod scheduling can be fine-tuned with the following two ways:
1. Adjust `divisor` parameter in the scoring expression
2. Change the weight for the node scoring against the default by kube-scheduler

The scoring expression in `topolvm-scheduler` is as follows:
```
min(10, max(0, log2(capacity >> 30 / divisor)))
```
For example, the default of `divisor` is `1`, then if a node has the free disk capacity more than `1024GiB`, `topolvm-scheduler` scores the node as `10`. `divisor` should be adjusted to suit each environment. It can be specified the default value and values for each device-class in [scheduler-options.yaml](./manifests/overlays/daemonset-scheduler/scheduler-options.yaml) as follows:

```yaml
default-divisor: 1
divisors:
  ssd: 1
  hdd: 10
```

Besides, the scoring weight can be passed to kube-scheduler via [scheduler-config-v1beta1.yaml](./scheduler-config/scheduler-config-v1beta1.yaml). Almost all scoring algorithms in kube-scheduler are weighted as `"weight": 1`. So if you want to give a priority to the scoring by `topolvm-scheduler`, you have to set the weight as a value larger than one like as follows:
```yaml
apiVersion: kubescheduler.config.k8s.io/v1beta1
kind: KubeSchedulerConfiguration
leaderElection:
  leaderElect: true
clientConnection:
  kubeconfig: /etc/kubernetes/scheduler.conf
extenders:
- urlPrefix: "http://127.0.0.1:9251"
  filterVerb: "predicate"
  prioritizeVerb: "prioritize"
  nodeCacheCapable: false
  weight: 100 ## EDIT THIS FIELD ##
  managedResources:
  - name: "topolvm.cybozu.com/capacity"
    ignoredByScheduler: true
```

### Storage Capacity Tracking

topolvm supports [Storage Capacity Tracking](https://kubernetes.io/docs/concepts/storage/storage-capacity/).
You can enable Storage Capacity Tracking mode instead of using topolvm-scheduler.
You need to use Kubernetes Cluster v1.21 or later when using Storage Capacity Tracking with topolvm.

You can see the limitations of using Storage Capacity Tracking from [here](https://kubernetes.io/docs/concepts/storage/storage-capacity/#scheduling).

#### Use Storage Capacity Tracking

If you want to use Storage Capacity Tracking instead of using topolvm-scheduler, you need following:

- Add `--enable-capacity` and `--capacity-ownerref-level=2` to csi-provisoner arguments.
- Append `NAMESPACE` and `POD_NAME` environment variables to csi-provisioner.

The [example manifest](./manifests/overlays/storage-capcacity/controller.yaml) can be used almost as is.

Deploy example is below:

```
kustomize build ./deploy/manifests/overlays/storage-capcacity | kubectl apply -f -
```

Protect system namespaces from TopoLVM webhook
---------------------------------------------

TopoLVM installs a mutating webhook for Pods.  It may prevent Kubernetes from bootstrapping
if the webhook pods and the system pods are both missing.

To workaround the problem, add a label to system namespaces such as `kube-system` as follows:

```console
$ kubectl label ns kube-system topolvm.cybozu.com/webhook=ignore
```

This label was already applied to the `topolvm-system` namespace via the `deploy/manifests/base/namespace.yaml` manifest.

Apply remaining manifests for TopoLVM
-------------------------------------

Previous sections describe how to tune the manifest configurations, apply them now using Kustomize as follows:

### Running topolvm-scheduler using DaemonSet

If using `topolvm-scheduler` as a DaemonSet, run the following command: 

```console
kustomize build ./deploy/manifests/overlays/daemonset-scheduler | kubectl apply -f -
```

### Running topolvm-scheduler using Deployment and Service

If using `topolvm-scheduler` as a Deployment, run the following command: 

```console
kustomize build ./deploy/manifests/overlays/deployment-scheduler | kubectl apply -f -
```

### Generate cert-manager manifests

Unless you chose to generate and install certificates manually, run the following:

```console
kubectl apply -f ./deploy/manifests/base/certificates.yaml
```

Configure kube-scheduler
------------------------

`kube-scheduler` need to be configured to use `topolvm-scheduler` extender.

First you need to choose an appropriate `KubeSchdulerConfiguration` YAML file according to your Kubernetes version.

```console
cp ./deploy/scheduler-config/scheduler-config-v1beta1.yaml ./deploy/scheduler-config/scheduler-config.yaml
```

And then copy the `deploy/scheduler-config` directory to the hosts where `kube-scheduler`s run.

### For new clusters

If you are installing your cluster from scratch with `kubeadm`, you can use the following configuration:

```yaml
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
metadata:
  name: config
kubernetesVersion: v1.18.2
scheduler:
  extraVolumes:
    - name: "config"
      hostPath: /path/to/scheduler-config     # absolute path to ./scheduler-config directory
      mountPath: /var/lib/scheduler
      readOnly: true
  extraArgs:
    config: /var/lib/scheduler/scheduler-config.yaml
```

### For existing clusters

The changes to `/etc/kubernetes/manifests/kube-scheduler.yaml` that are affected by this are as follows:

1. Add a line to the `command` arguments array such as ```- --config=/var/lib/scheduler/scheduler-config.yaml```. Note that this is the location of the file **after** it is mapped to the `kube-scheduler` container, not where it exists on the node local filesystem.
2. Add a volume mapping to the location of the configuration on your node:

    ```yaml
      spec.volumes:
      - hostPath:
          path: /path/to/scheduler-config     # absolute path to ./scheduler-config directory
          type: Directory
        name: topolvm-config
    ```

3. Add a `volumeMount` for the scheduler container:

    ```yaml
      spec.containers.volumeMounts:
      - mountPath: /var/lib/scheduler
        name: topolvm-config
        readOnly: true
    ```

Prepare StorageClasses
----------------------

Finally, you need to create [StorageClasses](https://kubernetes.io/docs/concepts/storage/storage-classes/) for TopoLVM.

An example is available in [provisioner.yaml](./manifests/base/provisioner.yaml).

See [podpvc.yaml](../example/podpvc.yaml) for how to use TopoLVM provisioner.

[lvmd]: ../docs/lvmd.md
[cert-manager]: https://github.com/jetstack/cert-manager
[topolvm-scheduler]: ../docs/topolvm-scheduler.md
[topolvm-controller]: ../docs/topolvm-controller.md
