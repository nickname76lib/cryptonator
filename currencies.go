package cryptonator

type Currency string

const (
	// US dollar (USD) - "usd"
	CurrencyUSDollar Currency = "usd"
	// Euro (EUR) - "eur"
	CurrencyEuro Currency = "eur"
	// Russian ruble (RUB) - "rur"
	CurrencyRURuble Currency = "rur"

	// Bitcoin (BTC) - "bitcoin"
	CurrencyBitcoin Currency = "bitcoin"
	// Bitcoin Cash (BCH) - "bitcoincash"
	CurrencyBitcoinCash Currency = "bitcoincash"
	// Dash (DASH) - "dash"
	CurrencyDash Currency = "dash"
	// Dogecoin (DOGE) - "dogecoin"
	CurrencyDogecoin Currency = "dogecoin"
	// Ethereum (ETH) - "ethereum"
	CurrencyEthereum Currency = "ethereum"
	// Litecoin (LTC) - "litecoin"
	CurrencyLitecoin Currency = "litecoin"
	// Monero (XMR) - "monero"
	CurrencyMonero Currency = "monero"
	// Ripple (XRP) - "ripple"
	CurrencyRipple Currency = "ripple"
	// TetherUS (USDT) - "usdt"
	CurrencyTetherUS Currency = "usdt"
	// Zcash (ZEC) - "zcash"
	CurrencyZcash Currency = "zcash"
)

// Cryptonator invoice currency
type InvoiceCurrency Currency

const (
	// US dollar (USD) - "usd"
	InvoiceCurrencyUSDollar = InvoiceCurrency(CurrencyUSDollar)
	// Euro (EUR) - "eur"
	InvoiceCurrencyEuro = InvoiceCurrency(CurrencyEuro)
	// Russian ruble (RUB) - "rur"
	InvoiceCurrencyRURuble = InvoiceCurrency(CurrencyRURuble)

	// Bitcoin (BTC) - "bitcoin"
	InvoiceCurrencyBitcoin = InvoiceCurrency(CurrencyBitcoin)
	// Bitcoin Cash (BCH) - "bitcoincash"
	InvoiceCurrencyBitcoinCash = InvoiceCurrency(CurrencyBitcoinCash)
	// Dash (DASH) - "dash"
	InvoiceCurrencyDash = InvoiceCurrency(CurrencyDash)
	// Dogecoin (DOGE) - "dogecoin"
	InvoiceCurrencyDogecoin = InvoiceCurrency(CurrencyDogecoin)
	// Ethereum (ETH) - "ethereum"
	InvoiceCurrencyEthereum = InvoiceCurrency(CurrencyEthereum)
	// Litecoin (LTC) - "litecoin"
	InvoiceCurrencyLitecoin = InvoiceCurrency(CurrencyLitecoin)
	// Monero (XMR) - "monero"
	InvoiceCurrencyMonero = InvoiceCurrency(CurrencyMonero)
	// Ripple (XRP) - "ripple"
	InvoiceCurrencyRipple = InvoiceCurrency(CurrencyRipple)
	// TetherUS (USDT) - "usdt"
	InvoiceCurrencyTetherUS = InvoiceCurrency(CurrencyTetherUS)
	// Zcash (ZEC) - "zcash"
	InvoiceCurrencyZcash = InvoiceCurrency(CurrencyZcash)
)

// Cryptonator checkout currency
type CheckoutCurrency Currency

const (
	// Bitcoin (BTC) - "bitcoin"
	CheckoutCurrencyBitcoin = CheckoutCurrency(CurrencyBitcoin)
	// Bitcoin Cash (BCH) - "bitcoincash"
	CheckoutCurrencyBitcoinCash = CheckoutCurrency(CurrencyBitcoinCash)
	// Dash (DASH) - "dash"
	CheckoutCurrencyDash = CheckoutCurrency(CurrencyDash)
	// Dogecoin (DOGE) - "dogecoin"
	CheckoutCurrencyDogecoin = CheckoutCurrency(CurrencyDogecoin)
	// Ethereum (ETH) - "ethereum"
	CheckoutCurrencyEthereum = CheckoutCurrency(CurrencyEthereum)
	// Litecoin (LTC) - "litecoin"
	CheckoutCurrencyLitecoin = CheckoutCurrency(CurrencyLitecoin)
	// Monero (XMR) - "monero"
	CheckoutCurrencyMonero = CheckoutCurrency(CurrencyMonero)
	// Ripple (XRP) - "ripple"
	CheckoutCurrencyRipple = CheckoutCurrency(CurrencyRipple)
	// TetherUS (USDT) - "usdt"
	CheckoutCurrencyTetherUS = CheckoutCurrency(CurrencyTetherUS)
	// Zcash (ZEC) - "zcash"
	CheckoutCurrencyZcash = CheckoutCurrency(CurrencyZcash)
)
