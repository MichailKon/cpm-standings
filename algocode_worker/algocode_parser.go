package algocode_worker

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"os"
	"strconv"
)

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	GroupShort string `json:"group_short"`
}

type Problem struct {
	Id    int    `json:"id"`
	Long  string `json:"long"`
	Short string `json:"short"`
	Index int    `json:"index"`
}

type UserSubmit struct {
	Score   int    `json:"score"`
	Penalty int    `json:"penalty"`
	Verdict string `json:"verdict"`
	Time    int    `json:"time"`
}

type Contest struct {
	Id          int                      `json:"id"`
	Date        string                   `json:"date"`
	EjudgeId    int                      `json:"ejudge_id"`
	Title       string                   `json:"title"`
	Coefficient float64                  `json:"coefficient"`
	Problems    []*Problem               `json:"problems"`
	Users       map[string][]*UserSubmit `json:"users"`
}

type SubmitsData struct {
	Users    []*User    `json:"users"`
	Contests []*Contest `json:"contests"`
}

func GetSubmitsData(url string) (data *SubmitsData) {
	client := resty.New()
	res, err := client.R().SetResult(&data).Get(url)
	if err != nil {
		slog.Warn("Error while querying algocode:", err.Error())
		return nil
	}
	if res.StatusCode() != 200 {
		slog.Warn("Algocode returned code", res.StatusCode())
		return nil
	}
	return
}

func LoadSubmitsData(filepath string) (data *SubmitsData) {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Warn("Error while opening standings file:", err.Error())
		return nil
	}
	defer file.Close()
	parser := json.NewDecoder(file)
	if err := parser.Decode(&data); err != nil {
		slog.Warn("Error while loading standings from file:", err.Error())
		return nil
	}
	return
}

func CreateUser2Id(data *SubmitsData) (res map[string]int) {
	res = make(map[string]int)
	for _, user := range data.Users {
		res[user.Name] = user.Id
	}
	return
}

func CreateId2User(data *SubmitsData) (res map[int]string) {
	res = make(map[int]string)
	for _, user := range data.Users {
		res[user.Id] = user.Name
	}
	return
}

func CreateContest2Id(data *SubmitsData) (res map[string]int) {
	res = make(map[string]int)
	for _, contest := range data.Contests {
		res[contest.Title] = contest.Id
	}
	return
}

func CreateId2Contest(data *SubmitsData) (res map[int]string) {
	res = make(map[int]string)
	for _, contest := range data.Contests {
		res[contest.Id] = contest.Title
	}
	return
}

// SolvedTasks is a list of solved tasks (bruh)
type SolvedTasks []string

// ContestsData ContestId -> SolvedTasks
type ContestsData map[int]SolvedTasks

// SolvedTable UserId -> ContestsData
type SolvedTable map[int]ContestsData

func UnpackSubmitsData(data *SubmitsData) SolvedTable {
	res := make(SolvedTable)
	for _, contest := range data.Contests {
		for userId, submits := range contest.Users {
			uid, _ := strconv.Atoi(userId)
			for i, submit := range submits {
				if submit.Score > 0 {
					if _, ok := res[uid]; !ok {
						res[uid] = make(ContestsData)
					}
					if _, ok := res[uid][contest.Id]; !ok {
						res[uid][contest.Id] = make(SolvedTasks, 0)
					}
					res[uid][contest.Id] = append(res[uid][contest.Id], contest.Problems[i].Short)
				}
			}
		}
	}
	return res
}
