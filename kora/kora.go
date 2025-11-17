package kora

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ojiehdavid5/campusbyte/model"
)

// CreateKoraPayment initializes a payment and returns the checkout URL
func CreateKoraPayment(user model.User, amount float64, purpose string) (string, error) {
	apiKey := os.Getenv("KORA_SECRET_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("KORA_SECRET_KEY not set in environment")
	}

	// Unique reference (important for webhook verification)
	reference := fmt.Sprintf("CB-%d-%d", user.ID, time.Now().Unix())

payload := map[string]interface{}{
	"amount":              6000,
	"redirect_url":        "https://korapay.com?order=1234",
	"currency":            "NGN",
	"reference":           reference,
	"narration":           "Payment for product Y",
	"merchant_bears_cost": false,
	"customer": map[string]interface{}{
		"name":  "Chigozie Madubuko",
		"email": "chigoziemadubuko@gmail.com",
	},
	"notification_url": "https://webhook.site/8d321d8d-397f-4bab-bf4d-7e9ae3afbd50",
}

	data, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.korapay.com/merchant/api/v1/charges/initialize",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Kora API: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("🔍 Kora API Response:", string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Kora API returned status: %s", resp.Status)
	}

	var res struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Reference   string `json:"reference"`
			CheckoutURL string `json:"checkout_url"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("failed to decode Kora response: %w", err)
	}

	if !res.Status || res.Data.CheckoutURL == "" {
		return "", fmt.Errorf("Kora API error: %s", res.Message)
	}

	fmt.Printf("✅ Kora payment initialized for %s (₦%.2f)\nCheckout URL: %s\n",
		user.FirstName, amount, res.Data.CheckoutURL)

	return res.Data.CheckoutURL, nil
}
