package e2e

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topolvm/topolvm"
	topolvmv1 "github.com/topolvm/topolvm/api/v1"
	"google.golang.org/grpc"
	"sigs.k8s.io/yaml"
)

type cleanup struct {
	// key is volumeID, value is target path
	volumes map[string]string
}

func (c cleanup) register(volumeID, targetPath string) {
	By("[cleanup] register")
	if c.volumes == nil {
		c.volumes = make(map[string]string)
	}
	c.volumes[volumeID] = targetPath
}

func (c cleanup) unregister(volumeID, targetPath string) {
	By("[cleanup] unregister")
	if c.volumes != nil {
		delete(c.volumes, volumeID)
	}
}

func (c cleanup) unpublishVolumes(nc csi.NodeClient) {
	By("[cleanup] unpublishVolumes")
	for volumeID, targetPath := range c.volumes {
		req := &csi.NodeUnpublishVolumeRequest{
			VolumeId:   volumeID,
			TargetPath: targetPath,
		}
		_, err := nc.NodeUnpublishVolume(context.Background(), req)
		if err != nil {
			fmt.Printf("failed to unpublish volume: %v", req)
		}
	}
	c.volumes = nil
}

func testPublishVolume() {
	var (
		cl   cleanup
		nc   csi.NodeClient
		conn *grpc.ClientConn
	)
	nodeSocket := "/tmp/topolvm/worker1/plugins/topolvm.cybozu.com/node/csi-topolvm.sock"

	var lvmCountBefore int
	var capacitiesBefore map[string]map[string]string

	BeforeEach(func() {
		dialer := &net.Dialer{}
		dialFunc := func(ctx context.Context, a string) (net.Conn, error) {
			return dialer.DialContext(ctx, "unix", a)
		}
		var err error
		conn, err = grpc.Dial(nodeSocket, grpc.WithInsecure(), grpc.WithContextDialer(dialFunc))
		Expect(err).ShouldNot(HaveOccurred())

		nc = csi.NewNodeClient(conn)

		lvmCountBefore, err = countLVMs()
		Expect(err).ShouldNot(HaveOccurred())

		capacitiesBefore, err = getNodeAnnotationMapWithPrefix(topolvm.CapacityKeyPrefix)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		cl.unpublishVolumes(nc)
		if conn != nil {
			conn.Close()
			conn = nil
		}

		if !CurrentGinkgoTestDescription().Failed {
			Eventually(func() error {
				lvmCountAfter, err := countLVMs()
				if err != nil {
					return err
				}
				if lvmCountBefore != lvmCountAfter {
					return fmt.Errorf("lvm num mismatched. before: %d, after: %d", lvmCountBefore, lvmCountAfter)
				}

				stdout, stderr, err := kubectl("get", "node", "-o", "json")
				if err != nil {
					return fmt.Errorf("stdout=%s, stderr=%s", stdout, stderr)
				}

				capacitiesAfter, err := getNodeAnnotationMapWithPrefix(topolvm.CapacityKeyPrefix)
				if err != nil {
					return err
				}
				if diff := cmp.Diff(capacitiesBefore, capacitiesAfter); diff != "" {
					return fmt.Errorf("capacities on nodes should be same before and after the test: diff=%q", diff)
				}
				return nil
			}).Should(Succeed())
		}
	})

	It("should publish filesystem", func() {
		mountTargetPath := "/mnt/csi-node-test"

		By("creating a logical volume resource")
		lvYaml := []byte(`apiVersion: topolvm.cybozu.com/v1
kind: LogicalVolume
metadata:
  name: csi-node-test-fs
spec:
  deviceClass: ssd
  name: csi-node-test-fs
  nodeName: topolvm-e2e-worker
  size: 1Gi
`)

		_, _, err := kubectlWithInput(lvYaml, "apply", "-f", "-")
		Expect(err).ShouldNot(HaveOccurred())

		var volumeID string
		Eventually(func() error {
			stdout, stderr, err := kubectl("get", "logicalvolume", "csi-node-test-fs", "-o", "yaml")
			if err != nil {
				return fmt.Errorf("failed to get logical volume. stdout: %s, stderr: %s, err: %v", stdout, stderr, err)
			}

			var lv topolvmv1.LogicalVolume
			err = yaml.Unmarshal(stdout, &lv)
			if err != nil {
				return err
			}

			if len(lv.Status.VolumeID) == 0 {
				return errors.New("VolumeID is not set")
			}
			volumeID = lv.Status.VolumeID
			return nil
		}).Should(Succeed())

		cl.register(volumeID, mountTargetPath)

		By("creating Filesystem volume")
		mountVolCap := &csi.VolumeCapability{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{FsType: "btrfs"},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
		}

		req := &csi.NodePublishVolumeRequest{
			PublishContext:   map[string]string{},
			TargetPath:       mountTargetPath,
			VolumeCapability: mountVolCap,
			VolumeId:         volumeID,
		}
		resp, err := nc.NodePublishVolume(context.Background(), req)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp).ShouldNot(BeNil())

		By("publishing Filesystem volume again to check idempotency")
		resp, err = nc.NodePublishVolume(context.Background(), req)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp).ShouldNot(BeNil())

		By("publishing volume on same target path, but requested volume and existing one are different")
		_, err = nc.NodePublishVolume(context.Background(), &csi.NodePublishVolumeRequest{
			PublishContext: map[string]string{},
			TargetPath:     mountTargetPath,
			VolumeCapability: &csi.VolumeCapability{
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{FsType: "ext4"},
				},
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
				},
			},
			VolumeId: volumeID,
		})
		Expect(err).Should(HaveOccurred())

		By("unpublishing the volume")
		unpubReq := csi.NodeUnpublishVolumeRequest{
			VolumeId:   volumeID,
			TargetPath: mountTargetPath,
		}
		unpubResp, err := nc.NodeUnpublishVolume(context.Background(), &unpubReq)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(unpubResp).ShouldNot(BeNil())

		By("unpublishing the volume again to check idempotency")
		unpubResp, err = nc.NodeUnpublishVolume(context.Background(), &unpubReq)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(unpubResp).ShouldNot(BeNil())

		cl.unregister(volumeID, mountTargetPath)

		By("cleaning logicalvolume")
		stdout, stderr, err := kubectl("delete", "logicalvolume", "csi-node-test-fs")
		Expect(err).ShouldNot(HaveOccurred(), "stdout=%s, stderr=%s", stdout, stderr)
	})

	It("should be worked NodePublishVolume successfully to create a block device", func() {
		deviceTargetPath := "/dev/csi-node-test"

		By("creating a logical volume resource")
		lvYaml := []byte(`apiVersion: topolvm.cybozu.com/v1
kind: LogicalVolume
metadata:
  name: csi-node-test-block
spec:
  deviceClass: ssd
  name: csi-node-test-block
  nodeName: topolvm-e2e-worker
  size: 1Gi
`)

		_, _, err := kubectlWithInput(lvYaml, "apply", "-f", "-")
		Expect(err).ShouldNot(HaveOccurred())

		var volumeID string
		Eventually(func() error {
			stdout, stderr, err := kubectl("get", "logicalvolume", "csi-node-test-block", "-o", "yaml")
			if err != nil {
				return fmt.Errorf("failed to get logical volume. stdout: %s, stderr: %s, err: %v", stdout, stderr, err)
			}

			var lv topolvmv1.LogicalVolume
			err = yaml.Unmarshal(stdout, &lv)
			if err != nil {
				return err
			}

			if len(lv.Status.VolumeID) == 0 {
				return errors.New("VolumeID is not set")
			}
			volumeID = lv.Status.VolumeID
			return nil
		}).Should(Succeed())

		cl.register(volumeID, deviceTargetPath)

		By("creating raw block volume")
		blockVolCap := &csi.VolumeCapability{
			AccessType: &csi.VolumeCapability_Block{
				Block: &csi.VolumeCapability_BlockVolume{},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
		}

		req := &csi.NodePublishVolumeRequest{
			PublishContext:   map[string]string{},
			TargetPath:       deviceTargetPath,
			VolumeCapability: blockVolCap,
			VolumeId:         volumeID,
		}
		resp, err := nc.NodePublishVolume(context.Background(), req)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp).ShouldNot(BeNil())

		By("creating raw block volume again to check idempotency")
		resp, err = nc.NodePublishVolume(context.Background(), req)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp).ShouldNot(BeNil())

		By("creating volume on the same target path, but requested volume and existing one are different")
		_, err = nc.NodePublishVolume(context.Background(), &csi.NodePublishVolumeRequest{
			PublishContext:   map[string]string{},
			TargetPath:       deviceTargetPath,
			VolumeCapability: blockVolCap,
			VolumeId:         volumeID + "-invalid",
		})
		Expect(err).Should(HaveOccurred())

		By("unpublishing the volume")
		unpubReq := csi.NodeUnpublishVolumeRequest{
			VolumeId:   volumeID,
			TargetPath: deviceTargetPath,
		}
		unpubResp, err := nc.NodeUnpublishVolume(context.Background(), &unpubReq)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(unpubResp).ShouldNot(BeNil())

		By("deleting the volume again to check idempotency")
		unpubResp, err = nc.NodeUnpublishVolume(context.Background(), &unpubReq)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(unpubResp).ShouldNot(BeNil())

		cl.unregister(volumeID, deviceTargetPath)

		By("cleaning logicalvolume")
		stdout, stderr, err := kubectl("delete", "logicalvolume", "csi-node-test-block")
		Expect(err).ShouldNot(HaveOccurred(), "stdout=%s, stderr=%s", stdout, stderr)
	})
}
