package form

type SaveExamination struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	StartTime        string `json:"start_time"`
	Remark           string `json:"remark"`
	ClassCurriculums []*struct {
		ClassId     int `json:"class_id"`
		Curriculums []*struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"curriculums"`
	} `json:"class_curriculum"`
}

type ListExamination struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Name      string `form:"name"`
}

type SaveExaminationClass struct {
	ExaminationId    int    `json:"id"`
	Name             string `json:"name"`
	StartTime        string `json:"start_time"`
	Remark           string `json:"remark"`
	ClassCurriculums []*struct {
		ClassId     int `json:"class_id"`
		Curriculums []*struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"curriculums"`
	} `json:"class_curriculum"`
}
