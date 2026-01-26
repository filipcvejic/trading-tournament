package tradingaccount

type Crypto interface {
	Encrypt(plain string) (string, error)
	Decrypt(cipher string) (string, error)
}
