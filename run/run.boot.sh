#!/usr/bin/env bash

set -x -e
# PATH: ${PKG_FILE_SERVER}${NAMESPACE}/public/run/2.0/run.replace.sh


########################
#
#  Author: aoxn
#  Usage:
#     export ROLE=WORKER TOKEN=aacdeb.61c3abfd1eac6pbc INTRANET_LB=172.16.9.43; bash run.replace.sh
########################

# change to /root/ dir
WORKDIR=/root/
cd ${WORKDIR}

function detecos() {
    export OS=centos
    export ARCH=amd64
}

function validatedefault() {
    # detect os arch
    detecos

    if [[ "$REGION" == "" ]];
    then
        # https://github.com/koalaman/shellcheck/wiki/SC2155
        REGION="$(curl 100.100.100.200/latest/meta-data/region-id)"; export REGION
    fi
    if [[ "$ENDPOINT" == "" ]];
    then
        export ENDPOINT=http://127.0.0.1:32443
    fi
    if [[ "$NET_MASK" == "" ]];
    then
        export NET_MASK=25
    fi

    if [[ "$INTRANET_LB" == "" ]];
    then
        echo "INTRANET_LB must be specified"; exit 1
    fi

    if [[ "$INTERNET_LB" == "" ]];
    then
        echo "Warning: INTERNET_LB was not specified";
    fi

    if [[ "$CLUSTER_ID" == "" ]];
    then
        export CLUSTER_ID=kubernetes-clusterid-demo
    fi

    if [[ "$REGISTRY" == "" ]];
    then
        export REGISTRY=registry-vpc.$REGION.aliyuncs.com/acs
    fi

    if [[ "$ROLE" == "" ]];
    then
        echo "ROLE must be provided, one of BOOTSTRAP|MASTER|WORKER"; exit 1
    fi
    if [[ "$NAMESPACE" == "" ]];
    then
        NAMESPACE=default
    fi
    if [[ "$OOC_VERSION" == "" ]];
    then
        export OOC_VERSION=0.1.1
    fi

    if [[ "$CLOUD_TYPE" == "" ]];
    then
        export CLOUD_TYPE=public
    fi

    if [[ "$TOKEN" == "" ]];
    then
        echo "TOKEN must be provided"; exit 1
    fi

    if [[ "$PKG_FILE_SERVER" == "" ]];
    then
        PKG_FILE_SERVER="http://host-oc-$REGION.oss-$REGION-internal.aliyuncs.com"
        echo "empty PKG_FILE_SERVER, using default ${PKG_FILE_SERVER}"
    fi

    echo "using beta version: [${NAMESPACE}]"
    wget --tries 10 --no-check-certificate -q \
        -O /tmp/ooc.${ARCH}\
        ${PKG_FILE_SERVER}/ack/${NAMESPACE}/${CLOUD_TYPE}/ooc/${OS}/${ARCH}/ooc-${OOC_VERSION}.${ARCH}
    chmod +x /tmp/ooc.${ARCH} ; mv /tmp/ooc.${ARCH} /usr/local/bin/ooc
}


function config() {
    mkdir -p /etc/ooc
    if [[ -f /etc/ooc/ooc.cfg.gen ]];
    then
        cp -rf /etc/ooc/ooc.cfg.gen /etc/ooc/ooc.cfg
        echo "Found /etc/ooc/ooc.cfg.gen, use previous one"; return
    fi
    cat > /etc/ooc/ooc.cfg << EOF
clusterid: ${CLUSTER_ID}
iaas:
  image: abclid.vxd
  disk:
    size: 40G
    type: cloudssd
registry: ${REGISTRY}
namespace: ${NAMESPACE}
cloudType: ${CLOUD_TYPE}
kubernetes:
  name: kubernetes
  version: 1.16.9-aliyun.1
  kubeadmToken: "${TOKEN}"
etcd:
  name: etcd
  version: v3.3.8
  endpoints: "${ETCD_ENDPOINTS}"
runtime:
  name: docker
  version: 19.03.5
  para:
    key1: value
    key2: value2
sans:
  - 192.168.0.1
network:
  mode: ipvs
  podcidr: 172.16.0.1/16
  svccidr: 172.19.0.1/20
  domain: cluster.domain
  netMask: "${NET_MASK}"
endpoint:
  intranet: "${INTRANET_LB}"
  internet: "${INTERNET_LB}"
EOF
    # start in backgroud
}

function bootstrap() {
    echo run bootstrap init
    # run bootsrap init
    nohup ooc bootstrap --token "${TOKEN}" --bootcfg /etc/ooc/ooc.cfg &
}

function init() {
    echo run master init
    # run master init
    ooc init --role "${ROLE}" --token "${TOKEN}" --config /etc/ooc/ooc.cfg
}

function join() {
    echo run worker init
    # run master init
    ooc init --role Worker --token "${TOKEN}" --config /etc/ooc/ooc.cfg
}

function postcheck() {

    echo 'Check ROS notify server health, and notify to ROS notify server if its healthy.'
    set +e
    for ((i=1; i<=5; i ++));
    do
        cnt=$(curl -s http://100.100.100.110/health-check | grep ok | wc -l)
        echo "wait for ros notify server to be healthy cnt=$cnt, this is round $i"
        if curl -s http://100.100.100.110/health-check | grep ok ;
        then
            echo "the ros notify server is healthy"; break
        fi
        sleep 2
    done
    if ! curl -s http://100.100.100.110/health-check | grep ok ;
    then
        echo "wait for ros notify server to be healthy failed."; exit 2
    fi
    set -e
}

# validate default parameter first
validatedefault
config

case ${ROLE} in
    "Hybrid")
        echo "join master"
        init
    ;;
    "Master")
        echo "join master"
        init
    ;;
    "Worker")
        echo "join worker"
        join
    ;;
    "Bootstrap")
        echo "bootstrap master"
        bootstrap; init
    ;;
    *)
        echo "unrecognized role"
    ;;
esac

postcheck
