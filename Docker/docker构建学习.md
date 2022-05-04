# Docker指令学习

## 1.Dockerfile注意事项

### 1.1 dockerfile的最基本的格式

```docker
FROM busybox
COPY imfromCTX.txt .
COPY html/hello.html .
```

之后执行dockerbuild来根据dockerfile来构建镜像:

`docker build -t hello:v1 .`

`.` 代表将当前目录作为上下文目录传给daemon,可以替换为路径;`-t`代表指定镜像名字和版本;`-f`指定Dockerfile的路径;  
执行build时不加`-f`是默认从传入的上下文目录中寻找Dockerfile文件的;  
Dockerfile中的所有操作都是**基于传入的上下文目录的**

### 1.2 Dokcerfile命令

#### COPY和ADD

格式:  

- `COPY [--chown=<user>:<group>] <源路径>... <目标路径>`

- `COPY [--chown=<user>:<group>] ["<源路径1>",... "<目标路径>"]`

都可以用于从上下文目录中拷贝文件到镜像中,COPY是单纯的拷贝,而ADD则增加更多功能,**ADD可以在拷贝压缩文件时,将会自动解压**;  
当copy的是一个目录时,只会拷贝该目录的所有文件;目标路径可以是镜像的绝对路径,也可以是相对于工作目录的路径(WORKDIR指定)  

#### CMD 容器启动命令

CMD和RUN类似,同样有2种语法格式:  

- `CMD <命令>`

- `CMD ["命令","参数1","参数2"]`  

- `CMD ["参数1","参数2"]` 这种方式是在用ENTRYPOINT指定命令后, 用CMD来指定参数  

Docker的容器不是虚拟机,而是进程,进程在运行时就需要指定要运行的程序和参数,就像Ubuntu镜像,默认的指令是/bin/bash;**CMD指令就是用于指定默认的容器主进程的启动命令的。**  
指令格式上更推荐第二种,此种会被解析为json格式,命令和参数一定要用"包裹好,而不是单引号;  

##### 容器和虚拟机

提到 CMD 就不得不提容器中应用在前台执行和后台执行的问题。这是初学者常出现的一个混淆。
Docker 不是虚拟机，容器中的应用都应该以前台执行，而不是像虚拟机、物理机里面那样，用 systemd 去启动后台服务，容器内没有后台服务的概念。  
容器是为了主进程而存在的,主进程都退出了容器也没有存在的意义;因此容器内的程序一定要以前台形式存在!  
如:

`CMD service nginx start`

会发现启动容器时马上就会退出;使用 service nginx start 命令，则是希望 upstart 来以后台守护进程形式启动 nginx 服务。而刚才说了 CMD service nginx start 会被理解为 CMD [ "sh", "-c", "service nginx start"]，因此主进程实际上是 sh。那么当 service nginx start 命令结束后，sh 也就结束了，sh 作为主进程退出了，自然就会令容器退出。  
正确的做法是之间运行nginx(使用第二种格式):

`CMD ["nginx", "-g", "daemon off;"]`

***CMD和RUN的区别?***

RUN是指定Dockerbuild中所需要执行的命令,而CMD是在docker run启动容器时执行;CMD的命令是可以被docker run时指定的要运行的程序的命令所**覆盖**的;如果存在多个CMD,则仅最后一个会生效!

#### ENTRYPOINT 入口点

ENTRYPOINT 的目的和 CMD 一样，都是在指定容器启动程序及参数。ENTRYPOINT 在运行时也可以替代，不过比 CMD 要略显繁琐，需要通过 docker run 的参数 --entrypoint 来指定。  

##### 场景一:将镜像作为命令使用

如,构建一个显示当前ip的镜像:

```docker
FROM alpine
RUN apk add curl
CMD ["curl","-s","http://myip.ipip.net"]
```

在docker run时可直接显示结果:

``` bash
➜  nginx git:(main) ✗ docker run --rm showip:v1
当前 IP：117.176.187.27  来自于：中国 四川 成都  移动
```

但是如果我们要修改参数怎么办呢,比如显示头信息,加入-h参数?此时只能修改Dockerfile,显然不是一个好方法,**而使用ENTRYPOINT就可以解决这个问题**

```docker
FROM alpine
RUN apk add curl
ENTRYPOINT [ "curl","-s","http://myip.ipip.net"]
```

执行`dokcer run --rm showip:v2 -i`

