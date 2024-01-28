package data_worker

import (
	"cpm-standings/config"
	"cpm-standings/parser"
	codeforces_api "github.com/MichailKon/codeforces-api"
	"log/slog"
	"slices"
	"strings"
)

type StudentContestData struct {
	Solved int
	Mark   int
}

type Student struct {
	Name         string
	Handle       string
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

func ExportStudentsData(
	session *codeforces_api.CodeforcesSession,
	mapping config.StudentsHandlesMapping,
	criteria config.Criteria,
) (res *StudentsExportData) {
	res = &StudentsExportData{
		ContestTitles: make([]string, 0),
		Students:      make([]*Student, 0),
	}

	for handle, name := range mapping {
		res.Students = append(res.Students, &Student{
			Name:         name,
			Handle:       handle,
			ContestsData: make([]*StudentContestData, 0),
		})
	}

	table := parser.LoadData(session, criteria)

	for _, groups := range criteria {
		for _, group := range groups.Groups {
			slices.Sort(group.Tasks)
		}
	}
	for _, contestsData := range table {
		for _, tasks := range contestsData {
			slices.Sort(tasks)
		}
	}

	contestIds := make([]int, 0)
	for _, i := range criteria {
		contestIds = append(contestIds, i.ContestId)
	}
	slices.Sort(contestIds)
	slices.Reverse(contestIds)
	for _, contestId := range contestIds {
		contestTitle := criteria.GetContestName(contestId)
		taskGroup := criteria[contestTitle]
		slog.Info("Filling contest", contestId)
		for _, group := range taskGroup.Groups {
			slog.Info("Filling task group", group.Name)
			res.ContestTitles = append(res.ContestTitles, group.Name)
			for i, student := range res.Students {
				solved := intersection(table[student.Handle][contestId], group.Tasks)
				res.Students[i].ContestsData = append(res.Students[i].ContestsData, &StudentContestData{
					Solved: solved,
					Mark:   convertMark(10 * float64(solved) / float64(group.Norm)),
				})
			}
		}
	}
	slices.SortFunc(res.Students, func(a, b *Student) int {
		return strings.Compare(a.Name, b.Name)
	})
	return
}
