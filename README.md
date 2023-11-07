# ACES#JS——Web端短视频应用

- ACES-JS项目是基于restful接口规范开发的一款具有丰富样式和功能实现的Web端短视频应用，实现基础的视频播放/切换/分类功能，以及进阶的账户系统（注册/登录/注销/修改/历史记录/消息通知/视频上传/视频删除/视频置顶）、推荐系统（用户侧个性化推荐/视频侧相关推荐/视频推流加热/流量充值/广告投放）、交互系统（点赞/分享/关注/评论/弹幕/分享/搜索）等多达x项功能、设计大小不等的共y个界面。

## Install

克隆项目代码

```shell
git clone https://github.com/LabEnbug/ACES-JS.git
```

#### 启动后端

```shell
# make sure Golang installed
go

cd ./ACES-JS/backend
go run main.go
```

#### 启动前端

```shell
# 安装 npm 
参考：https://www.npmjs.com/
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
# 导出前端文件
yarn export
# 运行
yarn start -p XXXX(port)
```

#### 打开Web应用

浏览器即可打开应用

```
http://localhost:XXXX(port)/
```

#### Demo演示视频