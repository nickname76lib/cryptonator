package cryptonator

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

// Minimum acceptable amounts of cryptocurrencies by Cryptonator.
// Minimum invoce amounts of EUR, USD, and RUR can be calculated by
// multiplying price of single unit of checkout cryptocurrency in needed fiat on minum payement amount of checkout cryptocurrency (below).
//
// https://rf.cryptonator.com/fees
// https://cryptonator.zendesk.com/hc/en-us/articles/203324122-Is-there-an-incoming-transaction-minimum-
const (
	// Minimum acceptable amount by Cryptonator in Bitcoin (BTC)
	MinimumPaymentAmountBitcoin float64 = 0.0001
	// Minimum acceptable amount by Cryptonator in Bitcoin Cash (BCH)
	MinimumPaymentAmountBitcoinCash float64 = 0.0001
	// Minimum acceptable amount by Cryptonator in Dash (DASH)
	MinimumPaymentAmountDash float64 = 0.001
	// Minimum acceptable amount by Cryptonator in Dogecoin (DOGE)
	MinimumPaymentAmountDogecoin float64 = 10
	// Minimum acceptable amount by Cryptonator in Ethereum (ETH)
	MinimumPaymentAmountEthereum float64 = 0.001
	// Minimum acceptable amount by Cryptonator in Litecoin (LTC)
	MinimumPaymentAmountLitecoin float64 = 0.001
	// Minimum acceptable amount by Cryptonator in Monero (XMR)
	MinimumPaymentAmountMonero float64 = 0.001
	// Minimum acceptable amount by Cryptonator in Ripple (XRP)
	MinimumPaymentAmountRipple float64 = 1
	// Minimum acceptable amount by Cryptonator in TetherUS (USDT)
	MinimumPaymentAmountTetherUS float64 = 10
	// Minimum acceptable amount by Cryptonator in Zcash (ZEC)
	MinimumPaymentAmountZcash float64 = 0.0001

	// Minimum acceptable amount by Cryptonator in US Dollar (equivalent)
	MinimumPaymentAmountUSDollar float64 = 1.0
	// Minimum acceptable amount by Cryptonator in Euro (equivalent)
	MinimumPaymentAmountEuro float64 = 1.0
	// Minimum acceptable amount by Cryptonator in Russian Ruble (equivalent)
	MinimumPaymentAmountRURuble float64 = 10.0
)

type coingeckoFiatPricesResoponse struct {
	Bitcoin     coingeckoFiatPricesResoponsePrices `json:"bitcoin"`
	BitcoinCash coingeckoFiatPricesResoponsePrices `json:"bitcoin-cash"`
	Dash        coingeckoFiatPricesResoponsePrices `json:"dash"`
	Dogecoin    coingeckoFiatPricesResoponsePrices `json:"dogecoin"`
	Ripple      coingeckoFiatPricesResoponsePrices `json:"ripple"`
	Ethereum    coingeckoFiatPricesResoponsePrices `json:"ethereum"`
	Litecoin    coingeckoFiatPricesResoponsePrices `json:"litecoin"`
	Monero      coingeckoFiatPricesResoponsePrices `json:"monero"`
	TetherUS    coingeckoFiatPricesResoponsePrices `json:"tether"`
	Zcash       coingeckoFiatPricesResoponsePrices `json:"zcash"`
}

type coingeckoFiatPricesResoponsePrices struct {
	USDollar float64 `json:"usd"`
	Euro     float64 `json:"eur"`
	RURuble  float64 `json:"rub"`
}

