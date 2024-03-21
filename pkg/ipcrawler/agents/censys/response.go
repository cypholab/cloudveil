package censys

type ResponseHits struct {
	Ip string `json:"ip"`
}

type ResponseResult struct {
	Hits  []ResponseHits    `json:"hits"`
	Total int               `json:"total"`
	Links map[string]string `json:"links"`
}

type Response struct {
	Code   int            `json:"code"`
	Result ResponseResult `json:"result"`
}
