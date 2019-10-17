package listen

import (
	"errors"
	"unsafe"

	"github.com/joomcode/errorx"
	"github.com/xlab/closer"
	"github.com/xlab/pocketsphinx-go/sphinx"
	"github.com/xlab/portaudio-go/portaudio"
)

type PocketSphinx struct {
	decoder  *sphinx.Decoder
	stream   *portaudio.Stream
	channels int

	channel    chan string
	inSpeech   bool
	uttStarted bool
}

func (ps *PocketSphinx) Close() error {
	// terminate portaudio
	err := convertError(portaudio.Terminate())
	if err != nil {
		return errorx.Decorate(err, "cannot terminate PortAudio")
	}

	// destroy sphinx decoder
	if ps.decoder != nil {
		ps.decoder.Destroy()
	}

	// close portaudio stream
	if ps.stream != nil {
		err = convertError(portaudio.CloseStream(ps.stream))
		if err != nil {
			return errorx.Decorate(err, "cannot close PortAudio stream")
		}
		err = convertError(portaudio.StopStream(ps.stream))
		if err != nil {
			return errorx.Decorate(err, "cannot stop PortAudio stream")
		}
	}

	return nil
}

func (ps *PocketSphinx) Listen() string {
	return <-ps.channel
}

// paCallback: for simplicity reasons we process raw audio with sphinx in the this stream callback,
// never do that for any serious applications, use a buffered channel instead.
func (ps *PocketSphinx) paCallback(input unsafe.Pointer, _ unsafe.Pointer, sampleCount uint,
	_ *portaudio.StreamCallbackTimeInfo, _ portaudio.StreamCallbackFlags, _ unsafe.Pointer) int32 {

	const (
		statusContinue = int32(portaudio.PaContinue)
		statusAbort    = int32(portaudio.PaAbort)
	)

	in := (*(*[1 << 24]int16)(input))[:int(sampleCount)*ps.channels]
	// ProcessRaw with disabled search because callback needs to be relatime
	_, ok := ps.decoder.ProcessRaw(in, true, false)
	if !ok {
		return statusAbort
	}

	if ps.decoder.IsInSpeech() {
		ps.inSpeech = true
		if !ps.uttStarted {
			ps.uttStarted = true
		}
		return statusContinue
	}

	if ps.uttStarted {
		// speech -> silence transition, time to start new utterance
		ps.decoder.EndUtt()
		ps.uttStarted = false
		hyp, _ := ps.decoder.Hypothesis()
		if len(hyp) > 0 {
			ps.channel <- hyp
		}
		if !ps.decoder.StartUtt() {
			closer.Fatalln("Sphinx failed to start utterance")
		}
	}
	return statusContinue
}

func NewPocketSphinx(config ListenConfig) (*PocketSphinx, error) {
	err := convertError(portaudio.Initialize())
	if err != nil {
		return nil, errorx.Decorate(err, "cannot init PortAudio")
	}

	ps := PocketSphinx{channels: config.Channels}
	cfg := sphinx.NewConfig(
		sphinx.HMMDirOption(config.HMM),
		sphinx.DictFileOption(config.Dict),
		sphinx.LMFileOption(config.LM),
		sphinx.SampleRateOption(float32(config.SampleRate)),
	)
	sphinx.LogFileOption("/dev/null")(cfg)
	ps.decoder, err = sphinx.NewDecoder(cfg)
	if err != nil {
		return &ps, err
	}

	var stream *portaudio.Stream
	err = convertError(portaudio.OpenDefaultStream(
		&stream, int32(config.Channels), 0, portaudio.PaInt16, float64(config.SampleRate),
		uint(config.Samples), ps.paCallback, nil,
	))
	if err != nil {
		return &ps, errorx.Decorate(err, "cannot open PortAudio stream")
	}

	err = convertError(portaudio.StartStream(stream))
	if err != nil {
		return &ps, errorx.Decorate(err, "cannot start PortAudio stream")
	}

	if !ps.decoder.StartUtt() {
		return &ps, errorx.Decorate(err, "Sphinx failed to start utterance")
	}

	ps.channel = make(chan string)
	return &ps, nil
}

func convertError(err portaudio.Error) error {
	if portaudio.ErrorCode(err) == portaudio.PaNoError {
		return nil
	}
	return errors.New(portaudio.GetErrorText(err))
}
