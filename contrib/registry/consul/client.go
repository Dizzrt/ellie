package consul

type Datacenter string

const (
	SingleDatacenter Datacenter = "SINGLE"
	MultiDatacenter  Datacenter = "MULTI"
)

type Client struct {
}
