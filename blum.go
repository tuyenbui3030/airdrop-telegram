// File: go.mod

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// API response structures
type TokenResponse struct {
	Token struct {
		Access string `json:"access"`
	} `json:"token"`
}

type BalanceResponse struct {
	AvailableBalance     string `json:"availableBalance"`
	PlayPasses           int    `json:"playPasses"`
	IsFastFarmingEnabled bool   `json:"isFastFarmingEnabled"`
	Timestamp            int64  `json:"timestamp"`
	Farming              struct {
		StartTime    int64  `json:"startTime"`
		EndTime      int64  `json:"endTime"`
		EarningsRate string `json:"earningsRate"`
		Balance      string `json:"balance"`
	} `json:"farming"`
}

type UserResponse struct {
	Username string `json:"username"`
}

type FarmRewardResponse struct {
	Message string `json:"message"`
}

type FarmDailyResponse struct {
	Message string `json:"message"`
}

type FarmingResponse struct {
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	EarningsRate string `json:"earningsRate"`
	Balance      string `json:"balance"`
}

type GamePlayResponse struct {
	GameID string `json:"gameId"`
}

type GameClaimPayload struct {
	GameID string `json:"gameId"`
	Points int    `json:"points"`
}

type Task struct {
	ID                   string              `json:"id"`
	Kind                 string              `json:"kind"`
	Type                 string              `json:"type"`
	Status               string              `json:"status"`
	ValidationType       string              `json:"validationType"`
	IconFileKey          string              `json:"iconFileKey"`
	BannerFileKey        *string             `json:"bannerFileKey,omitempty"`
	Title                string              `json:"title"`
	ProductName          *string             `json:"productName,omitempty"`
	Description          *string             `json:"description,omitempty"`
	Reward               string              `json:"reward"`
	SubTasks             []SubTask           `json:"subTasks,omitempty"`
	ProgressTarget       *ProgressTarget     `json:"progressTarget,omitempty"`
	SocialSubscription   *SocialSubscription `json:"socialSubscription,omitempty"`
	IsHidden             bool                `json:"isHidden"`
	IsDisclaimerRequired bool                `json:"isDisclaimerRequired"`
}

type SubTask struct {
	ID                   string              `json:"id"`
	Kind                 string              `json:"kind"`
	Type                 string              `json:"type"`
	Status               string              `json:"status"`
	ValidationType       string              `json:"validationType"`
	IconFileKey          string              `json:"iconFileKey"`
	Title                string              `json:"title"`
	ProductName          *string             `json:"productName,omitempty"`
	Reward               string              `json:"reward"`
	SocialSubscription   *SocialSubscription `json:"socialSubscription,omitempty"`
	IsDisclaimerRequired bool                `json:"isDisclaimerRequired"`
}

type SocialSubscription struct {
	OpenInTelegram bool   `json:"openInTelegram"`
	URL            string `json:"url"`
}

type ProgressTarget struct {
	Target   string `json:"target"`
	Progress string `json:"progress"`
	Accuracy int    `json:"accuracy"`
	Postfix  string `json:"postfix"`
}

type Section struct {
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Tasks       []Task    `json:"tasks"`
	SubSections []Section `json:"subSections,omitempty"`
}

type TasksResponse struct {
	Sections []Section `json:"sections"`
}

