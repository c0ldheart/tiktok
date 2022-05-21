### 这个手册可能过时，有问题直接在群里问



1. 测试环境搭建：

   本地测试需要搭建 `mysql` 和 `redis`，可以从网上找教程，推荐使用 docker

   启动服务：需要进入到项目根目录下的 `cmd` 文件夹，执行 `go run .`

   ![image-20220521223054140](D:\misc\vscodego\h68u-tiktok-app\doc\imgs\image-20220521223054140.png)

   出现如上则说明服务已启动，上图所示，服务开在 localhost:8080

   可以使用 apifox 等接口测试工具，或者直接使用浏览器，在地址栏输入，如图应为:

   > localhost:8080/ping

   如果是浏览器，拿的结果应该是：

   ![image-20220521223352434](D:\misc\vscodego\h68u-tiktok-app\doc\imgs\image-20220521223352434.png)

   说明服务器启动成功

2. 参与项目（只是建议如此

   **先提一个 issue 说明自己打算负责的部分**

   将 github 上的项目 clone 到本地，并在本地建立自己的分支：

   > git checkout -b 分支名-dev
   >
   > 例如 git checkout -b slime-dev

   之后在自己的分支写就可以了，写完后执行：

   > git commit -am '这里是提交的信息，应言简意赅'

   这里可能会出现问题，及时在群里问

   然后将分支推送到项目仓库：

   > git push

   如果没有配置 ssh 可能会不成功，建议配置 ssh，当然，也可以多试几次，有一定概率成功 push

   成功后去 github 提一个 pr，建议找个人 review 下，如果把握十足或者没有项目结构上的大改动也可直接 merge

3.  有问题及时沟通x3