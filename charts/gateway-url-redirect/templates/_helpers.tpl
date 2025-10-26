{{/*
Expand the name of the chart.
*/}}
{{- define "gateway-url-redirect.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gateway-url-redirect.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "gateway-url-redirect.labels" -}}
helm.sh/chart: {{ include "gateway-url-redirect.chart" . }}
{{ include "gateway-url-redirect.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "gateway-url-redirect.selectorLabels" -}}
app.kubernetes.io/name: {{ include "gateway-url-redirect.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
