package main

type PopupRequest struct {
	Type           string
	popup_URL      string
	popup_ClientID string
	popup_AppName  string
	popup_Favicon  string
	popup_Title    string

	input_Type        string
	input_Placeholder string
	input_MaxLength   int
	trackingID        string
}

func CreateTypingRequest(input_Type string, input_Placeholder string, input_MaxLength int) PopupRequest {
	req := PopupRequest{input_Type: input_Type, input_Placeholder: input_Placeholder, input_MaxLength: input_MaxLength, Type: "input"}
	return req
}

func CreatePopupRequest(popup_URL string, popup_ClientID string, popup_AppName string, popup_Favicon string, popup_Title string) PopupRequest {
	req := PopupRequest{popup_URL: popup_URL, popup_ClientID: popup_ClientID, popup_AppName: popup_AppName, popup_Favicon: popup_Favicon, popup_Title: popup_Title, Type: "popup"}
	return req
}