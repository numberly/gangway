# Configuration du Chart Helm pour Gangway

Ce document décrit les paramètres configurables du chart Helm pour Gangway et leurs valeurs par défaut.

## Paramètres Globaux

- `global.imageRegistry` : Registre global d'images Docker (par défaut : `""`)
- `global.imagePullSecrets` : Noms des secrets de registre Docker global sous forme de tableau (par défaut : `[]`)
- `global.storageClass` : Classe de stockage globale pour le(s) volume(s) persistant(s) (par défaut : `""`)

## Paramètres Communs

- `kubeVersion` : Forcer la version cible de Kubernetes (par défaut : `""`)
- `nameOverride` : Chaîne pour remplacer partiellement le modèle `common.names.fullname` (par défaut : `""`)
- `fullnameOverride` : Chaîne pour remplacer entièrement le modèle `common.names.fullname` (par défaut : `""`)
- `commonLabels` : Étiquettes à ajouter à tous les objets déployés (par défaut : `{}`)
- `commonAnnotations` : Annotations à ajouter à tous les objets déployés (par défaut : `{}`)
- `clusterDomain` : Domaine de cluster Kubernetes par défaut (par défaut : `cluster.local`)
- `extraDeploy` : Tableau d'objets supplémentaires à déployer avec la release (par défaut : `[]`)

## Paramètres de l'Image Gangway

- `image.registry` : Registre d'image Gangway (par défaut : `docker.io`)
- `image.repository` : Dépôt d'image Gangway (par défaut : `numberlyinfra/gangway`)
- `image.tag` : Tag de l'image Gangway (par défaut : `master`)
- `image.digest` : Digest de l'image Gangway (par défaut : `""`)
- `image.pullPolicy` : Politique de téléchargement de l'image (par défaut : `Always`)
- `image.pullSecrets` : Spécifier les noms des secrets de registre Docker sous forme de tableau (par défaut : `[]`)
- `image.debug` : Activer le mode debug de l'image Gangway (par défaut : `false`)

## Paramètres de Déploiement de Gangway

- `replicaCount` : Nombre de répliques de Gangway à déployer (par défaut : `3`)
- `updateStrategy.type` : StrategyType peut être défini sur `RollingUpdate` ou `OnDelete` (par défaut : `RollingUpdate`)
- `podSecurityContext.enabled` : Activer le contexte de sécurité des pods Gangway (par défaut : `true`)
- `podSecurityContext.fsGroup` : ID de groupe pour les volumes du pod (par défaut : `1001`)
- `containerSecurityContext.enabled` : Activer le contexte de sécurité des conteneurs Gangway (par défaut : `true`)
- `containerSecurityContext.runAsUser` : ID utilisateur pour exécuter les conteneurs Gangway (par défaut : `1001`)
- `containerSecurityContext.runAsNonRoot` : Définir le contexte de sécurité du conteneur Gangway runAsNonRoot (par défaut : `true`)
- `resources.limits` : Les limites de ressources pour le conteneur Gangway (par défaut : `{}`)
- `resources.requests` : Les ressources demandées pour le conteneur Gangway (par défaut : `{ memory: "512Mi", cpu: "300m" }`)

## Probes

- `startupProbe.enabled` : Activer startupProbe (par défaut : `false`)
- `livenessProbe.enabled` : Activer livenessProbe (par défaut : `true`)
- `readinessProbe.enabled` : Activer readinessProbe (par défaut : `true`)

## Configuration du Service

- `service.type` : Type de service Kubernetes (par défaut : `ClusterIP`)
- `service.ports.http` : Port HTTP du service (par défaut : `8080`)
- `service.sessionAffinity` : Affinité de session pour le service Kubernetes (par défaut : `ClientIP`)

## Configuration de l'Ingress

- `ingress.enabled` : Définir sur true pour activer la génération de l'enregistrement d'ingress (par défaut : `false`)
- `ingress.hostname` : Hôte par défaut pour la ressource ingress (par défaut : `gangway.local`)
- `ingress.path` : Le chemin vers Gangway (par défaut : `/`)
- `ingress.annotations` : Annotations supplémentaires pour la ressource Ingress (par défaut : `{}`)

## Configuration Supplémentaire

- `extraEnvVars` : Variables d'environnement supplémentaires à définir sur le conteneur Gangway (par défaut : `[]`)
- `extraVolumes` : Spécifier éventuellement une liste supplémentaire de volumes pour les pods Gangway (par défaut : `[]`)
- `extraVolumeMounts` : Spécifier éventuellement une liste supplémentaire de volumeMounts pour le(s) conteneur(s) Gangway (par défaut : `[]`)

## Affinité et Tolérances

- `podAffinityPreset` : Préréglage d'affinité de pod (par défaut : `""`)
- `podAntiAffinityPreset` : Préréglage d'anti-affinité de pod (par défaut : `soft`)
- `nodeAffinityPreset.type` : Type de préréglage d'affinité de nœud (par défaut : `""`)
- `affinity` : Affinité pour l'assignation de pod (par défaut : `{}`)
- `nodeSelector` : Étiquettes de nœud pour l'assignation de pod (par défaut : `{}`)
- `tolerations` : Tolérances pour l'assignation de pod (par défaut : `[]`)

## Personnalisation

- `command` : Remplacer la commande par défaut du conteneur (par défaut : `[]`)
- `args` : Remplacer les arguments par défaut du conteneur (par défaut : `[]`)
- `lifecycleHooks` : pour le(s) conteneur(s) Gangway afin d'automatiser la configuration avant ou après le démarrage (par défaut : `{}`)
- `extraEnvVarsCM` : Nom du ConfigMap existant contenant des variables d'environnement supplémentaires (par défaut : `""`)
- `extraEnvVarsSecret` : Nom du Secret existant contenant des variables d'environnement supplémentaires (par défaut : `""`)

## Session et Configuration

- `sessionkey` : Clé de session (par défaut : `mySessionKey`)
- `sessionsalt` : Sel de session (par défaut : `mySessionSalt`)
- `configuration` : Configuration en ligne pour Gangway (par défaut : `""`)

## Exemple de configuration : 
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
To get more informations : [Configuration](configuration.md)