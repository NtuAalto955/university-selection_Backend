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

	//upperGpa := 0.0
	//upperPercentage := 0.0
	tx := global.G_DB.Table("offer_info_tab")
	// 分数范围筛选
	//if req.Grade.GpaScore != 0 {
	//	upperGpa = req.Grade.GpaScore + 0.1
	//	tx = tx.Where("gpa_grade<=?", upperGpa)
	//}
	//if req.Grade.PercentageScore != 0 {
	//	upperPercentage = req.Grade.PercentageScore + 5
	//	tx = tx.Where("gpa_percentage<=?", upperPercentage)
	//}
	// 地区筛选
	tx = tx.Where("region IN (?)", req.DestinationRegion)
	// 专业筛选
	tx = tx.Where("major = ?", req.Subject)
	err := tx.Find(&res).Error
	if err != nil {
		global.GLog.Error(fmt.Sprintf("select school failed, err :%v", err.Error()))
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	ret := biz.dataAggregation(res)
	return ret, nil
}
func (biz *SchoolFilterBiz) dataAggregation(offerList []*global.OfferInfo) *sysRequest.SchoolFilterRsp {
	schoolNameMap := make(map[string]*sysRequest.ApplyResults)
	offerInfoMap := make(map[string][]*global.OfferInfo)

	for _, data := range offerList {
		if _, ok := schoolNameMap[data.SchoolName]; !ok {
			schoolNameMap[data.SchoolName] = &sysRequest.ApplyResults{
				TotalResult: sysRequest.AdmissionResult{},
				AvgGrade:    sysRequest.Grade{},
			}
			offerInfoMap[data.SchoolName] = make([]*global.OfferInfo, 0)
		}
		offerInfoMap[data.SchoolName] = append(offerInfoMap[data.SchoolName], data)
		if data.OfferStatus == 1 {
			schoolNameMap[data.SchoolName].TotalResult.AcceptedNum += 1
		} else {
			schoolNameMap[data.SchoolName].TotalResult.RejectedNum += 1
		}

		if data.GpaGrade != 0 {
			schoolNameMap[data.SchoolName].AvgGrade.GpaScore += data.GpaGrade
			schoolNameMap[data.SchoolName].AvgGrade.GpaNum += 1
		}
		if data.GpaPercentage != 0 {
			schoolNameMap[data.SchoolName].AvgGrade.PercentageScore += data.GpaPercentage
			schoolNameMap[data.SchoolName].AvgGrade.PercentageNum += 1
		}

	}
	res := make([]sysRequest.ApplyResults, 0)
	for schoolName, applyResult := range schoolNameMap {
		applyResult.SchoolName = schoolName
		applyResult.Region = offerInfoMap[schoolName][0].Region
		if applyResult.AvgGrade.GpaNum != 0 {
			applyResult.AvgGrade.GpaScore = applyResult.AvgGrade.GpaScore / float64(applyResult.AvgGrade.GpaNum)
		}
		if applyResult.AvgGrade.PercentageNum != 0 {
			applyResult.AvgGrade.PercentageScore = applyResult.AvgGrade.PercentageScore / float64(applyResult.AvgGrade.PercentageNum)
		}
		applyResult.GpaRange = biz.StatGpa(offerInfoMap[schoolName])
		applyResult.PercentageRange = biz.StatPercentage(offerInfoMap[schoolName])
		applyResult.SchoolRange = biz.StatSchoolLevel(offerInfoMap[schoolName])
		res = append(res, *applyResult)
	}

	ret := &sysRequest.SchoolFilterRsp{
		ApplyResults: res,
	}
	return ret
}
func (biz *SchoolFilterBiz) StatGpa(data []*global.OfferInfo) []sysRequest.AdmissionResult {
	gpaResult := make([]sysRequest.AdmissionResult, 0)

	gpa0To2P6 := sysRequest.AdmissionResult{}
	gpa2P6To2P8 := sysRequest.AdmissionResult{}
	gpa2P8To3P0 := sysRequest.AdmissionResult{}
	gpa3P0To3P2 := sysRequest.AdmissionResult{}
	gpa3P2To3P4 := sysRequest.AdmissionResult{}
	gpa3P4To3P6 := sysRequest.AdmissionResult{}
	gpa3P6To3P8 := sysRequest.AdmissionResult{}
	gpa3P8To4P0 := sysRequest.AdmissionResult{}
	for _, row := range data {
		if row.GpaGrade < 2.6 {
			if row.OfferStatus == 1 {
				gpa0To2P6.AcceptedNum += 1
			} else {
				gpa0To2P6.RejectedNum += 1
			}
		}
		if row.GpaGrade < 2.8 && row.GpaGrade >= 2.6 {
			if row.OfferStatus == 1 {
				gpa2P6To2P8.AcceptedNum += 1
			} else {
				gpa2P6To2P8.RejectedNum += 1
			}
		}
		if row.GpaGrade < 3.0 && row.GpaGrade >= 2.8 {
			if row.OfferStatus == 1 {
				gpa2P8To3P0.AcceptedNum += 1
			} else {
				gpa2P8To3P0.RejectedNum += 1
			}
		}
		if row.GpaGrade < 3.2 && row.GpaGrade >= 3.0 {
			if row.OfferStatus == 1 {
				gpa3P0To3P2.AcceptedNum += 1
			} else {
				gpa3P0To3P2.RejectedNum += 1
			}
		}
		if row.GpaGrade < 3.4 && row.GpaGrade >= 3.2 {
			if row.OfferStatus == 1 {
				gpa3P2To3P4.AcceptedNum += 1
			} else {
				gpa3P2To3P4.RejectedNum += 1
			}
		}
		if row.GpaGrade < 3.6 && row.GpaGrade >= 3.4 {
			if row.OfferStatus == 1 {
				gpa3P4To3P6.AcceptedNum += 1
			} else {
				gpa3P4To3P6.RejectedNum += 1
			}
		}
		if row.GpaGrade < 3.8 && row.GpaGrade >= 3.6 {
			if row.OfferStatus == 1 {
				gpa3P6To3P8.AcceptedNum += 1
			} else {
				gpa3P6To3P8.RejectedNum += 1
			}
		}
		if row.GpaGrade >= 3.8 {
			if row.OfferStatus == 1 {
				gpa3P6To3P8.AcceptedNum += 1
			} else {
				gpa3P6To3P8.RejectedNum += 1
			}
		}
	}
	gpaResult = append(gpaResult, gpa0To2P6, gpa2P6To2P8, gpa2P8To3P0, gpa3P0To3P2, gpa3P2To3P4, gpa3P4To3P6, gpa3P6To3P8, gpa3P8To4P0)
	return gpaResult
}
func (biz *SchoolFilterBiz) StatPercentage(data []*global.OfferInfo) []sysRequest.AdmissionResult {
	percentageResult := make([]sysRequest.AdmissionResult, 0)

	percentage0To76 := sysRequest.AdmissionResult{}
	percentage76To78 := sysRequest.AdmissionResult{}
	percentage78To80 := sysRequest.AdmissionResult{}
	percentage80To82 := sysRequest.AdmissionResult{}
	percentage82To84 := sysRequest.AdmissionResult{}
	percentage84To86 := sysRequest.AdmissionResult{}
	percentage86To88 := sysRequest.AdmissionResult{}
	percentage88To90 := sysRequest.AdmissionResult{}
	percentage90To92 := sysRequest.AdmissionResult{}
	percentage92To94 := sysRequest.AdmissionResult{}
	percentage94To96 := sysRequest.AdmissionResult{}
	percentage96To98 := sysRequest.AdmissionResult{}
	percentage98To100 := sysRequest.AdmissionResult{}
	for _, row := range data {
		if row.GpaPercentage < 76 {
			if row.OfferStatus == 1 {
				percentage0To76.AcceptedNum += 1
			} else {
				percentage0To76.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 78 && row.GpaPercentage >= 76 {
			if row.OfferStatus == 1 {
				percentage76To78.AcceptedNum += 1
			} else {
				percentage76To78.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 80 && row.GpaPercentage >= 78 {
			if row.OfferStatus == 1 {
				percentage78To80.AcceptedNum += 1
			} else {
				percentage78To80.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 82 && row.GpaPercentage >= 80 {
			if row.OfferStatus == 1 {
				percentage80To82.AcceptedNum += 1
			} else {
				percentage80To82.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 84 && row.GpaPercentage >= 82 {
			if row.OfferStatus == 1 {
				percentage82To84.AcceptedNum += 1
			} else {
				percentage82To84.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 86 && row.GpaPercentage >= 84 {
			if row.OfferStatus == 1 {
				percentage84To86.AcceptedNum += 1
			} else {
				percentage84To86.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 88 && row.GpaPercentage >= 86 {
			if row.OfferStatus == 1 {
				percentage86To88.AcceptedNum += 1
			} else {
				percentage86To88.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 90 && row.GpaPercentage >= 88 {
			if row.OfferStatus == 1 {
				percentage88To90.AcceptedNum += 1
			} else {
				percentage88To90.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 92 && row.GpaPercentage >= 90 {
			if row.OfferStatus == 1 {
				percentage90To92.AcceptedNum += 1
			} else {
				percentage90To92.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 94 && row.GpaPercentage >= 92 {
			if row.OfferStatus == 1 {
				percentage92To94.AcceptedNum += 1
			} else {
				percentage92To94.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 96 && row.GpaPercentage >= 94 {
			if row.OfferStatus == 1 {
				percentage94To96.AcceptedNum += 1
			} else {
				percentage94To96.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 98 && row.GpaPercentage >= 96 {
			if row.OfferStatus == 1 {
				percentage96To98.AcceptedNum += 1
			} else {
				percentage96To98.RejectedNum += 1
			}
		}
		if row.GpaPercentage >= 98 {
			if row.OfferStatus == 1 {
				percentage98To100.AcceptedNum += 1
			} else {
				percentage98To100.RejectedNum += 1
			}
		}

	}
	percentageResult = append(percentageResult, percentage0To76, percentage76To78, percentage78To80, percentage80To82, percentage82To84,
		percentage84To86, percentage86To88, percentage88To90, percentage88To90, percentage90To92, percentage92To94,
		percentage94To96, percentage96To98, percentage98To100)
	return percentageResult
}

func (biz *SchoolFilterBiz) StatSchoolLevel(data []*global.OfferInfo) []sysRequest.AdmissionResult {
	schoolLevelResult := make([]sysRequest.AdmissionResult, 0)
	schoolFirstLevel := sysRequest.AdmissionResult{}
	schoolSecondLevel := sysRequest.AdmissionResult{}
	schoolThirdLevel := sysRequest.AdmissionResult{}
	schoolOther := sysRequest.AdmissionResult{}
	for _, row := range data {
		if row.SchoolLevel == "985/211" {
			if row.OfferStatus == 1 {
				schoolFirstLevel.AcceptedNum += 1
			} else {
				schoolFirstLevel.RejectedNum += 1
			}
		}
		if row.SchoolLevel == "双非" {
			if row.OfferStatus == 1 {
				schoolSecondLevel.AcceptedNum += 1
			} else {
				schoolSecondLevel.RejectedNum += 1
			}
		}
		if row.SchoolLevel == "海本" {
			if row.OfferStatus == 1 {
				schoolThirdLevel.AcceptedNum += 1
			} else {
				schoolThirdLevel.RejectedNum += 1
			}
		}
		if row.SchoolLevel == "其他" {
			if row.OfferStatus == 1 {
				schoolOther.AcceptedNum += 1
			} else {
				schoolOther.RejectedNum += 1
			}
		}
	}
	schoolLevelResult = append(schoolLevelResult, schoolFirstLevel, schoolSecondLevel, schoolThirdLevel, schoolOther)
	return schoolLevelResult
}
