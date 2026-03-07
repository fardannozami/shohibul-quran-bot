package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
)

type BotRepository struct {
	db *sql.DB
}

func NewBotRepository(db *sql.DB) *BotRepository {
	return &BotRepository{db: db}
}

func (r *BotRepository) InitDatabase(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			phone TEXT,
			name TEXT,
			xp INTEGER DEFAULT 0,
			level INTEGER DEFAULT 1,
			streak INTEGER DEFAULT 0,
			joined_at TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS reports (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			pages INTEGER,
			message TEXT,
			date TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS daily_progress (
			user_id TEXT,
			date TEXT,
			pages INTEGER,
			reports_count INTEGER,
			PRIMARY KEY (user_id, date)
		)`,
		`CREATE TABLE IF NOT EXISTS badges (
			user_id TEXT,
			badge TEXT,
			created_at TEXT
		)`,
	}
	for _, q := range queries {
		if _, err := r.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

// User methods
func (r *BotRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, phone, name, xp, level, streak, joined_at FROM users WHERE id = ?", id)
	var u domain.User
	var joinedAt string
	err := row.Scan(&u.ID, &u.Phone, &u.Name, &u.XP, &u.Level, &u.Streak, &joinedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.JoinedAt, _ = time.Parse(time.RFC3339, joinedAt)
	return &u, nil
}

func (r *BotRepository) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (id, phone, name, xp, level, streak, joined_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.Phone, user.Name, user.XP, user.Level, user.Streak, user.JoinedAt.Format(time.RFC3339))
	return err
}

func (r *BotRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET phone=?, name=?, xp=?, level=?, streak=? WHERE id=?",
		user.Phone, user.Name, user.XP, user.Level, user.Streak, user.ID)
	return err
}

func (r *BotRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, phone, name, xp, level, streak, joined_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*domain.User
	for rows.Next() {
		var u domain.User
		var joinedAt string
		if err := rows.Scan(&u.ID, &u.Phone, &u.Name, &u.XP, &u.Level, &u.Streak, &joinedAt); err != nil {
			return nil, err
		}
		u.JoinedAt, _ = time.Parse(time.RFC3339, joinedAt)
		users = append(users, &u)
	}
	return users, nil
}

func (r *BotRepository) ResolveLIDToPhone(ctx context.Context, lid string) string {
	var phone string
	err := r.db.QueryRowContext(ctx, "SELECT pn FROM whatsmeow_lid_map WHERE lid = ?", lid).Scan(&phone)
	if err == nil && phone != "" {
		return phone
	}
	return lid
}

// ReportLog methods
func (r *BotRepository) InsertReport(ctx context.Context, report *domain.ReportLog) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO reports (id, user_id, pages, message, date) VALUES (?, ?, ?, ?, ?)",
		report.ID, report.UserID, report.Pages, report.Message, report.Date.Format(time.RFC3339))
	return err
}

func (r *BotRepository) GetReportsByUser(ctx context.Context, userID string) ([]*domain.ReportLog, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, pages, message, date FROM reports WHERE user_id = ? ORDER BY date DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reports []*domain.ReportLog
	for rows.Next() {
		var log domain.ReportLog
		var date string
		if err := rows.Scan(&log.ID, &log.UserID, &log.Pages, &log.Message, &date); err != nil {
			return nil, err
		}
		log.Date, _ = time.Parse(time.RFC3339, date)
		reports = append(reports, &log)
	}
	return reports, nil
}

// DailyProgress methods
func (r *BotRepository) GetDailyProgress(ctx context.Context, userID string, date time.Time) (*domain.DailyProgress, error) {
	dateStr := date.Format("2006-01-02")
	row := r.db.QueryRowContext(ctx, "SELECT user_id, date, pages, reports_count FROM daily_progress WHERE user_id = ? AND date = ?", userID, dateStr)
	var dp domain.DailyProgress
	var dStr string
	err := row.Scan(&dp.UserID, &dStr, &dp.Pages, &dp.ReportsCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	dp.Date, _ = time.Parse("2006-01-02", dStr)
	return &dp, nil
}

func (r *BotRepository) UpsertDailyProgress(ctx context.Context, progress *domain.DailyProgress) error {
	dateStr := progress.Date.Format("2006-01-02")
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO daily_progress (user_id, date, pages, reports_count)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, date) DO UPDATE SET
			pages = excluded.pages,
			reports_count = excluded.reports_count
	`, progress.UserID, dateStr, progress.Pages, progress.ReportsCount)
	return err
}

func (r *BotRepository) GetTotalPagesInRange(ctx context.Context, start, end time.Time) (int, error) {
	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")
	row := r.db.QueryRowContext(ctx, "SELECT SUM(pages) FROM daily_progress WHERE date BETWEEN ? AND ?", startStr, endStr)
	var total sql.NullInt64
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return int(total.Int64), nil
}

// Badge methods
func (r *BotRepository) InsertBadge(ctx context.Context, badge *domain.BadgeLog) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO badges (user_id, badge, created_at) VALUES (?, ?, ?)",
		badge.UserID, badge.Badge, badge.CreatedAt.Format(time.RFC3339))
	return err
}

func (r *BotRepository) GetBadgesByUser(ctx context.Context, userID string) ([]*domain.BadgeLog, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id, badge, created_at FROM badges WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var badges []*domain.BadgeLog
	for rows.Next() {
		var b domain.BadgeLog
		var created string
		if err := rows.Scan(&b.UserID, &b.Badge, &created); err != nil {
			return nil, err
		}
		b.CreatedAt, _ = time.Parse(time.RFC3339, created)
		badges = append(badges, &b)
	}
	return badges, nil
}
