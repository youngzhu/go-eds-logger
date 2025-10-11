package logger

import (
	"testing"
	"time"
)

func TestEDSLogger_RetrieveWorkReportRandom(t *testing.T) {
	err := r.RetrieveWorkReportRandom()
	if err != nil {
		t.Errorf("RetrieveWorkReportRandom() error = %v", err)
	}

	if r.workReport.LastWeekWorkContent == "" {
		t.Error("LastWeekWorkContent should not be empty.")
	}
	if r.workReport.LastWeekStudyContent == "" {
		t.Error("LastWeekStudyContent should not be empty.")
	}
	if r.workReport.LastWeekSummary == "" {
		t.Error("LastWeekSummary should not be empty.")
	}
	if len(r.workReport.WorkPlan) == 0 {
		t.Error("WorkPlan should not be empty.")
	}
	if r.workReport.StudyPlan == "" {
		t.Error("StudyPlan should not be empty.")
	}
}

func TestEDSLogger_RetrieveWorkReportRandom_humanable(t *testing.T) {
	err := r.RetrieveWorkReportRandom()
	if err != nil {
		t.Errorf("RetrieveWorkReportRandom() error = %v", err)
	}

	// 检查周报内容
	t.Logf("上周工作任务完成情况: %s", r.workReport.LastWeekWorkContent)
	t.Logf("上周学习完成任务情况: %s", r.workReport.LastWeekStudyContent)
	t.Logf("经验和收获总结: %s", r.workReport.LastWeekSummary)
	t.Logf("本周学习计划: %s", r.workReport.StudyPlan)
	t.Logf("本周工作计划与重点: %s", r.workReport.workPlanWeekly())

	// 日报，周一至周五随机获取一条工作计划
	for i := 0; i < 5; i++ {
		// 加上休眠，避免随机数相同
		time.Sleep(1 * time.Second)
		t.Logf("第 %d 天: %s", i+1, r.workReport.workPlanDaily())
	}

}

func TestEDSLogger_RetrieveLogContentViaWeb(t *testing.T) {
	err := lg.RetrieveLogContentViaWeb()
	if err != nil {
		t.Errorf("RetrieveWorkReportRandom() error = %v", err)
	}

	// 检查周报内容
	t.Logf("上周工作任务完成情况: %s", r.workReport.LastWeekWorkContent)
	t.Logf("上周学习完成任务情况: %s", r.workReport.LastWeekStudyContent)
	t.Logf("经验和收获总结: %s", r.workReport.LastWeekSummary)
	t.Logf("本周学习计划: %s", r.workReport.StudyPlan)
	t.Logf("本周工作计划与重点: %s", r.workReport.workPlanWeekly())

	// 日报，周一至周五随机获取一条工作计划
	for i := 0; i < 5; i++ {
		// 加上休眠，避免随机数相同
		time.Sleep(1 * time.Second)
		t.Logf("第 %d 天: %s", i+1, r.workReport.workPlanDaily())
	}

}
