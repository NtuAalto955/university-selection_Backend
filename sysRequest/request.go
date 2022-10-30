package sysRequest

type SchoolFilterReq struct {
	SchoolLevel       string   `json:"school_level"`
	Major             int      `json:"major"`
	Grade             Grade    `json:"grade"`
	DestinationRegion []string `json:"destination_region"`
}
type Grade struct {
	PercentageScore float64 `json:"percentage_score"`
	PercentageNum   int     `json:"percentage_num"`
	GpaScore        float64 `json:"gpa_score"`
	GpaNum          int     `json:"gpa_num"`
}
type SchoolFilterRsp struct {
	ApplyResults []ApplyResults `json:"apply_results"`
}
type ApplyResults struct {
	SchoolName      string                   `json:"school_name"`
	AdmissionYear   map[int]*AdmissionDetail `json:"admission_year"`
	Region          string                   `json:"region"`
	Country         string                   `json:"country"`
	GpaRange        []AdmissionResult        `json:"gpa_range"`
	PercentageRange []AdmissionResult        `json:"percentage_range"`
	SchoolRange     []AdmissionResult        `json:"school_range"`
	TotalResult     AdmissionResult          `json:"total_result"`
	AvgGrade        Grade                    `json:"avg_grade"`
}
type AdmissionDetail struct {
	ApplyYear       int               `json:"apply_year"`
	GpaRange        []AdmissionResult `json:"gpa_range"`
	PercentageRange []AdmissionResult `json:"percentage_range"`
	SchoolRange     []AdmissionResult `json:"school_range"`
	TotalResult     AdmissionResult   `json:"total_result"`
	AvgGrade        Grade             `json:"avg_grade"`
}
type AdmissionResult struct {
	AcceptedNum int `json:"accepted_num"`
	RejectedNum int `json:"rejected_num"`
}
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Id       string `json:"id"`
	B64s     string `json:"b64s"`
}

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type DeleteUser struct {
	Username string `json:"username"`
}
type UpdateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Comment struct {
	CommentID string `gorm:"unique;not null"`
	Name      string `gorm:"comment:评论者用户名"`
	Content   string `gorm:"comment:评论内容"`
}

type PageInfo struct {
	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}
