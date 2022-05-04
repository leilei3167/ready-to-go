FROM busybox
COPY imfromCTX.txt .
COPY html/hello.html .
# 如果执行docker build -t hello:v1 context 将报错,因为Dockerfile不在上下文目录(context)中
# 执行docker build时,如果未指定dockerfile目录,将默认为上下文目录中寻找dockerfile

# -f指定Dockerfile的位置,docker build -f Dockerfile -t hello:v1 context
# 因此Dockerfile可以不在上下文目录中;构建镜像时会先将上下文目录发送给daemon;

# 如果以当前目录作为上下文目录发送 docker build -f Dockerfile -t hello:v2 .
# 将无法正确的copy,因为copy指令是基于传入的上下文目录,而当前目录下找不到该文件