会发现成功显示了头信息,这是因为,**当CMD存在时(-i就是CMD命令),会将CMD作为参数传入ENTRYPOINT**

##### 场景二 容器运行前的准备工作

启动容器就是启动进程,但可能会有额外的一些准备工作,如启动mysql前配置数据库等,这些肯定是在最终运行前准备的;  
此外，可能希望避免使用 root 用户去启动服务，从而提高安全性，而在启动服务前还需要以 root 身份执行一些必要的准备工作，最后切换到服务用户身份启动服务。或者除了服务外，其它命令依旧可以使用 root 身份执行，方便调试等。  
这些准备工作是和容器 CMD 无关的，无论 CMD 为什么，都需要事先进行一个预处理的工作。这种情况下，可以写一个脚本，然后放入 ENTRYPOINT 中去执行，而这个脚本会将接到的参数（也就是 CMD）作为命令，在脚本最后执行  
如官方Redis镜像是这么写的:

```docker
FROM alpine:3.4
...
RUN addgroup -S redis && adduser -S -G redis redis
...
ENTRYPOINT ["docker-entrypoint.sh"]

EXPOSE 6379
CMD [ "redis-server" ]
```

脚本的内容:

```bash
#!/bin/sh
...
# allow the container to be started with `--user`
if [ "$1" = 'redis-server' -a "$(id -u)" = '0' ]; then
 find . \! -user redis -exec chown redis '{}' +
 exec gosu redis "$0" "$@"
fi

exec "$@"
```

该脚本内容就是根据CMD的内容来判断,如果是 redis-server 的话，则切换到 redis 用户身份启动服务器，否则依旧使用 root 身份执行

#### ENV设置环境变量

格式:

- `ENV key value`
  
- `ENV key1=value1 key2=value2...`
  
在设置环境变量后,后续的指令都可以使用;例:

```docker
ENV NODE_VERSION 7.2.0

RUN curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.xz" \
  && curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/SHASUMS256.txt.asc" \
  && gpg --batch --decrypt --output SHASUMS256.txt SHASUMS256.txt.asc \
  && grep " node-v$NODE_VERSION-linux-x64.tar.xz\$" SHASUMS256.txt | sha256sum -c - \
  && tar -xJf "node-v$NODE_VERSION-linux-x64.tar.xz" -C /usr/local --strip-components=1 \
  && rm "node-v$NODE_VERSION-linux-x64.tar.xz" SHASUMS256.txt.asc SHASUMS256.txt \
  && ln -s /usr/local/bin/node /usr/local/bin/nodejs
```

RUN中多次使用了环境变量`NODE_VERSION`,这样在版本发生改变时,只需修改ENV命令中的数值,这样的话维护更容易  
下列指令可以支持环境变量展开： `ADD、COPY、ENV、EXPOSE、FROM、LABEL、USER、WORKDIR、VOLUME、STOPSIGNAL、ONBUILD、RUN`。
可以从这个指令列表里感觉到，环境变量可以使用的地方很多，很强大。通过环境变量，我们可以让一份 Dockerfile 制作更多的镜像，只需使用不同的环境变量即可。

#### ARG 构建参数

构建参数和 ENV 的效果一样，都是设置环境变量。所不同的是，ARG 所设置的构建环境的环境变量，在将来容器运行时是不会存在这些环境变量的。但是不要因此就使用 ARG 保存密码之类的信息，因为 docker history 还是可以看到所有值的。
Dockerfile 中的 ARG 指令是定义参数名称，以及定义其默认值。该默认值可以在构建命令 docker build 中用 --build-arg <参数名>=<值> 来覆盖。
灵活的使用 ARG 指令，能够在不修改 Dockerfile 的情况下，构建出不同的镜像。  
ARG是有作用域的:

```docker
ARG DOCKER_USERNAME=library

FROM ${DOCKER_USERNAME}/alpine

# 在FROM 之后使用变量，必须在每个阶段分别指定
ARG DOCKER_USERNAME=library

RUN set -x ; echo ${DOCKER_USERNAME}

FROM ${DOCKER_USERNAME}/alpine

# 在FROM 之后使用变量，必须在每个阶段分别指定
ARG DOCKER_USERNAME=library

RUN set -x ; echo ${DOCKER_USERNAME}
```

#### VOLUME 定义匿名卷

