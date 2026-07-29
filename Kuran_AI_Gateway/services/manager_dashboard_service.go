package services

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/rampesna/kuran-ai-gateway/database"
	kuranredis "github.com/rampesna/kuran-ai-gateway/redis"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const presenceWindow = 2 * time.Minute

type MetricBreakdown struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

type DashboardTrendPoint struct {
	Day           time.Time `json:"day"`
	NewUsers      int64     `json:"newUsers"`
	Questions     int64     `json:"questions"`
	Subscriptions int64     `json:"subscriptions"`
}

type ManagerUserRow struct {
	ID                 uint       `json:"id"`
	Name               string     `json:"name"`
	Email              string     `json:"email"`
	Username           string     `json:"username"`
	Status             int        `json:"status"`
	AvatarURL          string     `json:"avatarUrl"`
	Platform           string     `json:"platform"`
	DeviceInfo         string     `json:"deviceInfo"`
	LastSeenAt         *time.Time `json:"lastSeenAt"`
	CreatedAt          time.Time  `json:"createdAt"`
	SessionCount       int64      `json:"sessionCount"`
	QuestionCount      int64      `json:"questionCount"`
	PackageName        string     `json:"packageName"`
	SubscriptionStatus string     `json:"subscriptionStatus"`
	Provider           string     `json:"provider"`
	Online             bool       `json:"online"`
}

type ManagerPurchaseRow struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userId"`
	UserName    string    `json:"userName"`
	UserEmail   string    `json:"userEmail"`
	PackageName string    `json:"packageName"`
	PackageSlug string    `json:"packageSlug"`
	Provider    string    `json:"provider"`
	Billing     string    `json:"billing"`
	AmountTL    float64   `json:"amountTL"`
	StartedAt   time.Time `json:"startedAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	Status      string    `json:"status"`
	ProductID   string    `json:"productId"`
}

type ManagerQuestionRow struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"userId"`
	UserName      string    `json:"userName"`
	UserEmail     string    `json:"userEmail"`
	SessionID     uint      `json:"sessionId"`
	Content       string    `json:"content"`
	SourceType    string    `json:"sourceType"`
	CreatedAt     time.Time `json:"createdAt"`
	AssistantText string    `json:"assistantText"`
}

