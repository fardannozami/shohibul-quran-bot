package sqlite

import (
	"context"
	"database/sql"
	"strings"
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
	// 1. Create tables if they don't exist
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT,
			group_id TEXT DEFAULT 'main',
			phone TEXT,
			name TEXT,
			xp INTEGER DEFAULT 0,
			level INTEGER DEFAULT 1,
			streak INTEGER DEFAULT 0,
			joined_at TEXT,
			daily_target INTEGER DEFAULT 0,
			PRIMARY KEY (id, group_id)
		)`,
		`CREATE TABLE IF NOT EXISTS reports (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			group_id TEXT DEFAULT 'main',
			pages INTEGER,
			message TEXT,
			date TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS daily_progress (
			user_id TEXT,
			group_id TEXT DEFAULT 'main',
			date TEXT,
			pages INTEGER,
			reports_count INTEGER,
			PRIMARY KEY (user_id, date, group_id)
		)`,
		`CREATE TABLE IF NOT EXISTS badges (
			user_id TEXT,
			group_id TEXT DEFAULT 'main',
			badge TEXT,
			created_at TEXT
		)`,
	}
	for _, q := range queries {
		if _, err := r.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}

	// 2. Migration: Check if primary keys are correct by recreating tables if needed
	// This handles the case where group_id was added later but not included in PK

	// Migrate daily_progress if PK is missing group_id
	var dpSchema string
	err := r.db.QueryRowContext(ctx, "SELECT sql FROM sqlite_master WHERE type='table' AND name='daily_progress'").Scan(&dpSchema)
	if err == nil {
		dpSchema = strings.ReplaceAll(dpSchema, "\n", " ")
		dpSchema = strings.ReplaceAll(dpSchema, "\t", " ")
		if !strings.Contains(dpSchema, "PRIMARY KEY (user_id, date, group_id)") && !strings.Contains(dpSchema, "PRIMARY KEY(user_id, date, group_id)") {
			migrationQueries := []string{
				"ALTER TABLE daily_progress RENAME TO daily_progress_old",
				`CREATE TABLE daily_progress (
					user_id TEXT,
					group_id TEXT DEFAULT 'main',
					date TEXT,
					pages INTEGER,
					reports_count INTEGER,
					PRIMARY KEY (user_id, date, group_id)
				)`,
				"INSERT INTO daily_progress (user_id, group_id, date, pages, reports_count) SELECT user_id, IFNULL(group_id, 'main'), date, pages, reports_count FROM daily_progress_old",
				"DROP TABLE daily_progress_old",
			}
			for _, q := range migrationQueries {
				if _, err := r.db.ExecContext(ctx, q); err != nil {
					return err
				}
			}
		}
	}

	// Migrate users if PK is missing group_id
	var usersSchema string
	err = r.db.QueryRowContext(ctx, "SELECT sql FROM sqlite_master WHERE type='table' AND name='users'").Scan(&usersSchema)
	if err == nil {
		usersSchema = strings.ReplaceAll(usersSchema, "\n", " ")
		usersSchema = strings.ReplaceAll(usersSchema, "\t", " ")
		if !strings.Contains(usersSchema, "PRIMARY KEY (id, group_id)") && !strings.Contains(usersSchema, "PRIMARY KEY(id, group_id)") {
			migrationQueries := []string{
				"ALTER TABLE users RENAME TO users_old",
				`CREATE TABLE users (
					id TEXT,
					group_id TEXT DEFAULT 'main',
					phone TEXT,
					name TEXT,
					xp INTEGER DEFAULT 0,
					level INTEGER DEFAULT 1,
					streak INTEGER DEFAULT 0,
					joined_at TEXT,
					daily_target INTEGER DEFAULT 0,
					PRIMARY KEY (id, group_id)
				)`,
				"INSERT INTO users (id, group_id, phone, name, xp, level, streak, joined_at, daily_target) SELECT id, IFNULL(group_id, 'main'), phone, name, xp, level, streak, joined_at, IFNULL(daily_target, 0) FROM users_old",
				"DROP TABLE users_old",
			}
			for _, q := range migrationQueries {
				if _, err := r.db.ExecContext(ctx, q); err != nil {
					return err
				}
			}
		}
	}

	// Migration: Add group_id column if it doesn't exist (for tables not handled above)
	alterQueries := []string{
		"ALTER TABLE reports ADD COLUMN group_id TEXT",
		"ALTER TABLE badges ADD COLUMN group_id TEXT",
	}
	for _, q := range alterQueries {
		_, _ = r.db.ExecContext(ctx, q) // ignore error if column already exists
	}

	// Add daily_target to users if it doesn't exist (simpler check than recreating table if we just need the column)
	_, _ = r.db.ExecContext(ctx, "ALTER TABLE users ADD COLUMN daily_target INTEGER DEFAULT 0")

	return nil
}

// User methods
func (r *BotRepository) GetUser(ctx context.Context, id string, groupID string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, group_id, phone, name, xp, level, streak, joined_at, daily_target FROM users WHERE id = ? AND group_id = ?", id, groupID)
	var u domain.User
	var joinedAt string
	err := row.Scan(&u.ID, &u.GroupID, &u.Phone, &u.Name, &u.XP, &u.Level, &u.Streak, &joinedAt, &u.DailyTarget)
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
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (id, group_id, phone, name, xp, level, streak, joined_at, daily_target) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.GroupID, user.Phone, user.Name, user.XP, user.Level, user.Streak, user.JoinedAt.Format(time.RFC3339), user.DailyTarget)
	return err
}

func (r *BotRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET phone=?, name=?, xp=?, level=?, streak=?, daily_target=? WHERE id=? AND group_id=?",
		user.Phone, user.Name, user.XP, user.Level, user.Streak, user.DailyTarget, user.ID, user.GroupID)
	return err
}

func (r *BotRepository) GetAllUsers(ctx context.Context, groupID string) ([]*domain.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, group_id, phone, name, xp, level, streak, joined_at, daily_target FROM users WHERE group_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*domain.User
	for rows.Next() {
		var u domain.User
		var joinedAt string
		if err := rows.Scan(&u.ID, &u.GroupID, &u.Phone, &u.Name, &u.XP, &u.Level, &u.Streak, &joinedAt, &u.DailyTarget); err != nil {
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
	_, err := r.db.ExecContext(ctx, "INSERT INTO reports (id, user_id, group_id, pages, message, date) VALUES (?, ?, ?, ?, ?, ?)",
		report.ID, report.UserID, report.GroupID, report.Pages, report.Message, report.Date.Format(time.RFC3339))
	return err
}

func (r *BotRepository) GetReportsByUser(ctx context.Context, userID string, groupID string) ([]*domain.ReportLog, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, group_id, pages, message, date FROM reports WHERE user_id = ? AND group_id = ? ORDER BY date DESC", userID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reports []*domain.ReportLog
	for rows.Next() {
		var log domain.ReportLog
		var date string
		if err := rows.Scan(&log.ID, &log.UserID, &log.GroupID, &log.Pages, &log.Message, &date); err != nil {
			return nil, err
		}
		log.Date, _ = time.Parse(time.RFC3339, date)
		reports = append(reports, &log)
	}
	return reports, nil
}

// DailyProgress methods
func (r *BotRepository) GetDailyProgress(ctx context.Context, userID string, groupID string, date time.Time) (*domain.DailyProgress, error) {
	dateStr := date.Format("2006-01-02")
	row := r.db.QueryRowContext(ctx, "SELECT user_id, group_id, date, pages, reports_count FROM daily_progress WHERE user_id = ? AND group_id = ? AND date = ?", userID, groupID, dateStr)
	var dp domain.DailyProgress
	var dStr string
	err := row.Scan(&dp.UserID, &dp.GroupID, &dStr, &dp.Pages, &dp.ReportsCount)
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
		INSERT INTO daily_progress (user_id, group_id, date, pages, reports_count)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, date, group_id) DO UPDATE SET
			pages = excluded.pages,
			reports_count = excluded.reports_count
	`, progress.UserID, progress.GroupID, dateStr, progress.Pages, progress.ReportsCount)
	return err
}

