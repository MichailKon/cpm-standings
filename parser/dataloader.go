package parser

import (
	"cpm-standings/config"
	"errors"
	"fmt"
	codeforcesapi "github.com/MichailKon/codeforces-api"
	"github.com/MichailKon/codeforces-api/objects"
	"github.com/MichailKon/codeforces-api/utils"
	"log/slog"
	"net/http"
	"strings"
)

// SolvedTasks is a list of solved tasks (bruh)
type SolvedTasks []string

// ContestsData ContestId -> SolvedTasks
type ContestsData map[int]SolvedTasks

// SolvedTable handle -> ContestsData
type SolvedTable map[string]ContestsData

func LoadData(session *codeforcesapi.CodeforcesSession, criteria config.Criteria) SolvedTable {
	res := SolvedTable{}
	for _, curCriteria := range criteria {
		var standings *objects.ContestStandings
		asManager := false
		for {
			cur, err := session.ContestStandings(
				curCriteria.ContestId,
				utils.
					NewContestStandingsParams().
					WithShowUnofficial(true).
					WithAsManager(asManager),
			)
			if err != nil {
				var cerr codeforcesapi.CodeforcesApiError
				if errors.As(err, &cerr) {
					if cerr.StatusCode == http.StatusServiceUnavailable {
						continue
					}
					if err.Error() == "contestId: Standings are not available" {
						slog.Info(fmt.Sprintf("Getting contest %v as manager", curCriteria.ContestId))
						asManager = true
						continue
					}
				}
				slog.Error("Error while getting standings;",
					"contestId", curCriteria.ContestId,
					"error", err)
				break
			}
			standings = cur
			break
		}
		if standings == nil {
			continue
		}
		for _, row := range standings.Rows {
			handle := strings.ToLower(row.Party.Members[0].Handle)
			if _, ok := res[handle]; !ok {
				res[handle] = make(ContestsData)
			}
			for problemInd, result := range row.ProblemResults {
				if result.Points == 1 {
					if _, ok := res[handle][curCriteria.ContestId]; !ok {
						res[handle][curCriteria.ContestId] = make(SolvedTasks, 0)
					}
					res[handle][curCriteria.ContestId] = append(res[handle][curCriteria.ContestId],
						standings.Problems[problemInd].Index)
				}
			}
		}
	}
	return res
}
