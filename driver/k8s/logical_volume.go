package k8s

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/cybozu-go/topolvm"
	topolvmv1 "github.com/cybozu-go/topolvm/api/v1"
	"github.com/cybozu-go/topolvm/csi"
	"github.com/cybozu-go/topolvm/driver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type logicalVolumeService struct {
	client.Client
	mu sync.Mutex
}

const (
	indexFieldVolumeID = "status.volumeID"
)

var (
	scheme = runtime.NewScheme()
	logger = logf.Log.WithName("LogicalVolume")
)

// +kubebuilder:rbac:groups=topolvm.cybozu.com,resources=logicalvolumes,verbs=get;list;watch;create;delete
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch

// NewLogicalVolumeService returns LogicalVolumeService.
func NewLogicalVolumeService(mgr manager.Manager) (driver.LogicalVolumeService, error) {
	err := mgr.GetFieldIndexer().IndexField(&topolvmv1.LogicalVolume{}, indexFieldVolumeID,
		func(o runtime.Object) []string {
			return []string{o.(*topolvmv1.LogicalVolume).Status.VolumeID}
		})
	if err != nil {
		return nil, err
	}

	return &logicalVolumeService{Client: mgr.GetClient()}, nil
}

func (s *logicalVolumeService) CreateVolume(ctx context.Context, node string, name string, sizeGb int64, capabilities []*csi.VolumeCapability) (string, error) {
	logger.Info("k8s.CreateVolume called", "name", name, "node", node, "size_gb", sizeGb)
	s.mu.Lock()
	defer s.mu.Unlock()

	lv := &topolvmv1.LogicalVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "LogicalVolume",
			APIVersion: "topolvm.cybozu.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: topolvmv1.LogicalVolumeSpec{
			Name:     name,
			NodeName: node,
			Size:     *resource.NewQuantity(sizeGb<<30, resource.BinarySI),
		},
	}

	existingLV := new(topolvmv1.LogicalVolume)
	err := s.Get(ctx, client.ObjectKey{Name: name}, existingLV)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return "", err
		}

		err := s.Create(ctx, lv)
		if err != nil {
			return "", err
		}
		logger.Info("created LogicalVolume CRD", "name", name)
	} else {
		// LV with same name was found; check compatibility
		// skip check of capabilities because (1) we allow both of two access types, and (2) we allow only one access mode
		// for ease of comparison, sizes are compared strictly, not by compatibility of ranges
		if !existingLV.IsCompatibleWith(lv) {
			return "", status.Error(codes.AlreadyExists, "Incompatible LogicalVolume already exists")
		}
		// compatible LV was found
	}

	for {
		logger.Info("waiting for setting 'status.volumeID'", "name", name)
		select {
		case <-ctx.Done():
			return "", errors.New("timed out")
		case <-time.After(1 * time.Second):
		}

		var newLV topolvmv1.LogicalVolume
		err := s.Get(ctx, client.ObjectKey{Name: name}, &newLV)
		if err != nil {
			logger.Error(err, "failed to get LogicalVolume", "name", name)
			continue
		}
		if newLV.Status.VolumeID != "" {
			logger.Info("end k8s.LogicalVolume", "volume_id", newLV.Status.VolumeID)
			return newLV.Status.VolumeID, nil
		}
		if newLV.Status.Code != codes.OK {
			err := s.Delete(ctx, &newLV)
			if err != nil {
				// log this error but do not return this error, because newLV.Status.Message is more important
				logger.Error(err, "failed to delete LogicalVolume")
			}
			return "", status.Error(newLV.Status.Code, newLV.Status.Message)
		}
	}
}

func (s *logicalVolumeService) DeleteVolume(ctx context.Context, volumeID string) error {
	lvList := new(topolvmv1.LogicalVolumeList)
	err := s.List(ctx, lvList, client.MatchingFields{indexFieldVolumeID: volumeID})
	if err != nil {
		return err
	}
	if len(lvList.Items) == 0 {
		logger.Info("volume is not found", "volume_id", volumeID)
		return nil
	} else if len(lvList.Items) > 1 {
		return fmt.Errorf("multiple LogicalVolume is found for VolumeID %s", volumeID)
	}

	return s.Delete(ctx, &lvList.Items[0])
}

func (s *logicalVolumeService) getLogicalVolume(ctx context.Context, volumeID string) (*topolvmv1.LogicalVolume, error) {
	lvList := new(topolvmv1.LogicalVolumeList)
	err := s.List(ctx, lvList, client.MatchingFields{indexFieldVolumeID: volumeID})
	if err != nil {
		return nil, err
	}

	if len(lvList.Items) > 1 {
		return nil, errors.New("found multiple volumes with volumeID " + volumeID)
	}
	if len(lvList.Items) == 0 {
		return nil, driver.ErrVolumeNotFound
	}
	return &lvList.Items[0], nil
}

func (s *logicalVolumeService) VolumeExists(ctx context.Context, volumeID string) error {
	_, err := s.getLogicalVolume(ctx, volumeID)
	return err
}

func (s *logicalVolumeService) listNodes(ctx context.Context) (*corev1.NodeList, error) {
	nl := new(corev1.NodeList)
	err := s.List(ctx, nl)
	if err != nil {
		return nil, err
	}
	return nl, nil
}

