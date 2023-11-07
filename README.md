# ACES#JS——Web端短视频应用

- ACES-JS项目是基于restful接口规范开发的一款具有丰富样式和功能实现的Web端短视频应用，实现基础的视频播放/切换/分类功能，以及进阶的账户系统（注册/登录/注销/修改/历史记录/消息通知/视频上传/视频删除/视频置顶）、推荐系统（用户侧个性化推荐/视频侧相关推荐/视频推流加热/流量充值/广告投放）、交互系统（点赞/分享/关注/评论/弹幕/分享/搜索）等多达x项功能、设计大小不等的共y个界面。

## Install

克隆项目代码

```shell
git clone https://github.com/LabEnbug/ACES-JS.git
```

#### 安装Golang开发环境

1. 下载Golang安装包：
   访问Golang官方网站（https://golang.org/dl/）并下载适用于Windows的最新版本的Golang安装包。选择与操作系统架构（32位或64位）相对应的安装包。

2. 运行安装程序：
   打开下载的安装包，并按照安装向导的指示进行操作。默认情况下，Golang将安装到"C:\Go"目录下，建议选择其他目录。

3. 配置环境变量：
   在Windows 10上，需要配置环境变量，以便在命令提示符或PowerShell中使用Golang。以下是配置环境变量的步骤：

   - 右键点击"此电脑"（或"我的电脑"），然后选择"属性"。
   - 在打开的窗口中，点击"高级系统设置"。
   - 在"系统属性"窗口中，点击"环境变量"按钮。
   - 在"用户变量"部分，点击"新建"按钮。
   - 在"变量名"字段中输入"GOROOT"，在"变量值"字段中输入Golang的安装路径（例如："C:\Go"）。
   - 在"用户变量"部分，找到名为"Path"的变量，并点击"编辑"按钮。
   - 在"编辑环境变量"窗口中，点击"新建"按钮。
   - 在"变量值"字段中输入"%GOROOT%\bin"，然后点击"确定"按钮。
   - 关闭所有打开的窗口。

4. 验证安装：
   win+R输入cmd回车打开命令提示符，或打开PowerShell，运行以下命令来验证Golang是否正确安装：

   ```shell
   go version
   ```

   如果一切正常，你将看到Golang的版本信息。

#### 启动后端

```shell
go version
# go version go1.21.3 windows/amd64

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