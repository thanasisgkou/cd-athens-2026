//go:build js && wasm

package game

import (
	"syscall/js"
)

// SpeakLine triggers browser Web Speech API synthesis with custom language, pitch, and rate.
func SpeakLine(text string, lang string, pitch float64, rate float64) {
	defer func() {
		_ = recover()
	}()
	win := js.Global()
	if win.IsUndefined() || win.IsNull() {
		return
	}
	synth := win.Get("speechSynthesis")
	uttCtor := win.Get("SpeechSynthesisUtterance")
	if synth.IsUndefined() || synth.IsNull() || uttCtor.IsUndefined() || uttCtor.IsNull() {
		return
	}
	synth.Call("cancel")
	utt := uttCtor.New(text)
	utt.Set("lang", lang)
	utt.Set("pitch", pitch)
	utt.Set("rate", rate)
	utt.Set("volume", 0.95)
	synth.Call("speak", utt)
}
