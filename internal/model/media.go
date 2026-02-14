package model

// MediaUploadRequest holds the parameters needed to initialize a media upload.
type MediaUploadRequest struct {
	Owner string `json:"owner"`
	Type  string `json:"type"`
}

// MediaUpload contains the upload details returned by LinkedIn after
// initializing an upload.
type MediaUpload struct {
	UploadURL   string `json:"uploadUrl"`
	MediaURN    string `json:"mediaUrn"`
	UploadToken string `json:"uploadToken"`
}

// MediaStatus represents the processing status of an uploaded media asset.
type MediaStatus struct {
	URN    string `json:"urn"`
	Status string `json:"status"`
}
