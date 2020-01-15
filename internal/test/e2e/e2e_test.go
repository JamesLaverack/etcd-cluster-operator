package e2e

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	kindv1alpha3 "sigs.k8s.io/kind/pkg/apis/config/v1alpha3"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/create"
	"sigs.k8s.io/kind/pkg/container/cri"

	etcdv1alpha1 "github.com/improbable-eng/etcd-cluster-operator/api/v1alpha1"
	"github.com/improbable-eng/etcd-cluster-operator/internal/test/try"
)

const (
	expectedClusterSize = 3
)

var (
	fUseKind           = flag.Bool("kind", false, "Creates a Kind cluster to run the tests against.")
	fOutputDirectory   = flag.String("output-directory", "/tmp/etcd-e2e", "The absolute path to a directory where E2E results logs will be saved.")
	fKindClusterName   = flag.String("kind-cluster-name", "etcd-e2e", "The name of the Kind cluster to use or create")
	fUseCurrentContext = flag.Bool("current-context", false, "Runs the tests against the current Kubernetes context, the path to kube config defaults to ~/.kube/config, unless overridden by the KUBECONFIG environment variable.")
	fRepoRoot          = flag.String("repo-root", "", "The absolute path to the root of the etcd-cluster-operator git repository.")
	fCleanup           = flag.Bool("cleanup", true, "Tears down the Kind cluster once the test is finished.")

	etcdConfig = etcd.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: etcd.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
)

func objFromYaml(objBytes []byte) (runtime.Object, error) {
	scheme := runtime.NewScheme()
	if err := etcdv1alpha1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return runtime.Decode(
		serializer.NewCodecFactory(scheme).UniversalDeserializer(),
		objBytes,
	)
}

func objFromYamlPath(objPath string) (runtime.Object, error) {
	objBytes, err := ioutil.ReadFile(objPath)
	if err != nil {
		return nil, err
	}
	return objFromYaml(objBytes)
}

func getSpec(t *testing.T, o interface{}) interface{} {
	switch ot := o.(type) {
	case *etcdv1alpha1.EtcdCluster:
		return ot.Spec
	case *etcdv1alpha1.EtcdPeer:
		return ot.Spec
	default:
		require.Failf(t, "unknown type", "%#v", o)
	}
	return nil
}

// Starts a Kind cluster on the local machine, exposing port 2379 accepting ETCD connections.
func startKind(t *testing.T, ctx context.Context, stopped chan struct{}) (kind *cluster.Context, err error) {
	clusterExists := true
	// Only attempt to delete the cluster if we created it.
	// But always ensure that the stopped channel gets closed when the context
	// ends or is cancelled.
	go func() {
		defer close(stopped)
		<-ctx.Done()
		if kind == nil {
			return
		}
		t.Log("Collecting Kind logs.")
		err := kind.CollectLogs(path.Join(*fOutputDirectory, "kind"))
		assert.NoError(t, err, "failed to collect Kind logs")
		if !*fCleanup {
			t.Log("Not deleting Kind cluster because --cleanup=false")
			return
		}
		if clusterExists {
			t.Log("Not deleting Kind cluster because this was an existing cluster.")
			return
		}
		t.Log("Deleting Kind cluster.")
		err = kind.Delete()
		assert.NoError(t, err)
	}()

	clusters, err := cluster.List()
	if err != nil {
		return nil, err
	}
	for _, kind := range clusters {
		if kind.Name() == *fKindClusterName {
			t.Log("Found existing Kind cluster")
			return &kind, nil
		}
	}
	clusterExists = false
	t.Log("Starting new Kind cluster")
	kind = cluster.NewContext(*fKindClusterName)
	err = kind.Create(create.WithV1Alpha3(&kindv1alpha3.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: "kind.sigs.k8s.io/v1alpha3",
		},
		Nodes: []kindv1alpha3.Node{
			{
				Role: "control-plane",
				ExtraPortMappings: []cri.PortMapping{
					{
						ContainerPort: 32379,
						HostPort:      2379,
					},
				},
			},
		},
	}))
	if err != nil {
		return nil, err
	}
	return kind, nil
}

