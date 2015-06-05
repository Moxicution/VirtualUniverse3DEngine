// Copyright © 2013-2015 Galvanized Logic Inc.
// Use is governed by a BSD-style license found in the LICENSE file.

package audio

import (
	"fmt"
	"log"
	"strings"

	"github.com/gazed/vu/audio/al"
)

// 64-bit OpenAL may be difficult to locate for Windows machines. Try
//    http://kcat.strangesoft.net/openal.html/openal-soft-1.15.1-bin.zip.
// Extract Win64/soft_oal.dll from the zip to c:/Windows/System32/OpenAL32.dll

// openal provides sound support for the engine. It exposes the useful parts
// of the underlying OpenAL audio library as well as providing some sound
// utility methods.
type openal struct {
	dev al.Device  // created on initialization.
	ctx al.Context // created on initialization.
}

// Init runs the one time openal library initialization. It is expected to
// be called once by the engine on startup.
func (a *openal) Init() (err error) {
	al.Init()
	if err = a.validate(); err != nil {
		return
	}

	// create an openal context for all sounds.
	if a.dev = al.OpenDevice(""); a.dev != 0 {
		if a.ctx = al.CreateContext(a.dev, nil); a.ctx != 0 {
			al.MakeContextCurrent(a.ctx)
			return // success
		}
	}
	return fmt.Errorf("openal audio init failed.")
}

// validate that OpenAL is available. OSX has OpenAL.
func (a *openal) validate() error {
	if report := al.BindingReport(); len(report) > 0 {
		for _, line := range report {
			if strings.Contains(line, "[-]") {
				return fmt.Errorf("OpenAL uninitialized.")
			}
		}
	} else {
		return fmt.Errorf("OpenAL unavailable.")
	}
	return nil
}

// Shutdown closes down the openal library.  This is only expected to be called
// once by the engine when it is shutting down.
func (a *openal) Shutdown() {
	al.MakeContextCurrent(0)
	if a.ctx != 0 {
		al.DestroyContext(a.ctx)
	}
	if a.dev != 0 {
		al.CloseDevice(a.dev)
	}
}

// SetGain sets the listener gain to a value between 0 and 1.
// Values outside the 0 to 1 range are ignored.
func (a *openal) SetGain(zeroToOne float64) {
	if zeroToOne >= 0 && zeroToOne <= 1 {
		al.Listenerf(al.GAIN, float32(zeroToOne))
	}
}

// BindSound copies sound data to the sound card. If successfull then the
// sound reference, snd, and sound data buffer reference, buff are updated
// with valid references.
func (a *openal) BindSound(snd, buff *uint32, d *Data) (err error) {
	if alerr := al.GetError(); alerr != al.NO_ERROR {
		log.Printf("openal.BindSound need to find and fix prior error %X", alerr)
	}

	// create the sound buffer and copy the audio data into the buffer
	var format int32
	if format, err = a.format(d); err == nil {
		al.GenBuffers(1, buff)
		al.BufferData(*buff, format, al.Pointer(&(d.AudioData[0])), int32(d.DataSize), int32(d.Frequency))
		if alerr := al.GetError(); alerr != al.NO_ERROR {
			err = fmt.Errorf("Failed binding sound %s", d.Name)
		} else {
			al.GenSources(1, snd)
			al.Sourcei(*snd, al.BUFFER, int32(*buff))
		}
	}
	return err
}

// Implement Audio.
func (a *openal) PlaceListener(x, y, z float64) {
	al.Listener3f(al.POSITION, float32(x), float32(y), float32(z))
}

// Implement Audio.
func (a *openal) PlaySound(snd uint32, x, y, z float64) {
	al.Source3f(snd, al.POSITION, float32(x), float32(y), float32(z))
	al.SourcePlay(snd)
}

// Implement Audio.
func (a *openal) ReleaseSound(sound uint32) { al.DeleteSources(1, &sound) }

// format figures out which of the OpenAL formats to use based on the
// WAVE file information. A -1 value, and error, is returned if the format
// cannot be determined.
func (a *openal) format(d *Data) (format int32, err error) {
	format = -1
	if d.Channels == 1 && d.SampleBits == 8 {
		format = al.FORMAT_MONO8
	} else if d.Channels == 1 && d.SampleBits == 16 {
		format = al.FORMAT_MONO16
	} else if d.Channels == 2 && d.SampleBits == 8 {
		format = al.FORMAT_STEREO8
	} else if d.Channels == 2 && d.SampleBits == 16 {
		format = al.FORMAT_STEREO16
	}
	if format < 0 {
		err = fmt.Errorf("openal:format cannot recognize audio format")
	}
	return format, err
}