type ManagerConversationRow struct {
	ID           uint      `json:"id"`
	FirstMessage string    `json:"firstMessage"`
	SourceType   string    `json:"sourceType"`
	MessageCount int64     `json:"messageCount"`
	LastMessage  string    `json:"lastMessage"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ManagerMessageRow struct {
	ID         uint      `json:"id"`
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	SourceType string    `json:"sourceType"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ManagerDeviceRow struct {
	ID         uint      `json:"id"`
	Platform   string    `json:"platform"`
	DeviceInfo string    `json:"deviceInfo"`
	UserAgent  string    `json:"userAgent"`
	IPAddress  string    `json:"ipAddress"`
	CreatedAt  time.Time `json:"createdAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

type ManagerDashboardService struct{}

func NewManagerDashboardService() *ManagerDashboardService {
	return &ManagerDashboardService{}
}

func (s *ManagerDashboardService) Overview(ctx context.Context, days int) (map[string]interface{}, error) {
	if days < 7 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	var totalUsers, questionsToday, activeSubscriptions int64
	if err := database.DB.WithContext(ctx).Table("users").Where("deleted_at IS NULL").Count(&totalUsers).Error; err != nil {
		return nil, err
	}
	if err := database.DB.WithContext(ctx).Table("chat_messages").
		Where("deleted_at IS NULL AND role = ? AND created_at >= CURRENT_DATE", "user").
		Count(&questionsToday).Error; err != nil {
		return nil, err
	}
	if err := database.DB.WithContext(ctx).Table("user_subscriptions").
		Where("deleted_at IS NULL AND expires_at > NOW()").
		Distinct("user_id").Count(&activeSubscriptions).Error; err != nil {
		return nil, err
	}

	online, err := s.onlineUsers(ctx)
	if err != nil {
		return nil, err
	}

	platforms, err := s.platformBreakdown(ctx)
	if err != nil {
		return nil, err
	}
	providers, err := s.providerBreakdown(ctx)
	if err != nil {
		return nil, err
	}

	var trend []DashboardTrendPoint
	trendSQL := `
		WITH days AS (
			SELECT generate_series(CURRENT_DATE - (?::int - 1) * INTERVAL '1 day', CURRENT_DATE, INTERVAL '1 day')::date AS day
		)
		SELECT
			d.day,
			(SELECT COUNT(*) FROM users u WHERE u.deleted_at IS NULL AND u.created_at >= d.day AND u.created_at < d.day + INTERVAL '1 day') AS new_users,
			(SELECT COUNT(*) FROM chat_messages m WHERE m.deleted_at IS NULL AND m.role = 'user' AND m.created_at >= d.day AND m.created_at < d.day + INTERVAL '1 day') AS questions,
			(SELECT COUNT(*) FROM user_subscriptions s WHERE s.deleted_at IS NULL AND s.created_at >= d.day AND s.created_at < d.day + INTERVAL '1 day') AS subscriptions
		FROM days d ORDER BY d.day`
	if err := database.DB.WithContext(ctx).Raw(trendSQL, days).Scan(&trend).Error; err != nil {
		return nil, err
	}

	var recentUsers []ManagerUserRow
	if err := s.userBaseQuery(ctx).
		Order("u.created_at DESC").Limit(6).Scan(&recentUsers).Error; err != nil {
		return nil, err
	}
	applyPresence(recentUsers, online)

	var recentPurchases []ManagerPurchaseRow
	if err := s.purchaseBaseQuery(ctx).Order("us.created_at DESC").Limit(6).Scan(&recentPurchases).Error; err != nil {
		return nil, err
	}

	var grossRevenueTL float64
	if err := database.DB.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(CASE WHEN us.is_monthly THEN sp.monthly_price_tl ELSE sp.yearly_price_tl END), 0)
		FROM user_subscriptions us
		JOIN subscription_packages sp ON sp.id = us.package_id AND sp.deleted_at IS NULL
		WHERE us.deleted_at IS NULL`).Scan(&grossRevenueTL).Error; err != nil {
		return nil, err
	}

	onlineRows, err := s.onlineUserRows(ctx, online)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"metrics": map[string]interface{}{
			"totalUsers":          totalUsers,
			"onlineUsers":         len(online),
			"activeSubscriptions": activeSubscriptions,
			"questionsToday":      questionsToday,
			"grossRevenueTL":      grossRevenueTL,
		},
		"platforms":             platforms,
		"providers":             providers,
		"trend":                 trend,
		"recentUsers":           recentUsers,
		"recentPurchases":       recentPurchases,
		"onlineUsers":           onlineRows,
		"presenceWindowSeconds": int(presenceWindow.Seconds()),
	}, nil
}

func (s *ManagerDashboardService) ListUsers(ctx context.Context, page, pageSize int, search, platform, status string) (map[string]interface{}, error) {
	page, pageSize = normalizePage(page, pageSize)
	query := s.userBaseQuery(ctx)
	countQuery := database.DB.WithContext(ctx).Table("users u").Where("u.deleted_at IS NULL")

	if search = strings.TrimSpace(search); search != "" {
		needle := "%" + strings.ToLower(search) + "%"
		condition := "(LOWER(u.name) LIKE ? OR LOWER(u.email) LIKE ? OR LOWER(u.username) LIKE ?)"
		query = query.Where(condition, needle, needle, needle)
		countQuery = countQuery.Where(condition, needle, needle, needle)
	}
	rawPlatform := strings.ToLower(strings.TrimSpace(platform))
	if rawPlatform != "" && rawPlatform != "all" {
		platform = normalizeDeviceType(rawPlatform)
		condition := "LOWER(COALESCE(NULLIF(jt.device_type, ''), 'unknown')) = ?"
		query = query.Where(condition, platform)
		countQuery = countQuery.Joins(`
			LEFT JOIN LATERAL (
				SELECT device_type FROM jwt_tokens
				WHERE tokenable_type = 'user' AND tokenable_id = u.id AND deleted_at IS NULL
				ORDER BY created_at DESC LIMIT 1
			) jt ON TRUE`).Where(condition, platform)
	}
	if status == "active" || status == "inactive" {
		value := 1
		if status == "inactive" {
			value = 0
		}
		query = query.Where("u.status = ?", value)
		countQuery = countQuery.Where("u.status = ?", value)
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}
	var rows []ManagerUserRow
	if err := query.Order("u.created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Scan(&rows).Error; err != nil {
		return nil, err
	}
	online, err := s.onlineUsers(ctx)
	if err != nil {
		return nil, err
	}
	applyPresence(rows, online)
	return pageResult(rows, page, pageSize, total), nil
}

