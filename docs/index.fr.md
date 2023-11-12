# Documentation de Gangway

gangway
=======

_(noun): An opening in the bulwark of the ship to allow passengers to board or leave the ship._

Une application qui peut être utilisée pour faciliter les flux d'authentification via OIDC pour un cluster Kubernetes.
Kubernetes prend en charge les [jetons OpenID Connect](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) comme moyen d'identifier les utilisateurs qui accèdent au cluster.
Gangway a été amélioré et est maintenant capable de gérer plusieurs clusters.
Gangway permet aux utilisateurs de configurer eux-mêmes leur configuration `kubectl` en quelques étapes simples.

![gangway multicluster](images/gangway-multicluster.png)

Une fois authentifié pour l'un de vos clusters :

![gangway](images/gangway.png)
