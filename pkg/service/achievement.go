package service

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/quexer/utee"

	"liuyu/stu/pkg/datasource"
	"liuyu/stu/pkg/model"
	"liuyu/stu/pkg/ut"
	"liuyu/stu/pkg/web/hdl/form"
)

type Achievement struct {
	Ds *datasource.Ds
}

func (p *Achievement) List(in *form.ListAchievement) ([]*model.Achievement, int, error) {
	db := p.Ds.Db.Model(&model.Examination{}).
		Select("id examination_id, name examination_name, start_time examination_time").
		Where("start_time<?", utee.Tick())
	if in.Name != "" {
		db = db.Where("name like ?", "%"+in.Name+"%")
	}
	if in.StartTime != "" && in.EndTime != "" {
		start, err := time.ParseInLocation("2006-01-02", in.StartTime, time.Local)
		if err != nil {
			return nil, 0, ut.NewValidateError("开始时间格式不对")
		}
		end, err := time.ParseInLocation("2006-01-02", in.EndTime, time.Local)
		if err != nil {
			return nil, 0, ut.NewValidateError("结束时间格式不对")
		}
		db = db.Where("start_time>=? AND start_time<=?", utee.Tick(start), utee.Tick(end))
	}

	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var as []*model.Achievement
	_, limit, offset := ut.MakePager(in.Page, in.Limit, 10)
	if err := db.Limit(limit).Offset(offset).Order("id desc").
		Scan(&as).Error; err != nil {
		return nil, 0, err
	}

	for _, a := range as {
		if err := p.Ds.Db.Model(&model.ExaminationClass{}).
			Where("examination_id=?", a.ExaminationId).
			Count(&a.ExaminationClassCount).Error; err != nil {
			return nil, 0, err
		}
	}

	return as, count, nil
}

func (p *Achievement) ListStudentScore(in *form.ListAchievementStudentScore) (interface{}, error) {
	sql := `SELECT s.score,
				   s.mt    score_mt,
				   c.name  class_curriculum_name,
				   t.name  class_teacher_name,
				   t2.name curriculum_teacher_name
			FROM score s
					 LEFT JOIN class_curriculum cc ON cc.id = s.class_curriculum_id
					 LEFT JOIN curriculum c ON c.id = cc.curriculum_id
					 LEFT JOIN class_curriculum_year ccy ON ccy.id = cc.cc_year_id
					 LEFT JOIN class cl ON ccy.class_id = cl.id
					 LEFT JOIN teacher_class tc ON tc.class_id = cl.id
					 LEFT JOIN teacher t ON t.id = tc.teacher_id
					 LEFT JOIN teacher_curriculum tcc ON tcc.curriculum_id = cc.curriculum_id
					 LEFT JOIN teacher t2 ON t2.id = tcc.teacher_id
			where s.student_id = ? and s.examination_class_id = ? ORDER BY s.mt DESC`

	var res []*struct {
		Score                 int    `json:"score"`
		ScoreMt               int64  `json:"score_mt"`
		ClassCurriculumName   string `json:"class_curriculum_name"`
		ClassTeacherName      string `json:"class_teacher_name"`
		CurriculumTeacherName string `json:"curriculum_teacher_name"`
	}
	if err := p.Ds.Db.Raw(sql, in.StudentId, in.ExaminationClassId).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Achievement) ListCurriculumScore(in *form.ListAchievementCurriculumScore) (interface{}, error) {
	sql := `SELECT s.score,
				   s.mt     score_mt,
				   st.name  student_name
			FROM score s
					 LEFT JOIN student st ON s.student_id = st.id
					 LEFT JOIN class_curriculum cc ON cc.id = s.class_curriculum_id
					 LEFT JOIN curriculum c ON c.id = cc.curriculum_id
			where s.class_curriculum_id = ? and s.examination_class_id = ? ORDER BY s.mt DESC`

	var res []*struct {
		Score       int    `json:"score"`
		ScoreMt     int64  `json:"score_mt"`
		StudentName string `json:"student_name"`
	}
	if err := p.Ds.Db.Raw(sql, in.ClassCurriculumId, in.ExaminationClassId).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Achievement) ListClassScore(in *form.ListAchievementClassScore) (interface{}, error) {
	sql := `SELECT s.score,
				   s.class_curriculum_id,
				   s.examination_class_id,
				   s.student_id,
				   s.mt         score_mt,
				   st.name      student_name,
				   c.name       class_curriculum_name,
				   ccy.year     class_curriculum_year,
				   ccy.pos      class_curriculum_pos
			FROM score s
					 LEFT JOIN student st ON s.student_id = st.id
					 LEFT JOIN class_curriculum cc ON cc.id = s.class_curriculum_id
					 LEFT JOIN curriculum c ON c.id = cc.curriculum_id
					 LEFT JOIN class_curriculum_year ccy ON ccy.id = cc.cc_year_id
			where s.examination_class_id = ?`

	var res []*struct {
		Score               int       `json:"score"`
		ClassCurriculumId   int       `json:"class_curriculum_id"`
		ExaminationClassId  int       `json:"examination_class_id"`
		StudentId           int       `json:"student_id"`
		ScoreMt             int64     `json:"score_mt"`
		StudentName         string    `json:"student_name"`
		ClassCurriculumName string    `json:"class_curriculum_name"`
		ClassCurriculumYear string    `json:"-"`
		ClassCurriculumPos  model.Pos `json:"-"`
	}
	if err := p.Ds.Db.Raw(sql, in.ExaminationClassId).Scan(&res).Error; err != nil {
		return nil, err
	}

	for _, v := range res {
		v.ClassCurriculumName = ut.MakeClassName(v.ClassCurriculumName, v.ClassCurriculumYear, v.ClassCurriculumPos)
	}

	return res, nil
}

