package conf

type Config struct {
	Listen    map[string]string
	AddrGroup map[string][]string
	DnsGroup  []string
	Router    map[string][]string
}
