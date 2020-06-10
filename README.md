# APRSWeb

APRSWeb connects to Direwolf (or any AGWPE-compatible server, in
theory) and serves received APRS frames.

It uses Twirp to serialize those frames into JSON or protobuf objects.
# Building
You'll need `git`, `golang` `protoc-gen-twirp`, and `protoc` per the
[Twirp guide](https://github.com/twitchtv/twirp#installation).

You will also need `go-bindata`:
```
go get -u github.com/go-bindata/go-bindata/... 
```

Once you have everything installed, run `make`.

This will output `aprsweb` as an ARMv6 Linux binary, and
`aprsweb-native` as whatever architecture/OS you're compiling on now.


# Development
When hacking locally, run `bindata` with debug mode enabled so you can
hot-reload files.

```
go-bindata -pkg bindata -debug -o bindata/bindata.go -fs -prefix "static/" assets/...

```

Since we have `bindata/bindata.go` in `.gitignore` and `make` always
overwrites this, you should end up with an all-in-one binary that
includes all of the js/html assets.

# Running 
Set up [Direwolf](https://github.com/wb2osz/direwolf) with a radio or
an rtl_fm device. This exercise is left to the reader, but certainly
can be run on a raspi co-locating this service.

Run `./aprsweb-native -p port -s server` substituting port and server
for the port and server of your AGWPE server.

If you set the environment variable `DEBUG=1` aprsweb will log every
packet with base64

Note that we use `git rev-parse HEAD --short` upon build so APRSWeb
versions are exposed in HTTP headers.

An example Kubernetes deployment is included at `aprsweb.yml` for
those with letencrypt and nginx-ingress running.

# TODO
* Change the default server from `kd2aoy-pi`
* Provide an easier way to customize:
    * Map starting position
    * Map zoom
    * HTML template
* Provide a map provider picker for off-net usage
* Provide a filter for:
    * Time period to query for and display on map
    * Displaying a specific station's checkins only
* Automatically refresh the map every minute
* Remember map position in local storage so you don't get stuck in SF
* Write better docs :)

# Bugs
Mic-e is still kind of sketchy location wise.
