# Forwardly-Go

Super simple POST to GET forwarder using FastAPI. This was born from the need
of connecting Alertmanager and Uptime Kuma. This app gets the webhook
notification in POST and converts it to GET and fowards it to Uptime Kuma.
That's it.

This one is a Go rewrite from
[Python/FastAPI](https://github.com/hadret/forwardly) using
[Gin](https://gin-gonic.com).

## Testing

To test whether the `POST` payload works with the sample `alert.json` file run
this:

```shell
curl -XPOST -H "Content-Type: application/json" -d@alert.json localhost:8000/AAaaAaAaaa
```

This assumes that `.env-sample` file is unchanged and symlinked/moved to `.env`.
