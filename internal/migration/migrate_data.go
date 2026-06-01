package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// MigrateData migrates data from a Spring Boot SQLite database to a new Go SQLite database.
func MigrateData(sourcePath, targetPath string) error {
	// Open source database (Spring Boot format)
	sourceDB, err := sql.Open("sqlite", sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source database: %w", err)
	}
	defer sourceDB.Close()

	if err := sourceDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping source database: %w", err)
	}

	// Open target database (Go format - must have migrations already applied)
	targetDB, err := sql.Open("sqlite", targetPath+"?_pragma=foreign_keys(ON)")
	if err != nil {
		return fmt.Errorf("failed to open target database: %w", err)
	}
	defer targetDB.Close()

	if err := targetDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping target database: %w", err)
	}

	// Migrate entity by entity
	migrations := []entityMigration{
		{source: "basic_data", target: "basic_data", orderCol: "", orderDesc: false, columns: basicDataColumns},
		{source: "home", target: "home", orderCol: "", orderDesc: false, columns: homeColumns},
		{source: "label", target: "label", orderCol: "id", orderDesc: false, columns: labelColumns},
		{source: "image_url", target: "image_url", orderCol: "", orderDesc: false, columns: genericColumns},
		{source: "video_url", target: "video_url", orderCol: "", orderDesc: false, columns: genericColumns},
		{source: "blog", target: "blog", orderCol: "", orderDesc: false, columns: blogColumns},
		{source: "blog_type", target: "blog_type", orderCol: "id", orderDesc: false, columns: blogTypeColumns},
		{source: "skill_type", target: "skill_type", orderCol: "id", orderDesc: false, columns: skillTypeColumns},
		{source: "skill", target: "skill", orderCol: "id", orderDesc: false, columns: skillColumns},
		{source: "skill_son", target: "skill_son", orderCol: "id", orderDesc: false, columns: skillSonColumns},
		{source: "experience", target: "experience", orderCol: "year_start", orderDesc: true, columns: experienceColumns},
		{source: "education", target: "education", orderCol: "id", orderDesc: false, columns: educationColumns},
		{source: "futured_project", target: "futured_project", orderCol: "id", orderDesc: false, columns: futuredProjectColumns},
	}

	for _, m := range migrations {
		if err := migrateTable(sourceDB, targetDB, m); err != nil {
			return fmt.Errorf("failed to migrate %s: %w", m.target, err)
		}
	}

	// Migrate join tables
	joinMigrations := []struct{ source, target string }{
		{"home_label", "home_label"},
		{"skill_type_skill", "skill_type_skill"},
		{"skill_skill_son", "skill_skill_son"},
		{"experience_skill_son", "experience_skill_son"},
	}

	for _, jm := range joinMigrations {
		if err := migrateJoinTable(sourceDB, targetDB, jm.source, jm.target); err != nil {
			return fmt.Errorf("failed to migrate join table %s: %w", jm.target, err)
		}
	}

	// Migrate user_credentials
	if err := migrateUserCredentials(sourceDB, targetDB); err != nil {
		return fmt.Errorf("failed to migrate user credentials: %w", err)
	}

	// Validate counts
	if err := validateCounts(sourceDB, targetDB, migrations); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	log.Println("Migration completed successfully")
	return nil
}

type entityMigration struct {
	source    string
	target    string
	orderCol  string
	orderDesc bool
	columns   []string
}

func migrateTable(source, target *sql.DB, m entityMigration) error {
	// Check if source table exists
	if !tableExists(source, m.source) {
		log.Printf("Source table %s does not exist, skipping", m.source)
		return nil
	}

	query := fmt.Sprintf("SELECT %s FROM %s", buildColumnList(m.columns), m.source)
	rows, err := source.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query source %s: %w", m.source, err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	order := 1
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row from %s: %w", m.source, err)
		}

		// Build insert query - skip the order column if it's auto-assigned
		var insertCols, placeholders string
		var insertArgs []interface{}

		if m.orderCol != "" {
			// Add order column
			insertCols = buildColumnList(m.columns) + ", \"order\""
			placeholders = buildPlaceholderList(len(m.columns)) + ", ?"
			insertArgs = append(convertArgs(values), order)
			order++
		} else {
			insertCols = buildColumnList(m.columns)
			placeholders = buildPlaceholderList(len(m.columns))
			insertArgs = convertArgs(values)
		}

		insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			m.target, insertCols, placeholders)

		if _, err := target.Exec(insertQuery, insertArgs...); err != nil {
			return fmt.Errorf("failed to insert into %s: %w", m.target, err)
		}
	}

	return rows.Err()
}

