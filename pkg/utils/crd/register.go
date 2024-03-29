package crd

import (
	"fmt"
	apiext "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"reflect"

	//"context"
	v1 "github.com/aoxn/wdrip/pkg/apis/alibabacloud.com/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	//"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// InitializeCRD register crds from in cluster config file
func InitializeCRD(cfg *rest.Config) error { return doRegisterCRD(cfg) }

// RegisterFromKubeconfig register crds from kubeconfig file
func RegisterFromKubeconfig(name string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", name)
	if err != nil {
		return fmt.Errorf("register crd: build rest.config, %s", err.Error())
	}
	return doRegisterCRD(cfg)
}

func doRegisterCRD(cfg *rest.Config) error {
	extc, err := apiext.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("error create incluster client: %s", err.Error())
	}
	client := NewClient(extc)
	for _, crd := range []CRD{
		NewClusterCRD(client),
		NewMasterCRD(client),
		NewMasterSetCRD(client),
		NewNodePoolCRD(client),
	} {
		err := crd.Initialize()
		if err != nil {
			return fmt.Errorf("initialize crd: %s, %s", reflect.TypeOf(crd), err.Error())
		}
		klog.Infof("register crd: %s", reflect.TypeOf(crd))
	}
	return nil
}

type CRD interface {
	Initialize() error
	GetObject() runtime.Object
	GetListerWatcher() cache.ListerWatcher
}

// ClusterCRD is the cluster crd .
type ClusterCRD struct {
	crdc Interface
	//wdrip vcset.Interface
}

func NewClusterCRD(
	//wdripClient vcset.Interface,
	crdClient Interface,
) *ClusterCRD {
	return &ClusterCRD{
		crdc: crdClient,
		//wdrip: wdripClient,
	}
}

// podTerminatorCRD satisfies resource.crd interface.
func (p *ClusterCRD) Initialize() error {
	crd := Conf{
		Kind:       v1.ClusterKind,
		NamePlural: v1.ClusterNamePlural,
		Group:      v1.SchemeGroupVersion.Group,
		Version:    v1.SchemeGroupVersion.Version,
		Scope:      apiextv1beta1.ClusterScoped,
	}

	return p.crdc.EnsurePresent(crd)
}

// GetListerWatcher satisfies resource.crd interface (and retrieve.Retriever).
func (p *ClusterCRD) GetListerWatcher() cache.ListerWatcher {
	//return &cache.ListWatch{
	//	ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
	//		return p.wdrip.WdripV1().Clusters("").List(context.TODO(), options)
	//	},
	//	WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
	//		return p.wdrip.WdripV1().Clusters("").Watch(context.TODO(),options)
	//	},
	//}
	return nil
}

// GetObject satisfies resource.crd interface (and retrieve.Retriever).
func (p *ClusterCRD) GetObject() runtime.Object { return &v1.Cluster{} }

// MasterCRD is the cluster crd .
type MasterCRD struct {
	crdc Interface
	//wdrip vcset.Interface
}

func NewMasterCRD(
	//wdripClient vcset.Interface,
	crdClient Interface,
) *MasterCRD {
	return &MasterCRD{
		crdc: crdClient,
		//wdrip: wdripClient,
	}
}

// podTerminatorCRD satisfies resource.crd interface.
func (p *MasterCRD) Initialize() error {
	crd := Conf{
		Kind:       v1.MasterKind,
		NamePlural: v1.MasterNamePlural,
		Group:      v1.SchemeGroupVersion.Group,
		Version:    v1.SchemeGroupVersion.Version,
		Scope:      apiextv1beta1.ClusterScoped,
	}

	return p.crdc.EnsurePresent(crd)
}

// GetListerWatcher satisfies resource.crd interface (and retrieve.Retriever).
func (p *MasterCRD) GetListerWatcher() cache.ListerWatcher {
	//return &cache.ListWatch{
	//	ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
	//		return p.wdrip.WdripV1().Clusters("").List(context.TODO(), options)
	//	},
	//	WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
	//		return p.wdrip.WdripV1().Clusters("").Watch(context.TODO(),options)
	//	},
	//}
	return nil
}

// GetObject satisfies resource.crd interface (and retrieve.Retriever).
func (p *MasterCRD) GetObject() runtime.Object { return &v1.Master{} }

// MasterCRD is the cluster crd .
type MasterSetCRD struct {
	crdc Interface
	//wdrip vcset.Interface
}

func NewMasterSetCRD(
	//wdripClient vcset.Interface,
	crdClient Interface,
) *MasterSetCRD {
	return &MasterSetCRD{
		crdc: crdClient,
		//wdrip: wdripClient,
	}
}

// podTerminatorCRD satisfies resource.crd interface.
func (p *MasterSetCRD) Initialize() error {
	crd := Conf{
		Kind:       v1.MasterSetKind,
		NamePlural: v1.MasterSetNamePlural,
		Group:      v1.SchemeGroupVersion.Group,
		Version:    v1.SchemeGroupVersion.Version,
		Scope:      apiextv1beta1.ClusterScoped,
	}

	return p.crdc.EnsurePresent(crd)
}

// GetListerWatcher satisfies resource.crd interface (and retrieve.Retriever).
func (p *MasterSetCRD) GetListerWatcher() cache.ListerWatcher {
	//return &cache.ListWatch{
	//	ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
	//		return p.wdrip.WdripV1().Clusters("").List(context.TODO(), options)
	//	},
	//	WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
	//		return p.wdrip.WdripV1().Clusters("").Watch(context.TODO(),options)
	//	},
	//}
	return nil
}

// GetObject satisfies resource.crd interface (and retrieve.Retriever).
func (p *MasterSetCRD) GetObject() runtime.Object { return &v1.MasterSet{} }

type NodePoolCRD struct {
	crdc Interface
	//wdrip vcset.Interface
}

func NewNodePoolCRD(
	//wdripClient vcset.Interface,
	crdClient Interface,
) *NodePoolCRD {
	return &NodePoolCRD{
		crdc: crdClient,
		//wdrip: wdripClient,
	}
}

// podTerminatorCRD satisfies resource.crd interface.
func (p *NodePoolCRD) Initialize() error {
	crd := Conf{
		Kind:       v1.NodePoolKind,
		NamePlural: v1.NodePoolPlural,
		Group:      v1.SchemeGroupVersion.Group,
		Version:    v1.SchemeGroupVersion.Version,
		Scope:      apiextv1beta1.ClusterScoped,
	}

	return p.crdc.EnsurePresent(crd)
}

// GetListerWatcher satisfies resource.crd interface (and retrieve.Retriever).
func (p *NodePoolCRD) GetListerWatcher() cache.ListerWatcher {
	//return &cache.ListWatch{
	//	ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
	//		return p.wdrip.WdripV1().Clusters("").List(context.TODO(), options)
	//	},
	//	WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
	//		return p.wdrip.WdripV1().Clusters("").Watch(context.TODO(),options)
	//	},
	//}
	return nil
}

// GetObject satisfies resource.crd interface (and retrieve.Retriever).
func (p *NodePoolCRD) GetObject() runtime.Object { return &v1.NodePool{} }
