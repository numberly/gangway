# Déploiement de Gangway

Le déploiement de Gangway consiste à écrire un fichier de configuration, puis à déployer le service.
Le service est sans état, il est donc relativement facile à gérer sur Kubernetes.
La manière dont vous fournirez l'accès au service dépendra de votre configuration spécifique.

Gangway est maintenant conscient d'une configuration multi-cluster, nous avons développé un chart Helm qui facilitera le déploiement.

Voici un exemple de configuration :

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
