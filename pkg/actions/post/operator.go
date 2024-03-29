//go:build linux || darwin
// +build linux darwin

package post

import (
	"bytes"
	"fmt"
	"github.com/aoxn/wdrip"
	"github.com/aoxn/wdrip/pkg/actions"
	"github.com/aoxn/wdrip/pkg/actions/post/addons"
	v12 "github.com/aoxn/wdrip/pkg/apis/alibabacloud.com/v1"
	"github.com/aoxn/wdrip/pkg/context"
	"github.com/aoxn/wdrip/pkg/utils"
	"github.com/aoxn/wdrip/pkg/utils/crd"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"html/template"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
	"path/filepath"
	"strings"
	"time"
)

const (
	ObjectName        = "config"
	KUBELET_UNIT_FILE = "/etc/systemd/system/kubelet.service"
)

type ActionPost struct {
}

// NewAction returns a new ActionPost for post kubernetes install
func NewActionPost() actions.Action {
	return &ActionPost{}
}

// Execute runs the ActionPost
func (a *ActionPost) Execute(ctx *actions.ActionContext) error {
	// Addon was installed by operator
	adds := ctx.WdripFlags().Addons
	cfgadds := []addons.ConfigTpl{addons.KUBEPROXY_MASTER, addons.KUBEPROXY_WORKER}
	if adds == "*" {
		cfgadds = addons.AddonConfigsTpl()
	}
	err := addons.InstallAddons(ctx.ProviderCtx(), ctx.Config(), cfgadds)
	if err != nil {
		return fmt.Errorf("install addons: %s", err.Error())
	}

	err = crd.RegisterFromKubeconfig("/etc/kubernetes/admin.conf")
	if err != nil {
		return fmt.Errorf("register crds: %s", err.Error())
	}
	err = WriteClusterCR(ctx.NodeContext)
	if err != nil {
		return fmt.Errorf("write cluster cfg: %s", err.Error())
	}
	err = WritePublicInfo(ctx.NodeContext)
	if err != nil {
		return fmt.Errorf("write public cluster info")
	}
	// Run wdrip operator default
	return RunWdrip(ctx.Config())
}

func WritePublicInfo(ctx *context.NodeContext) error {
	cfg := ctx.NodeObject().Status.BootCFG
	bcfg := clientcmdapi.Config{
		APIVersion: "v1",
		Clusters: map[string]*clientcmdapi.Cluster{
			"": {
				Server:                   fmt.Sprintf("https://%s:6443", cfg.Spec.Endpoint.Intranet),
				CertificateAuthorityData: cfg.Spec.Kubernetes.RootCA.Cert,
			},
		},
	}
	data, err := clientcmd.Write(bcfg)
	if err != nil {
		return errors.Wrapf(err, "write config")
	}
	cm := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-info",
			Namespace: metav1.NamespacePublic,
		},
		Data: map[string]string{
			"kubeconfig": string(data),
		},
	}
	return utils.ApplyYaml(utils.PrettyYaml(cm), "cluster-info")
}

func WriteClusterCR(ctx *context.NodeContext) error {

	cfg := ctx.NodeObject().Status.BootCFG
	m := ctx.NodeObject()
	node := v12.Master{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Master",
			APIVersion: v12.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: m.Spec.ID,
		},
		Spec: m.Spec,
	}
	klog.Infof("bind to infra: [%s]", cfg.Spec.Bind.ResourceId)
	return utils.ApplyYaml(
		strings.Join(
			[]string{
				utils.PrettyYaml(cfg),
				utils.PrettyYaml(node),
			}, "---\n",
		), "cluster-crd",
	)
}

func doRunWdrip(ctx *v12.ClusterSpec) error {
	cfg, err := RenderWdripYaml(ctx)
	if err != nil {
		return fmt.Errorf("write wdrip yaml: %s", err.Error())
	}
	return wait.Poll(
		2*time.Second,
		1*time.Minute,
		func() (done bool, err error) {
			if err := BootCFG(ctx); err != nil {
				klog.Errorf("retry upload bootcfg fail: %s", err.Error())
				return false, nil
			}
			if err := utils.ApplyYaml(cfg, "wdrip"); err != nil {
				klog.Errorf("retry wait for wdrip addon: %s", err.Error())
				return false, nil
			}
			return true, nil
		},
	)
}

func RunWdrip(ctx *v12.ClusterSpec) error {
	err := doRunWdrip(ctx)
	if err != nil {
		return err
	}
	return RunMonitor(ctx)
}

func BootCFG(spec *v12.ClusterSpec) error {
	bootcfg, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("marshal bootcfg: %s", err.Error())
	}
	cm := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bootcfg",
			Namespace: "kube-system",
		},
		Data: map[string][]byte{
			"bootcfg": bootcfg,
		},
	}

	cmdata, err := yaml.Marshal(cm)
	if err != nil {
		return fmt.Errorf("marshal cm: %s", err.Error())
	}
	return utils.ApplyYaml(string(cmdata), "bootcfg")
}

func RenderWdripYaml(spec *v12.ClusterSpec) (string, error) {
	t, err := template.New("wdrip-file").Parse(wdripf)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse config template")
	}

	// execute the template
	var buff bytes.Buffer
	err = t.Execute(
		&buff,
		struct {
			Version  string
			Registry string
			UUID     string
		}{
			Version:  wdrip.Version,
			Registry: fmt.Sprintf("%s/aoxn", filepath.Dir(spec.Registry)),
			//Registry: "registry.cn-hangzhou.aliyuncs.com/aoxn",
			UUID: uuid.New().String(),
		},
	)
	return buff.String(), err
}

var (
	wdripf = `
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin
  namespace: kube-system
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: wdrip
  name: wdrip
  namespace: kube-system
spec:
  ports:
    - name: tcp
      nodePort: 32443
      port: 9443
      protocol: TCP
      targetPort: 443
  selector:
    app: wdrip
  sessionAffinity: None
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: wdrip
    random.uuid: "{{ .UUID }}"
  name: wdrip
  namespace: kube-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: wdrip
  template:
    metadata:
      labels:
        app: wdrip
        random.uuid: "{{ .UUID }}"
    spec:
      hostNetwork: true
      priorityClassName: system-node-critical
      serviceAccount: admin
      containers:
        - image: {{ .Registry }}/wdrip:{{ .Version }}
          imagePullPolicy: Always
          name: wdrip-net
          command:
            - /wdrip
            - operator
            # - --bootcfg=/etc/wdrip/boot.cfg
          volumeMounts:
            - name: bootcfg
              mountPath: /etc/wdrip/
              readOnly: true
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
      - effect: NoSchedule
        operator: Exists
        key: node-role.kubernetes.io/master
      - effect: NoSchedule
        operator: Exists
        key: node.cloudprovider.kubernetes.io/uninitialized
      - effect: NoSchedule
        key: node.kubernetes.io/not-ready
        operator: Exists
      - effect: NoSchedule
        key: node.kubernetes.io/unreachable
        operator: Exists
      volumes:
        - name: bootcfg
          secret:
            secretName: bootcfg
            items:
              - key: bootcfg
                path: boot.cfg
`
)
