{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "cb-dragonfly.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cb-dragonfly.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cb-dragonfly.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "cb-dragonfly.labels" -}}
helm.sh/chart: {{ include "cb-dragonfly.chart" . }}
{{ include "cb-dragonfly.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "cb-dragonfly.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cb-dragonfly.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "cb-dragonfly.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "cb-dragonfly.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the Cluster Role to use
*/}}
{{- define "cb-dragonfly.clusterRoleName" -}}
{{- if .Values.clusterRole.create -}}
    {{ default (include "cb-dragonfly.fullname" .) .Values.clusterRole.name }}
{{- else -}}
    {{ default "default" .Values.clusterRole.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the Role Binding to use
*/}}
{{- define "cb-dragonfly.roleBindingName" -}}
{{- if .Values.roleBinding.create -}}
    {{ default (include "cb-dragonfly.fullname" .) .Values.roleBinding.name }}
{{- else -}}
    {{ default "default" .Values.roleBinding.name }}
{{- end -}}
{{- end -}}
