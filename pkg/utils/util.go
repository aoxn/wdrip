package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	validate "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/klog/v2"
	"math/big"
	"net"
	"os"
	"text/template"
)

var (
	CA          = "ca"
	ETCD_PEER   = "etcd.peer"
	ETCD_CLIENT = "etcd.client"
)

type Errors []error

func (e Errors) Error() string {
	result := ""
	for _, err := range e {
		if err == nil {
			klog.Errorf("nil error")
			continue
		}
		result += err.Error() + "\n"
	}
	return result
}

func (e Errors) HasError() error {
	if len(e) != 0 {
		return e
	}
	return nil
}

func PrettyYaml(obj interface{}) string {
	bs, err := yaml.Marshal(obj)
	if err != nil {
		fmt.Errorf("failed to parse yaml, ' %s'", err.Error())
	}
	return string(bs)
}

func PrettyJson(obj interface{}) string {
	pretty := bytes.Buffer{}
	data, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("PrettyJson, mashal error: %s", err.Error())
		return ""
	}
	err = json.Indent(&pretty, data, "", "    ")

	if err != nil {
		fmt.Printf("PrettyJson, indent error: %s", err.Error())
		return ""
	}
	return pretty.String()
}

func FileExist(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetDNSIP(subnet string, index int) (net.IP, error) {
	_, cidr, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse service subnet CIDR %q: %v", subnet, err)
	}

	bip := big.NewInt(0).SetBytes(cidr.IP.To4())
	ip := net.IP(big.NewInt(0).Add(bip, big.NewInt(int64(index))).Bytes())
	if cidr.Contains(ip) {
		return ip, nil
	}
	return nil, fmt.Errorf("can't generate IP with "+
		"index %d from subnet. subnet too small. subnet: %q", index, subnet)
}

func RenderConfig(
	tplName string,
	tpl string,
	data interface{},
) (string, error) {
	t, err := template.New(tplName).Parse(tpl)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse config template")
	}

	// execute the template
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	return buff.String(), err
}

var KubeConfigTpl = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: {{ .AuthCA }}
    server: https://{{ .Address }}:6443
  name: kubernetes-clusterid-demo
contexts:
- context:
    cluster: kubernetes-clusterid-demo
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes-clusterid-demo
current-context: kubernetes-admin@kubernetes-clusterid-demo
kind: Config
preferences: {}
users:
- name: kubernetes-admin
  user:
    client-certificate-data: {{ .ClientCRT }}
    client-key-data: {{ .ClientKey }}
`

func GetNamePrefix(p string) string {
	// use the dash (if the name isn't too long) to make the pod name a bit prettier
	prefix := fmt.Sprintf("%s-", p)
	if len(validate.NameIsDNSSubdomain(prefix, true)) != 0 {
		prefix = prefix
	}
	return prefix
}

const NODE_MASTER_LABEL = "node-role.kubernetes.io/master"

func NodeIsMaster(node *corev1.Node) bool {
	labels := node.Labels
	if _, ok := labels[NODE_MASTER_LABEL]; ok {
		return true
	}
	return false
}
