# from指定基础镜像,必须是第一条指令;有一个特殊镜像叫scratch,表示一个空白镜像,完全从零开始
FROM nginx
# run 代表执行指令,有两种格式,另一种带方括号 RUN ["可执行文件","参数1","参数2"]
# 要特别注意,每一个指令都会创建一层,每一个RUN就相当于我们执行完一些指令后commit一次,因此应尽可能将安装
#指令放在同一个RUN命令之下
RUN echo '<h1>Im lerning how to build image!</h1>' /usr/share/nginx/html/index.html
# 并且每一个命令执行完毕后 一定要清除无关的内容
# 在本目录执行以下语句构建镜像
# docker build -t nginx:v3 .

# 其中的 . 要注意,是指的是上下文路径;docker实际是客户端,服务端交互模型,服务端也就是daemon引擎,保留了restAPI
# 供客户端(也就是docker命令)调用,实际执行build时是在daemon中构建的,那么服务器是如何获得本地的文件的呢?
# 这就引入了上下文的概念。当构建的时候，用户会指定构建镜像上下文的路径，docker build 命令得知这个路径后，会将路径下的所有内容打包，然后上传给 Docker 引擎。这样 Docker 引擎收到这个上下文包后，展开就会获得构建镜像所需的一切文件。
# 如 COPY ./package.json /app/ ,这条指令并非是要复制当前目录下的package.json文件进入镜像,而是值复制上下文目录中的文件

# -f可以指定Dockerfile的地址, . 可以更换为上下文目录,其指定的目录将会被打包传送给daemon,docker build中的操作指令都是基于此上下文路径