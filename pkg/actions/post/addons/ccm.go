package addons

var CCM = ConfigTpl{
	Name:         "ccm",
	Tpl:          ccmtpl,
	ImageVersion: "v1.9.3.380-gd6d0962-aliyun",
}

var ccmtpl = `
{{ if eq .Action "Install" }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloud-controller-manager
  namespace: kube-system
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:cloud-controller-manager
rules:
  - apiGroups:
      - ""
    resources:
      - persistentvolumes
      - services
      - secrets
      - endpoints
      - serviceaccounts
    verbs:
      - get
      - list
      - watch
      - create
      - update
      - patch
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - get
      - list
      - watch
      - delete
      - patch
      - update
  - apiGroups:
      - ""
    resources:
      - services/status
    verbs:
      - update
      - patch
  - apiGroups:
      - ""
    resources:
      - nodes/status
    verbs:
      - patch
      - update
  - apiGroups:
      - ""
    resources:
      - events
      - endpoints
    verbs:
      - create
      - patch
      - update
  - apiGroups:
      - coordination.k8s.io
    resources:
      - leases
    verbs:
      - get
      - update
      - create
      - delete
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:cloud-controller-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:cloud-controller-manager
subjects:
- kind: ServiceAccount
  name: cloud-controller-manager
  namespace: kube-system
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:shared-informers
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:cloud-controller-manager
subjects:
- kind: ServiceAccount
  name: shared-informers
  namespace: kube-system
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:cloud-node-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:cloud-controller-manager
subjects:
- kind: ServiceAccount
  name: cloud-node-controller
  namespace: kube-system
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:pvl-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:cloud-controller-manager
subjects:
- kind: ServiceAccount
  name: pvl-controller
  namespace: kube-system
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:route-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:cloud-controller-manager
subjects:
- kind: ServiceAccount
  name: route-controller
  namespace: kube-system
---
{{ end }}
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    app: cloud-controller-manager
    tier: control-plane
  name: cloud-controller-manager
  namespace: kube-system
spec:
  selector:
    matchLabels:
      app: cloud-controller-manager
      tier: control-plane
  updateStrategy:
    rollingUpdate:
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: cloud-controller-manager
        tier: control-plane
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      serviceAccountName: cloud-controller-manager
      tolerations:
      - operator: Exists
      nodeSelector:
        node-role.kubernetes.io/master: ""
      containers:
      - command:
        -  /cloud-controller-manager
        - --kubeconfig=/etc/kubernetes/cloud-controller-manager.conf
        - --address=127.0.0.1
{{ if eq .CIDR "" }}
        - --configure-cloud-routes=false
{{ else }}
        - --cluster-cidr={{.CIDR}}
        - --configure-cloud-routes=true
        - --route-reconciliation-period=3m
{{ end }}
        - --allow-untagged-cloud=true
        - --leader-elect=true
        - --cloud-provider=alicloud
        - --allocate-node-cidrs=true
        - --use-service-account-credentials=true
        - --v=5
        - --feature-gates=ServiceNodeExclusion=true
        image: registry-vpc.{{.Region}}.aliyuncs.com/acs/cloud-controller-manager-amd64:{{.ImageVersion}}
        livenessProbe:
          failureThreshold: 8
          httpGet:
            host: 127.0.0.1
            path: /healthz
            port: 10258
            scheme: HTTP
          initialDelaySeconds: 15
          timeoutSeconds: 15
        name: cloud-controller-manager
        resources:
          requests:
            cpu: 200m
        volumeMounts:
        - mountPath: /etc/kubernetes/
          name: k8s
          readOnly: true
        - mountPath: /etc/ssl/certs
          name: certs
        - mountPath: /etc/pki
          name: pki
      hostNetwork: true
      volumes:
      - hostPath:
          path: /etc/kubernetes
        name: k8s
      - hostPath:
          path: /etc/ssl/certs
        name: certs
      - hostPath:
          path: /etc/pki
        name: pki
`