func (p *Achievement) ListSearchScore(in *form.ListAchievementSearchSource) (interface{}, int, error) {
	var res []*struct {
		StudentName         string    `json:"student_name"`
		StudentCode         string    `json:"student_code"`
		ClassCurriculumName string    `json:"class_curriculum_name"`
		ClassCurriculumYear string    `json:"-"`
		ClassCurriculumPos  model.Pos `json:"-"`
		ExaminationName     string    `json:"examination_name"`
		ExaminationTime     int64     `json:"examination_time"`
		ClassName           string    `json:"class_name"`
		ClassYear           string    `json:"-"`
		ClassPos            model.Pos `json:"-"`
		Score               int       `json:"score"`
		ScoreMt             int64     `json:"score_mt"`
		ExaminationClassId  int       `json:"examination_class_id"`
		StudentId           int       `json:"student_id"`
		ClassCurriculumId   int       `json:"class_curriculum_id"`
	}

	db := p.Ds.Db.Table("score s").
		Select(`s.score,
					s.examination_class_id,
					s.student_id,
					s.class_curriculum_id,
					st.name      student_name,
					st.code      student_code,
       				c.name       class_curriculum_name,
       				ccy.year     class_curriculum_year,
       				ccy.pos      class_curriculum_pos,
       				e.name       examination_name,
       				e.start_time examination_time,
       				cl.name      class_name,
       				sy.year      class_year,
       				sy.pos       class_pos,
					s.mt		 score_mt`).
		Joins("LEFT JOIN student st ON s.student_id = st.id").
		Joins("LEFT JOIN class_curriculum cc ON cc.id = s.class_curriculum_id").
		Joins("LEFT JOIN curriculum c ON c.id = cc.curriculum_id").
		Joins("LEFT JOIN class_curriculum_year ccy ON ccy.id = cc.cc_year_id").
		Joins("LEFT JOIN class cl ON ccy.class_id = cl.id").
		Joins("LEFT JOIN school_year sy ON cl.school_year_id = sy.id").
		Joins("LEFT JOIN examination_class ec ON ec.id = s.examination_class_id").
		Joins("LEFT JOIN examination e ON e.id = ec.examination_id").
		Where("st.code!=''")

	if in.StartTime != "" && in.EndTime != "" {
		start, err := time.ParseInLocation("2006-01-02", in.StartTime, time.Local)
		if err != nil {
			return nil, 0, ut.NewValidateError("开始时间格式不对")
		}
		end, err := time.ParseInLocation("2006-01-02", in.EndTime, time.Local)
		if err != nil {
			return nil, 0, ut.NewValidateError("结束时间格式不对")
		}
		db = db.Where("e.start_time>=? AND e.start_time<=?", utee.Tick(start), utee.Tick(end))
	}
	if in.StudentName != "" {
		db = db.Where("st.name like ?", "%"+in.StudentName+"%")
	}
	if in.CurriculumName != "" {
		db = db.Where("c.name like ?", "%"+in.CurriculumName+"%")
	}
	if in.ExaminationName != "" {
		db = db.Where("e.name like ?", "%"+in.ExaminationName+"%")
	}
	if in.ClassName != "" {
		db = db.Where("cl.name like ?", "%"+in.ClassName+"%")
	}

	var count int
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	_, limit, offset := ut.MakePager(in.Page, in.Limit, 10)

	db = db.Limit(limit).Offset(offset)
	if in.Sort != "" && in.Order != "" {
		db = db.Order(in.Sort + " " + in.Order)
	} else {
		db = db.Order("s.mt desc")
	}
	if err := db.Scan(&res).Error; err != nil {
		return nil, 0, err
	}

	for _, v := range res {
		v.ClassName = ut.MakeClassName(v.ClassName, v.ClassYear, v.ClassPos)
		v.ClassCurriculumName = ut.MakeClassName(v.ClassCurriculumName, v.ClassCurriculumYear, v.ClassCurriculumPos)
	}

	return res, count, nil
}

