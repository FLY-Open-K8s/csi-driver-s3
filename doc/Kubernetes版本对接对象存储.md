# Kubernetes版本对接对象存储

[toc]

## 对象存储

> 想要通过创建 PersistentVolume（PV）/PersistentVolumeClaim（PVC），并为工作负载挂载数据卷的方式使用云对象存储 COS

对于对象存储，在k8s中不需要PV/PVC来做资源抽象，应用可以直接访问和使用，如果需要在k8s支持访问，就需要相应的csi插件，支持对象存储转成文件存储
例如： 几个大厂的方式
阿里：https://help.aliyun.com/document_detail/130911.html
腾讯：https://cloud.tencent.com/document/product/457/44232
华为：https://support.huaweicloud.com/intl/zh-cn/usermanual-cce/cce_01_0267.html

如果对接文件存储，可以使用使用S3fs可以把Bucket当成一个文件夹挂载到Linux系统内部，当成一个系统文件夹使用。

## 社区方案

### 方案1：Object Storage API (COSI)

> 目前开发中，不成熟

> **官方文档**
>
> https://container-object-storage-interface.github.io/docs/
>
> https://container-object-storage-interface.github.io/
>
> https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/1979-object-storage-support#provisionergetinfo
>
> [Kubernetes v1.23 正式发布，有哪些增强？](https://zhuanlan.zhihu.com/p/442623224)



Kubernetes 现在对文件和块存储都有了较好的扩展支持 (CSI)，但是这些并不能很好的支持对象存储，原因如下：

- oss 以桶 (bucket) 来组织分配存储单元而不是文件系统挂载或是块设备
- oss 服务的访问是通过网络调用而不是本地的 POSIX 调用
- oss 不试用 csi 定义的 Attach/Detach 逻辑 (无需挂载/卸载)

COSI 是关于如何为容器化的工作负载 (Pod) 提供对象存储服务 (oss) 的标准协议。与 csi 和 cni 类似，Kubernetes 旨在通过定义一些标准的接口与第三方 oss 服务提供方解耦。



![img](https://cdn.jsdelivr.net/gh/Fly0905/note-picture@main/img/202204252104317.jpeg)

### **[方案2：csi-s3](https://github.com/ctrox/csi-s3)**

> 可以使用，但需要多测试

> 官方文档
>
> https://github.com/CTrox/csi-s3
>
> https://github.com/majst01/csi-driver-s3
>
> [使用s3(minio)为kubernetes提供pv存储](https://www.lishuai.fun/2021/01/07/k8s-pv-s3/)

#### Kubernetes 要求

- Kubernetes 1.13+（CSI v1.0.0 兼容性）
- <font color=red>**Kubernetes 必须允许特权容器**</font>
- Docker 守护进程必须允许共享挂载（systemd 标志`MountFlags=shared`）
- <font color=red>**是非常实验性的，还未在任何生产环境中使用。根据使用的挂载程序和 S3 存储后端，可能会发生意外的数据丢失。**</font>

##### MountFlags

在 18.09 之前的 Docker 版本中，containerd 由 Docker 引擎守护进程管理。

在 Docker Engine 18.09 中，containerd 由 systemd 管理。由于 containerd 由 systemd 管理，因此任何`docker.service`更改挂载设置的 systemd 配置的自定义配置（例如，`MountFlags=slave`）都会破坏 Docker Engine 守护进程和 containerd 之间的交互，并且您将无法启动容器。

运行以下命令以获取 的`MountFlags`属性的当前值`docker.service`：

```bash
$ sudo systemctl show --property=MountFlags docker.service
MountFlags=
```

如果此命令为 打印非空值，请更新您的配置`MountFlags`，然后重新启动 docker 服务。

## 为什么要将S3 以文件存储的方式挂载到 Kubernetes 平台?

对于原来使用本地目录访问数据的应用程序，比如使用本地磁盘或网络共享盘保存数据的应用系统，如果用户希望把数据放到S3上，则需要修改数据的访问方式，比如修改为使用SDK 或CLI访问S3中存储的数据。

同时实现kubernetes集群是很多用户的需求，无论是使用托管服务还是自建kubernetes集群，存储都是kubernetes集群搭建的重点。并且docker的部署方式也让客户程序减少了对于底层环境的依赖。为了让用户原来的应用系统能在不做修改的情况下直接使用S3服务，**需要把S3存储桶作为目录挂载到用户kubernetes集群中的worker节点上**。

> 利用S3fs将S3存储桶在kubernetes平台上以`sidecar`方式挂载到kubernetes集群的worker实例上的pod中，挂载后需要读写此存储桶的pod都可以对此桶进行读写，以实现共享存储功能。

## 什么是 S3FS ？

**S3fs是基于FUSE的文件系统**，允许Linux和Mac Os X 挂载S3的存储桶在本地文件系统，S3fs能够保持对象原来的格式,S3FS是POSIX的大子集，包括读/写文件、目录、符号链接、模式、uid/gid和扩展属性，**与**AmazonS3、Google云存储和其他**基于S3的对象存储兼容**。关于S3fs的详细介绍，请参见：https://github.com/s3fs-fuse/s3fs-fuse

## 后续

1. 验证特权容器，升级了多少特权

## 参考链接

1. [使用s3(minio)为kubernetes提供pv存储](https://www.lishuai.fun/2021/01/07/k8s-pv-s3/)
2. [基于openshift+华为对象存储的CSI开发](https://cloud.tencent.com/developer/article/1608177)
3. [利用 S3FS 将 S3 作为共享存储挂载到 Kubernetes Pod](https://aws.amazon.com/cn/blogs/china/use-u3fs-as-shared-storage-to-kubernetes-pod/)
4. [使用S3fs在Linux实例上挂载Bucket](https://docs.jdcloud.com/cn/object-storage-service/s3fs)
5. [S3FS：基于对象存储的文件系统](https://blog.shunzi.tech/post/s3fs/)
6. https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#podsecuritycontext-v1-core

