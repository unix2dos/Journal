### 安装启动

``` 
yum -y install docker 安装docker
service docker start 启动docker
docker run hello-world 测试安装成功
```

### 基本命令

```
docker ps 			//查看正在运行的容器  -a显示所有的容器

docker start con_name //启动容器名为con_name的容器(已终止的)。
docker stop con_name //停止容器名为con_name的容器。
docker rm con_name //删除容器名为con_name的容器。
docker rename old_name new_name //重命名一个容器。

docker images 		//列出镜像
docker image rm     //删除镜像
docker search mysql //查找镜像
docker pull mysql   //拉取镜像
docker run mysql  	//运行镜像容器
docker rmi img_name //删除镜像名为img_name的镜像。

docker exec -it webserver bash //进入容器

```



### 获取镜像

>docker pull [选项] [Docker Registry 地址[:端口号]/]仓库名[:标签]

也可以不pull,直接run, 没有会自动pull, 也可以 docker seath xxx查找镜像


### 运行容器
> docker run -it --rm ubuntu:16.04 bash

--rm 是运行完毕就删除容器

--it i是交互式, t是终端

bash 是启动容器运行的命令


### commit镜像

>docker commit [选项] <容器ID或容器名> [<仓库名>[:<标签>]]

```
docker commit \
    --author "levonfly@gmail.com" \
    --message "修改了默认网页" \
    webserver \
    nginx:v2
```

一般情况下不要使用commit


# DockerFile

 
### FROM 来源于什么

而 FROM 就是指定基础镜像，因此一个 Dockerfile 中 FROM 是必备的指令，并且必须是第一条指令。

在 Docker Store 上有非常多的高质量的官方镜像，有可以直接拿来使用的服务类的镜像，如 nginx、redis、mongo、mysql、httpd、php、tomcat 等

除了选择现有镜像为基础镜像外，Docker 还存在一个特殊的镜像，名为 scratch。这个镜像是虚拟的概念，并不实际存在，它表示一个空白的镜像。


### RUN  是在镜像前执行

RUN 指令是用来执行命令行命令的。


##### shell 格式 RUN <命令>

```
RUN echo "<h1>Hello, Docker!</h1>" > 
/usr/share/nginx/html/index.html
```
##### exec 格式：RUN ["可执行文件", "参数1", "参数2"]，这更像是函数调用中的格式。


---

Dockerfile 中每一个指令都会建立一层，RUN 也不例外。每一个 RUN 的行为，就和刚才我们手工建立镜像的过程一样：新建立一层，在其上执行这些命令，执行结束后，commit 这一层的修改，构成新的镜像。


```
FROM debian:jessie

RUN apt-get update
RUN apt-get install -y gcc libc6-dev make
RUN wget -O redis.tar.gz "http://download.redis.io/releases/redis-3.2.5.tar.gz"
RUN mkdir -p /usr/src/redis
RUN tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1
RUN make -C /usr/src/redis
RUN make -C /usr/src/redis install
```
上面的这种写法，创建了 7 层镜像。这是完全没有意义的

```
FROM debian:jessie

RUN buildDeps='gcc libc6-dev make' \
    && apt-get update \
    && apt-get install -y $buildDeps \
    && wget -O redis.tar.gz "http://download.redis.io/releases/redis-3.2.5.tar.gz" \
    && mkdir -p /usr/src/redis \
    && tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1 \
    && make -C /usr/src/redis \
    && make -C /usr/src/redis install \
    && rm -rf /var/lib/apt/lists/* \
    && rm redis.tar.gz \
    && rm -r /usr/src/redis \
    && apt-get purge -y --auto-remove $buildDeps
```

在撰写 Dockerfile 的时候，要经常提醒自己，这并不是在写 Shell 脚本，而是在定义每一层该如何构建。

### COPY 复制文件
和 RUN 指令一样，也有两种格式，一种类似于命令行，一种类似于函数调用。

