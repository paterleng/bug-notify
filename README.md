# redmine通过钉钉进行消息通知
## 介绍
本项目独立于你的项目之外，不影响原系统的使用，通过使用go-mysql库来检测mysql binlog日志变化，触发go-mysql库中的事件回调来处理mysql表中的数据变化。
## 部署
通过Dockerfile打包成docker镜像，然后通过编写的yaml文件部署到kubernetes集群中运行
### kubernetes集群中部署
```
kubectl apply -f yaml文件路径
```
### docker部署
```
docker build
docker run 
```
### 二进制运行
```
cd 二进制包所在目录
./main
```

## 使用
在部署前你可以编辑你的conf.yaml文件来监控你的多个数据库中的多张表的变化，钉钉机器人的webhook地址及密钥需要定义一个自定义字段配置使用， 地址与密钥通过@符分隔，同时需要在项目层面配置自定义字段用于存储用户的手机号
## 源代码获取
```
cd existing_repo
git remote add origin http://192.168.10.6/liangjunhao/bug-notify.git
git branch -M master
git push -uf origin master
```
