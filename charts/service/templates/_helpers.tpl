{{/* Create chart name and version as used by the chart label. */}}
{{- define "service.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Common labels. */}}
{{- define "service.labels" -}}
helm.sh/chart: {{ include "service.chart" . | quote }}
{{ include "service.selectorLabels" . }}
app.kubernetes.io/version: {{ .Values.image.tag | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
app.nicklasfrahm.dev/cluster: {{ .Values.platform.cluster | quote }}
app.nicklasfrahm.dev/environment: {{ .Values.platform.environment | quote }}
app.nicklasfrahm.dev/location: {{ .Values.platform.location | quote }}
app.nicklasfrahm.dev/tenant: {{ .Values.platform.tenant | quote }}
{{- if .Values.commonLabels }}
{{ toYaml .Values.commonLabels }}
{{- end }}
{{- end }}

{{/* Selector labels. */}}
{{- define "service.selectorLabels" -}}
app.kubernetes.io/name: {{ .Release.Name | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end }}
