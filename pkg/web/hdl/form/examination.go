package form

type SaveExamination struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	Remark    string `json:"remark"`
}

type ListExamination struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Name      string `form:"name"`
}

type SaveExaminationClass struct {
	ExaminationId int `json:"examination_id"`
	ClassId       int `json:"class_id"`
}

type ListExaminationClass struct {
	ExaminationId int `form:"examination_id"`
}

type SaveExaminationClassCurriculum struct {
	ExaminationClassId int   `json:"examination_class_id"`
	ClassCurriculumIds []int `json:"class_curriculum_ids"`
}
