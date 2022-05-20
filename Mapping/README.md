
# PortScan

## TODO

- [x] 接收上传文件并实现少量IPv4解析并构建任务

- [x] 解析CIDR和IP段为IP切片

- [ ] Node的超时退出机制,目前在master超时取消后,Node仍然会继续执行

- [ ] 在构建任务前,进行主机存活性检测,未存活主机不进行扫描

    - 使用PING(最简单,可靠性低)
    - 其他方法实现(ICMP)

- [ ] 超大量的IP,构建任务的方式必须采用缓冲的形式,即解析一定量IP就发送给节点开始扫描,扫描结果返回前不再解析
    
    - 3个节点的情况下,大概能承受多少? 任务数量为IP数量的Node倍;512个ip就需要发送1536个任务,总共会进行RPC调用1536次  
  站在node的角度,这一次扫描就被调用512次!假设1W个ip,每个node会被调用三千多次
        

### client
主要作为交互层,仅仅负责将任务发送至master,或展示结果
目前是作为http服务,用postman模拟上传扫描请求


### master
业务的核心,主要负责任务的分解,未来需要承担node节点的管理,对于client发送来的任务进行拆分并
发送至各个node节点,并处理返回的结果,后续可能拓展将结果存入数据库供client查询等功能

    TODO:优化超时控制机制:目前是写死的10s,如何根据任务数量(如发生重试 超时时间延迟)?





### node
主要的业务处理层,只负责处理master传入的任务包,执行完毕后返回结果,node节点间互相不感知,
可支持横向扩容

    TODO:node的并发控制:
    Node作为rpc服务提供节点,对每一次的请求G的数量限制为1W,但是还为做总体并发数量限制
    请求数量和单请求并发数量为乘积关系,考虑引入有效并发控制机制(第三方协程池Ants)
    TODO:node增加SYN扫描功能


## Deployment
#### 1.安装部署consul
项目以consul为服务发现,各个节点会注册至consul,master会从consul中获取在线的node节点执行任务分发,
因此,需要先部署consul(主机暴露8500端口)
```bash
consul agent -dev -ui=true -client 0.0.0.0
```

#### 2.编译master和node
cmd目录下分别进入node和master目录运行:

```bash
#windows下请交叉编译 SET GOOS=linux
  go build main.go
  chmod +x main
```

#### 3.在各节点运行master或node
master和consul可以在同一台主机,node建议分布式部署;
node需将自己的公网ip注册到consul(便于调用),因此要确保consul,node对应的端口处于可访问状态

运行node:
```bash
./main -s 0.0.0.0:8081(本机监听端口) -c (consul的主机地址和端口) -l (本机公网ip端口)
```

#### 4.使用api调试工具测试
使用postman或apifox,对master的地址发起post请求:
POST

    http://144.34.167.185:8080/v1/scanport

```bash
{
"targets":["8.8.8.8"],
"scanType:"top1000"

}
```