作用:**创建一个匿名数据卷挂载点**  
**容器运行时应该尽量保持容器存储层不发生写操作，对于数据库类需要保存动态数据的应用，其数据库文件应该保存于卷(volume)中**，为了防止运行时用户忘记将动态文件所保存目录挂载为卷，在 Dockerfile 中，我们可以事先指定某些目录挂载为匿名卷，这样在运行时如果用户不指定挂载，其应用也可以正常运行，不会向容器存储层写入大量数据。

```docker
VOLUME /data
```

这里的 /data 目录就会在容器运行时自动挂载为匿名卷，任何向 /data 中写入的信息都不会记录进容器存储层，从而保证了容器存储层的无状态化。当然，**运行容器时可以覆盖这个挂载设置**。比如：

```docker
docker run -d -v mydata:/data xxxx
```

使用了 mydata 这个**命名卷**挂载到了 /data 这个位置，替代了 Dockerfile 中定义的匿名卷的挂载配置。

#### EXPOSE 暴露端口

只是声明,而不是直接使用该端口,有2个好处:  

一个是帮助镜像使用者理解这个镜像服务的守护端口，以方便配置映射；另一个用处则是在运行时使用随机端口映射时，也就是 docker run -P 时，会自动随机映射 EXPOSE 的端口。  

要将 EXPOSE 和在运行时使用 -p <宿主端口>:<容器端口> 区分开来。-p，是映射宿主端口和容器端口，换句话说，就是将容器的对应端口服务公开给外界访问，**而 EXPOSE 仅仅是声明容器打算使用什么端口而已，并不会自动在宿主进行端口映射**。

#### WORKDIR 指定工作目录

使用 WORKDIR 指令可以来指定工作目录（或者称为当前目录），以后各层的当前目录就被改为指定的目录，如该目录不存在，WORKDIR 会帮你建立目录。

因为每一个指令都是完全新的一层,其上一层状态不会保留(如cd到某个目录),因此如果需要改变以后各层的工作目录的位置，那么应该使用 WORKDIR 指令。

#### USER 指定当前用户

USER 指令和 WORKDIR 相似，都是改变环境状态并影响以后的层。WORKDIR 是改变工作目录，USER 则是改变之后层的执行 RUN, CMD 以及 ENTRYPOINT 这类命令的身份。
注意，USER 只是帮助你切换到指定用户而已，这个用户必须是事先建立好的，否则无法切换。

#### HEALTHCHECK 健康检查

自 1.12 之后，Docker 提供了 HEALTHCHECK 指令，通过该指令指定一行命令，用这行命令来判断容器主进程的服务状态是否还正常，从而比较真实的反应容器实际状态。
当在一个镜像指定了 HEALTHCHECK 指令后，用其启动容器，初始状态会为 starting，在 HEALTHCHECK 指令检查成功后变为 healthy，如果连续一定次数失败，则会变为 unhealthy.  

假设我们有个镜像是个最简单的 Web 服务，我们希望增加健康检查来判断其 Web 服务是否在正常工作，我们可以用 curl 来帮助判断，其 Dockerfile 的 HEALTHCHECK 可以这么写：

```docker
FROM nginx
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
HEALTHCHECK --interval=5s --timeout=3s \
  CMD curl -fs http://localhost/ || exit 1
```

这里我们设置了**每 5 秒检查一次**（这里为了试验所以间隔非常短，实际应该相对较长），如果健康检查**命令超过 3 秒没响应就视为失败**，并且使用 curl -fs http://localhost/ || exit 1 作为健康检查命令。

#### ONBUILD

此命令较特殊,它后面跟的是其它指令，比如 RUN, COPY 等，而这些指令，在当前镜像构建时并不会被执行。只有当以当前镜像为基础镜像，去构建下一级镜像的时候才会被执行。

```docker
FROM node:slim
RUN mkdir /app
WORKDIR /app
ONBUILD COPY ./package.json /app
ONBUILD RUN [ "npm", "install" ]
ONBUILD COPY . /app/
CMD [ "npm", "start" ]
```

#### LABEL 为镜像添加元数据

LABEL 指令用来给镜像以键值对的形式添加一些元数据（metadata）。

`LABEL <key>=<value> <key>=<value> <key>=<value> ...`

可用于申明作者,文档地址等:

```docker
LABEL org.opencontainers.image.authors="yeasy"

LABEL org.opencontainers.image.documentation="https://yeasy.gitbooks.io"
```

#### SHELL 指令

SHELL 指令可以指定 RUN ENTRYPOINT CMD 指令的 shell，Linux 中默认为 ["/bin/sh", "-c"]