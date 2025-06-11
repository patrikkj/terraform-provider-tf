package subprovider

var registry = []Subprovider[any]{}

func RegisterSubprovider(subprovider Subprovider[any]) {
	registry = append(registry, subprovider)
}