func buildOperator(t *testing.T, ctx context.Context) (imageTar string, err error) {
	t.Log("Building the operator")
	// Tag for running this test, for naming resources.
	operatorImage := "etcd-cluster-operator:test"

	// Build the operator.
	out, err := exec.CommandContext(ctx, "docker", "build", "--target=debug", "-t", operatorImage, *fRepoRoot).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w Output: %s", err, out)
	}

	// Bundle the image to a tar.
	tmpDir, err := ioutil.TempDir("", "etcd-cluster-operator-e2e-test")
	if err != nil {
		return "", err
	}

	imageTar = filepath.Join(tmpDir, "etcd-cluster-operator.tar")

	t.Log("Exporting the operator image")
	out, err = exec.CommandContext(ctx, "docker", "save", "-o", imageTar, operatorImage).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w Output: %s", err, out)
	}
	return imageTar, nil
}

func buildRestoreAgent(t *testing.T, ctx context.Context) (imageTar string, err error) {
	t.Log("Building the restore agent")
	// Tag for running this test, for naming resources.
	restoreagentImage := "etcd-restore-agent:test"

	// Build the agent.
	t.Logf("Repo Root %s", *fRepoRoot)
	cmd := exec.CommandContext(ctx, "docker", "build", "--file=restoreagent.Dockerfile", "-t", restoreagentImage, *fRepoRoot)
	cmd.Dir = *fRepoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w Output: %s", err, out)
	}

	// Bundle the image to a tar.
	tmpDir, err := ioutil.TempDir("", "etcd-cluster-operator-e2e-test")
	if err != nil {
		return "", err
	}

	imageTar = filepath.Join(tmpDir, "etcd-resotreagent.tar")

	t.Log("Exporting the agent image")
	out, err = exec.CommandContext(ctx, "docker", "save", "-o", imageTar, restoreagentImage).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w Output: %s", err, out)
	}
	return imageTar, nil
}

func installOperator(t *testing.T,
	kubectl *kubectlContext,
	kind *cluster.Context,
	imageTar string,
	restoreAgentImageTar string,
) {
	t.Log("Installing cert-manager")
	err := kubectl.Apply("--validate=false", "--filename=https://github.com/jetstack/cert-manager/releases/download/v0.11.0/cert-manager.yaml")
	require.NoError(t, err)

	// Ensure CRDs exist in the cluster.
	t.Log("Applying CRDs")
	err = kubectl.Apply("--kustomize", filepath.Join(*fRepoRoot, "config", "crd"))
	require.NoError(t, err)

	imageFile, err := os.Open(imageTar)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, imageFile.Close(), "failed to close operator image tar")
	}()
	// Load the built image into the Kind cluster.
	t.Log("Loading image in to Kind cluster")
	nodes, err := kind.ListInternalNodes()
	require.NoError(t, err)
	for _, node := range nodes {
		err := node.LoadImageArchive(imageFile)
		require.NoError(t, err)
	}

	restoreAgentImageFile, err := os.Open(restoreAgentImageTar)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, restoreAgentImageFile.Close(), "failed to close restoreagent image tar")
	}()
	// Load the built image into the Kind cluster.
	t.Log("Loading image in to Kind cluster")
	require.NoError(t, err)
	for _, node := range nodes {
		err := node.LoadImageArchive(restoreAgentImageFile)
		require.NoError(t, err)
	}

	t.Log("Waiting for cert-manager to be ready")
	err = kubectl.Wait("--for=condition=Available", "--timeout=300s", "apiservice", "v1beta1.webhook.cert-manager.io")
	require.NoError(t, err)

	// Deploy the operator.
	t.Log("Applying operator")
	err = kubectl.Apply("--kustomize", filepath.Join(*fRepoRoot, "config", "test"))
	require.NoError(t, err)

	// Ensure the operator starts.
	err = try.Eventually(func() error {
		out, err := kubectl.Get("--namespace", "eco-system", "deploy", "eco-controller-manager", "-o=jsonpath='{.status.readyReplicas}'")
		if err != nil {
			return err
		}
		if out != "'1'" {
			return errors.New("expected exactly 1 replica of the operator to be available, got: " + out)
		}
		return nil
	}, time.Minute, time.Second*5)
	require.NoError(t, err)
}

