package payment

import (
	"dvalyayevkbtu/my-booking/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type InvoiceCreated struct {
	Reference string `json:"reference"`
	Volume    string `json:"volume"`
	Currency  string `json:"currency"`
}

type Invoice struct {
	Reference       string       `json:"reference"`
	Volume          string       `json:"volume"`
	Currency        string       `json:"currency"`
	VolumeFulfilled string       `json:"volumeFulfilled"`
	Status          string       `json:"status"`
	Confirments     []Confirment `json:"confirments"`
}

type Confirment struct {
	Reference   string `json:"reference"`
	Volume      string `json:"volume"`
	Currency    string `json:"currency"`
	AccountCode string `json:"accountCode"`
}

type Payment struct {
	baseUrl string
	client  *http.Client
}

func CreatePayment(conf config.PaymentConfig) *Payment {
	return &Payment{conf.URL, &http.Client{}}
}

func (p *Payment) CreateInvoice(reference, volume, currency string) error {
	reqStr, jsonErr := json.Marshal(InvoiceCreated{reference, volume, currency})
	if jsonErr != nil {
		return jsonErr
	}

	url := fmt.Sprintf("%s/payment", p.baseUrl)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(reqStr)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("status is not accepted")
	}
	return nil
}

func (p *Payment) CheckPayment(reference string) (bool, error) {
	url := fmt.Sprintf("%s/payment/%s", p.baseUrl, reference)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, errors.New("unsuccessful status code of check payment response")
	}

	str, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var invoice Invoice
	err = json.Unmarshal(str, &invoice)
	if err != nil {
		return false, err
	}
	return invoice.Status == "FULFILLED", nil
}
