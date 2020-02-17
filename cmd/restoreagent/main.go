package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	pb "github.com/improbable-eng/etcd-cluster-operator/api/proxy"
	"google.golang.org/grpc"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/otiai10/copy"
	"github.com/spf13/pflag"
	"go.etcd.io/etcd/clientv3/snapshot"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

func main() {

	etcdPeerName := pflag.String("etcd-peer-name",
		"",
		"Name of this peer, must match the name of the eventual peer")
	fmt.Printf("Using etcd peer name %s\n", *etcdPeerName)

	etcdClusterName := pflag.String("etcd-cluster-name",
		"",
		"Name of this cluster, must match the name of the eventual cluster")
	fmt.Printf("Using etcd cluster name %s\n", *etcdClusterName)

	etcdInitialCluster := pflag.String("etcd-initial-cluster",
		"",
		"Comma separated list of the peer advertise URLs of the complete eventual cluster, including our own.")
	fmt.Printf("Using etcd initial cluster %s\n", *etcdInitialCluster)

	etcdAdvertiseURL := pflag.String("etcd-peer-advertise-url",
		"",
		"The peer advertise URL of *this* peer, must match the one for the eventual cluster")
	fmt.Printf("Using advertise URL %s\n", *etcdAdvertiseURL)

	etcdDataDir := pflag.String("etcd-data-dir",
		"/var/etcd",
		"Location of the etcd data directory to restore into.")
	fmt.Printf("Using etcd data directory %s\n", *etcdDataDir)

	snapshotDir := pflag.String("snapshot-dir",
		"/tmp/snapshot",
		"Location of a temporary directory to make the backup into")
	fmt.Printf("Using snapshot directory %s\n", *snapshotDir)

	uploadProxyURL := pflag.String("upload-proxy-url",
		"",
		"URL of the proxy server to use to download the backup from remote storage.")
	fmt.Printf("Using Bucket URL %s\n", *uploadProxyURL)

	backupClusterIdentifier := pflag.String("backup-cluster-identifier",
		"",
		"Identifier for the cluster that the backup was taken from.")
	fmt.Printf("Requesting backup from %s\n", *backupClusterIdentifier)

	backupIdentifier := pflag.String("backup-identifier",
		"",
		"Identifier for the backup.")
	fmt.Printf("Requesting backup from %s\n", *backupIdentifier)

	timeoutSeconds := pflag.Int64("timeout-seconds",
		300,
		"Timeout, in seconds, of the whole restore operation.")
	fmt.Printf("Using %d seconds timeout`n", timeoutSeconds)

	// Pull the object from cloud storage into the snapshot directory.
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*timeoutSeconds))
	defer ctxCancel()

	conn, err := grpc.Dial(*uploadProxyURL)
	if err != nil {
		panic(err)
	}
	c := pb.NewProxyClient(conn)
	r, err := c.Download(ctx, &pb.DownloadRequest{
		ClusterIdentifier: *backupClusterIdentifier,
		BackupIdentifier:  *backupIdentifier,
	})
	err = conn.Close()
	if err != nil {
		panic(err)
	}

	snapshotFilePath := filepath.Join(*snapshotDir, "snapshot.db")
	fmt.Printf("Saving Object to local storage location %s\n", snapshotFilePath)
	snapshotFile, err := os.Create(snapshotFilePath)
	if err != nil {
		panic(err)
	}
	snapshotFileWriter := bufio.NewWriter(snapshotFile)
	_, err = io.Copy(snapshotFileWriter, bytes.NewReader(r.Backup))
	if err != nil {
		panic(err)
	}
	err = snapshotFileWriter.Flush()
	if err != nil {
		panic(err)
	}

	restoreDir := filepath.Join(*snapshotDir, "data-dir")

	restoreConfig := snapshot.RestoreConfig{
		SnapshotPath:        snapshotFilePath,
		Name:                *etcdPeerName,
		OutputDataDir:       restoreDir,
		OutputWALDir:        filepath.Join(restoreDir, "member", "wal"),
		PeerURLs:            []string{*etcdAdvertiseURL},
		InitialCluster:      *etcdInitialCluster,
		InitialClusterToken: *etcdClusterName,
		SkipHashCheck:       false,
	}

	client := snapshot.NewV3(nil)
	fmt.Printf("Executing restore\n")
	err = client.Restore(restoreConfig)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Copying restored data directory %s into correct PV on path %s\n", restoreDir, etcdDataDir)
	err = copy.Copy(restoreDir, *etcdDataDir)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Restore complete\n")
}