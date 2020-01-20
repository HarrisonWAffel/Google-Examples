package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

//This example uses 'malgo', the mini-audio library for golang.

func main() {
	Recognize()
}

//This function calls the google speech to text api
//We use gstreamer to get audio input from the rpi
//we then pass that audio stream to the Stdin of another script
//which then formats the audio stream and passes it to the google cloud platform
//We continue to stream to the GCP until we receive a non-empty response body
//at which point we return the contents and kill the streaming process.
func Recognize() (string, float64, error) {

	//First we need to craft the command we want to execute
	cmdName := "SpeechToTextExamples/scripts/recognize"
	cmdArgs := []string{""}
	cmd := exec.Command(cmdName, cmdArgs...)
	//We need to create a reader for the stdout of this script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	//If we want to return the values that are returned from the above script
	//we need to declare the return values at a higher scope
	transcription := ""
	confidence := 0.0

	//A scanner is created to read the stdout of the above command
	scanner := bufio.NewScanner(cmdReader)

	//A new go thread is created to handle the audio streaming and subsequent response bodies
	go func() {
		for scanner.Scan() {

			fmt.Println("Response Recognized...")
			//A Third party can interrupt this streaming process by simply saying "stop"
			//useful when you want to stop the test, but don't want orphan processes
			if strings.Contains(scanner.Text(), "stop") {
				if err := cmd.Process.Kill(); err != nil {
					log.Fatal("failed to kill process: ", err)
				}
			}

			//Regular expressions are used to parse out the transcription and confidence score
			//from the return body of our API request
			transcriptRegx := regexp.MustCompile("(\"([^\"]|\"\")*\")")
			match := transcriptRegx.FindStringSubmatch(scanner.Text())

			confRegx := regexp.MustCompile("([+-]?[0-9]*\\.[0-9]*)")
			conf := confRegx.FindStringSubmatch(scanner.Text())

			//We need to ensure that an empty response body doesn't stop our recognition
			if len(match) != 0 {
				transcription = match[1]
				confidence, err = strconv.ParseFloat(conf[1], 64)
				if err != nil {
					fmt.Println("Could not convert confidence from string to float64")
					fmt.Println(conf[1])
					fmt.Println(match[1])
				}

				//Now that we have our transcription we can stop the recognition process
				if err := cmd.Process.Kill(); err != nil {
					log.Fatal("failed to kill process: ", err)
				}
			}
		}
	}()

	//We need to start our goroutine from the main thread
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting Cmd", err)
		os.Exit(1)
	}

	//We need to wait for a transcription before we can return said transcription
	err = cmd.Wait()
	var out bytes.Buffer
	if err != nil && transcription == "" {
		fmt.Println("Recognition crash")
		fmt.Println("tried to run the command : ./scripts/recognize")
		fmt.Println(fmt.Sprint(err) + ": " + out.String())
		return transcription, confidence, err
	}
	return transcription, confidence, nil
}

/// Supporting Funcs ///

func c(err error) {
	if err != nil {
		panic(err)
	}

}