```
COPY package.json /usr/src/app
```
<目标路径> 可以是容器内的绝对路径，也可以是相对于工作目录的相对路径（工作目录可以
用 WORKDIR 指令来指定）。

### ADD 更高级的复制文件

ADD 指令和 COPY 的格式和性质基本一致。但是在 COPY 基础上增加了一些功能。

比如 <源路径> 可以是一个 URL ，这种情况下，Docker 引擎会试图去下载这个链接的文件放到 <目标路径> 去。

如果下载的是个压缩包，需要解压缩，也一样
还需要额外的一层 RUN 指令进行解压缩。所以不如直接使用 RUN 指令，然后使用 wget 或者 curl 工具下载，处理权限、解压缩、然后清理无用文件更合理。因此，这个功能其实并不实用，而且不推荐使用。

### CMD 容器启动命令 (容器启动后执行命令, 只能一个cmd,)

CMD 指令的格式和 RUN 相似，也是两种格式：


CMD 指令就是用于指定默认的容器主进程的启动命令的。


ubuntu 镜像默认的CMD 是 /bin/bash ，如果我们直接 docker run -it ubuntu 的话，会直接进入 bash 。



Docker 不是虚拟机，容器中的应用都应该以前台执行，而不是像虚拟机、物理机里面那样，
用 upstart/systemd 去启动后台服务，容器内没有后台服务的概念。


`CMD service nginx start` 发现容器执行后就立即退出了。

使用 service nginx start 命令，则是希望 upstart 来以后台守护进程形式启动 nginx 服务。而刚才说了 CMD service nginx start 会被理解为 CMD [ "sh", "-c", "service nginx start"] ，因此主进程实际上是 sh 。那么当 service nginx start 命令结束后， sh 也就结束了， sh 作为主进程退出了，自然就会令容器退出。

正确的做法 `CMD ["nginx", "-g", "daemon off;"]`

### ENTRYPOINT 入口点
ENTRYPOINT 的格式和 RUN 指令格式一样，分为 exec 格式和 shell 格式。

ENTRYPOINT 的目的和 CMD 一样，都是在指定容器启动程序及参数。 ENTRYPOINT 在运行时也可以替代，不过比 CMD 要略显繁琐，需要通过 docker run 的参数 --entrypoint 来指定。

+ 传递参数

```
 CMD [ "curl", "-s", "http://ip.cn" ]
 docker run myip -i //错误
 
 ENTRYPOINT [ "curl", "-s", "http://ip.cn" ]
 docker run myip -i //正确
```

+ 执行脚本

可以写一个脚本，然后放入 ENTRYPOINT 中去执行，而这个脚本会将接到
的参数（也就是 <CMD> ）作为命令，在脚本最后执行。


### ENV 设置环境变量
格式有两种：

```
ENV <key> <value>

ENV <key1>=<value1> <key2>=<value2>...
```
下列指令可以支持环境变量展开：
ADD 、 COPY 、 ENV 、 EXPOSE 、 LABEL 、 USER 、 WORKDIR 、 VOLUME 、 STOPSIGNAL 、 ONBUILD 。

### ARG 构建参数

```
格式： ARG <参数名>[=<默认值>]
```
构建参数和 ENV 的效果一样，都是设置环境变量。所不同的是， ARG 所设置的构建环境的
环境变量，在将来容器运行时是不会存在这些环境变量的。

### VOLUME 定义匿名卷
格式为：

```
格式为：
VOLUME ["<路径1>", "<路径2>"...]
VOLUME <路径>
```
容器运行时应该尽量保持容器存储层不发生写操作，对于数据库类需要保存
动态数据的应用，其数据库文件应该保存于卷(volume)中，


`VOLUME /data` 这里的 /data 目录就会在运行时自动挂载为匿名卷，任何向 /data 中写入的信息都不会记录进容器存储层，从而保证了容器存储层的无状态化。


### EXPOSE 声明端口

