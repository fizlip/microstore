package api

type Review struct {
	PID     string `json:"pid"`
	Title   string `json:"title"`
	Value   string `json:"value"`
	Content string `json:"content"`
	Creator string `json:"creator"`
	Date    string `json:"date"`
}
