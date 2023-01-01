Welcome to JTunnel .

Source code is at `https://github.com/manoj-inukolunu/jtunnel-go`

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

The endpoint you get is `https://myhost.lc-algorithms.com`

All the requests to `https://myhost.lc-algorithms.com` will now

be routed to your server running on port `3131`

