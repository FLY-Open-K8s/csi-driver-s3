# CSI for S3

This is a Container Storage Interface ([CSI](https://github.com/container-storage-interface/spec/blob/master/spec.md)) for S3 (or S3 compatible) storage. This can dynamically allocate buckets and mount them via a fuse mount into any container.

## Status

This is still very experimental and should not be used in any production environment. Unexpected data loss could occur depending on what mounter and S3 storage backend is being used.

## Kubernetes installation

### Requirements

* Kubernetes 1.16+ (CSI v1.0.0 compatibility)
* Kubernetes has to allow privileged containers
* Docker daemon must allow shared mounts (Systemd flag `MountFlags=shared`)

### Create a secret with your S3 credentials

```yaml
# deploy/kubernetes/secret-minio.yaml
apiVersion: v1
kind: Secret
metadata:
  name: csi-driver-s3-secret
  namespace: kube-system
stringData:
  accessKeyID: <YOUR_ACCESS_KEY_ID>
  secretAccessKey: <YOUR_SECRET_ACCES_KEY>
  # For AWS set it to "https://s3.<region>.amazonaws.com"
  endpoint: <S3_ENDPOINT_URL>
  # If not on S3, set it to ""
  region: <S3_REGION>
```

The region can be empty if you are using some other S3 compatible storage.

### Deploy the driver

```bash
kubectl apply -f deploy/kubernetes
```

### Test the S3 driver

Create a PVC

```bash
kubectl apply -f deploy/sample/pvc.yaml
```

Check if the PVC has been bound:

```bash
$ kubectl get pvc
NAME                STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
csi-driver-s3-pvc   Bound     pvc-c5d4634f-8507-11e8-9f33-0e243832354b   5Gi        RWO            csi-driver-s3  9s
```

Create a test pod which mounts your volume:

```bash
kubectl apply -f deploy/sample/pod.yaml
```

If the pod can start, everything should be working.

Test the mount

```bash
$ kubectl exec -ti csi-driver-s3-test-nginx bash
$ mount | grep fuse
s3fs on /var/lib/www/html type fuse.s3fs (rw,nosuid,nodev,relatime,user_id=0,group_id=0,allow_other)
$ touch /var/lib/www/html/hello_world
```

If something does not work as expected, check the troubleshooting section below.

## Additional configuration

### Mounter

As S3 is not a real file system there are some limitations to consider here. Depending on what mounter you are using, you will have different levels of POSIX compatibility. Also depending on what S3 storage backend you are using there are not always [consistency guarantees](https://github.com/gaul/are-we-consistent-yet#observed-consistency).

The csi-driver-s3 uses s3fs to mount the bucket into the PVC.

#### s3fs

* Large subset of POSIX
* Files can be viewed normally with any S3 client
* Does not support appends or random writes
* [s3fs](https://github.com/s3fs-fuse/s3fs-fuse)

## Troubleshooting

### Issues while creating PVC

Check the logs of the provisioner:

```bash
kubectl logs -l app=csi-provisioner-s3 -c csi-driver-s3
```

### Issues creating containers

1. Ensure feature gate `MountPropagation` is not set to `false`
2. Check the logs of the s3-driver:

```bash
kubectl logs -l app=csi-driver-s3 -c csi-driver-s3
```

## Development

This project can be built like any other go application.

```bash
go get -u github.com/majst01/csi-driver-s3
```

### Build executable

```bash
make build
```

### Tests

Currently, the driver is tested by the [CSI Sanity Tester](https://github.com/kubernetes-csi/csi-test/tree/master/pkg/sanity). As end-to-end tests require S3 storage and a mounter like s3fs, this is best done in a docker container. A Dockerfile and the test script are in the `test` directory. The easiest way to run the tests is to just use the make command:

```bash
make test
```
