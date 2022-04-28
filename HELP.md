## 注意项

- 项目需要在 Linux 下运行，golang 开发版本为 1.18，ubuntu 20.04 下测试
- Docker 已经将 aufs 改为 Overlay
- 本实例还是使用 aufs，看系统是否支持 `cat /proc/filesystems`
- ubuntu 20.04 需要安装 aufs `apt-get install aufs-tools`
- 所有操作在 `/root` 目录下

### 项目测试需使用 docker busybox 镜像

```bash
# 下载 busybox
docker pull busybox
# 运行
docker run -d busybox top -b
# 导出
docker export -o busybox.tar (容器ID)
# 解压到 /root 文件夹下
cd /root
mkdir busybox
tar -xvf busybox.tar -C busybox/
```

### 操作指南

```bash
# 编译
go build .

# 启动一个容器, busybox为镜像名，存放在 /root/busybox.tar
./simple-docker run -ti --name test busybox sh

# 后台启动
./simple-docker run -d --name test busybox sh

# 挂载文件
./simple-docker run -d -v /root/test:/test --name test busybox sh

# 进入容器
./simple-docker exec test sh

# 查看容器日志
./simple-docker logs test

# 查看容器列表
./simple-docker ps

# 停止容器
./simple-docker stop test

# 删除容器
./simple-docker rm test
```

### 指令小记

- 查看 Linux 程序父进程

```bash
pstree -pl | grep main
```

- 查看进程id

```bash
echo $$
```

- 查看进程的 uts

```bash
readling /proc/进程id/ns/uts
```

- 修改hostname

```bash
hostname -b 新名称
```

- 常看当前用户和用户组

```bash
id
```

- 创建并挂载一个 hierarchy

> 在这个文件夹下面创建新的文件夹，会被 kernel 标记为该 cgroup 的子 cgroup

```bash
mkdir cgroup-test
mount -t cgroup -o none,name=cgroup-test cgroup-test ./cgroup-test
```

- 将其他进程移动到其他的 cgroup 中

> 只要将该进程的 ID 放到其 cgroup 的 tasks 里面即可

```bash
echo "进程ID" >> cgroup/tasks
```

- 导出容器

```bash
docker export -o busybox.tar 45c98e055883(容器ID)
```

- 移除mount

```bash
unshare -m
```