# !/bin/bash
source ./script/config.sh

###
EXPECT_INFO=${NFS_HOST}:${NFS_DIR}

# mount data
if [ ! -d ../data ]; then
   mkdir -p ../data
fi

MOUNT_INFO=$(df ../data -kh | cut -d" " -f 1 | tail -n 1)
if [ ${MOUNT_INFO} = ${EXPECT_INFO} ]; then 
   echo "The data is already mount."
else
  sudo mount ${EXPECT_INFO} ../data
fi
