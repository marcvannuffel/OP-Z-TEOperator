# teoperator (OP-Z Edition)

*teoperator* lets you easily make drum and synth patches for the teenage engineering op-1 or op-z. 

This is an extended fork of the original [teoperator](https://github.com/schollz/teoperator) by infinitedigits, specifically optimized for the **teenage engineering OP-Z**. It adds powerful batch processing features, automatic kit splitting, and a unique "compress" mode to bypass the OP-Z's strict sample length limits.

## New OP-Z Features

### 1. Multi-File Upload & Batch Processing
Upload as many audio files as you want at once. The tool automatically processes them based on your selected patch type:

*   **Drum Kits:** The OP-Z limits drum kits to a maximum of 24 splices and 12 seconds total duration. When you upload multiple files (e.g., 60 drum samples), the tool automatically groups them and generates as many `.aif` drum kits as needed, perfectly respecting both limits. You get multiple ready-to-use kits to drop onto your OP-Z.
*   **Synth Patches:** Each uploaded file is converted into its own individual `.aif` synth patch. All patches are then bundled into a single `.zip` file for easy downloading.

### 2. Compress Mode (Double Playback Length)
The OP-Z has strict sample length limits: 12 seconds for a drum kit and 5.75 seconds for a synth patch. The **Compress Mode** allows you to effectively double these limits for any sample longer than 1 second.

When enabled, the tool halves the sample rate of the audio (from 44.1 kHz to 22.05 kHz). This makes the sample play twice as fast and one octave higher, fitting twice as much audio into the same time limit. 

**How to use it on the OP-Z:**
1. Load the compressed patch onto your OP-Z.
2. Set the pitch parameter of the track to **-1 octave**.
3. The original pitch is restored, and you effectively get **double the playback length** out of each slot (up to 24 seconds for drum kits, and ~11.5 seconds for synth patches). This is incredibly useful for longer loops, pads, or textures.

### 3. Mobile-First UI Redesign
The entire user interface has been redesigned from the ground up with a mobile-first approach. It features a clean, dark-mode aesthetic, large touch targets, clear section labels, and a prominent multi-upload area. The results page now includes a modern waveform display and an integrated audio player for previewing your patches before downloading.

## Installation

*teoperator* requires ffmpeg. First, [install ffmpeg](https://ffmpeg.org/download.html). 

To use as a command line program you first need to [install Go](https://golang.org/doc/install) and then in a terminal run:

```
go get -v github.com/schollz/teoperator@latest
```

That will install `teoperator` on your system.

## Usage (Command Line)

You can use *teoperator* to create drum patches or sample-based synth patches for the op-1 or op-z. The resulting file is a `.aif` converted to mono 44.1khz with metadata representing key-assignment information for the op-1 or op-z. You can use any kind of input music file (wav, aif, mp3, flac, etc.).

### Make synth sample patches

To make a synth patch just type:

```
teoperator synth piano.wav
```

Optionally, you can include the base frequency information which can be used on the op-1/opz to convert to the right pitch:

```
teoperator synth --freq 220 piano.wav
```

### Make a drum kit patch

To make a drumkit patch you can convert multiple files and splice points will be set at the boundaries of each individual file:

```
teoperator drum kick.wav snare.wav openhat.wav closedhat.wav
```

### Make a drum sample patch

To make a sample patch you can convert one sample and splice points will be automatically determined by transients:

```
teoperator drum vocals.wav
```

### Make a drum loop patch

To make a drum loop patch you can convert one sample and define splice points to be equally spaced along the sample:

```
teoperator drum --slices 16 drumloop.wav
```

## Web server

The webserver requires a few more dependencies. You can install them via `apt`:

```
$ sudo apt install imagemagick 
$ sudo add-apt-repository ppa:chris-needham/ppa
$ sudo apt-get update
$ sudo apt-get install audiowaveform
$ sudo -H python3 -m pip install youtube-dl
```

And then you can run the server via

```
$ git clone https://github.com/schollz/teoperator
$ cd teoperator && go build -v
$ ./teoperator server
```

Then open a browser to `localhost:8053`!

# License

MIT license

Please note THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
