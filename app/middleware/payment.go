package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var SnapClient snap.Client

func InitMidtrans() {
	SnapClient.New("SB-Mid-server-cf2nmb6kJ3mbVZ-R3XpLTOuq", midtrans.Sandbox)
}

func VerifyMidtransTrx(trxId string) (*MidtransResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "https://api.sandbox.midtrans.com/v2/" + trxId + "/status"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("SB-Mid-server-cf2nmb6kJ3mbVZ-R3XpLTOuq", "")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("midtrans error: %s", string(bodyBytes))
	}

	var midtransResponse MidtransResponse
	if err := json.NewDecoder(resp.Body).Decode(&midtransResponse); err != nil {
		return nil, err
	}

	return &midtransResponse, nil
}

type MidtransResponse struct {
	TransactionStatus string `json:"transaction_status"`
}