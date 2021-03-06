/*
Copyright 2013 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"os"

	"code.google.com/p/portaudio-go/portaudio"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	e := NewEngine()

	oscMod := NewOsc()
	oscMod.Input("pitch", Value(-0.1))

	oscModAmp := NewAmp()
	oscModAmp.Input("car", oscMod)
	oscModAmp.Input("mod", Value(0.1))

	osc := NewOsc()
	osc.Input("pitch", oscModAmp)

	envMod := NewOsc()
	envMod.Input("pitch", Value(-1))

	envModAmp := NewAmp()
	envModAmp.Input("car", envMod)
	envModAmp.Input("mod", Value(0.02))

	envModSum := NewSum()
	envModSum.Input("car", envModAmp)
	envModSum.Input("mod", Value(0.021))

	env := NewEnv()
	env.Input("att", Value(0.0001))
	env.Input("dec", envModSum)

	amp := NewAmp()
	amp.Input("car", osc)
	amp.Input("mod", env)

	ampAmp := NewAmp()
	ampAmp.Input("car", amp)
	ampAmp.Input("mod", Value(0.5))

	e.Input("root", ampAmp)

	if err := e.Start(); err != nil {
		log.Println(err)
		return
	}

	os.Stdout.Write([]byte("Press enter to stop...\n"))
	os.Stdin.Read([]byte{0})

	if err := e.Stop(); err != nil {
		log.Println(err)
	}
}
