package main

import (
	"context"
	"fmt"
	"gocloud.dev/blob"
	"io/ioutil"
	"log"
	"net"
	"net/url"

	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"

	pb "github.com/improbable-eng/etcd-cluster-operator/api/proxy/v1"
)

type proxyServer struct {
	pb.UnimplementedProxyServer
	BucketURL string
}

// Turn a full object URL like `gs://my-bucket/my-dir/my-obj.db` into a bucket URL (`gs://my-bucket`) and an object path
// (`/my-dir/my-obj.db`).
func parseBackupUrl(backupUrl string) (string, string, error) {
	u, err := url.Parse(backupUrl)
	if err != nil {
		return "", "", err
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host), u.Path, nil
}

func (ps *proxyServer) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadReply, error) {
	bucketName, objectPath, err := parseBackupUrl(req.BackupUrl)
	if err != nil {
		return nil, err
	}
	bucket, err := blob.OpenBucket(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	blobReader, err := bucket.NewReader(ctx, objectPath, nil)
	if err != nil {
		return nil, err
	}

	// Here we read the entire contents of the backup into memory. In theory these could be quite big (multiple
	// gigabytes). So we're actually taking a risk that the backup could be *too big* for our available memory.
	backup, err := ioutil.ReadAll(blobReader)
	if err != nil {
		return nil, err
	}

	return &pb.DownloadReply{Backup: backup}, nil
}

func main() {
	// Setup defaults for expected configuration keys
	var apiPort = flag.Int("api-port", 8080, "Port to serve the API on.")
	flag.Parse()

	// Launch gRPC server
	grpcAddress := fmt.Sprintf(":%d", *apiPort)
	log.Printf("Using %q as listen address for proxy server", grpcAddress)
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Starting etcd-cluster-controller upload/download Proxy service")
	srv := grpc.NewServer()
	pb.RegisterProxyServer(srv, &proxyServer{})
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