func (p *Achievement) ListSource(in *form.ListAchievementSource) (interface{}, error) {

	var res []*struct {
		ExaminationClassId  int    `json:"examination_class_id"`
		ClassCurriculumId   int    `json:"class_curriculum_id"`
		ClassCurriculumName string `json:"class_curriculum_name"`
		StudentId           int    `json:"student_id"`
		StudentName         string `json:"student_name"`
		Score               int    `json:"score"`
	}

	sql := `select s.examination_class_id, s.class_curriculum_id, c.name class_curriculum_name, s.student_id, st.name student_name, s.score
			from score s
         	left join class_curriculum cc on cc.id = s.class_curriculum_id
         	left join curriculum c on cc.curriculum_id = c.id
         	left join student st on s.student_id = st.id
			where s.examination_class_id=?`
	if err := p.Ds.Db.Raw(sql, in.ExaminationClassId).Scan(&res).Error; err != nil {
		return nil, err
	}

	if len(res) != 0 {
		return res, nil
	}

	sql = `select ec.id examination_class_id, cc.id class_curriculum_id, c.name class_curriculum_name, s.id student_id, s.name student_name 
			from examination_class_curriculum ecc
			left join examination_class ec on ec.id=ecc.examination_class_id
			left join class_curriculum cc on cc.id=ecc.class_curriculum_id
			left join curriculum c on cc.curriculum_id = c.id
			left join student s on ec.class_id = s.class_id
			where ecc.examination_class_id=?`
	if err := p.Ds.Db.Raw(sql, in.ExaminationClassId).Scan(&res).Error; err != nil {
		return nil, err
	}

	err := datasource.RunTransaction(p.Ds.Db, func(tx *gorm.DB) error {
		tick := utee.Tick()
		for _, v := range res {

			if err := tx.Exec(`insert into score 
				(ct, mt, examination_class_id, student_id, class_curriculum_id, score) 
				values (?, ?, ?, ?, ?, ?)`, tick, 0, v.ExaminationClassId, v.StudentId, v.ClassCurriculumId, 0).
				Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Achievement) SaveScore(in []*form.SaveAchievementScore) error {
	return datasource.RunTransaction(p.Ds.Db, func(tx *gorm.DB) error {
		sql := `UPDATE score SET score=? WHERE examination_class_id=? AND student_id=? AND class_curriculum_id=?`
		for _, v := range in {
			if err := tx.Exec(sql, v.Score, v.ExaminationClassId, v.StudentId, v.ClassCurriculumId).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *Achievement) Archived(examinationClassId int) error {
	var res []*struct {
		ExaminationClassId  int    `json:"examination_class_id"`
		ClassCurriculumId   int    `json:"class_curriculum_id"`
		ClassCurriculumName string `json:"class_curriculum_name"`
		StudentId           int    `json:"student_id"`
		StudentName         string `json:"student_name"`
		Score               int    `json:"score"`
	}

	sql := `select s.examination_class_id, s.class_curriculum_id, c.name class_curriculum_name, s.student_id, st.name student_name, s.score
			from score s
         	left join class_curriculum cc on cc.id = s.class_curriculum_id
         	left join curriculum c on cc.curriculum_id = c.id
         	left join student st on s.student_id = st.id
			where s.examination_class_id=?`
	if err := p.Ds.Db.Raw(sql, examinationClassId).Scan(&res).Error; err != nil {
		return err
	}

	if len(res) == 0 {
		return ut.NewValidateError("没找到考试成绩单")
	}

	if err := p.Ds.Db.Model(&model.ExaminationClass{}).Where("id=?", examinationClassId).Updates(map[string]interface{}{
		"mt":     utee.Tick(),
		"status": model.ECStatus_Archived,
	}).Error; err != nil {
		return err
	}

	return nil
}
