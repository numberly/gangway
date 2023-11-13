# How It Works

Kubernetes supports OpenID Connect (OIDC) as a user authentication mechanism. OIDC is an
authentication protocol that allows servers to verify the identity of a user by way of an ID Token.

When using OIDC to authenticate with Kubernetes, the client (e.g. `kubectl`) sends the ID token
alongside all requests to the API server. On the server side, the Kubernetes API server verifies the
token to ensure it is valid and has not expired. Once verified, the API server extracts username and
group membership information from the token, and continues processing the request.

In order to obtain the ID token, the user must go through the OIDC authentication process. This is
where Gangway comes in. Gangway is a web application that enables the OIDC authentication flow which
results in the minting of the ID Token.

Gangway is configured as a client of an upstream Identity Service that speaks OIDC. To obtain the ID
token, the user accesses Gangway, initiates the OIDC flow by clicking the "Log In" button, and
completes the flow by authenticating with the upstream Identity Service. The user's credentials are
never shared with Gangway.

Once the authentication flow is complete, the user is redirected to a Gangway page that provides
instructions on how to configure `kubectl` to use the ID token.

The following sequence diagram details the authentication flow:

![gangway](images/gangway-sequence-diagram.png)