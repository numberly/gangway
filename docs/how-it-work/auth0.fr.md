# Connexion de Gangway à Auth0

1. Créez un compte pour Auth0 et connectez-vous.
2. Depuis le tableau de bord, cliquez sur "Nouvelle Application".
3. Entrez un nom et choisissez "Applications Web à Page Unique".
4. Cliquez sur "Paramètres" et rassemblez les informations pertinentes, puis mettez à jour le fichier `docs/yaml/02-configmap.yaml` et appliquez-le au cluster.
5. Mettez à jour les "URL de Rappel Autorisées" pour correspondre au paramètre "redirectURL" dans la configmap configurée précédemment.
6. Cliquez sur "Enregistrer les Changements".
7. Ajoutez une règle pour ajouter des métadonnées de groupe en cliquant sur "Règles" dans le menu.
8. Donnez un nom à la règle et copiez/collez ce qui suit :

    ```go
    function (user, context, callback) {
        if (user.app_metadata && 'groups' in user.app_metadata) {
            context.idToken.groups = user.app_metadata.groups;
        } else {
            context.idToken.groups = [];
        }

    callback(null, user, context);
    }
    ```

9. Configurez le serveur API avec la configuration suivante en remplaçant les valeurs issuer-url et client-id :

    ```
    --oidc-issuer-url=https://example.auth0.com/
    --oidc-client-id=<clientid>
    --oidc-username-claim=email
    --oidc-groups-claim=groups
    ```

## Exemple

Une configuration typique de gangway pour Auth0 :

```yaml
clusterName: "VotreCluster"
providerURL: "https://example.auth0.com"
clientID: "<votre ID client>"
clientSecret: "<votre secret client>"
audience: "https://example.auth0.com/userinfo"
redirectURL: "https://gangway.example.com/callback"
scopes: ["openid", "profile", "email", "offline_access"]
usernameClaim: "sub"
emailClaim: "email"
apiServerURL: "https://kube-apiserver.votrecluster.com"
```