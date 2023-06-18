# MigTunnel



Welcome to MigTunnel .

To create a new tunnel

Make a `POST` request to `http://127.0.0.1:1234/create`
with the payload

```
{
    "HostName":"myhost",
    "TunnelName":"Tunnel Name",
    "localServerPort":"3131"
}

```

The endpoint you get is `https://myhost.migtunnel.net`

All the requests to `https://myhost.migtunnel.net` will now

be routed to your server running on port `3131`

