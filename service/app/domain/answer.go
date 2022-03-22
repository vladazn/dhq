package domain

type Answer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Action struct {
	Event string  `json:"event"`
	Data  *Answer `json:"data,omitempty"`
}
