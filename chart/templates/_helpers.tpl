{{/*
Copyright VMware, Inc.
SPDX-License-Identifier: APACHE-2.0
*/}}

{{/* vim: set filetype=mustache: */}}

{{/*
Get the user defined LoadBalancerIP for this release.
Note, returns 127.0.0.1 if using ClusterIP.
*/}}
{{- define "gangway.serviceIP" -}}
{{- if eq .Values.service.type "ClusterIP" -}}
127.0.0.1
{{- else -}}
{{- .Values.service.loadBalancerIP | default "" -}}
{{- end -}}
{{- end -}}

{{/*
Gets the host to be used for this application.
If not using ClusterIP, or if a host or LoadBalancerIP is not defined, the value will be empty.
When using Ingress, it will be set to the Ingress hostname.
*/}}
{{- define "gangway.host" -}}
{{- $host := index .Values (printf "%sHost" .Chart.Name) | default "" -}}
{{- if .Values.ingress.enabled }}
{{- $host := .Values.ingress.hostname | default "" -}}
{{- end -}}
{{- default (include "gangway.serviceIP" .) $host -}}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "gangway.imagePullSecrets" -}}
{{- include "common.images.pullSecrets" (dict "images" (list .Values.image ) "global" .Values.global) -}}
{{- end -}}

{{/*
Return the proper gangway image name
*/}}
{{- define "gangway.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.image "global" .Values.global) }}
{{- end -}}

{{/*
Check if there are rolling tags in the images
*/}}
{{- define "gangway.checkRollingTags" -}}
{{- include "common.warnings.rollingTag" .Values.image }}
{{- end -}}

{{/*
Validate values of password and existing Secret mutually exclusive
*/}}
{{- define "gangway.validateValues.passwordSecretMutExclusive" -}}
{{- if and .Values.gangwayPassword .Values.gangwaySecret -}}
gangway: gangwaySecret
    You cannot provide both gangwayPassword and gangwaySecret
{{- end -}}
{{- end -}}

{{/*
Compile all warnings into a single message, and call fail.
*/}}
{{- define "gangway.validateValues" -}}
{{- $messages := list -}}
{{- $messages := append $messages (include "gangway.validateValues.passwordSecretMutExclusive" .) -}}
{{- $messages := without $messages "" -}}
{{- $message := join "\n" $messages -}}

{{- if $message -}}
{{-   printf "\nVALUES VALIDATION:\n%s" $message | fail -}}
{{- end -}}
{{- end -}}