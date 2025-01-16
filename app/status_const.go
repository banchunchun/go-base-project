package app

type TaskStatusType int32

const (
	TaskStatusPending   TaskStatusType = 0
	TaskStatusRunning   TaskStatusType = 1
	TaskStatusCompleted TaskStatusType = 2
	TaskStatusCancelled TaskStatusType = 3
	TaskStatusError     TaskStatusType = 4
	ThreeTaskPending    TaskStatusType = 10
)

func CheckCompleted(status TaskStatusType) bool {
	switch status {
	case TaskStatusCompleted, TaskStatusCancelled, TaskStatusError:
		return true
	default:
		return false
	}
}

func ErrorStatus(status TaskStatusType) bool {
	switch status {
	case TaskStatusCancelled, TaskStatusError:
		return true
	default:
		return false
	}
}