func (s *ManagerDashboardService) UserDetail(ctx context.Context, userID uint) (map[string]interface{}, error) {
	var user ManagerUserRow
	if err := s.userBaseQuery(ctx).Where("u.id = ?", userID).Limit(1).Scan(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	online, err := s.onlineUsers(ctx)
	if err != nil {
		return nil, err
	}
	rows := []ManagerUserRow{user}
	applyPresence(rows, online)
	user = rows[0]

	var devices []ManagerDeviceRow
	if err := database.DB.WithContext(ctx).Raw(`
		SELECT id,
			LOWER(COALESCE(NULLIF(device_type, ''), 'unknown')) AS platform,
			device_info, user_agent, ip_address, created_at, expires_at
		FROM jwt_tokens
		WHERE tokenable_type = 'user' AND tokenable_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 25`, userID).Scan(&devices).Error; err != nil {
		return nil, err
	}
	for i := range devices {
		devices[i].Platform = normalizeDeviceType(devices[i].Platform)
	}

	var subscriptions []ManagerPurchaseRow
	if err := s.purchaseBaseQuery(ctx).Where("us.user_id = ?", userID).
		Order("us.created_at DESC").Scan(&subscriptions).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user":          user,
		"devices":       devices,
		"subscriptions": subscriptions,
	}, nil
}

