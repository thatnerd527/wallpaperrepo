package main

type PopupResult struct {
	Type           string
	trackingID	   string
	popup_ResultData string
	input_ResultData string
	cancelled bool
}

func CreatePopupResult(popup_ResultData string, trackingID string, cancelled bool) PopupResult {
	return PopupResult{popup_ResultData: popup_ResultData, trackingID: trackingID, Type: "popup", cancelled: cancelled}
}

func CreateInputResult(input_ResultData string, trackingID string, cancelled bool) PopupResult {
	return PopupResult{input_ResultData: input_ResultData, trackingID: trackingID, Type: "input", cancelled: cancelled}
}