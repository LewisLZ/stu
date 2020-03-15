package form

type ListAchievement struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Name      string `form:"name"`
}

type ListAchievementSource struct {
	ExaminationClassId int `form:"examination_class_id"`
}

type ListAchievementSearchSource struct {
	Page            int    `form:"page"`
	Limit           int    `form:"limit"`
	StudentName     string `form:"student_name"`
	CurriculumName  string `form:"curriculum_name"`
	ExaminationName string `form:"examination_name"`
	ClassName       string `form:"class_name"`
	StartTime       string `form:"start_time"`
	EndTime         string `form:"end_time"`
	Sort            string `form:"sort"`
	Order           string `form:"order"`
}

type SaveAchievementScore struct {
	ExaminationClassId int `json:"examination_class_id"`
	ClassCurriculumId  int `json:"class_curriculum_id"`
	StudentId          int `json:"student_id"`
	Score              int `json:"score"`
}

type ListAchievementStudentScore struct {
	ExaminationClassId int `form:"examination_class_id"`
	StudentId          int `form:"student_id"`
}

type ListAchievementCurriculumScore struct {
	ExaminationClassId int `form:"examination_class_id"`
	ClassCurriculumId  int `form:"class_curriculum_id"`
}

type ListAchievementClassScore struct {
	ExaminationClassId int `form:"examination_class_id"`
}
