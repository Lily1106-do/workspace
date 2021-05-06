## 华链时代yardman服务用教程
### 介绍
```
yardman部署调度区块链服务，并提供标准的restful接口供上层服务调用
```
### 前提条件
```
1. 部署服务器需要安装python3
2. 需要安装kubectl， 并能访问kubernetes集群
3. 需要能访问网络存储（NFS）服务器
```

### 部署流程

#### 1. 从仓库拉取部署运行程序
```
git clone http://git.timechainer.com/blockchain/yardman-deploy.git
```
#### 2. 修改配置文件， 按实际的配置信息修改，示例如下
```
vim config.sh
```
![image.png](https://i.loli.net/2020/11/19/47WyUiN86qnHseE.png)

#### 3. 修改yardman-deplyment.yaml
```
vim yardman-deployment.yaml
```
##### 3.1 nfs 模块, 示例如下
![image.png](https://i.loli.net/2020/11/19/eJTMWkbDBCL31hH.png)

##### 3.2 k8s config模块, 示例如下
##### 把接入k8s的config文件内容贴到如下红色方框，yardman服务就可以访问控制k8s集群
![image.png](https://i.loli.net/2020/11/19/oju3hI5seVwprFK.png)


#### 4.启动yarman集群
bash create-cluster.sh
#### 5.删除yardman集群，及部署区块链集群
bash delete-cluster.sh