package crontab

import "testing"

var schTest = []struct {
	s   string
	cnt [5]int
}{
	{"* * * * *", [5]int{60, 24, 31, 12, 7}},
	{"*/2 * * * *", [5]int{30, 24, 31, 12, 7}},
	{"*/10 * * * *", [5]int{6, 24, 31, 12, 7}},
	{"* * * * */2", [5]int{60, 24, 0, 12, 4}},
	{"5,8,9 */2 2,3 * */2", [5]int{3, 12, 2, 12, 4}},
	{"* 5-11 2-30/2 * *", [5]int{60, 7, 15, 12, 0}},
	{"1,2,5-8 * * */3 *", [5]int{6, 24, 31, 4, 7}},
}

// TestSchedule parse the crontab syntax and compare number of target min/hour/days/month with expected ones
func TestSchedule(t *testing.T) {
	for _, sch := range schTest {
		j, err := parseSchedule(sch.s)
		if err != nil {
			t.Error(err)
		}

		if len(j.min) != sch.cnt[0] {
			t.Error(sch.s, "min count expected to be", sch.cnt[0], "result", len(j.min), j.min)
		}

		if len(j.hour) != sch.cnt[1] {
			t.Error(sch.s, "hour count expected to be", sch.cnt[1], "result", len(j.hour), j.hour)
		}

		if len(j.day) != sch.cnt[2] {
			t.Error(sch.s, "day count expected to be", sch.cnt[2], "result", len(j.day), j.day)
		}

		if len(j.month) != sch.cnt[3] {
			t.Error(sch.s, "month count expected to be", sch.cnt[3], "result", len(j.month), j.month)
		}

		if len(j.dayOfWeek) != sch.cnt[4] {
			t.Error(sch.s, "dayOfWeek count expected to be", sch.cnt[4], "result", len(j.dayOfWeek), j.dayOfWeek)
		}
	}
}

var schErrorTest = []string{
	"* * * * * *",
	"1 2 3 4 5 6",
	"*/ 2 * * * *",
	"1,2,3/10 * * * *",
	"1,2,3,1-15/10 * * * *",
	"a b c d e",
}

// TestScheduleError tests crontab syntax which should not be accepted
func TestScheduleError(t *testing.T) {
	for _, s := range schErrorTest {
		if _, err := parseSchedule(s); err == nil {
			t.Error(s, "should be error")
		}
	}
}