// API functions
func getToken() (string, error) {
	queryId := os.Getenv("QUERY_ID")
	url := "https://user-domain.blum.codes/api/v1/auth/provider/PROVIDER_TELEGRAM_MINI_APP"
	payload := map[string]string{
		"query":         queryId,
		"referralToken": "vTHusRz4j0",
	}

	resp, err := postJSON(url, payload)
	if err != nil {
		return "", err
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(resp, &tokenResp)
	if err != nil {
		return "", err
	}

	return "Bearer " + tokenResp.Token.Access, nil
}

func getUsername(token string) (string, error) {
	url := "https://gateway.blum.codes/v1/user/me"
	resp, err := getWithAuth(url, token)
	if err != nil {
		return "", err
	}

	var userResp UserResponse
	err = json.Unmarshal(resp, &userResp)
	if err != nil {
		return "", err
	}

	return userResp.Username, nil
}

func getBalance(token string) (*BalanceResponse, error) {
	url := "https://game-domain.blum.codes/api/v1/user/balance"
	resp, err := getWithAuth(url, token)
	if err != nil {
		return nil, err
	}

	var balanceResp BalanceResponse
	err = json.Unmarshal(resp, &balanceResp)
	if err != nil {
		return nil, err
	}
	return &balanceResp, nil
}

func claimFarmReward(token string) (string, error) {
	url := "https://game-domain.blum.codes/api/v1/farming/claim"
	resp, err := postWithAuth(url, token, nil)
	if err != nil {
		return "", err
	}

	var rewardResp FarmRewardResponse
	err = json.Unmarshal(resp, &rewardResp)
	if err != nil {
		return "", err
	}

	return rewardResp.Message, nil
}

func startFarmingSession(token string) (*FarmingResponse, error) {
	url := "https://game-domain.blum.codes/api/v1/farming/start"
	resp, err := postWithAuth(url, token, nil)
	if err != nil {
		return nil, err
	}

	var startFarmingResp FarmingResponse
	err = json.Unmarshal(resp, &startFarmingResp)
	if err != nil {
		return nil, err
	}

	return &startFarmingResp, nil
}

func getIdGame(token string) (*GamePlayResponse, error) {
	url := "https://game-domain.blum.codes/api/v1/game/play"
	resp, err := postWithAuth(url, token, nil)
	if err != nil {
		return nil, err
	}

	var gamePlayResp GamePlayResponse
	err = json.Unmarshal(resp, &gamePlayResp)
	if err != nil {
		return nil, err
	}
	return &gamePlayResp, nil
}

func claimGamePoins(token string, gameId string, points int) (string, error) {
	url := "https://game-domain.blum.codes/api/v1/game/claim"

	payload := GameClaimPayload{
		GameID: gameId,
		Points: points,
	}

	resp, err := postWithAuth(url, token, payload)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func claimDailyReward(token string) (*FarmDailyResponse, error) {
	url := "https://game-domain.blum.codes/api/v1/daily-reward?offset=-420"
	resp, err := postWithAuth(url, token, nil)
	if err != nil {
		return nil, err
	}

	var farmDailyResponse FarmDailyResponse
	err = json.Unmarshal(resp, &farmDailyResponse)
	if err != nil {
		return nil, err
	}

	return &farmDailyResponse, nil
}

func getTasks(token string) ([]Section, error) {
	url := "https://earn-domain.blum.codes/api/v1/tasks"
	resp, err := getWithAuth(url, token)
	if err != nil {
		return nil, err
	}

	var sectionsResp []Section
	err = json.Unmarshal(resp, &sectionsResp)
	if err != nil {
		return nil, err
	}

	return sectionsResp, nil

	//for _, section := range tasksResp.Sections {
	//	fmt.Printf("Section Title: %s\n", section.Title)
	//	if section.Description != nil {
	//		fmt.Printf("Section Description: %s\n", *section.Description)
	//	}
	//	for _, task := range section.Tasks {
	//		fmt.Printf("Task ID: %s\n", task.ID)
	//		fmt.Printf("  Kind: %s\n", task.Kind)
	//		fmt.Printf("  Type: %s\n", task.Type)
	//		fmt.Printf("  Status: %s\n", task.Status)
	//		fmt.Printf("  Validation Type: %s\n", task.ValidationType)
	//		fmt.Printf("  Icon File Key: %s\n", task.IconFileKey)
	//		if task.BannerFileKey != nil {
	//			fmt.Printf("  Banner File Key: %s\n", *task.BannerFileKey)
	//		}
	//		fmt.Printf("  Title: %s\n", task.Title)
	//		if task.ProductName != nil {
	//			fmt.Printf("  Product Name: %s\n", *task.ProductName)
	//		}
	//		if task.Description != nil {
	//			fmt.Printf("  Description: %s\n", *task.Description)
	//		}
	//		fmt.Printf("  Reward: %s\n", task.Reward)
	//		for _, subTask := range task.SubTasks {
	//			fmt.Printf("  SubTask ID: %s\n", subTask.ID)
	//			fmt.Printf("    Kind: %s\n", subTask.Kind)
	//			fmt.Printf("    Type: %s\n", subTask.Type)
	//			fmt.Printf("    Status: %s\n", subTask.Status)
	//			fmt.Printf("    Validation Type: %s\n", subTask.ValidationType)
	//			fmt.Printf("    Icon File Key: %s\n", subTask.IconFileKey)
	//			fmt.Printf("    Title: %s\n", subTask.Title)
	//			if subTask.ProductName != nil {
	//				fmt.Printf("    Product Name: %s\n", *subTask.ProductName)
	//			}
	//			fmt.Printf("    Reward: %s\n", subTask.Reward)
	//		}
	//		if task.ProgressTarget != nil {
	//			fmt.Printf("  Progress Target - Target: %s, Progress: %s, Accuracy: %d, Postfix: %s\n",
	//				task.ProgressTarget.Target, task.ProgressTarget.Progress, task.ProgressTarget.Accuracy, task.ProgressTarget.Postfix)
	//		}
	//		if task.SocialSubscription != nil {
	//			fmt.Printf("  Social Subscription - Open In Telegram: %v, URL: %s\n",
	//				task.SocialSubscription.OpenInTelegram, task.SocialSubscription.URL)
	//		}
	//		fmt.Printf("  Is Hidden: %v\n", task.IsHidden)
	//		fmt.Printf("  Is Disclaimer Required: %v\n", task.IsDisclaimerRequired)
	//	}
	//}

}

//func claimFarmReward(token string) error {
//	url := "https://game-domain.blum.codes/api/v1/farming/claim"
//	resp, err := postWithAuth(url, token, nil)
//	if err != nil {
//		return fmt.Errorf("error claiming farm reward: %w", err)
//	}
//
//	var rewardResp FarmRewardResponse
//
//	err = json.Unmarshal(resp, &rewardResp)
//	if err != nil {
//		return fmt.Errorf("error unmarshaling farm reward response: %w", err)
//	}
//	fmt.Printf(rewardResp.Message)
//
//	if rewardResp.Message == "It's too early to claim" {
//		return fmt.Errorf("🚨 Claim failed! It's too early to claim")
//	} else if rewardResp.Message == "Need to start farm" {
//		return fmt.Errorf("🚨 Claim failed! Need to start farm")
//	}
//
//	fmt.Println("✅ Farm reward claimed successfully!")
//	return nil
//}

// Helper functions for HTTP requests
func postJSON(url string, payload interface{}) ([]byte, error) {
	client := &http.Client{}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Handle the error from closing the response body
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()

	return io.ReadAll(resp.Body)
}

func getWithAuth(url, token string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Handle the error from closing the response body
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()

	return io.ReadAll(resp.Body)
}

func postWithAuth(url, token string, payload interface{}) ([]byte, error) {
	client := &http.Client{}
	var req *http.Request
	var err error

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	} else {
		req, err = http.NewRequest("POST", url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Main function
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get token
	token, err := getToken()
	if err != nil {
		log.Fatal("Error getting token:", err)
	}

	// Get username and balance
	username, err := getUsername(token)
	if err != nil {
		log.Fatal("Error getting username:", err)
	}
	//
	//balance, err := getBalance(token)
	//if err != nil {
	//	log.Fatal("Error getting balance:", err)
	//}
	//
	//claimResult, err := claimFarmReward(token)
	//if err != nil {
	//	log.Fatal("Error claim farm reward:", err)
	//}
	//
	//dailyResult, err := claimDailyReward(token)
	//if err != nil {
	//	log.Fatal("🚨 Error daily farming:", err)
	//}

	fmt.Printf("👋 Hello, %s!\n", username)
	//fmt.Printf("💰 Your current BLUM balance is: %s\n", balance.AvailableBalance)
	//fmt.Printf("🎮 Your chances to play the game: %d\n", balance.PlayPasses)
	//fmt.Printf("🌾 Your farming balance: %s\n", balance.Farming.Balance)
	//fmt.Printf("🍄 Your claim farm reward: %s\n", claimResult)
	//fmt.Printf("🔃 Daily farming: %s\n", dailyResult.Message)
	//
	//if claimResult == "Need to start farm" {
	//	farmingResult, _ := startFarmingSession(token)
	//	fmt.Println(farmingResult.Balance)
	//}
	//
	//if balance.PlayPasses > 0 {
	//	infoGame, err := getIdGame(token)
	//	if err != nil {
	//		log.Fatal("🚨 Error getting idgame info:", err)
	//	}
	//
	//	source := rand.NewSource(time.Now().UnixNano())
	//	r := rand.New(source)
	//
	//	minValue := 200
	//	maxValue := 240
	//	points := r.Intn(maxValue-minValue+1) + minValue
	//
	//	fmt.Printf("💳 Your GameID: %s\n", infoGame.GameID)
	//	fmt.Printf("🪙 Your Points: %d\n", points)
	//	time.Sleep(60 * time.Second)
	//	status, err := claimGamePoins(token, infoGame.GameID, points)
	//	if err != nil {
	//		log.Fatal("Error getting status info:", err)
	//	}
	//	fmt.Printf("⌛ Status Game: %s\n", status)
	//} else {
	//	fmt.Println("🎰 Turn over")
	//}

	tasksData, err := getTasks(token)
	if err != nil {
		log.Fatal("🚨 Error get tasks: ", err)
	}

	for _, categoryTask := range tasksData {
		if len(categoryTask.Tasks) > 0 && len(categoryTask.Tasks[0].SubTasks) > 0 {
			//for (const task of category.tasks[0].subTasks)

			fmt.Printf("🚨Category: %s\n", categoryTask.Title)
			for _, task := range categoryTask.Tasks[0].SubTasks {
				//if (task.status === 'FINISHED') {
				//	console.log(
				//		`⏭️  Task "${task.title}" is already completed.`.cyan
				//	);
				//}
				if task.Status == "FINISHED" {
					fmt.Printf("⏭️ %s already completed.\n", task.Title)
				} else if task.Status == "NOT_STARTED" {
					fmt.Printf("⏭️ %s not completed.\n", task.Title)
				} else if task.Status == "STARTED" || task.Status == "READY_FOR_CLAIM" {
					fmt.Printf("✅ %s has been claimed!\n", task.Title)
				}
			}
		}
	}
}
