package processor

type Response struct {
	ID      string    `json:"id"`
	Choices []Choice  `json:"choices"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Object  string    `json:"object"`
	Usage   UsageInfo `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Logprobs     *string `json:"logprobs"` // 可以是 null，所以用指针
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
}

type UsageInfo struct {
	PromptTokens      int          `json:"prompt_tokens"`
	CompletionTokens  int          `json:"completion_tokens"`
	TotalTokens       int          `json:"total_tokens"`
	CompletionDetails TokenDetails `json:"completion_tokens_details"`
	PromptDetails     TokenDetails `json:"prompt_tokens_details"`
	SystemFingerprint *string      `json:"system_fingerprint"` // 可以是 null
}

type TokenDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
	CachedTokens    int `json:"cached_tokens"`
}