func (s *ManagerDashboardService) Conversations(ctx context.Context, userID uint, page, pageSize int) (map[string]interface{}, error) {
	page, pageSize = normalizePage(page, pageSize)
	base := database.DB.WithContext(ctx).Table("chat_sessions cs").
		Where("cs.deleted_at IS NULL AND cs.auth_type = 'user' AND cs.auth_id = ?", userID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	var rows []ManagerConversationRow
	err := base.Select(`
			cs.id, cs.first_message, cs.created_at, cs.updated_at,
			COALESCE(stats.message_count, 0) AS message_count,
			COALESCE(stats.source_type, '') AS source_type,
			COALESCE(stats.last_message, '') AS last_message`).
		Joins(`
			LEFT JOIN LATERAL (
				SELECT COUNT(*) AS message_count,
					(ARRAY_AGG(source_type ORDER BY created_at ASC))[1] AS source_type,
					(ARRAY_AGG(content ORDER BY created_at DESC))[1] AS last_message
				FROM chat_messages
				WHERE chat_session_id = cs.id AND deleted_at IS NULL
			) stats ON TRUE`).
		Order("cs.updated_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return pageResult(rows, page, pageSize, total), nil
}

func (s *ManagerDashboardService) ConversationMessages(ctx context.Context, sessionID uint) (map[string]interface{}, error) {
	var session struct {
		ID           uint      `json:"id"`
		UserID       uint      `json:"userId"`
		UserName     string    `json:"userName"`
		UserEmail    string    `json:"userEmail"`
		FirstMessage string    `json:"firstMessage"`
		CreatedAt    time.Time `json:"createdAt"`
	}
	if err := database.DB.WithContext(ctx).Raw(`
		SELECT cs.id, cs.auth_id AS user_id, u.name AS user_name, u.email AS user_email,
			cs.first_message, cs.created_at
		FROM chat_sessions cs
		LEFT JOIN users u ON cs.auth_type = 'user' AND u.id = cs.auth_id AND u.deleted_at IS NULL
		WHERE cs.id = ? AND cs.deleted_at IS NULL`, sessionID).Scan(&session).Error; err != nil {
		return nil, err
	}
	if session.ID == 0 {
		return nil, fmt.Errorf("conversation not found")
	}
	var messages []ManagerMessageRow
	if err := database.DB.WithContext(ctx).Raw(`
		SELECT id, role, content, source_type, created_at
		FROM chat_messages WHERE chat_session_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC`, sessionID).Scan(&messages).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{"session": session, "messages": messages}, nil
}

func (s *ManagerDashboardService) Purchases(ctx context.Context, page, pageSize int, search, provider, status string) (map[string]interface{}, error) {
	page, pageSize = normalizePage(page, pageSize)
	query := s.purchaseBaseQuery(ctx)
	countQuery := database.DB.WithContext(ctx).Table("user_subscriptions us").
		Joins("JOIN users u ON u.id = us.user_id AND u.deleted_at IS NULL").
		Where("us.deleted_at IS NULL")
	if search = strings.TrimSpace(search); search != "" {
		needle := "%" + strings.ToLower(search) + "%"
		condition := "(LOWER(u.name) LIKE ? OR LOWER(u.email) LIKE ?)"
		query = query.Where(condition, needle, needle)
		countQuery = countQuery.Where(condition, needle, needle)
	}
	provider = strings.ToLower(strings.TrimSpace(provider))
	if provider == "apple" || provider == "google" {
		query = query.Where("LOWER(us.provider) = ?", provider)
		countQuery = countQuery.Where("LOWER(us.provider) = ?", provider)
	}
	if status == "active" {
		query = query.Where("us.expires_at > NOW()")
		countQuery = countQuery.Where("us.expires_at > NOW()")
	} else if status == "expired" {
		query = query.Where("us.expires_at <= NOW()")
		countQuery = countQuery.Where("us.expires_at <= NOW()")
	}
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}
	var rows []ManagerPurchaseRow
	if err := query.Order("us.created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return pageResult(rows, page, pageSize, total), nil
}

func (s *ManagerDashboardService) Questions(ctx context.Context, page, pageSize int, search, sourceType string) (map[string]interface{}, error) {
	page, pageSize = normalizePage(page, pageSize)
	base := database.DB.WithContext(ctx).Table("chat_messages m").
		Joins("JOIN chat_sessions cs ON cs.id = m.chat_session_id AND cs.deleted_at IS NULL AND cs.auth_type = 'user'").
		Joins("JOIN users u ON u.id = cs.auth_id AND u.deleted_at IS NULL").
		Where("m.deleted_at IS NULL AND m.role = 'user'")
	if search = strings.TrimSpace(search); search != "" {
		needle := "%" + strings.ToLower(search) + "%"
		base = base.Where("(LOWER(m.content) LIKE ? OR LOWER(u.name) LIKE ? OR LOWER(u.email) LIKE ?)", needle, needle, needle)
	}
	if sourceType = strings.ToLower(strings.TrimSpace(sourceType)); sourceType != "" && sourceType != "all" {
		base = base.Where("LOWER(m.source_type) = ?", sourceType)
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	var rows []ManagerQuestionRow
	err := base.Select(`
			m.id, cs.auth_id AS user_id, u.name AS user_name, u.email AS user_email,
			m.chat_session_id AS session_id, m.content, m.source_type, m.created_at,
			COALESCE(answer.content, '') AS assistant_text`).
		Joins(`
			LEFT JOIN LATERAL (
				SELECT content FROM chat_messages reply
				WHERE reply.chat_session_id = m.chat_session_id
					AND reply.role = 'assistant'
					AND reply.deleted_at IS NULL
					AND reply.created_at >= m.created_at
				ORDER BY reply.created_at ASC LIMIT 1
			) answer ON TRUE`).
		Order("m.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return pageResult(rows, page, pageSize, total), nil
}

func (s *ManagerDashboardService) userBaseQuery(ctx context.Context) *gorm.DB {
	return database.DB.WithContext(ctx).Table("users u").
		Select(`
			u.id, u.name, u.email, u.username, u.status, u.avatar_url, u.created_at,
			LOWER(COALESCE(NULLIF(jt.device_type, ''), 'unknown')) AS platform,
			COALESCE(jt.device_info, '') AS device_info,
			jt.created_at AS last_seen_at,
			COALESCE(stats.session_count, 0) AS session_count,
			COALESCE(stats.question_count, 0) AS question_count,
			COALESCE(sub.package_name, 'Ücretsiz') AS package_name,
			COALESCE(sub.provider, '') AS provider,
			CASE WHEN sub.expires_at > NOW() THEN 'active' WHEN sub.expires_at IS NOT NULL THEN 'expired' ELSE 'free' END AS subscription_status`).
		Joins(`
			LEFT JOIN LATERAL (
				SELECT device_type, device_info, created_at
				FROM jwt_tokens WHERE tokenable_type = 'user' AND tokenable_id = u.id AND deleted_at IS NULL
				ORDER BY created_at DESC LIMIT 1
			) jt ON TRUE`).
		Joins(`
			LEFT JOIN LATERAL (
				SELECT COUNT(DISTINCT cs.id) AS session_count,
					COUNT(cm.id) FILTER (WHERE cm.role = 'user') AS question_count
				FROM chat_sessions cs
				LEFT JOIN chat_messages cm ON cm.chat_session_id = cs.id AND cm.deleted_at IS NULL
				WHERE cs.auth_type = 'user' AND cs.auth_id = u.id AND cs.deleted_at IS NULL
			) stats ON TRUE`).
		Joins(`
			LEFT JOIN LATERAL (
				SELECT sp.name AS package_name, us.provider, us.expires_at
				FROM user_subscriptions us
				JOIN subscription_packages sp ON sp.id = us.package_id AND sp.deleted_at IS NULL
				WHERE us.user_id = u.id AND us.deleted_at IS NULL
				ORDER BY us.expires_at DESC LIMIT 1
			) sub ON TRUE`).
		Where("u.deleted_at IS NULL")
}

func (s *ManagerDashboardService) purchaseBaseQuery(ctx context.Context) *gorm.DB {
	return database.DB.WithContext(ctx).Table("user_subscriptions us").
		Select(`
			us.id, us.user_id, u.name AS user_name, u.email AS user_email,
			sp.name AS package_name, sp.slug AS package_slug,
			LOWER(COALESCE(NULLIF(us.provider, ''), 'unknown')) AS provider,
			CASE WHEN us.is_monthly THEN 'monthly' ELSE 'yearly' END AS billing,
			CASE WHEN us.is_monthly THEN sp.monthly_price_tl ELSE sp.yearly_price_tl END AS amount_tl,
			us.started_at, us.expires_at,
			CASE WHEN us.expires_at > NOW() THEN 'active' ELSE 'expired' END AS status,
			CASE WHEN LOWER(us.provider) = 'google' THEN us.google_product_id ELSE us.apple_product_id END AS product_id`).
		Joins("JOIN users u ON u.id = us.user_id AND u.deleted_at IS NULL").
		Joins("JOIN subscription_packages sp ON sp.id = us.package_id AND sp.deleted_at IS NULL").
		Where("us.deleted_at IS NULL")
}

type presenceEntry struct {
	Platform string
	LastSeen time.Time
}

func (s *ManagerDashboardService) onlineUsers(ctx context.Context) (map[uint]presenceEntry, error) {
	now := time.Now()
	values, err := kuranredis.Client.ZRangeByScoreWithScores(ctx, kuranredis.KeyUserPresence, &goredis.ZRangeBy{
		Min: strconv.FormatInt(now.Add(-presenceWindow).Unix(), 10),
		Max: "+inf",
	}).Result()
	if err != nil && err != goredis.Nil {
		return nil, err
	}
	result := make(map[uint]presenceEntry)
	for _, value := range values {
		member, ok := value.Member.(string)
		if !ok {
			continue
		}
		parts := strings.Split(member, "|")
		if len(parts) < 2 {
			continue
		}
		id, parseErr := strconv.ParseUint(parts[0], 10, 64)
		if parseErr != nil {
			continue
		}
		entry := presenceEntry{Platform: normalizeDeviceType(parts[1]), LastSeen: time.Unix(int64(value.Score), 0)}
		existing, exists := result[uint(id)]
		if !exists || entry.LastSeen.After(existing.LastSeen) {
			result[uint(id)] = entry
		}
	}
	return result, nil
}

func (s *ManagerDashboardService) onlineUserRows(ctx context.Context, online map[uint]presenceEntry) ([]ManagerUserRow, error) {
	if len(online) == 0 {
		return []ManagerUserRow{}, nil
	}
	ids := make([]uint, 0, len(online))
	for id := range online {
		ids = append(ids, id)
	}
	var rows []ManagerUserRow
	if err := s.userBaseQuery(ctx).Where("u.id IN ?", ids).Scan(&rows).Error; err != nil {
		return nil, err
	}
	applyPresence(rows, online)
	return rows, nil
}

func (s *ManagerDashboardService) platformBreakdown(ctx context.Context) ([]MetricBreakdown, error) {
	var rows []MetricBreakdown
	err := database.DB.WithContext(ctx).Raw(`
		SELECT LOWER(COALESCE(NULLIF(device_type, ''), 'unknown')) AS key, COUNT(*) AS count
		FROM (
			SELECT DISTINCT ON (tokenable_id) tokenable_id, device_type
			FROM jwt_tokens
			WHERE tokenable_type = 'user' AND deleted_at IS NULL
			ORDER BY tokenable_id, created_at DESC
		) latest GROUP BY key`).Scan(&rows).Error
	return completeBreakdown(rows, []string{"ios", "android", "web", "unknown"}), err
}

func (s *ManagerDashboardService) providerBreakdown(ctx context.Context) ([]MetricBreakdown, error) {
	var rows []MetricBreakdown
	err := database.DB.WithContext(ctx).Raw(`
		SELECT LOWER(COALESCE(NULLIF(provider, ''), 'unknown')) AS key, COUNT(*) AS count
		FROM user_subscriptions WHERE deleted_at IS NULL GROUP BY key`).Scan(&rows).Error
	return completeBreakdown(rows, []string{"apple", "google", "unknown"}), err
}

func applyPresence(rows []ManagerUserRow, online map[uint]presenceEntry) {
	for i := range rows {
		rows[i].Platform = normalizeDeviceType(rows[i].Platform)
		if entry, ok := online[rows[i].ID]; ok {
			rows[i].Online = true
			rows[i].Platform = entry.Platform
			lastSeen := entry.LastSeen
			rows[i].LastSeenAt = &lastSeen
		}
	}
}

func completeBreakdown(rows []MetricBreakdown, keys []string) []MetricBreakdown {
	counts := make(map[string]int64, len(rows))
	for _, row := range rows {
		counts[normalizeDeviceTypeOrProvider(row.Key)] += row.Count
	}
	result := make([]MetricBreakdown, 0, len(keys))
	for _, key := range keys {
		result = append(result, MetricBreakdown{Key: key, Count: counts[key]})
	}
	return result
}

func normalizeDeviceTypeOrProvider(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "apple" || value == "google" {
		return value
	}
	return normalizeDeviceType(value)
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func pageResult[T any](items []T, page, pageSize int, total int64) map[string]interface{} {
	return map[string]interface{}{
		"items":      items,
		"page":       page,
		"pageSize":   pageSize,
		"total":      total,
		"totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
	}
}
