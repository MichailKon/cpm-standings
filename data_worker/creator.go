package data_worker

import (
	"cpm-standings/algocode_worker"
	"cpm-standings/config"
	"log/slog"
	"slices"
)

type StudentContestData struct {
	Solved int
	Mark   int
}

type Student struct {
	Name         string
	ContestsData []*StudentContestData
}

type StudentsExportData struct {
	ContestTitles []string
	Students      []*Student
}

func intersection(a, b []string) (res int) {
	l, r := 0, 0
	for l < len(a) && r < len(b) {
		if a[l] == a[r] {
			res++
			l++
			r++
		} else if a[l] < a[r] {
			l++
		} else {
			r++
		}
	}
	return
}

func convertMark(mark float64) int {
	if mark < 3.5 {
		return 2
	} else if mark < 5.5 {
		return 3
	} else if mark < 7.5 {
		return 4
	} else {
		return 5
	}
}

func ExportStudentsData(data *algocode_worker.SubmitsData, criteria config.Criteria) (res *StudentsExportData) {
	res = &StudentsExportData{
		ContestTitles: make([]string, 0),
		Students:      make([]*Student, 0),
	}
	for _, user := range data.Users {
		res.Students = append(res.Students, &Student{
			Name:         user.Name,
			ContestsData: make([]*StudentContestData, 0),
		})
	}

	table := algocode_worker.UnpackSubmitsData(data)
	contest2id := algocode_worker.CreateContest2Id(data)
	user2id := algocode_worker.CreateUser2Id(data)

	for _, groups := range criteria {
		for _, group := range groups {
			slices.Sort(group.Tasks)
		}
	}
	for _, contestsData := range table {
		for _, tasks := range contestsData {
			slices.Sort(tasks)
		}
	}

	keys := make([]string, 0)
	for contestTitle := range criteria {
		keys = append(keys, contestTitle)
	}
	slices.Sort(keys)
	for _, contestTitle := range keys {
		taskGroup := criteria[contestTitle]
		slog.Info("Filling contest", contestTitle)
		for _, group := range taskGroup {
			slog.Info("Filling task group", group.Name)
			res.ContestTitles = append(res.ContestTitles, group.Name)
			contestId := contest2id[contestTitle]
			for i, student := range res.Students {
				studentId := user2id[student.Name]
				solved := intersection(table[studentId][contestId], group.Tasks)
				res.Students[i].ContestsData = append(res.Students[i].ContestsData, &StudentContestData{
					Solved: solved,
					Mark:   convertMark(10 * float64(solved) / float64(group.Norm)),
				})
			}
		}
	}
	return
}
