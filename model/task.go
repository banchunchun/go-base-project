package model

import "com.banxiaoxiao.server/app"

type Task struct {
	Model
	TaskId            int64              `gorm:"index" json:"taskId"`
	Name              string             `json:"name"`
	FileName          string             `gorm:"index" json:"fileName"`
	LiveMode          int32              `json:"liveMode"`
	FileType          string             `json:"fileType"`
	Features          string             `json:"features"`
	AppId             int64              `json:"appId"`
	Priority          int32              `gorm:"index" json:"priority"`
	MediaType         string             `json:"mediaType"`
	MediaTag          string             `json:"mediaTag"`
	FingerServerName  string             `json:"fingerServerName"`
	CallbackUrl       string             `json:"callbackUrl"`
	CallbackContent   string             `gorm:"size:65535" json:"callbackContent"`
	Published         bool               `gorm:"index" json:"published"`
	PublishStatus     app.TaskStatusType `gorm:"index" json:"publishStatus"`
	StartOffset       int64              `json:"startOffset"`
	EndOffset         int64              `json:"endOffset"`
	TaskStatus        app.TaskStatusType `gorm:"index" json:"taskStatus"`
	CreatedTime       int64              `json:"createTime"`
	StartTime         int64              `json:"startTime"`
	EndTime           int64              `json:"endTime"`
	OcrSkipFrames     int64              `json:"ocrSkipFrames"`
	FaceSkipFrames    int64              `json:"faceSkipFrames"`
	LibraryIdList     string             `json:"libraryIdList"`
	RetryTimes        int32              `json:"retryTimes"`
	ExtraData         string             `json:"extra_data"`
	TencentSearchType int32              `json:"tencentSearchType"`
}
