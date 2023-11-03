# Configuration de Gangway

Gangway lit un fichier de configuration au démarrage. Le chemin vers le fichier de configuration doit être défini en utilisant le drapeau `--config`.

Le fichier de configuration doit être au format YAML et contenir un dictionnaire (alias hash ou map) de paires clé/valeur. Les options disponibles sont décrites ci-dessous.

## Options de Configuration

Les options suivantes peuvent être définies via le fichier de configuration YAML.

### Configuration Générale

| Clé                      | Description                                                                                                   |
|--------------------------|---------------------------------------------------------------------------------------------------------------|
| `host`                   | L'adresse d'écoute. Par défaut à `0.0.0.0` (toutes les interfaces).                                           |
| `port`                   | Le port d'écoute. Par défaut à `8080`.                                                                        |
| `serveTLS`               | Gangway doit-il utiliser TLS au lieu de HTTP simple ? Par défaut à `false`.                                    |
| `certFile`               | Le fichier de certificat public à utiliser lors de l'utilisation de TLS. Par défaut à `/etc/gangway/tls/tls.crt`.|
| `keyFile`                | Le fichier de clé privée lors de l'utilisation de TLS. Par défaut à `/etc/gangway/tls/tls.key`.                |
| `trustedCAPath`          | Le chemin vers un CA racine à faire confiance pour les certificats auto-signés aux URL Oauth2.                |
| `httpPath`               | Le chemin utilisé par gangway pour créer des URL. Par défaut à `""`.                                           |
| `sessionSecurityKey`     | La clé de sécurité de session.                                                                                 |
| `sessionSalt`            | Le sel de la session. Valeur codée par défaut à `MkmfuPNHnZBBivy0L0aW`.                                       |
| `customHTMLTemplatesDir` | Le chemin vers un répertoire contenant des modèles HTML personnalisés.                                        |
| `customAssetsDir`        | Le chemin vers un répertoire contenant des ressources.                                                         |

### Configuration Multi-Cluster

La configuration multi-cluster permet d'avoir des configurations spécifiques pour chaque cluster dans un seul fichier.

| Clé                   | Description                                                                                                   |
|-----------------------|---------------------------------------------------------------------------------------------------------------|
| `clusters`            | Un dictionnaire de configurations par nom de cluster, où chaque cluster a sa propre liste d'options.          |

### Configuration de Cluster Spécifique

Chaque cluster peut avoir les configurations suivantes :

| Clé                          | Description                                                                                                   |
|------------------------------|---------------------------------------------------------------------------------------------------------------|
| `clusterName`                | Le nom du cluster. Utilisé dans l'UI et les instructions de configuration de kubectl.                         |
| `providerURL`                | URL du fournisseur OAuth2. Doit offrir un point de terminaison `$providerURL/.well-known/openid-configuration` pour la découverte.|
| `clientID`                   | ID client de l'API tel qu'indiqué par le fournisseur d'identité.                                              |
| `clientSecret`               | Secret client de l'API tel qu'indiqué par le fournisseur d'identité.                                          |
| `allowEmptyClientSecret`     | Certains fournisseurs d'identité acceptent un secret client vide, ce qui n'est généralement pas une bonne idée. Si vous devez utiliser un secret vide et accepter les risques qui l'accompagnent, alors vous pouvez le définir sur `true`. Par défaut à `false`.|
| `audience`                   | Point de terminaison qui fournit des informations sur le profil utilisateur [optionnel]. Non requis par tous les fournisseurs.|
| `scopes`                     | Utilisé pour spécifier la portée de l'autorisation OAuth demandée. Par défaut à `["openid", "profile", "email", "offline_access"]`.|
| `redirectURL`                | Où rediriger après l'authentification. Cela doit être une URL où gangway est accessible. Typiquement, cela doit aussi être enregistré dans l'application OAuth avec le fournisseur OAuth.|
| `usernameClaim`              | La revendication JWT à utiliser comme nom d'utilisateur. Ceci est utilisé dans l'UI. Combiné avec le clusterName pour la partie "utilisateur" de kubeconfig. Par défaut à `nickname`.|
| `emailClaim`                 | Obsolète. Par défaut à `email`.                                                                              |
| `apiServerURL`               | Le point de terminaison du serveur API utilisé pour configurer kubectl.                                       |
| `clusterCAPath`              | Le chemin pour trouver le paquet CA pour le serveur API. Utilisé pour configurer kubectl. Ce chemin est typiquement monté dans l'emplacement par défaut pour les charges de travail fonctionnant sur un cluster Kubernetes et n'a généralement pas besoin d'être défini. Par défaut à `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt`.|
| `showClaims`                 | Affiche les revendications reçues. Par défaut à `true`.                                                       |

### Validation et Préfixes

Les configurations de cluster spécifiques peuvent être validées et des préfixes d'environnement peuvent être appliqués pour permettre la configuration à partir de variables d'environnement.

## Fonctions Supplémentaires

- `NewMultiClusterConfig`: Crée une nouvelle instance de configuration multi-cluster à partir d'un fichier de configuration sérialisé.
- `Validate`: Vérifie toutes les propriétés de la structure de configuration pour s'assurer qu'elles sont initialisées.
- `GetRootPathPrefix`: Renvoie '/' si aucun préfixe n'est spécifié, sinon renvoie le chemin configuré.
- `loadCerts`: Charge les certificats pour les configurations de cluster à partir des chemins spécifiés.

## Utilisation des Variables d'Environnement

Les variables d'environnement peuvent être utilisées pour outrepasser les configurations spécifiées dans le fichier YAML en utilisant le préfixe `GANGWAY` suivi du nom de la clé correspondante en majuscules et des scores pour les espaces.

Exemple: Pour outrepasser `clientSecret`, utilisez la variable d'environnement `GANGWAY_CLIENT_SECRET`.

## Example de configuration Yaml

```yaml
host: "0.0.0.0"
port: 8080
serveTLS: false
certFile: "/chemin/vers/tls.crt"
keyFile: "/chemin/vers/tls.key"
clusters:
  production:
    - clusterName: "production"
      providerURL: "https://oauth.example.com"
      clientID: "mon-client-id"
      clientSecret: "mon-secret-client"
      ...
  staging:
    - clusterName: "staging"
      providerURL: "https://oauth-staging.example.com"
      clientID: "mon-client-id-staging"
      clientSecret: "mon-secret-client-staging"
      ...
httpPath: "/gangway"
customHTMLTemplatesDir: "/chemin/vers/templates"
customAssetsDir: "/chemin/vers/assets"

```