func setupKind(t *testing.T, ctx context.Context, stopped chan struct{}) *kubectlContext {
	ctx, cancel := context.WithCancel(ctx)
	var (
		kind                 *cluster.Context
		imageTar             string
		restoreAgentImageTar string
		wg                   sync.WaitGroup
	)
	stoppedKind := make(chan struct{})
	go func() {
		<-stoppedKind
		close(stopped)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		imageTar, err = buildOperator(t, ctx)
		if err != nil {
			assert.NoError(t, err)
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		restoreAgentImageTar, err = buildRestoreAgent(t, ctx)
		if err != nil {
			assert.NoError(t, err)
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		kind, err = startKind(t, ctx, stoppedKind)
		if err != nil {
			assert.NoError(t, err)
			cancel()
		}
	}()
	wg.Wait()
	require.NoError(t, ctx.Err())
	kubectl := &kubectlContext{
		t:          t,
		configPath: kind.KubeConfigPath(),
	}

	installOperator(t, kubectl, kind, imageTar, restoreAgentImageTar)

	return kubectl
}

func setupCurrentContext(t *testing.T) *kubectlContext {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	configPath := filepath.Join(home, ".kube", "config")
	if path, found := os.LookupEnv("KUBECONFIG"); found {
		configPath = path
	}
	return &kubectlContext{
		t:          t,
		configPath: configPath,
	}
}

func TestE2E(t *testing.T) {
	var kubectl *kubectlContext
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()
	switch {
	case *fUseKind:
		stoppedKind := make(chan struct{})
		defer func() { <-stoppedKind }()
		defer cancel()
		kubectl = setupKind(t, ctx, stoppedKind)
	case *fUseCurrentContext:
		kubectl = setupCurrentContext(t)
	default:
		t.Skip("Supply either --kind or --current-context to run E2E tests")
	}

	// Delete all existing test namespaces, to free up resources before running new tests.
	DeleteAllTestNamespaces(t, kubectl)

	sampleClusterPath := filepath.Join(*fRepoRoot, "config", "samples", "etcd_v1alpha1_etcdcluster.yaml")

	// Pre-flight check that we can submit etcd API resources, before continuing
	// with the remaining tests.
	// Because the Etcd mutating and validating webhook service may not
	// immediately be responding.
	kubectl = kubectl.WithDefaultNamespace("default")
	var out string
	err := try.Eventually(func() (err error) {
		out, err = kubectl.DryRun(sampleClusterPath)
		return err
	}, time.Second*5, time.Second*1)
	require.NoError(t, err, out)

	// This outer function is needed because call to t.Run does not block if
	// t.Parallel is used in its test function.
	// See https://github.com/golang/go/issues/17791#issuecomment-258527390
	t.Run("Parallel", func(t *testing.T) {
		/*		t.Run("SampleCluster", func(t *testing.T) {
					t.Parallel()
					kubectl := kubectl.WithT(t)
					ns, cleanup := NamespaceForTest(t, kubectl)
					defer cleanup()
					sampleClusterTests(t, kubectl.WithDefaultNamespace(ns), sampleClusterPath)
				})
				t.Run("Webhooks", func(t *testing.T) {
					t.Parallel()
					kubectl := kubectl.WithT(t)
					ns, cleanup := NamespaceForTest(t, kubectl)
					defer cleanup()
					webhookTests(t, kubectl.WithDefaultNamespace(ns))
				})
				t.Run("Persistence", func(t *testing.T) {
					t.Parallel()
					kubectl := kubectl.WithT(t)
					ns, cleanup := NamespaceForTest(t, kubectl)
					defer cleanup()
					persistenceTests(t, kubectl.WithDefaultNamespace(ns))
				})
				t.Run("ScaleDown", func(t *testing.T) {
					t.Parallel()
					kubectl := kubectl.WithT(t)
					ns, cleanup := NamespaceForTest(t, kubectl)
					defer cleanup()
					scaleDownTests(t, kubectl.WithDefaultNamespace(ns))
				})*/
		t.Run("Backup", func(t *testing.T) {
			t.Parallel()
			ns, cleanup := NamespaceForTest(t, kubectl)
			if *fCleanup {
				defer cleanup()
			} else {
				t.Logf("Not performing namespace cleanup for %s because --cleanup=false", ns)
			}
			backupRestoreTests(t, kubectl.WithT(t).WithDefaultNamespace(ns))
		})
	})
}

func backupRestoreTests(t *testing.T, kubectl *kubectlContext) {
	t.Log("Given a one node cluster.")
	err := kubectl.Apply("--filename", filepath.Join(*fRepoRoot, "config", "test", "e2e", "backup", "etcdcluster.yaml"))
	require.NoError(t, err)

	clusterName := "e2e-backup-cluster"

	t.Log("Containing data.")
	out, err := eventuallyInCluster(kubectl,
		"set-etcd-value",
		time.Minute*5,
		"quay.io/coreos/etcd:v3.3.17",
		map[string]string{"ETCDCTL_API": "3"},
		"etcdctl", "--insecure-discovery", fmt.Sprintf("--discovery-srv=%s", clusterName), "put", "foo", "bar")
	require.NoError(t, err, out)

	t.Log("A backup can be taken to a local disk.")
	err = kubectl.Apply("--filename", filepath.Join(*fRepoRoot, "config", "test", "e2e", "backup", "etcdbackup.yaml"))
	require.NoError(t, err)

	kubectlSystem := kubectl.WithDefaultNamespace("eco-system")
	var podNames string
	err = try.Eventually(func() (err error) {
		podNames, err = kubectlSystem.Get("pods", "--output=name")
		return err
	}, time.Minute, time.Second*5)
	require.NoError(t, err)
	podNameList := strings.Split(strings.TrimSpace(podNames), "\n")
	require.Len(t, podNameList, 1)
	podName := strings.TrimPrefix(podNameList[0], "pod/")

	t.Log("And it will be persisted in the expected location.")
	err = try.Eventually(func() (err error) {
		out, err = kubectlSystem.Exec(podName, "ls", "/tmp/backups", "-c", "manager")
		return err
	}, time.Minute*2, time.Second*10)
	require.NoError(t, err)
	t.Log(string(out))
	require.Len(t, strings.Split(string(out), "\n"), 2)

	// So we need to be certain that when we bring the database back with a restore that its actually a restore and not
	// a resumption of the existing database by accident. So we'll set the 'foo' key to something else *after* we've
	// taken our backup.
	t.Log("And there are modifications made after the backup was taken")
	_, err = eventuallyInCluster(kubectl,
		"set-etcd-value-after-backup",
		time.Minute*2,
		"quay.io/coreos/etcd:v3.3.17",
		map[string]string{"ETCDCTL_API": "3"},
		"etcdctl", "--insecure-discovery", fmt.Sprintf("--discovery-srv=%s", clusterName), "put", "foo", "key-changed-after-backup")
	require.NoError(t, err, out)

	t.Log("And the cluster is deleted")
	err = kubectl.Delete("etcdcluster", clusterName)
	require.NoError(t, err)
	err = try.Eventually(func() (err error) {
		out, err = kubectl.Get("etcdpeer", "-o=jsonpath='{.items[*].metadata.name}'")
		if strings.Contains(out, clusterName) {
			return errors.New(fmt.Sprintf("EtcdPeers not deleted, was %s", out))
		}
		return err
	}, time.Minute*2, time.Second*10)
	require.NoError(t, err)
	err = try.Eventually(func() (err error) {
		out, err = kubectl.Get("pod",
			fmt.Sprintf("-l etcd.improbable.io/cluster-name=%s,app.kubernetes.io/name=etcd", clusterName),
			"-o=jsonpath='{.items[*].metadata.name}'")
		if strings.Contains(out, clusterName) {
			return errors.New("Pods not deleted")
		}
		return err
	}, time.Minute*2, time.Second*10)
	require.NoError(t, err)

	t.Log("And the PVCs are deleted")
	err = try.Eventually(func() (err error) {
		return kubectl.Delete("pvc",
			fmt.Sprintf("-l etcd.improbable.io/cluster-name=%s,app.kubernetes.io/name=etcd", clusterName))
	}, time.Minute*2, time.Second*10)
	require.NoError(t, err)
	err = try.Eventually(func() (err error) {
		out, err = kubectl.Get("pvc",
			fmt.Sprintf("-l etcd.improbable.io/cluster-name=%s,app.kubernetes.io/name=etcd", clusterName),
			"-o=jsonpath='{.items[*].metadata.name}'")
		if strings.Contains(out, clusterName) {
			return errors.New("PVCs not deleted")
		}
		return err
	}, time.Minute*2, time.Second*10)
	require.NoError(t, err)

	// At this point the cluster should be well-and-truly dead. So do the restore
	t.Log("We restore the cluster")
	_, err = kubectl.Create("secret", "generic", "backup-gcs-cred", fmt.Sprintf("--from-file=%s", filepath.Join(*fRepoRoot, "config", "test", "e2e", "backup", "gcp.json")))
	require.NoError(t, err)
	err = kubectl.Apply("--filename", filepath.Join(*fRepoRoot, "config", "test", "e2e", "backup", "etcdrestore.yaml"))
	require.NoError(t, err)

	// Wait for our cluster to appear
	err = try.Eventually(func() error {
		t.Log("")
		members, err := kubectl.Get("etcdcluster", clusterName, "-o=jsonpath='{.status.members...name}'")
		if err != nil {
			return err
		}
		// Don't assert on exact members, just that we have one of them.
		numMembers := len(strings.Split(members, " "))
		if numMembers != 1 {
			return errors.New(fmt.Sprintf("Expected etcd member list to have three members. Had %d.", numMembers))
		}
		return nil
	}, time.Minute*2, time.Second*10)
	require.NoError(t, err)

	t.Log("And our data is still there")
	out, err = eventuallyInCluster(kubectl,
		"get-etcd-value",
		time.Minute*2,
		"quay.io/coreos/etcd:v3.3.17",
		map[string]string{"ETCDCTL_API": "3"},
		"etcdctl", "--insecure-discovery", fmt.Sprintf("--discovery-srv=%s", clusterName), "get", "foo")
	require.NoError(t, err, out)
	// It'll output the key name, then a newline, then the value, then a newline.
	assert.Equal(t, "foo\nbar\n", out)
}

func webhookTests(t *testing.T, kubectl *kubectlContext) {
	t.Run("Defaulting", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			name string
			path string
		}{
			{
				name: "EtcdCluster",
				path: "defaulting/etcdcluster.yaml",
			},
			{
				name: "EtcdPeer",
				path: "defaulting/etcdpeer.yaml",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				kubectl := kubectl.WithT(t)
				// local object
				lPath := filepath.Join(*fRepoRoot, "config/test/e2e", tc.path)
				l, err := objFromYamlPath(lPath)
				require.NoError(t, err)

				rBytes, err := kubectl.DryRun(lPath)
				require.NoError(t, err, rBytes)

				// Remote object
				r, err := objFromYaml([]byte(rBytes))
				require.NoError(t, err)

				if diff := cmp.Diff(getSpec(t, l), getSpec(t, r)); diff == "" {
					assert.Failf(t, "defaults were not applied to: %s", lPath)
				} else {
					t.Log(diff)
					assert.Contains(t, diff, `AccessModes:`)
					assert.Contains(t, diff, `VolumeMode:`)
				}
			})

		}
	})

	t.Run("Validation", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			name string
			path string
		}{
			{
				name: "EtcdCluster",
				path: "validation/etcdcluster_missing_storageclassname.yaml",
			},
			{
				name: "EtcdPeer",
				path: "validation/etcdpeer_missing_storageclassname.yaml",
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				kubectl := kubectl.WithT(t)
				lPath := filepath.Join(*fRepoRoot, "config/test/e2e", tc.path)
				out, err := kubectl.DryRun(lPath)
				assert.Regexp(
					t,
					`^Error from server \(spec.storage.volumeClaimTemplate.storageClassName: Required value\):`,
					err,
				)
				assert.Empty(t, out)
			})
		}
	})

}

