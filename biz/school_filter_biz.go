package biz

import (
	"admin_project/global"
	"admin_project/sysRequest"
	"fmt"
)

type SchoolFilterBiz struct{}

var schoolFilter = SchoolFilterBiz{}

func RefSchoolFilterBiz() *SchoolFilterBiz {
	return &schoolFilter
}
func (biz *SchoolFilterBiz) FilterSchool(req *sysRequest.SchoolFilterReq) (*sysRequest.SchoolFilterRsp, error) {
	res := make([]*global.OfferInfo, 0)

	upperGpa := 0.0
	upperPercentage := 0.0
	tx := global.G_DB.Table("offer_info_tab")
	// 分数范围筛选
	if req.Grade.GpaScore != 0 {
		upperGpa = req.Grade.GpaScore + 0.1
		tx = tx.Where("gpa_grade<=?", upperGpa)
	}
	if req.Grade.PercentageScore != 0 {
		upperPercentage = req.Grade.PercentageScore + 5
		tx = tx.Where("gpa_percentage<=?", upperPercentage)
	}
	// 地区筛选
	tx = tx.Where("region IN (?)", req.DestinationRegion)
	// 专业筛选
	tx = tx.Where("candidate_major = ?", req.Subject)
	err := tx.Find(&res).Error
	if err != nil {
		global.GLog.Error(fmt.Sprintf("select school failed, err :%v", err.Error()))
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	accept := 0
	reject := 0
	subjectList := make([]string, 0)
	subjectMap := make(map[string]interface{})
	avgGrade := 0
	gradeNum := 0
	avgPercentage := 0
	percentageNum := 0
	for _, data := range res {
		if data.OfferStatus == "Offer" || data.OfferStatus == "AD小奖" || data.OfferStatus == "AD无奖" {
			accept += 1
		} else {
			reject += 1
		}
		if data.GpaGrade != 0 {
			avgGrade += data.GpaGrade
			gradeNum += 1
		}
		if data.GpaPercentage != 0 {
			avgPercentage += data.GpaPercentage
			percentageNum += 1
		}

		if _, ok := subjectMap[data.Major]; !ok {
			subjectMap[data.Major] = true
			subjectList = append(subjectList, data.Major)
		}
	}
	avg := sysRequest.Grade{}
	if percentageNum != 0 {
		avg.PercentageScore = float64(avgPercentage) / float64(percentageNum)
	}
	if gradeNum != 0 {
		avg.GpaScore = float64(avgGrade) / float64(gradeNum)
	}
	ret := &sysRequest.SchoolFilterRsp{
		SchoolName:  req.SchoolName,
		AcceptedNum: accept,
		DeclinedNum: reject,
		SubjectList: subjectList,
		AvgGrade:    avg,
	}
	return ret, nil
}
