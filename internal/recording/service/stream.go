package service

type Stream struct {
	ID            int `json:"id"`
	VideoRecorder *VideoRecorder
}

func InitNewStream(id int, vr *VideoRecorder) *Stream {
	return &Stream{
		ID:            id,
		VideoRecorder: vr,
	}
}