func (r *BotRepository) GetTotalPagesInRange(ctx context.Context, start, end time.Time, groupID string) (int, error) {
	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")
	row := r.db.QueryRowContext(ctx, "SELECT SUM(pages) FROM daily_progress WHERE group_id = ? AND date BETWEEN ? AND ?", groupID, startStr, endStr)
	var total sql.NullInt64
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return int(total.Int64), nil
}

// Badge methods
func (r *BotRepository) InsertBadge(ctx context.Context, badge *domain.BadgeLog) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO badges (user_id, group_id, badge, created_at) VALUES (?, ?, ?, ?)",
		badge.UserID, badge.GroupID, badge.Badge, badge.CreatedAt.Format(time.RFC3339))
	return err
}

func (r *BotRepository) GetBadgesByUser(ctx context.Context, userID string, groupID string) ([]*domain.BadgeLog, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id, group_id, badge, created_at FROM badges WHERE user_id = ? AND group_id = ?", userID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var badges []*domain.BadgeLog
	for rows.Next() {
		var b domain.BadgeLog
		var created string
		if err := rows.Scan(&b.UserID, &b.GroupID, &b.Badge, &created); err != nil {
			return nil, err
		}
		b.CreatedAt, _ = time.Parse(time.RFC3339, created)
		badges = append(badges, &b)
	}
	return badges, nil
}