```
格式为 EXPOSE <端口1> [<端口2>...] 。
```

EXPOSE 指令是声明运行时容器提供服务端口，这只是一个声明，在运行时并不会因为这个声
明应用就会开启这个端口的服务。

在 Dockerfile 中写入这样的声明有两个好处，一个是帮助镜像使用者理解这个镜像服务的守护端口，以方便配置映射；另一个用处则是在运行时使用随机端口映射时，也就是 docker run -P 时，会自动随机映射 EXPOSE 的端口。


要将 EXPOSE 和在运行时使用 -p <宿主端口>:<容器端口> 区分开来。 -p ，是映射宿主端口和
容器端口，换句话说，就是将容器的对应端口服务公开给外界访问，而 EXPOSE 仅仅是声明
容器打算使用什么端口而已，并不会自动在宿主进行端口映射。


docker run -p 81:80 使用宿主的81端口 容器是80端口 file里 `EXPOSE 80`


### WORKDIR 指定工作目录

格式为 `WORKDIR <工作目录路径> `


使用 WORKDIR 指令可以来指定工作目录（或者称为当前目录），以后各层的当前目录就被改
为指定的目录，如该目录不存在， WORKDIR 会帮你建立目录。


```
RUN cd /app
RUN echo "hello" > world.txt
```

如果将这个 Dockerfile 进行构建镜像运行后，会发现找不到 /app/world.txt 文件，或者其内容不是 hello 。原因其实很简单，在 Shell 中，连续两行是同一个进程执行环境，因此前一个命令修改的内存状态，会直接影响后一个命令；而在 Dockerfile 中，这两行 RUN 命令的执行环境根本不同，是两个完全不同的容器。这就是对 Dockerfile 构建分层存储的概念不了解所导致的错误。


### USER 指定当前用户

格式： `USER <用户名>`

USER 指令和 WORKDIR 相似，都是改变环境状态并影响以后的层。 WORKDIR 是改变工作目录， USER 则是改变之后层的执行 RUN , CMD 以及 ENTRYPOINT 这类命令的身份。

### HEALTHCHECK 健康检查

```
HEALTHCHECK [选项] CMD <命令> ：设置检查容器健康状况的命令
HEALTHCHECK NONE ：如果基础镜像有健康检查指令，使用这行可以屏蔽掉其健康检查指
令
```

HEALTHCHECK 支持下列选项：

```
--interval=<间隔> ：两次健康检查的间隔，默认为 30 秒；
--timeout=<时长> ：健康检查命令运行超时时间，如果超过这个时间，本次健康检查就被
视为失败，默认 30 秒；
--retries=<次数> ：当连续失败指定次数后，则将容器状态视为 unhealthy ，默认 3
次。
```
和 CMD , ENTRYPOINT 一样， HEALTHCHECK 只可以出现一次，如果写了多个，只有最后一个生效。


### ONBUILD 为他人做嫁衣裳

格式： `ONBUILD <其它指令> 。`

ONBUILD 是一个特殊的指令，它后面跟的是其它指令，比如 RUN , COPY 等，而这些指令，在当前镜像构建时并不会被执行。只有当以当前镜像为基础镜像，去构建下一级镜像的时候才会被执行。

# 构建镜像

>docker build [选项] <上下文路径/URL/->

```
docker build -t nginx:v3 .
```
如果注意，会看到 docker build 命令最后有一个 .。. 表示当前目录，而 Dockerfile 就在当前目录，因此不少初学者以为这个路径是在指定 Dockerfile 所在路径，这么理解其实是不准确的


+ 上下文

docker build 命令构建镜像，其实并非在本地构建，而是在服务端，也就是 Docker 引擎中构建的。那么在这种客户端/服务端的架构中，如何才能让服务端获得本地文件呢？

docker build 命令得知这个路径后，会将路径下的所有内容打包，然后上传给 Docker 引擎。这样 Docker 引擎收到这个上下文包后，展开就会获得构建镜像所需的一切文件。

