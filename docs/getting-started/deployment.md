# Deploying Gangway

Deploying Gangway consists of writing a config file and then deploying the service.
The service is stateless so it is relatively easy to manage on Kubernetes.
How you provide access to the service is going to be dependent on your specific configuration.

Gangway is now aware of a multi-cluster configuration, we have developped a helm chart that will made deployment easily.

Here is a configration example : 

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
