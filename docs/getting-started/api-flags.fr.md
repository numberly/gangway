## Paramètres du serveur API

Gangway nécessite que le serveur API Kubernetes soit configuré pour OIDC :

[https://kubernetes.io/docs/admin/authentication/#configuring-the-api-server](https://kubernetes.io/docs/admin/authentication/#configuring-the-api-server)

```bash
kube-apiserver
...
--oidc-issuer-url="https://example.auth0.com/"
--oidc-client-id=3YM4ue8MoXgBkvCIHh00000000000
--oidc-username-claim=email
--oidc-groups-claim=groups
```