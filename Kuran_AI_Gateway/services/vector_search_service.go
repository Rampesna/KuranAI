package services

import (
	"context"

	"github.com/pgvector/pgvector-go"
	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
	"gorm.io/gorm/clause"
)


const defaultSimilarityLimit = 10

// VectorSearchService benzer vektör arama servisi.
type VectorSearchService struct{}

// NewVectorSearchService VectorSearchService oluşturur.
func NewVectorSearchService() *VectorSearchService {
	return &VectorSearchService{}
}

// SearchMealSimilar meal_embeddings tablosunda benzer vektörleri arar.
func (s *VectorSearchService) SearchMealSimilar(ctx context.Context, embedding []float32, limit int) ([]models.MealEmbedding, error) {
	if limit <= 0 {
		limit = defaultSimilarityLimit
	}
	vec := pgvector.NewVector(embedding)
	var results []models.MealEmbedding
	err := database.DB.WithContext(ctx).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "embedding <=> ?", Vars: []interface{}{vec}},
		}).
		Limit(limit).
		Find(&results).Error
	return results, err
}

// SearchAyahExegesisSimilar ayah_exegesis_embeddings tablosunda benzer vektörleri arar.
func (s *VectorSearchService) SearchAyahExegesisSimilar(ctx context.Context, embedding []float32, limit int) ([]models.AyahExegesisEmbedding, error) {
	if limit <= 0 {
		limit = defaultSimilarityLimit
	}
	vec := pgvector.NewVector(embedding)
	var results []models.AyahExegesisEmbedding
	err := database.DB.WithContext(ctx).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "embedding <=> ?", Vars: []interface{}{vec}},
		}).
		Limit(limit).
		Find(&results).Error
	return results, err
}

// SearchBySourceType sourceType'a göre ilgili tabloda benzer vektör arar.
// meal -> meal_embeddings, tefsir -> ayah_exegesis_embeddings, hadis -> boş (nil slice).
func (s *VectorSearchService) SearchBySourceType(ctx context.Context, sourceType string, embedding []float32, limit int) ([]map[string]interface{}, error) {
	switch sourceType {
	case "meal":
		items, err := s.SearchMealSimilar(ctx, embedding, limit)
		if err != nil {
			return nil, err
		}
		out := make([]map[string]interface{}, 0, len(items))
		for _, m := range items {
			out = append(out, map[string]interface{}{
				"surah": m.Surah, "ayah": m.Ayah, "data": m.Data,
			})
		}
		return out, nil
	case "tefsir":
		items, err := s.SearchAyahExegesisSimilar(ctx, embedding, limit)
		if err != nil {
			return nil, err
		}
		out := make([]map[string]interface{}, 0, len(items))
		for _, a := range items {
			out = append(out, map[string]interface{}{
				"surah": a.Surah, "ayah": a.Ayah, "data": a.Data,
			})
		}
		return out, nil
	case "hadis":
		return []map[string]interface{}{}, nil
	default:
		return []map[string]interface{}{}, nil
	}
}