```
COPY ./package.json /app/	
```
这并不是要复制执行 docker build 命令所在的目录下的 package.json，也不是复制 Dockerfile 所在目录下的 package.json，而是复制 上下文（context） 目录下的 package.json。


COPY 这类指令中的源文件的路径都是相对路径。这也是初学者经常会问的为什么 COPY ../package.json /app 或者 COPY /opt/xxxx /app 无法工作的原因，因为这些路径已经超出了上下文的范围，Docker 引擎无法获得这些位置的文件。如果真的需要那些文件，应该将它们复制到上下文目录中去。

一般来说，应该会将 Dockerfile 置于一个空目录下，或者项目根目录下。如果该目录下没有所需文件，那么应该把所需文件复制一份过来。如果目录下有些东西确实不希望构建时传给 Docker 引擎，那么可以用 .gitignore 一样的语法写一个 .dockerignore，该文件是用于剔除不需要作为上下文传递给 Docker 引擎的。


+ 从github构建

```
docker build https://github.com/twang2218/gitlab-ce-zh.git#:8.14
```
+ 从tar构建

```
docker build http://server/context.tar.gz
```


# 容器

### 启动容器

启动容器有两种方式，一种是基于镜像新建一个容器并启动，另外一个是将在终止状态
（ stopped ）的容器重新启动。
因为 Docker 的容器实在太轻量级了，很多时候用户都是随时删除和新创建容器。


+ 当利用 docker run 来创建容器时，Docker 在后台运行的标准操作包括：

```
检查本地是否存在指定的镜像，不存在就从公有仓库下载
利用镜像创建并启动一个容器
分配一个文件系统，并在只读的镜像层外面挂载一层可读写层
从宿主主机配置的网桥接口中桥接一个虚拟接口到容器中去
从地址池配置一个 ip 地址给容器
执行用户指定的应用程序
执行完毕后容器被终止
```


+ 后台运行容器  加上-d

```
docker run -d ubuntu:17.10 /bin/sh -c "while true; do echo hello world; sleep 1; don
e"
```

消息通过 docker logs 可以查看

### 进入容器

+ docker attach
+  docker exec 

	-i -t 参数 docker exec 后边可以跟多个参数，这里主要说明 -i -t 参数。只用 -i 参数时，由于没有分配伪终端，界面没有我们熟悉的 Linux 命令提示符，但命令执行结果仍然可以返回。当 -i -t 参数一起使用时，则可以看到我们熟悉的 Linux 命令提示符。


### 导出和导入容器

+ 导出容器

```
	docker export 7691a814370e > ubuntu.tar
```

+ 导入容器快照
可以使用 docker import 从容器快照文件中再导入为镜像

```
cat ubuntu.tar | docker import - test/ubuntu:v1.0
```

用户既可以使用 docker load 来导入镜像存储文件到本地镜像库，也可以使用 docker
import 来导入一个容器快照到本地镜像库。这两者的区别在于容器快照文件将丢弃所有的历
史记录和元数据信息（即仅保存容器当时的快照状态），而镜像存储文件将保存完整记录，
体积也要大。此外，从容器快照文件导入时可以重新指定标签等元数据信息。


### 删除容器

`docker container rm trusting_newton`


### 清理所有处于终止状态的容器

`docker container prune`


# 仓库 TODO:
仓库（ Repository ）是集中存放镜像的地方。


# TODO: Docker-Compose



# 运行

### mysql

```
//运行
docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=12345 -d mysql

//连接到mysql
docker run -it --link some-mysql:mysql --rm mysql sh -c 'exec mysql -h"$MYSQL_PORT_3306_TCP_ADDR" -P"$MYSQL_PORT_3306_TCP_PORT" -uroot -p"$MYSQL_ENV_MYSQL_ROOT_PASSWORD"'

//进入bash 连接mysql
docker exec -it some-mysql bash
mysql -uroot -p
```






