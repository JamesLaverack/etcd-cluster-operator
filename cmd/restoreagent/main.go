package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3/snapshot"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

func main() {
	viper.SetEnvPrefix("restore")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	fmt.Printf("BEES\n")

	etcdPeerName := viper.GetString("etcd-peer-name")
	fmt.Printf("Using etcd peer name %s\n", etcdPeerName)

	etcdClusterName := viper.GetString("etcd-cluster-name")
	fmt.Printf("Using etcd cluster name %s\n", etcdClusterName)

	etcdInitialCluster := viper.GetString("etcd-initial-cluster")
	fmt.Printf("Using etcd initial cluster %s\n", etcdInitialCluster)

	etcdAdvertiseURL := viper.GetString("etcd-advertise-url")
	fmt.Printf("Using advertise URL %s\n", etcdAdvertiseURL)

	etcdDataDir := viper.GetString("etcd-data-dir")
	fmt.Printf("Using etcd data directory %s\n", etcdDataDir)

	snapshotDir := viper.GetString("snapshot-dir")
	fmt.Printf("Using snapshot directory %s\n", snapshotDir)

	bucketURL := viper.GetString("bucket-url")
	fmt.Printf("Using Bucket URL %s\n", bucketURL)

	objectPath := viper.GetString("object-path")
	fmt.Printf("Using Object Path %s\n", objectPath)

	// Pull the object from cloud storage into the snapshot directory.
	fmt.Printf("Pulling object from %s/%s\n", bucketURL, objectPath)
	ctx := context.Background()
	bucket, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		panic(err)
	}
	defer bucket.Close()
	blobReader, err := bucket.NewReader(ctx, objectPath, nil)
	if err != nil {
		panic(err)
	}
	snapshotFilePath := filepath.Join(snapshotDir, "snapshot.db")
	fmt.Printf("Saving Object to local storage location %s\n", snapshotFilePath)
	snapshotFile, err := os.Create(snapshotFilePath)
	if err != nil {
		panic(err)
	}
	snapshotFileWriter := bufio.NewWriter(snapshotFile)
	_, err = io.Copy(snapshotFileWriter, blobReader)
	if err != nil {
		panic(err)
	}
	err = snapshotFileWriter.Flush()
	if err != nil {
		panic(err)
	}

	restoreDir := filepath.Join(snapshotDir, "data-dir")

	restoreConfig := snapshot.RestoreConfig{
		SnapshotPath:        snapshotFilePath,
		Name:                etcdPeerName,
		OutputDataDir:       restoreDir,
		OutputWALDir:        filepath.Join(restoreDir, "member", "wal"),
		PeerURLs:            []string{etcdAdvertiseURL},
		InitialCluster:      etcdInitialCluster,
		InitialClusterToken: etcdClusterName,
		SkipHashCheck:       false,
	}

	// We've got a problematic Catch-22 here. We need to set the peer advertise URL to be the real URL that the eventual
	// peer will use, and we need to set the initial cluster list to contain all the peer advertise URLs for the future
	// cluster. We know those ahead of time (and the restore controller has provided them in environment variables).
	// However without the Kubernetes service to provide DNS and the Hostnames set on the pods they can't actually be
	// resolved. Unfortunately, the etcd API will try to resolve the advertise URL to check it's real.
	// So, we need to fake it. I'm not looking for judgement.
	u, err := url.Parse(etcdAdvertiseURL)
	if err != nil {
		panic(err)
	}
	uParts := strings.Split(u.Host, ":")
	f, err := os.OpenFile("/etc/hosts",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	hostsFileEntry := fmt.Sprintf("127.0.0.1 %s", uParts[0])
	fmt.Printf("Appending '%s' to hosts file\n", hostsFileEntry)
	if _, err := f.WriteString(hostsFileEntry); err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

	client := snapshot.NewV3(nil)
	fmt.Printf("Executing restore\n")
	err = client.Restore(restoreConfig)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Copying restored data directory %s into correct PV on path %s\n", restoreDir, etcdDataDir)
	err = copy.Copy(restoreDir, etcdDataDir)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Restore complete\n")
}
