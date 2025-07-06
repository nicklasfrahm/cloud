{{/* Create chart name and version as used by the chart label. */}}
{{- define "gateway.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Common labels. */}}
{{- define "gateway.labels" -}}
helm.sh/chart: {{ include "gateway.chart" . | quote }}
{{ include "gateway.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- if .Values.commonLabels }}
{{ toYaml .Values.commonLabels }}
{{- end }}
{{- end }}

{{/* Selector labels. */}}
{{- define "gateway.selectorLabels" -}}
app.kubernetes.io/name: {{ .Chart.Name | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end }}

{{/* Create secret names for an FQDN */}}
{{- define "gateway.secretName" -}}
{{ . | trimSuffix "." | replace "." "-" | printf "%s-tls" | quote }}
{{- end }}
