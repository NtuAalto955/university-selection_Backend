package biz

import (
	"admin_project/global"
	"admin_project/sysRequest"
	"fmt"
	"github.com/shopspring/decimal"
	"regexp"
	"sort"
	"strings"
)

type SchoolFilterBiz struct{}

var schoolFilter = SchoolFilterBiz{}

func RefSchoolFilterBiz() *SchoolFilterBiz {
	return &schoolFilter
}
func (biz *SchoolFilterBiz) FilterSchool(req *sysRequest.SchoolFilterReq) (*sysRequest.SchoolFilterRsp, error) {
	res := make([]*global.OfferInfo, 0)

	tx := global.G_DB.Table("offer_info_tab")

	// 地区筛选
	tx = tx.Where("region IN (?)", req.DestinationRegion)
	// 专业筛选
	tx = tx.Where("major_type = ?", req.Major)
	err := tx.Find(&res).Error
	if err != nil {
		global.GLog.Error(fmt.Sprintf("select school failed, err :%v", err.Error()))
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	ret := biz.dataAggregation(res, req.SchoolLevel)
	return ret, nil
}
func (biz *SchoolFilterBiz) dataAggregation(offerList []*global.OfferInfo, schoolLevel string) *sysRequest.SchoolFilterRsp {
	// region - country -result
	result := make(map[string]map[string][]*sysRequest.Results)

	// schoolName - year - offerInfo
	offerInfoMap := make(map[string]map[int][]*global.OfferInfo)
	// schoolName - year - 平均分
	offerList = biz.SchoolLevelFilter(offerList, schoolLevel)
	// 总数据
	for _, data := range offerList {
		// 针对每一个学校
		if _, ok := offerInfoMap[data.SchoolName]; !ok {

			if _, ok := result[data.Region]; !ok {
				result[data.Region] = make(map[string][]*sysRequest.Results)
			}
			if _, ok := result[data.Region][data.SchoolCountry]; !ok {
				result[data.Region][data.SchoolCountry] = make([]*sysRequest.Results, 0)
			}

			offerInfoMap[data.SchoolName] = make(map[int][]*global.OfferInfo)
		}

		offerInfoMap[data.SchoolName][data.ApplyYear] = append(offerInfoMap[data.SchoolName][data.ApplyYear], data)
	}
	// 处理重复学校 比如boston university 和boston university 波士顿大学
	schoolNameMap := make(map[string]*sysRequest.Results)
	// 分年份处理数据
	for schoolName, applyResult := range offerInfoMap {
		var schoolResult *sysRequest.Results
		key := checkIsSchoolProcessed(schoolName, schoolNameMap)
		if key == "" {
			schoolResult = &sysRequest.Results{
				SchoolName:      schoolName,
				AdmissionYear:   make(map[int]*sysRequest.AdmissionDetail),
				GpaRange:        make([]sysRequest.AdmissionResult, 8),
				PercentageRange: make([]sysRequest.AdmissionResult, 9),
				SchoolRange:     make([]sysRequest.AdmissionResult, 4),
				TotalResult:     sysRequest.AdmissionResult{},
				AvgGrade:        sysRequest.Grade{},
			}
			schoolNameMap[schoolName] = schoolResult
		} else {
			schoolResult = schoolNameMap[key]
		}

		for year, offerListPerYear := range applyResult {
			gpaRange := biz.StatGpa(offerListPerYear)
			PercentageRange := biz.StatPercentage(offerListPerYear)
			SchoolRange := biz.StatSchoolLevel(offerListPerYear)
			if _, ok := schoolResult.AdmissionYear[year]; !ok {
				schoolResult.AdmissionYear[year] = &sysRequest.AdmissionDetail{
					GpaRange:        gpaRange,
					PercentageRange: PercentageRange,
					SchoolRange:     SchoolRange,
					TotalResult:     sysRequest.AdmissionResult{},
					AvgGrade:        sysRequest.Grade{},
				}
			} else {
				schoolResult.AdmissionYear[year].GpaRange = mergeRange(schoolResult.AdmissionYear[year].GpaRange, gpaRange)
				schoolResult.AdmissionYear[year].PercentageRange = mergeRange(schoolResult.AdmissionYear[year].PercentageRange, PercentageRange)
				schoolResult.AdmissionYear[year].SchoolRange = mergeRange(schoolResult.AdmissionYear[year].SchoolRange, SchoolRange)

			}

			schoolResult.GpaRange = mergeRange(schoolResult.GpaRange, gpaRange)
			schoolResult.PercentageRange = mergeRange(schoolResult.PercentageRange, PercentageRange)
			schoolResult.SchoolRange = mergeRange(schoolResult.SchoolRange, SchoolRange)

			for _, data := range offerListPerYear {
				if schoolResult.Region == "" {
					schoolResult.Region = data.Region
					schoolResult.Country = data.SchoolCountry
				}

				if IsOfferAdmitted(data) {
					schoolResult.AdmissionYear[year].TotalResult.AcceptedNum += 1
					schoolResult.TotalResult.AcceptedNum += 1
				} else {
					schoolResult.AdmissionYear[year].TotalResult.RejectedNum += 1
					schoolResult.TotalResult.RejectedNum += 1
				}

				if data.GpaGrade != 0 {
					schoolResult.AdmissionYear[year].AvgGrade.GpaScore += data.GpaGrade
					schoolResult.AdmissionYear[year].AvgGrade.GpaNum += 1
					schoolResult.AvgGrade.GpaScore += data.GpaGrade
					schoolResult.AvgGrade.GpaNum += 1
				}
				if data.GpaPercentage != 0 {
					schoolResult.AdmissionYear[year].AvgGrade.PercentageScore += data.GpaPercentage
					schoolResult.AdmissionYear[year].AvgGrade.PercentageNum += 1
					schoolResult.AvgGrade.PercentageScore += data.GpaPercentage
					schoolResult.AvgGrade.PercentageNum += 1
				}
			}

		}
		if key == "" {
			result[schoolResult.Region][schoolResult.Country] = append(result[schoolResult.Region][schoolResult.Country], schoolResult)

		}
	}
	// 按学校名称排序
	result, regionCountryList, countrySchoolList := calculateAvg(result)
	for _, countryResult := range result {
		for _, res := range countryResult {
			sort.Slice(res, func(i, j int) bool {
				isChinese := regexp.MustCompile("^[\u4e00-\u9fa5]+") //中文开头默认放到后面
				if isChinese.MatchString(res[i].SchoolName) {
					return false
				}
				return res[i].SchoolName < res[j].SchoolName
			})
		}
	}
	for _, res := range regionCountryList {
		sort.Slice(res, func(i, j int) bool {
			return res[i] < res[j]
		})
	}
	for _, res := range countrySchoolList {
		sort.Slice(res, func(i, j int) bool {
			isChinese := regexp.MustCompile("^[\u4e00-\u9fa5]+") //中文开头默认放到后面
			if isChinese.MatchString(res[i]) {
				return false
			}
			return res[i][0] < res[j][0]
		})
	}
	ret := &sysRequest.SchoolFilterRsp{
		ApplyResults:      result,
		CountrySchoolList: countrySchoolList,
		RegionCountryList: regionCountryList,
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
			if IsOfferAdmitted(row) {
				gpa0To2P6.AcceptedNum += 1
			} else {
				gpa0To2P6.RejectedNum += 1
			}
		}
		if row.GpaGrade < 2.8 && row.GpaGrade >= 2.6 {
			if IsOfferAdmitted(row) {
				gpa2P6To2P8.AcceptedNum += 1
			} else {
				gpa2P6To2P8.RejectedNum += 1
			}
		}
		if row.GpaGrade < 3.0 && row.GpaGrade >= 2.8 {
			if IsOfferAdmitted(row) {
				gpa2P8To3P0.AcceptedNum += 1
			} else {
				gpa2P8To3P0.RejectedNum += 1
			}
		}
		if row.GpaGrade < 3.2 && row.GpaGrade >= 3.0 {
			if IsOfferAdmitted(row) {
				gpa3P0To3P2.AcceptedNum += 1
			} else {
				gpa3P0To3P2.RejectedNum += 1
			}
		}
		if row.GpaGrade < 3.4 && row.GpaGrade >= 3.2 {
			if IsOfferAdmitted(row) {
				gpa3P2To3P4.AcceptedNum += 1
			} else {
				gpa3P2To3P4.RejectedNum += 1
			}
		}
		if row.GpaGrade < 3.6 && row.GpaGrade >= 3.4 {
			if IsOfferAdmitted(row) {
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
			if IsOfferAdmitted(row) {
				gpa3P8To4P0.AcceptedNum += 1
			} else {
				gpa3P8To4P0.RejectedNum += 1
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
	percentage90To100 := sysRequest.AdmissionResult{}
	for _, row := range data {
		if row.GpaPercentage < 76 {
			if IsOfferAdmitted(row) {
				percentage0To76.AcceptedNum += 1
			} else {
				percentage0To76.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 78 && row.GpaPercentage >= 76 {
			if IsOfferAdmitted(row) {
				percentage76To78.AcceptedNum += 1
			} else {
				percentage76To78.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 80 && row.GpaPercentage >= 78 {
			if IsOfferAdmitted(row) {
				percentage78To80.AcceptedNum += 1
			} else {
				percentage78To80.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 82 && row.GpaPercentage >= 80 {
			if IsOfferAdmitted(row) {
				percentage80To82.AcceptedNum += 1
			} else {
				percentage80To82.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 84 && row.GpaPercentage >= 82 {
			if IsOfferAdmitted(row) {
				percentage82To84.AcceptedNum += 1
			} else {
				percentage82To84.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 86 && row.GpaPercentage >= 84 {
			if IsOfferAdmitted(row) {
				percentage84To86.AcceptedNum += 1
			} else {
				percentage84To86.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 88 && row.GpaPercentage >= 86 {
			if IsOfferAdmitted(row) {
				percentage86To88.AcceptedNum += 1
			} else {
				percentage86To88.RejectedNum += 1
			}
		}
		if row.GpaPercentage < 90 && row.GpaPercentage >= 88 {
			if IsOfferAdmitted(row) {
				percentage88To90.AcceptedNum += 1
			} else {
				percentage88To90.RejectedNum += 1
			}
		}
		if row.GpaPercentage <= 100 && row.GpaPercentage >= 90 {
			if IsOfferAdmitted(row) {
				percentage90To100.AcceptedNum += 1
			} else {
				percentage90To100.RejectedNum += 1
			}
		}

	}
	percentageResult = append(percentageResult, percentage0To76, percentage76To78, percentage78To80, percentage80To82, percentage82To84,
		percentage84To86, percentage86To88, percentage88To90, percentage90To100)
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
			if IsOfferAdmitted(row) {
				schoolFirstLevel.AcceptedNum += 1
			} else {
				schoolFirstLevel.RejectedNum += 1
			}
		}
		if row.SchoolLevel == "双非" {
			if IsOfferAdmitted(row) {
				schoolSecondLevel.AcceptedNum += 1
			} else {
				schoolSecondLevel.RejectedNum += 1
			}
		}
		if row.SchoolLevel == "海本" {
			if IsOfferAdmitted(row) {
				schoolThirdLevel.AcceptedNum += 1
			} else {
				schoolThirdLevel.RejectedNum += 1
			}
		}
		if row.SchoolLevel == "其他" {
			if IsOfferAdmitted(row) {
				schoolOther.AcceptedNum += 1
			} else {
				schoolOther.RejectedNum += 1
			}
		}
	}
	schoolLevelResult = append(schoolLevelResult, schoolFirstLevel, schoolSecondLevel, schoolThirdLevel, schoolOther)
	return schoolLevelResult
}
func (biz *SchoolFilterBiz) SchoolLevelFilter(offerList []*global.OfferInfo, schoolLevel string) []*global.OfferInfo {
	res := make([]*global.OfferInfo, 0)

	schoolMap := make(map[string]map[string]bool)
	// 过滤掉这个学校level从没申请记录的学校
	for _, offer := range offerList {
		if _, ok := schoolMap[offer.SchoolName]; !ok {
			schoolMap[offer.SchoolName] = make(map[string]bool)
		}
		schoolMap[offer.SchoolName][offer.SchoolLevel] = true
	}
	for _, offer := range offerList {
		if schoolMap[offer.SchoolName][schoolLevel] {
			res = append(res, offer)
		}
	}
	return res

}
func IsOfferAdmitted(data *global.OfferInfo) bool {
	if data.OfferStatus == 1 || data.OfferStatus == 2 || data.OfferStatus == 3 {
		return true
	} else {
		return false
	}
}
func mergeRange(target, source []sysRequest.AdmissionResult) []sysRequest.AdmissionResult {
	for index, result := range source {
		target[index].AcceptedNum += result.AcceptedNum
		target[index].RejectedNum += result.RejectedNum
	}
	return target

}
func checkIsSchoolProcessed(schoolName string, schoolNameMap map[string]*sysRequest.Results) string {
	for key := range schoolNameMap {
		if strings.Contains(schoolName, key) || strings.Contains(key, schoolName) {
			return key
		}
	}
	return ""
}
func calculateAvg(data map[string]map[string][]*sysRequest.Results) (map[string]map[string][]*sysRequest.Results, map[string][]string, map[string][]string) {
	regionCountryList := make(map[string][]string)
	countrySchoolList := make(map[string][]string)
	for region, regionValue := range data {
		for country, countryValue := range regionValue {
			if country == "nan" {
				country = "Others"
			}
			regionCountryList[region] = append(regionCountryList[region], country)
			for _, value := range countryValue {

				countrySchoolList[country] = append(countrySchoolList[country], value.SchoolName)
				// 计算每年平均值
				for _, perYearValue := range value.AdmissionYear {
					if perYearValue.AvgGrade.GpaNum != 0 {
						perYearValue.AvgGrade.GpaScore = perYearValue.AvgGrade.GpaScore / float64(perYearValue.AvgGrade.GpaNum)
						perYearValue.AvgGrade.GpaScore, _ = decimal.NewFromFloat(perYearValue.AvgGrade.GpaScore).Round(2).Float64()

					}
					if perYearValue.AvgGrade.PercentageNum != 0 {
						perYearValue.AvgGrade.PercentageScore = perYearValue.AvgGrade.PercentageScore / float64(perYearValue.AvgGrade.PercentageNum)
						perYearValue.AvgGrade.PercentageScore, _ = decimal.NewFromFloat(perYearValue.AvgGrade.PercentageScore).Round(2).Float64()
					}
				}
				//计算总平均
				if value.AvgGrade.GpaNum != 0 {
					value.AvgGrade.GpaScore = value.AvgGrade.GpaScore / float64(value.AvgGrade.GpaNum)
					value.AvgGrade.GpaScore, _ = decimal.NewFromFloat(value.AvgGrade.GpaScore).Round(2).Float64()

				}
				if value.AvgGrade.PercentageNum != 0 {
					value.AvgGrade.PercentageScore = value.AvgGrade.PercentageScore / float64(value.AvgGrade.PercentageNum)
					value.AvgGrade.PercentageScore, _ = decimal.NewFromFloat(value.AvgGrade.PercentageScore).Round(2).Float64()
				}
			}
		}
	}
	return data, regionCountryList, countrySchoolList
}
