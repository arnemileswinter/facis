package review_task_state

type ReviewTaskState string

const (
	Open     ReviewTaskState = "OPEN"
	Approved ReviewTaskState = "APPROVED"
)

var validState = map[ReviewTaskState]bool{
	Open:     true,
	Approved: true,
}

func NewReviewTaskState(state ReviewTaskState) ReviewTaskState {
	return ReviewTaskState(state)
}
