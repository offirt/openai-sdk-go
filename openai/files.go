package openai

import (
	"fmt"
)

const FilesEndpointPath = "/files/"

// Files Endpoint
//
// Files are used to upload documents that can be used with features like [Fine-tuning]: https://platform.openai.com/docs/api-reference/fine-tunes.
type FilesEndpoint struct {
	*endpoint
}

// Files Endpoint
func (c *Client) Files() *FilesEndpoint {
	return &FilesEndpoint{newEndpoint(c, FilesEndpointPath)}
}

type File struct {
	Id        string   `json:"id"`
	Object    string   `json:"object"`
	Bytes     int      `json:"bytes"`
	CreatedAt int      `json:"created_at"`
	Filename  string   `json:"filename"`
	Purpose   FineTune `json:"fine-tune"`
}

type Files struct {
	Object string `json:"object"`
	Data   []File `json:"data"`
}

// Returns a list of files that belong to the user's organization.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/files
func (e *FilesEndpoint) ListFiles() ([]File, error) {
	var files Files
	err := e.do(e, "GET", "", nil, &files)
	if err == nil && files.Object != "list" {
		err = fmt.Errorf("expected 'list' object type, got %s", files.Object)
	}
	return files.Data, err
}

type UploadFileRequest struct {
	// The audio file to transcribe, in one of these formats: mp3, mp4, mpeg, mpga, m4a, wav, or webm.
	File string `json:"file" binding:"required"`
	// 	ID of the model to use. Only whisper-1 is currently available.
	Purpose string `json:"purpose" binding:"required"`
}

// Upload a file that contains document(s) to be used across various endpoints/features.
// Currently, the size of all the files uploaded by one organization can be up to 1 GB.
// Please contact us if you need to increase the storage limit.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/files
func (e *FilesEndpoint) UploadFile(req *UploadFileRequest) (File, error) {
	var file File
	err := e.do(e, "POST", "", req, &file)
	return file, err
}

type DeleteFileResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// Delete a file.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/files
func (e *FilesEndpoint) DeleteFile(fileId string) (bool, error) {
	var resp DeleteFileResponse
	err := e.do(e, "POST", fileId, nil, &resp)
	if err != nil {
		return false, err
	}
	return resp.Deleted, nil
}

// Returns information about a specific file.
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/files
func (e *FilesEndpoint) RetrieveFile(fileId string) (*File, error) {
	var file File
	err := e.do(e, "GET", fileId, nil, &file)
	return &file, err
}

// Returns the contents of the specified file
// [OpenAI Documentation]: https://platform.openai.com/docs/api-reference/files
func (e *FilesEndpoint) RetrieveFileContent(fileId string) (*string, error) {
	var data string
	err := e.do(e, "GET", fileId+"/content", nil, &data)
	return &data, err
}
