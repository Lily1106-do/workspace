#! /bin/bash

sudo apt-get install -y nfs-kernel-server
echo "/opt/share *(rw,sync,no_root_squash,no_subtree_check)" >> /etc/exports
# 增加权限
exportfs -rv 
sudo chmod 777 /opt/share -R
#重启服务
sudo /etc/init.d/rpcbind restart
sudo /etc/init.d/nfs-kernel-server restart
#显示目前共享目录
showmount -e
#测试挂载
# sudo mount -t nfs localhost:/opt/share /mntsudo mount 127.0.0.1:/opt/share /mnt