// Returns minimum fiat (usd, eur, rub) invoice amounts for each supported cryptocurrency
func GetMinimumInvoiceAmountsInFiat() (map[CheckoutCurrency]map[InvoiceCurrency]float64, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin%2Cbitcoin-cash%2Cdash%2Cdogecoin%2Cethereum%2Clitecoin%2Cmonero%2Cripple%2Ctether%2Czcash&vs_currencies=usd%2Ceur%2Crub")
	if err != nil {
		return nil, fmt.Errorf("GetMinimumInvoiceAmountsInFiat: %w", err)
	}
	defer resp.Body.Close()

	responseData := &coingeckoFiatPricesResoponse{}

	err = json.NewDecoder(resp.Body).Decode(responseData)
	if err != nil {
		return nil, fmt.Errorf("GetMinimumInvoiceAmountsInFiat: %w", err)
	}

	amountsMap := map[CheckoutCurrency]map[InvoiceCurrency]float64{}

	amountsMap[CheckoutCurrencyBitcoin] = calcMinimumFiatPrices(MinimumPaymentAmountBitcoin, responseData.Bitcoin)
	amountsMap[CheckoutCurrencyBitcoinCash] = calcMinimumFiatPrices(MinimumPaymentAmountBitcoinCash, responseData.Bitcoin)
	amountsMap[CheckoutCurrencyDash] = calcMinimumFiatPrices(MinimumPaymentAmountDash, responseData.Dash)
	amountsMap[CheckoutCurrencyDogecoin] = calcMinimumFiatPrices(MinimumPaymentAmountDogecoin, responseData.Dogecoin)
	amountsMap[CheckoutCurrencyEthereum] = calcMinimumFiatPrices(MinimumPaymentAmountEthereum, responseData.Ethereum)
	amountsMap[CheckoutCurrencyLitecoin] = calcMinimumFiatPrices(MinimumPaymentAmountLitecoin, responseData.Litecoin)
	amountsMap[CheckoutCurrencyMonero] = calcMinimumFiatPrices(MinimumPaymentAmountMonero, responseData.Monero)
	amountsMap[CheckoutCurrencyRipple] = calcMinimumFiatPrices(MinimumPaymentAmountRipple, responseData.Ripple)
	amountsMap[CheckoutCurrencyTetherUS] = calcMinimumFiatPrices(MinimumPaymentAmountTetherUS, responseData.TetherUS)
	amountsMap[CheckoutCurrencyZcash] = calcMinimumFiatPrices(MinimumPaymentAmountZcash, responseData.Zcash)

	return amountsMap, nil
}

func calcMinimumFiatPrices(minCryptocurPayment float64, prices coingeckoFiatPricesResoponsePrices) map[InvoiceCurrency]float64 {
	m := map[InvoiceCurrency]float64{
		InvoiceCurrencyUSDollar: calcMinimumPrice(minCryptocurPayment, prices.USDollar),
		InvoiceCurrencyEuro:     calcMinimumPrice(minCryptocurPayment, prices.Euro),
		InvoiceCurrencyRURuble:  calcMinimumPrice(minCryptocurPayment, prices.RURuble),
	}
	// Min. 1 USD: https://cryptonator.zendesk.com/hc/en-us/articles/208377889-Minimum-transactions
	if m[InvoiceCurrencyUSDollar] < MinimumPaymentAmountUSDollar {
		m[InvoiceCurrencyUSDollar] = MinimumPaymentAmountUSDollar
	}
	// Min. 1 EUR: https://cryptonator.zendesk.com/hc/en-us/articles/208377889-Minimum-transactions
	if m[InvoiceCurrencyEuro] < MinimumPaymentAmountEuro {
		m[InvoiceCurrencyEuro] = MinimumPaymentAmountEuro
	}
	// Min. 10 RUB: https://cryptonator.zendesk.com/hc/en-us/articles/208377889-Minimum-transactions
	if m[InvoiceCurrencyRURuble] < MinimumPaymentAmountRURuble {
		m[InvoiceCurrencyRURuble] = MinimumPaymentAmountRURuble
	}
	return m
}

func calcMinimumPrice(minCryptocurPayment float64, cryptoCurFiatPrice float64) float64 {
	return math.Ceil(cryptoCurFiatPrice*minCryptocurPayment*100) / 100
}
