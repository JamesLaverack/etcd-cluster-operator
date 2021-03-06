package controllers

const (
	etcdClientPort    = 2379
	etcdPeerPort      = 2380
	etcdDataMountPath = "/var/lib/etcd"
	appName           = "etcd"
	appLabel          = "app.kubernetes.io/name"
	clusterLabel      = "etcd.improbable.io/cluster-name"
)
