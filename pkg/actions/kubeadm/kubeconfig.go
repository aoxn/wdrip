//go:build linux || darwin || windows
// +build linux darwin windows

package kubeadm

import (
	"encoding/base64"
	"fmt"
	"github.com/aoxn/ovm/pkg/actions"
	"github.com/aoxn/ovm/pkg/utils"
	"github.com/aoxn/ovm/pkg/utils/sign"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
)

type ActionKubeAuth struct {
}

// NewActionKubeAuth returns a new ActionKubeAuth for kubeadm init
func NewActionKubeAuth() actions.Action {
	return &ActionKubeAuth{}
}

// Execute runs the NewActionKubeAuth
func (a *ActionKubeAuth) Execute(ctx *actions.ActionContext) error {

	node := ctx.NodeObject()
	if node == nil {
		return fmt.Errorf("node info nil: ActionKubelet")
	}
	klog.Info("try write admin.local auth config")
	err := os.MkdirAll("/etc/kubernetes/", 0755)
	if err != nil {
		return fmt.Errorf("ensure dir /etc/kubernetes for admin.local:%s", err.Error())
	}

	key, crt, err := sign.SignKubernetesClient(
		node.Status.BootCFG.Spec.Kubernetes.RootCA.Cert,
		node.Status.BootCFG.Spec.Kubernetes.RootCA.Key, []string{},
	)
	if err != nil {
		return fmt.Errorf("sign kubernetes client crt: %s", err.Error())
	}
	err = os.MkdirAll("/etc/ovm", 0755)
	if err != nil {
		return fmt.Errorf("make ovm dir: %s", err.Error())
	}
	err = ioutil.WriteFile(
		"/etc/ovm/ovm.cfg.gen",
		[]byte(utils.PrettyYaml(ctx.Config())), 0755,
	)
	if err != nil {
		klog.Warningf("write bach config failed: %s", err.Error())
	}
	cfg, err := utils.RenderConfig(
		"admin.authconfig",
		utils.KubeConfigTpl,
		struct {
			AuthCA    string
			Address   string
			ClientCRT string
			ClientKey string
		}{
			AuthCA:    base64.StdEncoding.EncodeToString(node.Status.BootCFG.Spec.Kubernetes.RootCA.Cert),
			Address:   node.Spec.IP,
			ClientCRT: base64.StdEncoding.EncodeToString(crt),
			ClientKey: base64.StdEncoding.EncodeToString(key),
		},
	)
	if err != nil {
		return fmt.Errorf("render admin.local config error: %s", err.Error())
	}
	return ioutil.WriteFile(utils.AUTH_FILE, []byte(cfg), 0755)
}