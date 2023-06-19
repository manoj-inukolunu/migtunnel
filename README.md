# MigTunnel

## Installation

```shell
brew tap manoj-inukolunu/homebrew-tools
brew install migtunnel
```

## Usage

### Starting migtunnel

&nbsp;

```shell
migtunnel start
```

&nbsp;

After starting open up a new terminal window or a new tab and then

&nbsp;

```shell
migtunnel register --adminPort 1234 --host subdomain --port 3030
```

&nbsp;

The endpoint you get is `https://subdomain.migtunnel.net`

All the requests to `https://subdomain.migtunnel.net` will now

be routed to your server running on port `3030`

&nbsp;

Local Server can be running tls as well.

```shell
migtunnel register tls --adminPort 1234 --host tlstest --port 8080 --server localhost
```



