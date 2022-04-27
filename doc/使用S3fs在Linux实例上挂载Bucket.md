# 使用S3fs在Linux实例上挂载Bucket

S3fs是基于FUSE的文件系统，允许Linux 挂载Bucket在本地文件系统，S3fs能够保持对象原来的格式。使用S3fs可以把Bucket当成一个文件夹挂载到Linux系统内部，当成一个系统文件夹使用。

## 环境安装以及配置参考官方说明

https://github.com/s3fs-fuse/s3fs-fuse

**1.安装依赖包**

On CentOS 7:

```
sudo yum install automake fuse fuse-devel gcc-c++ git libcurl-devel libxml2-devel make openssl-devel
```

On Ubuntu 16.04:

```
sudo apt-get install automake autotools-dev fuse g++ git libcurl4-openssl-dev libfuse-dev libssl-dev libxml2-dev make pkg-config
```

**2.安装**

```
git clone https://github.com/s3fs-fuse/s3fs-fuse.git
cd s3fs-fuse
./autogen.sh
./configure
make
sudo make install
```

**3.创建密码文件**

```
echo Access_Key_ID:Access_Key_Secret > ~/.passwd-s3fs
chmod 600 ~/.passwd-s3fs
```

说明

Access_Key_ID:Access_Key_Secret获取方式：https://uc.jdcloud.com/account/accessKey

chmod 600：设置密钥文件只能被当前用户访问。

**4.挂载对象存储到本地目录/new**

```
mkdir /new
s3fs bucketname /new -o passwd_file=~/.passwd-s3fs -o url="https://s3.cn-north-1.jdcloud-oss.com"
```

说明

mkdir：创建new文件夹作为本地挂载目录

s3fs：手动挂载命令，其中bucketname为bucket名称、/new是本地挂载路径、passwd_file为密码文件位置、url为[京东云对象存储兼容S3域名](https://docs.jdcloud.com/cn/object-storage-service/regions-and-endpoints)（请输入bucket所在区域的服务域名）

![image-20220426142505954](https://cdn.jsdelivr.net/gh/Fly0905/note-picture@main/img/202204261425159.png)

**5.查看挂载结果**

```
df -h
```

[![查看挂载结果](https://jdcloud-portal.oss.cn-north-1.jcloudcs.com/cn/image/Object-Storage-Service/OSS-072.png)](https://jdcloud-portal.oss.cn-north-1.jcloudcs.com/cn/image/Object-Storage-Service/OSS-072.png)

**6.进入目录可以查看到object文件**

[![查看object文件](https://cdn.jsdelivr.net/gh/Fly0905/note-picture@main/img/202204261132736.png)](https://jdcloud-portal.oss.cn-north-1.jcloudcs.com/cn/image/Object-Storage-Service/OSS-073.png)

**Tips：**

1.如果您使用的CentOS 6，请参考以下步骤安装fuse依赖：

```
yum install automake gcc-c++ git libcurl-devel libxml2-devel make openssl-devel

wget https://github.com/libfuse/libfuse/releases/download/fuse_2_9_4/fuse-2.9.2.tar.gz
tar -zxvf fuse-2.9.2.tar.gz
cd fuse-2.9.2
./configure --prefix=/usr
make
make install
export PKG_CONFIG_PATH=/usr/lib/pkgconfig:/usr/lib64/pkgconfig/
ldconfig
```

2.使用Mac OS安装S3fs，请参考以下步骤：

```
git clone https://github.com/s3fs-fuse/s3fs-fuse.git
cd s3fs-fuse
./autogen.sh
./configure --prefix=/usr/local
PKG_CONFIG_PATH="/usr/local/opt/openssl/lib/pkgconfig:/usr/local/opt/libxml2/lib/pkgconfig"
make
sudo make install
```

注：--prefix=/usr/local非必须；PKG_CONFIG_PATH必须，/usr/local/要替换成用户本地路径。

3.如果您在挂载Bucket时使用的非root账号，请在指定挂载命令是指定当前账户的uid及gid。如下示例：

```
sudo s3fs bucketname /new -o passwd_file=~/.passwd-s3fs -o url="https://s3.cn-north-1.jdcloud-oss.com" -o uid=11111 -o gid=11111
```

4.使用s3fs-fuse工具挂载京东云对象存储，通过cp命令拷贝文件时，若遇到文件无content-type的问题，可通过如下方式解决：

- 使用`cp`命令拷贝文件，`s3fs-fuse`工具底层进行的操作依赖于`/etc/mime.types`文件，这个文件决定了`cp`命令目的文件的mime-type属性。请查看您目录下是否有该文件。
- 默认情况下，京东云的centos7版本并不包含`/etc/mime.types`文件，所以需要通过拷贝，或者安装`httpd`获得，安装命令为`yum install httpd`
- 对于已经通过`s3fs`命令挂载的目录，需要先`umount`，然后再次执行`s3fs`命令才能生效。

5.如果您在使用S3fs挂载Bucket之前开启了静态网站托管，会导致挂载失败；如果您使用S3fs挂载Bucket之后开启了静态网站托管，会导致文件操作失效。

6.如需使用Debug模式，请挂载时指定参数`-o dbglevel=dbg -f -o curldbg`