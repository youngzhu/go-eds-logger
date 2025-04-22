package reportor

import "testing"

func Test_retrieveWorkReportFromInternet(t *testing.T) {
	workReport, err := retrieveWorkReportFromInternet()

	if err != nil {
		t.Error(err)
	}

	if workReport.LastWeekWorkContent == "" {
		t.Error("LastWeekWorkContent should not empty.")
	}
}

func Test_retrieveWorkReportFromLocal(t *testing.T) {
	workReport := retrieveWorkReportFromLocal()

	if workReport.LastWeekWorkContent == "" {
		t.Error("LastWeekWorkContent should not empty.")
	}
}
