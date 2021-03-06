package model

// 成绩
type Score struct {
	Base
	ExaminationClassId int `gorm:"not null;unique_index:uidx_ecid_sid_ccid_score" json:"examination_class_id"` // 考试的班Id
	StudentId          int `gorm:"not null;unique_index:uidx_ecid_sid_ccid_score" json:"student_id"`           // 学生Id
	ClassCurriculumId  int `gorm:"not null;unique_index:uidx_ecid_sid_ccid_score" json:"class_curriculum_id"`  // 班级的课Id
	Score              int `gorm:"not null;size:5" json:"score"`                                               // 成绩
}
