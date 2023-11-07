# ACES#JS——Web端短视频应用

- ACES-JS项目是基于restful接口规范开发的一款具有丰富样式和功能实现的Web端短视频应用实现基础的视频播放/切换/分类功能，以及进阶的账户系统（注册/登录/注销/修改/历史记录/消息通知/视频上传/视频删除/视频置顶）、推荐系统（用户侧个性化推荐/视频侧相关推荐/视频推流加热/流量充值/广告投放）、交互系统（点赞/分享/关注/评论/弹幕/分享/搜索）等多达五十余项功能、设计大小不等的共十几个界面。

#### Demo演示视频

https://www.bilibili.com/video/BV14G411Q7CN/

#### Demo网站

http://101.133.129.34/

### 使用方法

##### 克隆项目代码

```shell
git clone https://github.com/LabEnbug/ACES-JS.git
```

##### 启动后端

```shell
go version
# go version go1.21.3 windows/amd64
# 进入后端目录
cd ./ACES-JS/backend
# 启动后端，默认端口为8051
go run main.go
```

##### 启动前端

```shell
# 安装 npm 
# 参考：https://www.npmjs.com/
# 用以下命令查看npm版本，保证输出版本大于等于 10.2.1
npm -v
# 安装 yarn
npm install -g yarn
# 进入前端目录
cd ./ACES-JS/acro_frontend
# 安装依赖
yarn install
# 编译前端
yarn build
# 启动前端，假设端口设为3000
yarn start -p 3000

# 如果需要导出静态文件到nginx之类的运行，可通过该命令导出
yarn export
```

##### 打开Web前端

浏览器即可打开应用

```
http://localhost:3000/
```
