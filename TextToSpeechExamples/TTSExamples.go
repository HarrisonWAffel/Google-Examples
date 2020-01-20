package main

//  If there is any code duplication below it is intentional.
// The purpose of the code is to clearly demonstrate how to work with GCP TTS, not to be production ready.

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"context"
	"fmt"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"io/ioutil"
	"strings"
)

/// Entry Point ///
func main() {
	p := SpeechRequest{
		Text:         "<speak> Hello World, <emphasis level=\"strong\"> let's get to work </emphasis> </speak>",
		LanguageCode: "en-US",
		SsmlGender:   "FEMALE",
		VoiceName:    "en-US-Wavenet-C",
	}

	s := p.SpeakToStream()

	err := ioutil.WriteFile("test.mp3", s, 0644)
	checkErr(err)
}

/// Structs ///

type SpeechRequest struct {
	Text         string
	LanguageCode string
	SsmlGender   string
	VoiceName    string
}

type SpeechExampleError struct {
	Message string
}

/// Core func's ///

//SpeakToFile uses the values within the receiver to open a connection to GCP, create a request, and then take the response and put it into a file.
func (st *SpeechRequest) SpeakToFile(outputFile string) {

	//Create a go context, a key component of nearly all golang web requests
	ctx := context.Background()

	//create a client connection to the GCP TTS backend
	client, err := texttospeech.NewClient(ctx)
	checkErr(err)

	// we 'defer' the closing of the client connection until the function has exited
	defer client.Close()

	//Craft a request using the parameters specified
	req, serr := st.CraftTextSpeechRequest()
	checkSpeechErr(serr)

	//Receive a response from GCP TTS
	resp, err := client.SynthesizeSpeech(ctx, &req)
	checkErr(err)

	//Write the contents of the response body to a file

	err = ioutil.WriteFile(outputFile, resp.AudioContent, 0644)
	checkErr(err)

	fmt.Printf("TTS Successfully written to %s", outputFile)
}

func (pt *SpeechRequest) SpeakFromFileToFile(inputFile string, outputFile string) {

	//Read the contents of the input file and pass them to the struct
	content, err := ioutil.ReadFile(inputFile)
	checkErr(err)
	pt.Text = string(content)

	//Proceed to translate the contents
	pt.SpeakToFile(outputFile)
}

func (st *SpeechRequest) SpeakToStream() []byte {

	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	checkErr(err)
	defer client.Close()

	req, er := st.CraftTextSpeechRequest()
	checkSpeechErr(er)

	resp, err := client.SynthesizeSpeech(ctx, &req)
	checkErr(err)

	return resp.AudioContent
}

/// Supporting funcs ///

func (st *SpeechRequest) CraftTextSpeechRequest() (texttospeechpb.SynthesizeSpeechRequest, SpeechExampleError) {

	//make sure st has the required values

	if st.Text == "" {
		return texttospeechpb.SynthesizeSpeechRequest{}, SpeechExampleError{Message: "TTS Request Has Empty Text"}
	}

	if st.LanguageCode == "" {
		return texttospeechpb.SynthesizeSpeechRequest{}, SpeechExampleError{Message: "TTS Request Has Empty Language Code"}
	}

	if st.SsmlGender == "" {
		return texttospeechpb.SynthesizeSpeechRequest{}, SpeechExampleError{Message: "TTS Request Has Empty Ssml Gender"}
	}

	if st.VoiceName == "" {
		return texttospeechpb.SynthesizeSpeechRequest{}, SpeechExampleError{Message: "TTS Request Has Empty Voice Name"}
	}

	// convert input strings to the proper types
	gender := texttospeechpb.SsmlVoiceGender_FEMALE
	if strings.Contains(st.SsmlGender, "MALE") {
		gender = texttospeechpb.SsmlVoiceGender_MALE
	}
	if st.SsmlGender == "" {
		gender = texttospeechpb.SsmlVoiceGender_NEUTRAL
	}

	input := &texttospeechpb.SynthesisInput{InputSource: &texttospeechpb.SynthesisInput_Text{Text: st.Text}}
	if strings.Contains(st.Text, "<speak>") {
		input = &texttospeechpb.SynthesisInput{InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: st.Text}}
	}

	//create and return the request
	return texttospeechpb.SynthesizeSpeechRequest{

		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},

		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: st.LanguageCode,
			Name:         st.VoiceName,
			SsmlGender:   gender,
		},

		Input: input,
	}, SpeechExampleError{}

}

//short hand for quick error checks
func checkErr(e error) {
	if e != nil {
		fmt.Print(e)
	}
}

//short hand for speechExampleError checks
func checkSpeechErr(exampleError SpeechExampleError) {
	if exampleError.Message != "" {
		fmt.Println(exampleError.Message)
	}
}
