package dns

type DNS interface {
	Resolve(domain string) (string, error)
}