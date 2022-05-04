# FROM alpine
# RUN apk add curl
# CMD ["curl","-s","http://myip.ipip.net"]
# 构建一个用于显示当前ip的镜像,容器启动时就会执行CMD的命令,但无法修改run时的参数

FROM alpine
RUN apk add curl
ENTRYPOINT [ "curl","-s","http://myip.ipip.net"]