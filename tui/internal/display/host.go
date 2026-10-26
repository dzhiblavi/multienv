package display

import "hash/fnv"

func HostAlias(host, alias string) string {
	colors := [...]string{"6", "3", "5", "2", "4", "1"}
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(host))

	return Style(Clean(alias), colors[hash.Sum32()%uint32(len(colors))])
}