func sampleClusterTests(t *testing.T, kubectl *kubectlContext, sampleClusterPath string) {
	err := kubectl.Apply(
		"--filename", filepath.Join(*fRepoRoot, "internal", "test", "e2e", "fixtures", "cluster-client-service.yaml"),
		"--filename", sampleClusterPath,
	)
	require.NoError(t, err)

	t.Run("EtcdClusterAvailability", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()

		etcdClient, err := etcd.New(etcdConfig)
		require.NoError(t, err)

		err = try.Eventually(func() error {
			members, err := etcd.NewMembersAPI(etcdClient).List(ctx)
			if err != nil {
				return err
			}

			if len(members) != expectedClusterSize {
				return errors.New(fmt.Sprintf("expected %d etcd peers, got %d", expectedClusterSize, len(members)))
			}
			return nil
		}, time.Minute*2, time.Second*10)
		require.NoError(t, err)
	})

	t.Run("EtcdClusterStatus", func(t *testing.T) {
		kubectl := kubectl.WithT(t)
		err := try.Eventually(func() error {
			t.Log("")
			members, err := kubectl.Get("etcdcluster", "my-cluster", "-o=jsonpath='{.status.members...name}'")
			if err != nil {
				return err
			}
			// Don't assert on exact members, just that we have three of them.
			numMembers := len(strings.Split(members, " "))
			if numMembers != 3 {
				return errors.New(fmt.Sprintf("Expected etcd member list to have three members. Had %d.", numMembers))
			}
			return nil
		}, time.Minute*2, time.Second*10)
		require.NoError(t, err)
	})

	t.Run("ScaleUp", func(t *testing.T) {
		kubectl := kubectl.WithT(t)
		// Attempt to scale up to five nodes
		err = kubectl.Scale("etcdcluster/my-cluster", 5)
		require.NoError(t, err)

		etcdClient, err := etcd.New(etcdConfig)
		require.NoError(t, err)

		err = try.Eventually(func() error {
			t.Log("Listing etcd members")
			members, err := etcd.NewMembersAPI(etcdClient).List(context.Background())
			if err != nil {
				return err
			}

			if len(members) != 5 {
				return errors.New(fmt.Sprintf("expected %d etcd peers, got %d", 5, len(members)))
			}

			for _, member := range members {
				if len(member.ClientURLs) == 0 {
					return errors.New("peer has no client URLs")
				}
				if len(member.PeerURLs) == 0 {
					return errors.New("peer has no peer URLs")
				}
				if member.ID == "" {
					return errors.New("peer has no ID")
				}
				if member.Name == "" {
					return errors.New("peer has no Name")
				}
			}
			return nil
		}, time.Minute*2, time.Second*10)
		require.NoError(t, err)
	})
}

