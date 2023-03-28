# Forwardly-Go

Super simple POST to GET forwarder using FastAPI. This was born from the need
of connecting Alertmanager and Uptime Kuma. This app gets the webhook
notification in POST and converts it to GET and fowards it to Uptime Kuma.
That's it.

This one is a rewrite from [Python/FastAPI](https://github.com/hadret/forwardly).
