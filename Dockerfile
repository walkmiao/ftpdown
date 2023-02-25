FROM alpine:latest
MAINTAINER lch

# 在容器根目录 创建一个 apps 目录
WORKDIR /apps

# 拷贝当前目录下 go_docker_demo1 可以执行文件
COPY dist/linux/backup /apps

# 拷贝配置文件到容器中
COPY config.yaml /apps/config.yaml

# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8

# 暴露端口
EXPOSE 5000

# 运行golang程序的命令
ENTRYPOINT ["/apps/backup"]