func migrateJoinTable(source, target *sql.DB, sourceTable, targetTable string) error {
	if !tableExists(source, sourceTable) {
		log.Printf("Source join table %s does not exist, skipping", sourceTable)
		return nil
	}

	rows, err := source.Query(fmt.Sprintf("SELECT * FROM %s", sourceTable))
	if err != nil {
		return fmt.Errorf("failed to query %s: %w", sourceTable, err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row from %s: %w", sourceTable, err)
		}
		insertQuery := fmt.Sprintf("INSERT INTO %s VALUES (%s)",
			targetTable, buildPlaceholderList(len(cols)))
		if _, err := target.Exec(insertQuery, convertArgs(values)...); err != nil {
			return fmt.Errorf("failed to insert into %s: %w", targetTable, err)
		}
	}

	return rows.Err()
}

func migrateUserCredentials(source, target *sql.DB) error {
	if !tableExists(source, "user_credentials") {
		log.Println("Source user_credentials does not exist, skipping")
		return nil
	}

	var username, passwordHash string
	err := source.QueryRow("SELECT username, password_hash FROM user_credentials LIMIT 1").
		Scan(&username, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user credentials found, skipping")
			return nil
		}
		return fmt.Errorf("failed to read user_credentials: %w", err)
	}

	_, err = target.Exec(
		"INSERT INTO user_credentials (username, password_hash, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))",
		username, passwordHash)
	return err
}

func validateCounts(source, target *sql.DB, migrations []entityMigration) error {
	for _, m := range migrations {
		if !tableExists(source, m.source) {
			continue
		}
		var srcCount, tgtCount int
		source.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", m.source)).Scan(&srcCount)
		target.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", m.target)).Scan(&tgtCount)
		if srcCount != tgtCount {
			return fmt.Errorf("count mismatch for %s: source=%d target=%d", m.target, srcCount, tgtCount)
		}
		log.Printf("Table %s: source=%d target=%d ✓", m.target, srcCount, tgtCount)
	}
	return nil
}

func tableExists(db *sql.DB, table string) bool {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?",
		table).Scan(&count)
	return err == nil && count > 0
}

func buildColumnList(cols []string) string {
	result := ""
	for i, col := range cols {
		if i > 0 {
			result += ", "
		}
		result += col
	}
	return result
}

func buildPlaceholderList(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			result += ", "
		}
		result += "?"
	}
	return result
}

func convertArgs(values []interface{}) []interface{} {
	args := make([]interface{}, len(values))
	for i, v := range values {
		if v == nil {
			args[i] = nil
		} else {
			args[i] = v
		}
	}
	return args
}

// Column definitions for each entity
var basicDataColumns = []string{
	"id", "first_name", "others_name", "first_surname", "others_surname",
	"date_birth", "located", "located_eng", "start_working_date",
	"greeting", "greeting_eng", "email", "instagram", "linkedin",
	"x", "github", "description", "description_eng",
	"description_pdf", "description_pdf_eng", "wrapper", "wrapper_eng",
}

var homeColumns = []string{
	"id", "greeting", "greeting_eng", "image_url_id",
	"button_work_label", "button_work_label_eng",
	"button_contact_label", "button_contact_label_eng",
}

var labelColumns = []string{"id", "name", "name_eng"}
var genericColumns = []string{"id", "name", "name_eng", "url"}
var blogColumns = []string{
	"id", "title", "title_eng", "clean_url_title",
	"description_short", "description", "description_short_eng", "description_eng",
	"image_url_id", "video_url_id", "blog_type_id",
}
var blogTypeColumns = []string{"id", "name", "name_eng"}
var skillTypeColumns = []string{"id", "name", "name_eng"}
var skillColumns = []string{"id", "name", "name_eng"}
var skillSonColumns = []string{"id", "name", "name_eng"}
var experienceColumns = []string{
	"id", "year_start", "year_end", "company", "location", "location_eng",
	"position", "position_eng", "summary", "summary_eng",
	"summary_pdf", "summary_pdf_eng",
	"description_items_pdf", "description_items_pdf_eng",
}
var educationColumns = []string{
	"id", "institution", "area", "area_eng", "degree", "degree_eng",
	"start_date", "end_date", "location", "location_eng",
	"highlights", "highlights_eng",
}
var futuredProjectColumns = []string{
	"id", "name", "name_eng", "description_short", "description",
	"description_short_eng", "description_eng",
	"experience_id", "image_list_url_id", "image_url_id",
}

// Run main migration command
func Run() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: migrate --source <old.db> --target <new.db>")
	}

	sourcePath := ""
	targetPath := ""

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--source":
			if i+1 < len(os.Args) {
				sourcePath = os.Args[i+1]
				i++
			}
		case "--target":
			if i+1 < len(os.Args) {
				targetPath = os.Args[i+1]
				i++
			}
		}
	}

	if sourcePath == "" || targetPath == "" {
		log.Fatal("Both --source and --target are required")
	}

	if err := MigrateData(sourcePath, targetPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
