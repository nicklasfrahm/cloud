{{- define "talos.config" }}
machine:
  type: {{ .MachineType }}
  kubelet:
    nodeIP:
      validSubnets:
        {{- toYaml .Values.advertisedSubnets | nindent 8 }}
  install:
    {{- (include "talm.discovered.disks_info" .) | nindent 4 }}
    disk: {{ include "talm.discovered.system_disk_name" . | quote }}
  network:
    hostname: {{ include "talm.discovered.hostname" . | quote }}
    nameservers: {{ include "talm.discovered.default_resolvers" . }}
    {{- (include "talm.discovered.physical_links_info" .) | nindent 4 }}
    interfaces:
    - deviceSelector:
        {{- include "talm.discovered.default_link_selector_by_gateway" . | nindent 8 }}
      addresses: {{ include "talm.discovered.default_addresses_by_gateway" . }}
      routes:
        - network: 0.0.0.0/0
          gateway: {{ include "talm.discovered.default_gateway" . }}
      {{- with .Values.floatingIP }}
      vip:
        ip: {{ . }}
      {{- end }}

cluster:
  allowSchedulingOnControlPlanes: true
  network:
    podSubnets:
      {{- toYaml .Values.podSubnets | nindent 6 }}
    serviceSubnets:
      {{- toYaml .Values.serviceSubnets | nindent 6 }}
    {{- if .Values.cilium.enabled }}
    cni:
      name: none
    {{- end }}
  clusterName: "{{ .Chart.Name }}"
  controlPlane:
    endpoint: "{{ .Values.endpoint }}"
  {{- if eq .MachineType "controlplane" }}
  etcd:
    advertisedSubnets:
      {{- toYaml .Values.advertisedSubnets | nindent 6 }}
  {{- end }}
  {{- if .Values.cilium.enabled }}
  inlineManifests:
    - name: cilium-install
      contents: |
        ---
        apiVersion: rbac.authorization.k8s.io/v1
        kind: ClusterRoleBinding
        metadata:
          name: cilium-install
        roleRef:
          apiGroup: rbac.authorization.k8s.io
          kind: ClusterRole
          name: cluster-admin
        subjects:
        - kind: ServiceAccount
          name: cilium-install
          namespace: kube-system
        ---
        apiVersion: v1
        kind: ServiceAccount
        metadata:
          name: cilium-install
          namespace: kube-system
        ---
        apiVersion: batch/v1
        kind: Job
        metadata:
          name: cilium-install
          namespace: kube-system
        spec:
          backoffLimit: 10
          template:
            metadata:
              labels:
                app: cilium-install
            spec:
              restartPolicy: OnFailure
              tolerations:
                - operator: Exists
                - effect: NoSchedule
                  operator: Exists
                - effect: NoExecute
                  operator: Exists
                - effect: PreferNoSchedule
                  operator: Exists
                - key: node-role.kubernetes.io/control-plane
                  operator: Exists
                  effect: NoSchedule
                - key: node-role.kubernetes.io/control-plane
                  operator: Exists
                  effect: NoExecute
                - key: node-role.kubernetes.io/control-plane
                  operator: Exists
                  effect: PreferNoSchedule
              affinity:
                nodeAffinity:
                  requiredDuringSchedulingIgnoredDuringExecution:
                    nodeSelectorTerms:
                      - matchExpressions:
                          - key: node-role.kubernetes.io/control-plane
                            operator: Exists
              serviceAccount: cilium-install
              serviceAccountName: cilium-install
              hostNetwork: true
              containers:
                - name: cilium-install
                  image: alpine/helm:3.12.3
                  env:
                    - name: KUBERNETES_SERVICE_HOST
                      valueFrom:
                        fieldRef:
                          apiVersion: v1
                          fieldPath: status.podIP
                    - name: KUBERNETES_SERVICE_PORT
                      value: "6443"
                  command:
                    - /bin/sh
                    - -c
                    - |
                      helm repo add cilium https://helm.cilium.io
                      helm repo update
                      helm upgrade --install cilium cilium/cilium \
                        --namespace kube-system \
                        --set ipam.mode=kubernetes \
                        --set kubeProxyReplacement=true \
                        --set securityContext.capabilities.ciliumAgent="{CHOWN,KILL,NET_ADMIN,NET_RAW,IPC_LOCK,SYS_ADMIN,SYS_RESOURCE,DAC_OVERRIDE,FOWNER,SETGID,SETUID}" \
                        --set securityContext.capabilities.cleanCiliumState="{NET_ADMIN,SYS_ADMIN,SYS_RESOURCE}" \
                        --set cgroup.autoMount.enabled=false \
                        --set cgroup.hostRoot=/sys/fs/cgroup \
                        --set k8sServiceHost=localhost \
                        --set k8sServicePort=7445 \
                        --set operator.replicas={{ .Values.cilium.operator.replicas }}
  {{- end }}
{{- end }}
