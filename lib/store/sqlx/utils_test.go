package sqlx

import (
	"fmt"
	"testing"
)

func TestFormatQuery(t *testing.T) {
	query, err := formatQuery("update daily_account_num set total=10000 where total > 10000")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(query)
}