func (s *logicalVolumeService) ExpandVolume(ctx context.Context, volumeID string, sizeGb int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	lv, err := s.getLogicalVolume(ctx, volumeID)
	if err != nil {
		logger.Error(err, "failed to get logical volume", "name", lv.Name)
		return err
	}

	targetNodeName := lv.Spec.NodeName
	node, err := s.getNode(ctx, targetNodeName)
	if err != nil {
		logger.Error(err, "failed to get node", "name", targetNodeName)
		return err
	}

	cap := s.getNodeCapacity(*node)
	cs, err := s.GetCurrentSize(ctx, volumeID)
	if err != nil {
		logger.Error(err, "failed to get current size", "name", lv.Name)
		return err
	}

	if cap < (sizeGb<<30 - cs) {
		err := errors.New("not enough space")
		logger.Error(err, "not enough space is left", "name", targetNodeName)
		return err
	}

	lv2 := lv.DeepCopy()
	lv2.Spec.Size = *resource.NewQuantity(sizeGb<<30, resource.BinarySI)
	// clear Status.Code not to be rejected at topolvm-node
	lv2.Status.Code = codes.OK
	lv2.Status.Message = ""
	patch := client.MergeFrom(lv)
	if err := s.Patch(ctx, lv2, patch); err != nil {
		logger.Error(err, "failed to patch .spec.size", "name", lv.Name)
		return err
	}

	// wait until topolvm-node expands the target volume
	lvName := lv.Name
	for {
		logger.Info("waiting for update of 'status.currentSize'", "name", lvName)
		select {
		case <-ctx.Done():
			return errors.New("timed out")
		case <-time.After(1 * time.Second):
		}

		var changedLV topolvmv1.LogicalVolume
		err := s.Get(ctx, client.ObjectKey{Name: lvName}, &changedLV)
		if err != nil {
			logger.Error(err, "failed to get LogicalVolume", "name", lvName)
			continue
		}
		if changedLV.Status.CurrentSize == nil {
			return errors.New("status.currentSize should not be nil")
		}
		if changedLV.Status.CurrentSize.Value() != changedLV.Spec.Size.Value() {
			logger.Info("failed to match current size and requested size", "current", changedLV.Status.CurrentSize.Value(), "requested", changedLV.Spec.Size.Value())
			continue
		}

		if changedLV.Status.Code != codes.OK {
			return status.Error(changedLV.Status.Code, changedLV.Status.Message)
		}

		return nil
	}
}

func (s *logicalVolumeService) GetCapacity(ctx context.Context, requestNodeNumber string) (int64, error) {
	nl, err := s.listNodes(ctx)
	if err != nil {
		return 0, err
	}

	capacity := int64(0)
	if len(requestNodeNumber) == 0 {
		for _, node := range nl.Items {
			capacity += s.getNodeCapacity(node)
		}
		return capacity, nil
	}

	for _, node := range nl.Items {
		if nodeNumber, ok := node.Labels[topolvm.TopologyNodeKey]; ok {
			if requestNodeNumber != nodeNumber {
				continue
			}
			c, ok := node.Annotations[topolvm.CapacityKey]
			if !ok {
				return 0, fmt.Errorf("%s is not found", topolvm.CapacityKey)
			}
			return strconv.ParseInt(c, 10, 64)
		}
	}

	return 0, errors.New("capacity not found")
}

func (s *logicalVolumeService) GetMaxCapacity(ctx context.Context) (string, int64, error) {
	nl, err := s.listNodes(ctx)
	if err != nil {
		return "", 0, err
	}
	var nodeName string
	var maxCapacity int64
	for _, node := range nl.Items {
		c := s.getNodeCapacity(node)

		if maxCapacity < c {
			maxCapacity = c
			nodeName = node.Name
		}
	}
	return nodeName, maxCapacity, nil
}

func (s *logicalVolumeService) getNodeCapacity(node corev1.Node) int64 {
	c, ok := node.Annotations[topolvm.CapacityKey]
	if !ok {
		return 0
	}
	val, _ := strconv.ParseInt(c, 10, 64)
	return val
}

func (s *logicalVolumeService) getNode(ctx context.Context, name string) (*corev1.Node, error) {
	nl, err := s.listNodes(ctx)
	if err != nil {
		return nil, err
	}
	for _, node := range nl.Items {
		if node.Name == name {
			return &node, nil
		}
	}
	return nil, errors.New("node not found")
}

func (s *logicalVolumeService) GetCurrentSize(ctx context.Context, volumeID string) (int64, error) {
	lv, err := s.getLogicalVolume(ctx, volumeID)
	if err != nil {
		return 0, err
	}

	if lv.Status.CurrentSize != nil {
		return lv.Status.CurrentSize.Value(), nil
	}

	// set `status.currentSize` to the value of `spec.size`
	lv2 := lv.DeepCopy()
	lv2.Status.CurrentSize = &lv2.Spec.Size
	patch := client.MergeFrom(lv)
	if err := s.Patch(ctx, lv2, patch); err != nil {
		logger.Error(err, "failed to patch .status.currentSize", "name", lv.Name)
		return 0, err
	}

	return lv2.Status.CurrentSize.Value(), nil
}
