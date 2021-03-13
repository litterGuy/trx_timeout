
#简介

检测连接trx官方节点rpc的连接成功率，以及trx其他操作的测试代码

#启动

```shell script
nohup ./main &
```

#停止

```shell script
给父进程发送一个TERM信号，试图杀死它和它的子进程。
# kill -TERM PPID
```