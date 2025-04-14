package server

type Routes []Router

func (r Routes) Len() int {
	return len(r)
}

func (r Routes) Less(i, j int) bool {
	return r[i].pattern < r[j].pattern
}

func (r Routes) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}