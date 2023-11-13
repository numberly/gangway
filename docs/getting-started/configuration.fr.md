# Configuration de Gangway

Gangway lit un fichier de configuration au démarrage. Le chemin vers le fichier de configuration doit être défini à l'aide du drapeau `--config`.

Le fichier de configuration doit être au format YAML et contenir un dictionnaire (alias hash ou map) de paires clé/valeur. Les options disponibles sont décrites ci-dessous.

## Options de Configuration

Les options suivantes peuvent être définies via le fichier de configuration YAML.

### Configuration Générale

| Clé                      | Description                                                                                                    |
|--------------------------|----------------------------------------------------------------------------------------------------------------|
| `host`                   | L'adresse sur laquelle écouter. Par défaut à `0.0.0.0` (toutes les interfaces).                                |
| `port`                   | Le port sur lequel écouter. Par défaut à `8080`.                                                               |
| `serveTLS`               | Gangway doit-il utiliser TLS au lieu de HTTP simple ? Par défaut à `false`.                                    |
| `certFile`               | Le fichier de certificat public à utiliser lors de l'utilisation de TLS. Par défaut à `/etc/gangway/tls/tls.crt`.|
| `keyFile`                | Le fichier de clé privée lors de l'utilisation de TLS. Par défaut à `/etc/gangway/tls/tls.key`.                |
| `trustedCAPath`          | Chemin vers une CA racine de confiance pour les certificats auto-signés aux URL Oauth2.                        |
| `httpPath`               | Le chemin utilisé par gangway pour créer des URL. Par défaut à `""`, en supprimant tout slash final.           |
| `sessionSecurityKey`     | La clé de sécurité de session.                                                                                  |
| `sessionSalt`            | Le sel de session. Valeur par défaut codée en dur `MkmfuPNHnZBBivy0L0aW`.                                      |
| `customHTMLTemplatesDir` | Chemin vers un répertoire contenant des modèles HTML personnalisés.                                            |
| `customAssetsDir`        | Chemin vers un répertoire contenant des actifs.                                                                |

### Configuration Multi-Cluster

La configuration multi-cluster permet des configurations spécifiques pour chaque cluster au sein d'un seul fichier.

#### Cluster de Production

- **EnvPrefix** : `kube01`
- **apiServerURL** : `https://kube01-api-url:443`
- **audience** : `xxxxxxxxxxx-xxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com`
- **providerURL** : `https://accounts.google.com`
- **clientID** : `xxxxxxxxxxx-xxxxxxxxxxxxxxxxxxx.apps.googleusercontent.com`
- **clientSecret** : `GXXXX-XXXXXXXXXXXXXXXXXXX`
- **clusterName** : `kube01`
- **emailClaim** : `email`
- **redirectURL** : `https://gangway.local/callback`
- **scopes** : `["openid", "profile", "email"]`
- **tokenURL** : `https://www.googleapis.com/oauth2/v4/token`
- **usernameClaim** : `email`

#### Clusters de Développement

- **Cluster 1 (kube02)**
  - **EnvPrefix** : `kube02`
  - **apiServerURL** : `https://kube02-api-url:443`
  - ... (identique au cluster de Production)

- **Cluster 2 (kube03)**
  - **EnvPrefix** : `kube03`
  - **apiServerURL** : `https://kube03-api-url:443`
  - **clusterCAPath** : `/etc/gangly/pki/kube03/ca.crt`
  - ... (identique au cluster de Production)

### Configuration Spécifique au Cluster

Chaque cluster peut avoir les configurations suivantes :

| Clé                          | Description                                                                                                    |
|------------------------------|----------------------------------------------------------------------------------------------------------------|
| `clusterName`                | Le nom du cluster. Utilisé dans l'UI et les instructions de configuration de kubectl.                          |
| `providerURL`                | URL du fournisseur OAuth2. Doit offrir un point de terminaison `$providerURL/.well-known/openid-configuration` pour la découverte.|
| `clientID`                   | ID client API tel que fourni par le fournisseur d'identité.                                                    |
| `clientSecret`               | Secret client API tel que fourni par le fournisseur d'identité.                                                |
| `allowEmptyClientSecret`     | Certains fournisseurs d'identité acceptent un secret client vide, ce qui n'est généralement pas une bonne idée. Si vous devez utiliser un secret vide et accepter les risques associés, alors vous pouvez le définir sur `true`. Par défaut à `false`.|
| `audience`                   | Le point de terminaison qui fournit des informations de profil utilisateur [optionnel]. Non requis par tous les fournisseurs.|
| `scopes`                     | Utilisé pour spécifier la portée de la demande d'autorisation OAuth. Par défaut à `["openid", "profile", "email", "offline_access"]`.|
| `redirectURL`                | Où rediriger après l'authentification. Cela devrait être une URL où Gangway est accessible. Typiquement, cela doit également être enregistré dans l'application OAuth avec le fournisseur OAuth.|
| `usernameClaim`              | La revendication JWT à utiliser comme nom d'utilisateur. Cela est utilisé dans l'UI. Combiné avec le clusterName pour la partie "utilisateur" de kubeconfig. Par défaut à `nickname`.|
| `emailClaim`                 | La revendication JWT à utiliser comme email. Par défaut à `email`.                                            |
| `apiServerURL`               | Le point de terminaison du serveur API utilisé pour configurer kubectl.                                        |
| `clusterCAPath`              | Chemin pour trouver le bundle CA pour le serveur API. Utilisé pour configurer kubectl. Ce chemin est généralement monté à l'emplacement par défaut pour les charges de travail fonctionnant sur un cluster Kubernetes et n'a généralement pas besoin d'être défini. Par défaut à `/var/run/secrets/kubernetes.io/serviceaccount/ca.crt`.|
| `showClaims`                 | Afficher les revendications reçues. Par défaut à `true`.                                                       |


## Fonctions Supplémentaires

- `NewMultiClusterConfig` : Crée une nouvelle instance de configuration multi-cluster à partir d'un fichier de configuration sérialisé.
- `Validate` : Vérifie toutes les propriétés de la structure de configuration pour s'assurer qu'elles sont initialisées.
- `GetRootPathPrefix` : Retourne '/' si aucun préfixe n'est spécifié, sinon retourne le chemin configuré.
- `loadCerts` : Charge les certificats pour les configurations de cluster à partir des chemins spécifiés.

## Utilisation des Variables d'Environnement

Les variables d'environnement peuvent être utilisées pour remplacer les configurations spécifiées dans le fichier YAML en utilis
