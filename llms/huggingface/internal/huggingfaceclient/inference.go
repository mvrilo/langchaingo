package huggingfaceclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// InferenceTask is the type of inference task to run.
type InferenceTask string

const (
	InferenceTaskTextGeneration      InferenceTask = "text-generation"
	InferenceTaskText2TextGeneration InferenceTask = "text2text-generation"
)

type inferencePayload struct {
	Model  string `json:"-"`
	Inputs string `json:"inputs"`
}

type (
	inferenceResponsePayload []inferenceResponse
	inferenceResponse        struct {
		Text string `json:"generated_text"`
	}
)

const hfInferenceAPI = "https://api-inference.huggingface.co/models/"

func (c *Client) runInference(ctx context.Context, payload *inferencePayload) (inferenceResponsePayload, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", hfInferenceAPI, payload.Model), body) //nolint:lll
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	// debug print the http request with httputil:

	// reqDump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Fprintf(os.Stderr, "%s", reqDump)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	// debug print the http response with httputil:
	// resDump, err := httputil.DumpResponse(r, true)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Fprintf(os.Stderr, "%s", resDump)

	var response inferenceResponsePayload
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
