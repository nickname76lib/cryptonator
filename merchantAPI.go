package cryptonator

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

const merchantAPIEndpointV1URL = "https://api.cryptonator.com/api/merchant/v1"

// Cryptonator API instance
type MerchantAPI struct {
	// Your merchant ID.
	ID string
	// Your merchant secret.
	Secret string
	// net/http client to be used to perform api requests.
	HttpClient *http.Client
}

// A wrapper for API instancing
//
// merchantID - your merchant ID. (REQUIRED)
// secret - your merchant secret. (REQUIRED)
func NewMerchantAPI(merchantID string, merchantSecret string) *MerchantAPI {
	return &MerchantAPI{
		ID:         merchantID,
		Secret:     merchantSecret,
		HttpClient: http.DefaultClient,
	}
}

var ErrCryptonatorAPIErr = errors.New("cryptonator api error")

// Cryptonator invoice status
type InvoiceStatus string

const (
	// Unpaid - "unpaid"
	InvoiceStatusUnpaid InvoiceStatus = "unpaid"
	// Confirming - "confirming"
	InvoiceStatusConfirming InvoiceStatus = "confirming"
	// Paid - "paid". Final.
	InvoiceStatusPaid InvoiceStatus = "paid"
	// Cancelled - "cancelled". Final.
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
	// Mispaid - "mispaid". Final.
	InvoiceStatusMispaid InvoiceStatus = "mispaid"
)

type CreateInvoiceOptions struct {
	// orderID - order ID for your accounting purposes.
	OrderID string `url:"order_id,omitempty"`
	// Name of an item or service. (REQUIRED)
	ItemName string `url:"item_name"`
	// Description of an item or service.
	ItemDescription string `url:"item_name,omitempty"`
	// Invoice currency. (REQUIRED)
	InvoiceCurrency InvoiceCurrency `url:"invoice_currency"`
	// InvoiceAmount - invoice amount. (REQUIRED)
	InvoiceAmount float64 `url:"invoice_amount"`
	// CheckoutCurrency - checkout currency. (REQUIRED)
	CheckoutCurrency CheckoutCurrency `url:"checkout_currency"`
	// An URL to which users will be redirected after a successful payment. If undefined, default setting is used.
	SuccessURL string `url:"success_url,omitempty"`
	// An URL to which users will be redirected after a cancelled or failed payment. If undefined, default setting is used.
	FailedURL string `url:"failed_url,omitempty"`
	// Language of the checkout page. If empty, English is used by default.
	Language Language `url:"language,omitempty"`
}

type CreateInvoiceReponseData struct {
	// Invoice ID
	InvoiceID string `json:"invoice_id"`
	// URL of the hosted payment page
	InvoiceURL string `json:"invoice_url"`
	// Creation timestamp UTC
	InvoiceCreated int64 `json:"invoice_created"`
	// Expiration timestamp UTC
	InvoiceExpires int64 `json:"invoice_expires"`
	// Checkout cryptocurrency
	CheckoutCurrency CheckoutCurrency `json:"checkout_currency"`
	// Amount due in cryptocurrency
	CheckoutAmount float64 `json:"checkout_amount"`
	// Cryptocurrency payment address
	CheckoutAddress string `json:"checkout_address"`

	// Internal field for api errors checking (DON NOT CHECK, IT IS CHECKED BY THE LIBRARY)
	Error string `json:"error"`
}

// Create new invoice with a predefined checkout currency.
//
// Users cannot choose different payment method from checkoutCurrency argument.
//
// A successful request will create a hosted checkout page
// https://www.cryptonator.com/merchant/invoice/<invoice_id> and
// an HTTP-notification with all the invoice parameters will be
// sent to the provided callback URL.
//
// See required values for CreateInvoiceOptions in it's commentaries.
func (mapi *MerchantAPI) CreateInvoice(options *CreateInvoiceOptions) (*CreateInvoiceReponseData, error) {
	requestURL := merchantAPIEndpointV1URL + "/createinvoice"

	// Errors can't occure in this call
	qv, err := query.Values(options)
	if err != nil {
		panic(err)
	}

	qv.Add("merchant_id", mapi.ID)

	secretHashDataOrder := [...]string{
		"merchant_id",
		"item_name",
		"order_id",
		"item_description",
		"checkout_currency",
		"invoice_amount",
		"invoice_currency",
		"success_url",
		"failed_url",
		"language",
	}

	secretHashData := make([]string, 0, len(secretHashDataOrder)+1)
	for _, paramName := range secretHashDataOrder {
		vs := qv[paramName]
		if vs == nil {
			secretHashData = append(secretHashData, "")
			continue
		}

		secretHashData = append(secretHashData, strings.Join(vs, ","))
	}
	secretHashData = append(secretHashData, mapi.Secret)

	h := sha1.New()
	h.Write([]byte(strings.Join(secretHashData, "&")))
	secretHash := hex.EncodeToString(h.Sum(nil))

	qv.Add("secret_hash", secretHash)

	resp, err := mapi.HttpClient.PostForm(requestURL, qv)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice: %w", err)
	}
	defer resp.Body.Close()

	responseData := &CreateInvoiceReponseData{}

	err = json.NewDecoder(resp.Body).Decode(responseData)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice: %w", err)
	}

	if responseData.Error != "" {
		return nil, fmt.Errorf("CreateInvoice: %w: %v %v", ErrCryptonatorAPIErr, resp.Status, responseData.Error)
	}

	responseData.InvoiceURL = "https://www.cryptonator.com/merchant/invoice/" + responseData.InvoiceID

	return responseData, nil
}

