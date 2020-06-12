# APRSWeb

APRSWeb connects to Direwolf (or any AGWPE-compatible server, in
theory) and rendered received APRS frames onto a map in your web
browser.  ![aprsweb screenshot](aprsweb.png?raw=true "APRSWeb")

Unlike other APRS software, APRSWeb's target is to render into a web
browser. You can think of APRSWeb as a self-hosted version of
<aprs.fi>, except it can run on hardware as modest as a Raspberry Pi.

It uses Twirp to serialize those frames into JSON or protobuf objects,
so it also doubles as an HTTP API for your APRS frames.

## Building
You'll need `git`, `golang` `protoc-gen-twirp`, and `protoc` per the
[Twirp guide](https://github.com/twitchtv/twirp#installation).

You will also need `go-bindata`:
```
go get -u github.com/go-bindata/go-bindata/... 
```

Once you have everything installed, run `make`.

This will output `aprsweb` as an ARMv6 Linux binary, and
`aprsweb-native` as whatever architecture/OS you're compiling on now.

## Running 
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

## Development
When hacking locally, run `bindata` with debug mode enabled so you can
hot-reload files.

```
go-bindata -pkg bindata -debug -o bindata/bindata.go -fs -prefix "static/" assets/...

```

Since we have `bindata/bindata.go` in `.gitignore` and `make` always
overwrites this, you should end up with an all-in-one binary that
includes all of the js/html assets.

## TODO
* Change the default server from `kd2aoy-pi`
* Provide an easier way to customize:
    * Map starting position
    * Map zoom
    * HTML template
* Provide a map provider picker for off-net usage
* Provide a filter for:
    * Displaying a specific station's checkins only
* Automatically refresh the map every minute or so
* Remember map position in local storage so you don't get stuck in SF
* Use specific icons for different station checkins
* Write better docs :)
