package models

import "time"

// VideoLesson rappresenta una videolezione programmata
type VideoLesson struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ScheduledAt time.Time `json:"scheduled_at"`
	TeacherID   int       `json:"teacher_id"`
}