func persistenceTests(t *testing.T, kubectl *kubectlContext) {
	t.Log("Given a 1-node cluster.")
	configPath := filepath.Join(*fRepoRoot, "config", "test", "e2e", "persistence")
	err := kubectl.Apply("--filename", configPath)
	require.NoError(t, err)

	t.Log("Containing data.")
	expectedValue := "foobarbaz"

	out, err := eventuallyInCluster(kubectl, "set-etcd-value", time.Minute*2, "quay.io/coreos/etcd:v3.2.27", nil, "etcdctl", "--insecure-discovery", "--discovery-srv=cluster1", "set", "--", "foo", expectedValue)
	require.NoError(t, err, out)

	t.Log("If the cluster is deleted.")
	err = kubectl.Delete("etcdcluster", "cluster1", "--wait")
	require.NoError(t, err)

	t.Log("And all the cluster pods terminate.")
	err = try.Eventually(func() error {
		out, err := kubectl.Get(
			"pods",
			"--selector", "etcd.improbable.io/cluster-name=cluster1",
			"--output", "go-template={{ len .items }}",
		)
		if err != nil {
			return err
		}
		if out != "0" {
			return errors.New("expected 0 pods, got: " + out)
		}
		return nil
	}, time.Minute, time.Second*5)
	require.NoError(t, err)

	t.Log("The cluster can be restored.")
	err = kubectl.Apply("--filename", configPath)
	require.NoError(t, err)

	t.Log("And the data is still available.")
	out, err = eventuallyInCluster(kubectl, "get-etcd-value", time.Minute*2, "quay.io/coreos/etcd:v3.2.27", nil, "etcdctl", "--insecure-discovery", "--discovery-srv=cluster1", "get", "--quorum", "--", "foo")
	require.NoError(t, err, out)
	assert.Equal(t, expectedValue+"\n", out)
}

