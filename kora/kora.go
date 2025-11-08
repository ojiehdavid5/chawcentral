package kora

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ojiehdavid5/campusbyte/model"
)

func CreateKoraPayment(reference string, user model.User, amount float64) (string, error) {
	apiKey := os.Getenv("KORA_SECRET_KEY")

	payload := map[string]interface{}{
		"amount":    amount,
		"currency":  "NGN",
		"reference": reference,
		"redirect_url": "https://campusbite.ng/payment/success",
		"customer": map[string]interface{}{
			"name":  fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			"email": fmt.Sprintf("%s@campusbite.ng", user.Username),
		},
	}

	data, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.korapay.com/merchant/api/v1/charges", bytes.NewBuffer(data))
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Data struct {
			CheckoutURL string `json:"checkout_url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Data.CheckoutURL, nil
}