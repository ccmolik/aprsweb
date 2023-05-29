# APRSWeb

APRSWeb connects to Direwolf (or any AGWPE-compatible server, in
theory) and rendered received APRS frames onto a map in your web
browser.  ![aprsweb screenshot](aprsweb.png?raw=true "APRSWeb")

Unlike other APRS software, APRSWeb's target is to render into a web
browser. You can think of APRSWeb as a self-hosted version of
<https://aprs.fi>, except it can run on hardware as modest as a
Raspberry Pi.

It uses [Twirp](https://github.com/twitchtv/twirp) to serialize those
frames into JSON or protobuf objects, so it also doubles as an HTTP
API for your APRS frames.

## Building
You'll need `git`, `golang` `protoc-gen-twirp`, and `protoc` per the
[Twirp guide](https://github.com/twitchtv/twirp#installation).

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

I pipe the TNC test CD to Direwolf to play back synthetic packets, such as:
```
$ cat ~/direwolf.dummy
ADEVICE - null
CHANNEL 0

$ ffmpeg -i ~/TNC_Test_Ver-1.1.00.wav -f s16le -acodec pcm_s16le -ac 1 pipe: | direwolf -t 0 -r44100 -c ~/direwolf.dummy -

```
## TODO
* Provide an easier way to customize:
    * Map starting position
    * Map zoom
    * HTML template
* Provide a map provider picker for off-net usage
* Provide a filter for:
    * Displaying a specific station's checkins only
* Automatically refresh the map every minute or so
* ~Use specific icons for different station checkins~ (In progress)
* Provide a way to swap out map API keys
* Integrate twirp / pb generation into Makefile better
* APRS test cd has some issues with Mic-E data, leading to screwed up lat/lng (specifically AC6VV-9)
* Write better docs :)