func scaleDownTests(t *testing.T, kubectl *kubectlContext) {
	t.Log("Given a 3-node cluster.")
	configPath := filepath.Join(*fRepoRoot, "config", "samples", "etcd_v1alpha1_etcdcluster.yaml")
	err := kubectl.Apply("--filename", configPath)
	require.NoError(t, err)

	t.Log("Where all the nodes are up")
	err = try.Eventually(
		func() error {
			out, err := kubectl.Get("etcdcluster", "my-cluster", "-o=jsonpath={.status.replicas}")
			if err != nil {
				return err
			}
			statusReplicas, err := strconv.Atoi(out)
			if err != nil {
				return err
			}
			if statusReplicas != 3 {
				return fmt.Errorf("unexpected status.replicas. Wanted: 3, Got: %d", statusReplicas)
			}
			return nil
		},
		time.Minute*2, time.Second*10,
	)
	require.NoError(t, err)

	t.Log("Which contains data")
	const expectedValue = "foobarbaz"

	out, err := eventuallyInCluster(kubectl, "set-etcd-value", time.Minute*2, "quay.io/coreos/etcd:v3.2.27", nil, "etcdctl", "--insecure-discovery", "--discovery-srv=my-cluster", "set", "--", "foo", expectedValue)
	require.NoError(t, err, out)

	t.Log("If the cluster is scaled down")
	const expectedReplicas = 1
	err = kubectl.Scale("etcdcluster/my-cluster", expectedReplicas)
	require.NoError(t, err)

	t.Log("The etcdcluster.status is updated when the cluster has been resized.")
	err = try.Eventually(
		func() error {
			out, err := kubectl.Get("etcdcluster", "my-cluster", "-o=jsonpath={.status.replicas}")
			require.NoError(t, err, out)
			statusReplicas, err := strconv.Atoi(out)
			require.NoError(t, err, out)
			if expectedReplicas != statusReplicas {
				return fmt.Errorf("unexpected status.replicas. Wanted: %d, Got: %d", expectedReplicas, statusReplicas)
			}
			return err
		},
		time.Minute*5, time.Second*10,
	)
	require.NoError(t, err)

	t.Log("And the data is still available.")
	out, err = eventuallyInCluster(kubectl, "get-etcd-value", time.Minute*2, "quay.io/coreos/etcd:v3.2.27", nil, "etcdctl", "--insecure-discovery", "--discovery-srv=my-cluster", "get", "--quorum", "--", "foo")
	require.NoError(t, err, out)
	assert.Equal(t, expectedValue+"\n", out)
}
