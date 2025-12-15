package mongo

type AchievementStatisticsResponse struct {
	TotalByType                    map[string]int `json:"totalByType"`
	TotalByPeriod                  map[string]int `json:"totalByPeriod"`
	TopStudents                    []map[string]interface{} `json:"topStudents"`
	CompetitionLevelDistribution   map[string]int `json:"competitionLevelDistribution"`
}