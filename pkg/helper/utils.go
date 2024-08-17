package helper

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewIncrementId(db *pgxpool.Pool, increment_column, tableName string, prifix string, Length int) (func() string, error) {
	var (
		id    sql.NullString
		query = fmt.Sprintf("SELECT %s FROM %s ORDER BY created_at DESC LIMIT 1", increment_column, tableName)
	)
	contex, cancelF := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*2))
	resp := db.QueryRow(contex, query)

	resp.Scan(&id)
	defer cancelF()

	if !id.Valid {
		id.String = ""
	}

	return func() string {
		idNumber := idToInt(id.String)
		idNumber++
		var (
			numberLength = idNumber
			count        = 0
		)

		for numberLength > 0 {
			numberLength /= 10
			count++
		}
		if count == 0 {
			count++
		}

		prifix = Prifix(prifix)

		id := fmt.Sprintf("%s-%s%d", prifix, strings.Repeat("0", Length-count), idNumber)
		return id

	}, nil
}

func Prifix(str string) string {
	prifix_Up := strings.ToUpper(str)
	prifix := ""
	for _, v := range strings.Split(prifix_Up, " ") {
		prifix += string(v[0])
	}

	return prifix_Up
}

func idToInt(DatabaseId string) int {
	pattern := regexp.MustCompile("[0-9]+")
	firstMatchSubstring := pattern.FindString(DatabaseId)
	id, err := strconv.Atoi(firstMatchSubstring)
	if err != nil {
		return 0
	}
	return id
}

// IfElse evaluates a condition, if true returns the first parameter otherwise the second
func IfElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}
