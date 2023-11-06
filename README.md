# Welcome to Numberly/Gangway a new fork with multi-cluster improvements

Gangway is EOL but still used at Numberly, we decided to fork and maintain it.

See the original repository : https://github.com/vmware-archive/gangway

This fork aims to continue development of Gangway by Numberly corporation.

gangway
=======

_(noun): An opening in the bulwark of the ship to allow passengers to board or leave the ship._

An application that can be used to easily enable authentication flows via OIDC for a kubernetes cluster.
Kubernetes supports [OpenID Connect Tokens](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) as a way to identify users who access the cluster.
Gangway has been improved and is now able to handle multiple clusters
Gangway allows users to self-configure their `kubectl` configuration in a few short steps.

![gangway multicluster](docs/images/gangway-multicluster.png)

Once authenticated for one of your cluster : 

![gangway](docs/images/gangway.png)

## Deployment

Instructions for deploying gangway for common cloud providers can be found [here](docs/README.md).

We can use our dedicated helm chart that will allow you to configure Gangway easily [here](chart/README.md)

## The multi-cluster way

At Numberly, initially, we was deploying 1 gangway per cluster which could lead to a lot of WebUI to connect
our clusters. One of the goal of the fork was also to permit to Gangway to manage any clusters to facilitate
our management. Gangway can now take a configuration with environment and a list of cluster for each environments.
You can read more about it [here](docs/configuration.md)

## How It Works

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

<p align="center">
    <img src="docs/images/gangway-sequence-diagram.png" width="600px" />
</p>

## API-Server flags

gangway requires that the Kubernetes API server is configured for OIDC:

https://kubernetes.io/docs/admin/authentication/#configuring-the-api-server

```bash
kube-apiserver
...
--oidc-issuer-url="https://example.auth0.com/"
--oidc-client-id=3YM4ue8MoXgBkvCIHh00000000000
--oidc-username-claim=email
--oidc-groups-claim=groups
```

## Build

Requirements for building

- Go (built with version >= 1.21)

A Makefile is provided for building tasks. The options are as follows

Getting started is as simple as:

```bash
go get -u github.com/soulkyu/gangway
cd $GOPATH/src/github.com/soulkyu/gangway
make setup
make
```
