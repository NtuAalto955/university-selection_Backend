package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var (
	GLog     *zap.Logger
	G_DB     *gorm.DB
	G_Viper  *viper.Viper
	G_Config Config
)

const (
	Priv_Admin   = 0
	Priv_User    = 1
	Priv_Visitor = 2
)

type OfferInfo struct {
	ID             uint   `gorm:"primarykey" json:"id"`
	OfferStatus    string `gorm:"not null;unique;comment:" json:"offer_status"`
	SchoolCountry  string `gorm:"not null;school_country" json:"school_country"`
	Region         string `gorm:"not null;comment:apply country region" json:"region"`
	Degree         string `gorm:"not null;comment:apply degree" json:"degree"`
	Major          string `gorm:"not null;comment:apply major" json:"username"`
	SchoolName     string `gorm:"not null;comment:apply schoolname" json:"major"`
	Semester       string `gorm:"not null;comment:apply semester" json:"semester"`
	IeltsScore     int    `gorm:"not null;comment:ielts score" json:"ielts_score"`
	ToeflScore     int    `gorm:"not null;comment:toefl score" json:"toefl_score"`
	GpaPercentage  int    `gorm:"not null;comment:gpa in percentage 85,95 etc" json:"gpa_percentage"`
	GpaGrade       int    `gorm:"not null;comment:gpa in grade 4.5,4.6 etc" json:"gpa_grade"`
	MaxGpaGrade    int    `gorm:"not null;comment:max gpa in grade 4,5 etc" json:"max_gpa_grade"`
	CandidateMajor string `gorm:"not null;comment:candidate's major" json:"candidate_major"`
	SchoolLevel    string `gorm:"not null;comment:school Level" json:"school_level"`
}
type User struct {
	ID        uint      `gorm:"primarykey"`
	Username  string    `gorm:"not null;unique;comment:用户账户" json:"username"`
	Password  string    `gorm:"comment:用户登录名" json:"password"`
	Priv      int       `gorm:"default:2" json:"privilege" `
	CreatedAt time.Time `gorm:"CreatedAt" json:"creattime"`
	UpdatedAt time.Time `gorm:"UpdatedAt" json:"updatetime"`
}
type Comment struct {
	CommentID uint      `gorm:"primarykey;unique;not null" json:"commentid"`
	Name      string    `gorm:"comment:评论者用户名"  json:"name"`
	Content   string    `gorm:"comment:评论内容"  json:"content"`
	CreatedAt time.Time `gorm:"CreatedAt" json:"creattime"`
	UpdatedAt time.Time `gorm:"UpdatedAt" json:"updatetime"`
}
type Config struct {
	Mysql Mysql
}
type Mysql struct {
	Path         string `yaml:"path"`
	Port         string `yaml:"port"`
	Config       string `yaml:"config"`
	Dbname       string `yaml:"dbname"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
	MaxOpenConns int    `yaml:"max-open-conns"`
	LogMode      int    `yaml:"logmode"`
	LogZap       bool   `yaml:"logzap"`
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}
