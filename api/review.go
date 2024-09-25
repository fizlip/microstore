package api

type Review struct {
  Title   string `json:"title"`
  Value   string `json:"value"`
  Content string `json:"content"`
  Creator string `json:"creator"`
  Date    string `json:"date"`
}
