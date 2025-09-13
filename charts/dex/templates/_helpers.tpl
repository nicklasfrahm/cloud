{{/* Create chart name and version as used by the chart label. */}}
{{- define "dex.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/* Common labels. */}}
{{- define "dex.labels" -}}
helm.sh/chart: {{ include "dex.chart" . | quote }}
{{ include "dex.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- if .Values.commonLabels }}
{{ toYaml .Values.commonLabels }}
{{- end }}
{{- end }}

{{/* Selector labels. */}}
{{- define "dex.selectorLabels" -}}
app.kubernetes.io/name: {{ .Chart.Name | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end }}

{{/* Create secret names for an FQDN */}}
{{- define "dex.secretName" -}}
{{ . | trimSuffix "." | replace "." "-" | printf "%s-tls" | quote }}
{{- end }}
