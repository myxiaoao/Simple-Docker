# Simple Docker

> Tips：用来学习 Docker

简单来说 `docker` 本质其实是一个特殊的进程，这个进程特殊在它被 `Namespace` 和 `CGroup` 技术做了装饰，`Namespace` 将该进程与 `Linux`
系统进行隔离开来，让该进程处于一个虚拟的沙盒中，而 `CGroup` 则对该进程做了一系列的资源限制，两者配合模拟出来一个沙盒的环境。

## Namespace

Linux 对线程提供了六种隔离机制，分别为：`uts`、 `pid`、 `user`、 `mount`、 `network`、 `ipc` ，它们的作用如下：

- `uts`: 用来隔离主机名
- `pid`: 用来隔离进程 PID 号
- `user`: 用来隔离用户
- `mount`: 用来隔离各个进程看到的挂载点视图
- `network`: 用来隔离网络
- `ipc`: 用来隔离 System V IPC 和 POSIX message queues

## CGroup

`Linux CGroup` 提供了对一组进程及子进程的资源限制，控制和统计的能力。这些资源包括 `CPU`、`内存`、`存储`、`网络`等。通过 `CGroup` 可以方便限制某个进程的资源占用，并且可以实时监控进程和统计信息。

CGroup 完成资源限制主要通过下面三个组件：

- `CGroup`: 是对进程分组管理的一种机制
- `subsystem`: 是一组资源控制的模块
- `hierarchy`: 把一组 `CGroup` 串成一个树状结构 (可让其实现继承)

> 主要实现方式是在 `/sys/fs/cgroup/` 文件夹下，根据限制的不同，创建一个新的文件夹即可，kernel 会将这个文件夹标记为它的 `子cgroup`
> 。比如要限制内存使用，则在 `/sys/fs/cgroup/memory/` 下创建 `test-limit-memory` 文件夹即可，将内存限制数写到该文件夹里面的 `memory.limit_in_bytes` 即可。

> [更多帮助](./HELP.md)

---

**Happy coding guys!!! :-)**
