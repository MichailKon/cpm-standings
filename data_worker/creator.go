package data_worker

import (
	"cpm-standings/config"
	"cpm-standings/parser"
	"cpm-standings/utils"
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
			Handle:       strings.ToLower(handle),
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
		slog.Info("Filling", "contestId", contestId)
		for _, group := range taskGroup.Groups {
			slog.Info("Filling", "taskGroup", group.Name)
			res.ContestTitles = append(res.ContestTitles, group.Name)
			for i, student := range res.Students {
				solved := utils.Intersection(table[student.Handle][contestId], group.Tasks)
				res.Students[i].ContestsData = append(res.Students[i].ContestsData, &StudentContestData{
					Solved: solved,
					Mark:   utils.ConvertMark(10 * float64(solved) / float64(group.Norm)),
				})
			}
		}
	}
	slices.SortFunc(res.Students, func(a, b *Student) int {
		return strings.Compare(a.Name, b.Name)
	})
	return
}
