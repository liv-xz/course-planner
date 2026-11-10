package service

import (
	"context"
	"sort"
	"strings"

	"sjqu-course-planner/internal/config"
	"sjqu-course-planner/internal/model"
	"sjqu-course-planner/internal/repository"
)

// RAGService 轻量知识检索（MySQL 关键词，无需向量库）
type RAGService struct {
	repo *repository.KnowledgeChunkRepo
	cfg  config.RAG
}

type RAGHit struct {
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Source  string  `json:"source"`
	Score   float64 `json:"score"`
}

func NewRAGService(repo *repository.KnowledgeChunkRepo, cfg config.RAG) *RAGService {
	return &RAGService{repo: repo, cfg: cfg}
}

func (s *RAGService) Search(ctx context.Context, query string) ([]RAGHit, error) {
	if !s.cfg.Enabled {
		return nil, nil
	}
	chunks, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	tokens := queryTokens(query)
	if len(tokens) == 0 {
		return nil, nil
	}
	hits := make([]RAGHit, 0)
	for _, c := range chunks {
		score := scoreChunk(c, tokens)
		if score < s.cfg.MinScore {
			continue
		}
		source := ""
		if c.Source != nil {
			source = *c.Source
		}
		hits = append(hits, RAGHit{Title: c.Title, Content: c.Content, Source: source, Score: score})
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
	topK := s.cfg.TopK
	if topK <= 0 {
		topK = 4
	}
	if len(hits) > topK {
		hits = hits[:topK]
	}
	return hits, nil
}

func scoreChunk(c model.KnowledgeChunk, tokens []string) float64 {
	text := strings.ToLower(c.Title + " " + c.Content)
	if c.Keywords != nil {
		text += " " + strings.ToLower(*c.Keywords)
	}
	score := 0.0
	for _, t := range tokens {
		if t == "" {
			continue
		}
		count := strings.Count(text, t)
		if count > 0 {
			score += float64(count)
			if strings.Contains(strings.ToLower(c.Title), t) {
				score += 1.5
			}
			if c.Keywords != nil && strings.Contains(strings.ToLower(*c.Keywords), t) {
				score += 1.0
			}
		}
	}
	return score
}

func queryTokens(q string) []string {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return nil
	}
	out := []string{}
	seen := map[string]bool{}
	add := func(t string) {
		if len([]rune(t)) < 2 || seen[t] {
			return
		}
		seen[t] = true
		out = append(out, t)
	}
	// 英文/数字词
	for _, w := range strings.FieldsFunc(q, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-')
	}) {
		add(w)
	}
	// 中文二元组
	runes := []rune(q)
	for i := 0; i+1 < len(runes); i++ {
		a, b := runes[i], runes[i+1]
		if isHan(a) && isHan(b) {
			add(string(runes[i : i+2]))
		}
	}
	return out
}

func isHan(r rune) bool { return r >= 0x4e00 && r <= 0x9fff }
