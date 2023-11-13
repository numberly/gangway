# Connexion de Gangway à Dex

[Dex](https://github.com/coreos/dex) est un outil pratique créé par CoreOS qui fournit un point de terminaison OIDC commun pour plusieurs fournisseurs d'identité.
Pour configurer Gangway afin de communiquer avec Dex, certaines informations doivent être collectées auprès de Dex.
Suivant la norme OIDC, Dex fournit une URL où sa configuration OIDC peut être rassemblée.
Cette URL se trouve à `.well-known/openid-configuration`.
Si Dex est configuré avec une URL d'émetteur de `http://app.example.com`, sa configuration OpenID peut être trouvée à `http://app.example.com/.well-known/openid-configuration`.
Un exemple de la configuration OpenID fournie par Dex :

 ```
 {
   "issuer": "http://app.example.com",
   "authorization_endpoint": "http://app.example.com/auth",
   "token_endpoint": "http://app.example.com/token",
   "jwks_uri": "http:/app.example.com/keys",
   "response_types_supported": [
     "code"
   ],
   "subject_types_supported": [
     "public"
   ],
   "id_token_signing_alg_values_supported": [
     "RS256"
   ],
   "scopes_supported": [
     "openid",
     "email",
     "groups",
     "profile",
     "offline_access"
   ],
   "token_endpoint_auth_methods_supported": [
     "client_secret_basic"
   ],
   "claims_supported": [
     "aud",
     "email",
     "email_verified",
     "exp",
     "iat",
     "iss",
     "locale",
     "name",
     "sub"
   ]
 }
 ```


En utilisant l'exemple de Gangway, il suffit de fournir votre installation Dex comme fournisseur. La configuration de Dex fournit une liste 
nommée `claims_supported` parmi laquelle vous pouvez choisir lors de la définition de `username_claim` et `email_claim`.
La revendication correcte à utiliser dépend du fournisseur d'identité en amont pour lequel Dex est configuré.
`client_id` et `client_secret` sont des chaînes de caractères qui peuvent avoir n'importe quelle valeur, mais elles doivent correspondre à l'ID client et au secret dans votre configuration Dex.
