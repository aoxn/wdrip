## Autopilot your kubernetes cluster

k8s提供了云原生时代下应用运行的标准环境，统一了应用管理的各个环节，但管理好k8s集群却是一件极其复杂的事情，如何初始化一个k8s集群（day 0 问题），如何运维一个k8s集群（day 1问题），集群故障了怎么办，高可用怎么做，节点容量不够了怎么办，如何配置集群才能达到性能最优，如何观测集群内部发生了哪些事件，如何监控系统运行状态？这些都是分布式系统中最为常见也是非常难解决的问题，需要投入昂贵的人力物力才能换来系统的平稳运行。

wdrip以一种更加高效且低成本的方式实现了k8s集群的自动驾驶能力，无人工干预下的k8s集群自治，如容量的自动配置，故障监测与自愈，组件自运维等等。 自动根据需要调节系统容量，如果集群陷入故障，能自动自我检测，并从故障中自行恢复，尤其是在面临灾难的情况下，仍然能够自我修复，全程甚至可以无需人工干预。wdrip实现了面向Cattle的运维，节点故障了，运维人员需要做的仅仅是干掉故障节点即可，系统会自动使用全新的节点补齐，这极大的缩短了从故障中恢复的时间。

[去试一试](./manage-cluster.md)