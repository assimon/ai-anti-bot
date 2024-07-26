package adapter

type _Option struct {
	Proxy  string `json:"proxy"`
	ApiKey string `json:"apiKey"`
	Model  string `json:"model"`
}

type Option _Option

type RecognizeResult struct {
	State        int    `json:"state"`          // 1: It's an advertisement 0: Not an advertisement
	SpamScore    int    `json:"spam_score"`     // Spam score 0-100
	SpamReason   string `json:"spam_reason"`    // Reasons for determining spam ads
	SpamMockText string `json:"spam_mock_text"` // Laugh at the other person
}
