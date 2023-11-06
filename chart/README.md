# Helm Chart Configuration for Gangway

This document outlines the configurable parameters of the Helm chart for Gangway and their default values.

## Global Parameters

- `global.imageRegistry`: Global Docker image registry (default: `""`)
- `global.imagePullSecrets`: Global Docker registry secret names as an array (default: `[]`)
- `global.storageClass`: Global StorageClass for Persistent Volume(s) (default: `""`)

## Common Parameters

- `kubeVersion`: Force target Kubernetes version (default: `""`)
- `nameOverride`: String to partially override `common.names.fullname` template (default: `""`)
- `fullnameOverride`: String to fully override `common.names.fullname` template (default: `""`)
- `commonLabels`: Labels to add to all deployed objects (default: `{}`)
- `commonAnnotations`: Annotations to add to all deployed objects (default: `{}`)
- `clusterDomain`: Default Kubernetes cluster domain (default: `cluster.local`)
- `extraDeploy`: Array of extra objects to deploy with the release (default: `[]`)

## Gangway Image Parameters

- `image.registry`: Gangway image registry (default: `docker.io`)
- `image.repository`: Gangway image repository (default: `numberlyinfra/gangway`)
- `image.tag`: Gangway image tag (default: `master`)
- `image.digest`: Gangway image digest (default: `""`)
- `image.pullPolicy`: Image pull policy (default: `Always`)
- `image.pullSecrets`: Specify docker-registry secret names as an array (default: `[]`)
- `image.debug`: Enable Gangway image debug mode (default: `false`)

## Gangway Deployment Parameters

- `replicaCount`: Number of Gangway replicas to deploy (default: `3`)
- `updateStrategy.type`: StrategyType can be set to `RollingUpdate` or `OnDelete` (default: `RollingUpdate`)
- `podSecurityContext.enabled`: Enable Gangway pods' Security Context (default: `true`)
- `podSecurityContext.fsGroup`: Group ID for the volumes of the pod (default: `1001`)
- `containerSecurityContext.enabled`: Enable Gangway containers' SecurityContext (default: `true`)
- `containerSecurityContext.runAsUser`: User ID to run Gangway containers (default: `1001`)
- `containerSecurityContext.runAsNonRoot`: Set Gangway container's Security Context runAsNonRoot (default: `true`)
- `resources.limits`: The resources limits for the Gangway container (default: `{}`)
- `resources.requests`: The requested resources for the Gangway container (default: `{ memory: "512Mi", cpu: "300m" }`)

## Probes

- `startupProbe.enabled`: Enable startupProbe (default: `false`)
- `livenessProbe.enabled`: Enable livenessProbe (default: `true`)
- `readinessProbe.enabled`: Enable readinessProbe (default: `true`)

## Service Configuration

- `service.type`: Kubernetes Service type (default: `ClusterIP`)
- `service.ports.http`: Service HTTP port (default: `8080`)
- `service.sessionAffinity`: Session Affinity for Kubernetes service (default: `ClientIP`)

## Ingress Configuration

- `ingress.enabled`: Set to true to enable ingress record generation (default: `false`)
- `ingress.hostname`: Default host for the ingress resource (default: `gangway.local`)
- `ingress.path`: The Path to Gangway (default: `/`)
- `ingress.annotations`: Additional annotations for the Ingress resource (default: `{}`)

## Additional Configuration

- `extraEnvVars`: Extra environment variables to be set on Gangway container (default: `[]`)
- `extraVolumes`: Optionally specify extra list of additional volumes for Gangway pods (default: `[]`)
- `extraVolumeMounts`: Optionally specify extra list of additional volumeMounts for Gangway container(s) (default: `[]`)

## Affinity and Tolerations

- `podAffinityPreset`: Pod affinity preset (default: `""`)
- `podAntiAffinityPreset`: Pod anti-affinity preset (default: `soft`)
- `nodeAffinityPreset.type`: Node affinity preset type (default: `""`)
- `affinity`: Affinity for pod assignment (default: `{}`)
- `nodeSelector`: Node labels for pod assignment (default: `{}`)
- `tolerations`: Tolerations for pod assignment (default: `[]`)

## Customization

- `command`: Override default container command (default: `[]`)
- `args`: Override default container args (default: `[]`)
- `lifecycleHooks`: for the Gangway container(s) to automate configuration before or after startup (default: `{}`)
- `extraEnvVarsCM`: Name of existing ConfigMap containing extra env vars (default: `""`)
- `extraEnvVarsSecret`: Name of existing Secret containing extra env vars (default: `""`)

## Session and Configuration

- `sessionkey`: Session key (default: `mySessionKey`)
- `sessionsalt`: Session salt (default: `mySessionSalt`)
- `configuration`: Inline configuration for Gangway (default: `""`)

## Example of configuration : 
```yaml
host: 0.0.0.0
port: 8080
serveTLS: false
clusters:
  Production:
  - EnvPrefix: kube01
    apiServerURL: https://kube01-api-url:443
    audience: xxxxxxxxxxx-xxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com
    providerURL: https://accounts.google.com
    clientID: xxxxxxxxxxx-xxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com
    clientSecret: GXXXX-XXXXXXXXXXXXXXXXXXX
    clusterName: kube01
    emailClaim: email
    redirectURL: https://gangway.local/callback
    scopes:
    - openid
    - profile
    - email
    tokenURL: https://www.googleapis.com/oauth2/v4/token
    usernameClaim: email
  Development:
  - EnvPrefix: kube02
    apiServerURL: https://kube02-api-url:443
    audience: xxxxxxxxxxx-xxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com
    providerURL: https://accounts.google.com
    clientID: xxxxxxxxxxx-xxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com
    clientSecret: GXXXX-XXXXXXXXXXXXXXXXXXX
    clusterName: kube02
    emailClaim: email
    redirectURL: https://gangway.local/callback
    scopes:
    - openid
    - profile
    - email
    tokenURL: https://www.googleapis.com/oauth2/v4/token
    usernameClaim: email
  - EnvPrefix: kube03
    apiServerURL: https://kube03-api-url:443
    audience: xxxxxxxxxxx-xxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com
    providerURL: https://accounts.google.com
    clientID: xxxxxxxxxxx-xxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com
    clientSecret: GXXXX-XXXXXXXXXXXXXXXXXXX
    clusterName: kube03
    emailClaim: email
    redirectURL: https://gangway.local/callback
    clusterCAPath: /etc/gangly/pki/kube03/ca.crt
    scopes:
    - openid
    - profile
    - email
    tokenURL: https://www.googleapis.com/oauth2/v4/token
    usernameClaim: email
```
To get more informations : [Configuration](../docs/configuration.md)