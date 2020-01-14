package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("restore")
	viper.AutomaticEnv()

	fmt.Printf("Using etcd data directory %s\n", viper.Get("etcd-data-dir"))
	fmt.Printf("Using snapshot directory %s\n", viper.Get("snapshot-dir"))
	fmt.Printf("Using restore type %s\n", viper.Get("type"))
	fmt.Printf("Using GCS Bucket Name %s\n", viper.Get("gcs-bucket-name"))
	fmt.Printf("Using GCS Object Path %s\n", viper.Get("gcs-object-path"))
	fmt.Printf("Using GCS Secret Key %s\n", viper.Get("gcs-secret-key"))
}