type GetInvoiceResponseData struct {
	// Your ID for this order
	OrderID string `json:"order_id"`
	// Invoice status
	InvoiceStatus InvoiceStatus `json:"status"`
	// Invoice currency
	InvoiceCurrency InvoiceCurrency `json:"currency"`
	// Invoice amount. API returns string, but it can always be parsed as float64 without errors.
	InvoiceAmount string `json:"amount"`

	// Internal field for api errors checking (DON NOT CHECK, IT IS CHECKED BY THE LIBRARY)
	Error string `json:"error"`
}

// Get information and status of an invoice.
//
// invoiceID - Invoice ID. Can't be empty.
func (mapi *MerchantAPI) GetInvoice(invoiceID string) (*GetInvoiceResponseData, error) {
	requestURL := merchantAPIEndpointV1URL + "/getinvoice"

	h := sha1.New()
	h.Write([]byte(mapi.ID + "&" + invoiceID + "&" + mapi.Secret))
	secretHash := hex.EncodeToString(h.Sum(nil))

	qv := url.Values{}
	qv.Set("merchant_id", mapi.ID)
	qv.Set("invoice_id", invoiceID)
	qv.Set("secret_hash", secretHash)

	resp, err := mapi.HttpClient.PostForm(requestURL, qv)
	if err != nil {
		return nil, fmt.Errorf("GetInvoice: %w", err)
	}
	defer resp.Body.Close()

	responseData := &GetInvoiceResponseData{}

	err = json.NewDecoder(resp.Body).Decode(responseData)
	if err != nil {
		return nil, fmt.Errorf("GetInvoice: %w", err)
	}

	if responseData.Error != "" {
		return nil, fmt.Errorf("GetInvoice: %w: %v %v", ErrCryptonatorAPIErr, resp.Status, responseData.Error)
	}

	return responseData, nil
}

type ListInvoicesResponseData struct {
	// Count of invoices in InvoiceList
	InvoiceCount int `json:"invoice_count"`
	// Invoice ID list
	InvoiceList []string `json:"invoice_list"`

	// Internal field for api errors checking (DON NOT CHECK, IT IS CHECKED BY LIBRARY)
	Error string `json:"error"`
}

// Get a list of all your invoices filtered by given parameters.
//
// invoiceStatus - invoice status filter. Use empty string to not use the filter.
// invoiceCurrency - invoice currency filter. Use empty string to not use the filter.
// checkoutCurrency - checkout currency filter. Use empty string to not use the filter.
func (mapi *MerchantAPI) ListInvoices(invoiceStatus InvoiceStatus, invoiceCurrency InvoiceCurrency, checkoutCurrency CheckoutCurrency) (*ListInvoicesResponseData, error) {
	requestURL := merchantAPIEndpointV1URL + "/listinvoices"

	h := sha1.New()
	h.Write([]byte(mapi.ID + "&" + string(invoiceStatus) + "&" + string(invoiceCurrency) + "&" + string(checkoutCurrency) + "&" + mapi.Secret))
	secretHash := hex.EncodeToString(h.Sum(nil))

	qv := url.Values{}
	qv.Set("merchant_id", mapi.ID)
	if invoiceStatus != "" {
		qv.Set("invoice_status", string(invoiceStatus))
	}
	if invoiceCurrency != "" {
		qv.Set("invoice_currency", string(invoiceCurrency))
	}
	if checkoutCurrency != "" {
		qv.Set("checkout_currency", string(checkoutCurrency))
	}
	qv.Set("secret_hash", secretHash)

	resp, err := mapi.HttpClient.PostForm(requestURL, qv)
	if err != nil {
		return nil, fmt.Errorf("ListInvoices: %w", err)
	}
	defer resp.Body.Close()

	responseData := &ListInvoicesResponseData{}

	err = json.NewDecoder(resp.Body).Decode(responseData)
	if err != nil {
		return nil, fmt.Errorf("ListInvoices: %w", err)
	}

	if responseData.Error != "" {
		return nil, fmt.Errorf("ListInvoices: %w: %v %v", ErrCryptonatorAPIErr, resp.Status, responseData.Error)
	}

	return responseData, nil
}
