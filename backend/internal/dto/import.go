package dto

// ImportError 导入行错误
type ImportError struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}

// ImportResult 导入结果
type ImportResult struct {
	Total   int           `json:"total"`
	Success int           `json:"success"`
	Failed  int           `json:"failed"`
	Errors  []ImportError `json:"errors"`
}
