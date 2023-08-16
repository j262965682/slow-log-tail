package parsefile

import (
	"fmt"
	"regexp"
	"testing"
)

func TestTimeRule(*testing.T) {
	timeReg := regexp.MustCompile(`:\s(.*)`)
	timeRule("# Time: 2022-02-23T02:00:02.845139+08:00", timeReg)
	timeRule("# Time: 2022-02-23T20:21:44.488349", timeReg)
}

func TestUserRule(*testing.T) {
	userReg := regexp.MustCompile(`\[(.*?)\]`)
	userRule("# User@Host: root[root] @ localhost []  Id: 1895435", userReg)
	userRule("# User@Host: db_item_user[db_item_user] @  [192.168.98.237]  Id: 1793146", userReg)
}

func TestRowRule(*testing.T) {
	rowRule("# Query_time: 0.016805  Lock_time: 0.000184 Rows_sent: 1  Rows_examined: 0")
	rowRule("# Query_time: 0.012647  Lock_time: 0.000000 Rows_sent: 0  Rows_examined: 0")
}

func TestTimestampRule(t *testing.T) {
	timestampRule("SET timestamp=1645552803;")
}

func TestDbRule(t *testing.T) {
	dbRule("use db_item;")
}

func TestMaoHaoChange(t *testing.T) {
	sql := `/* Traceid: ba69879fbd8bd74a4d174fbd2c0ec505 */ SELECT

                id,name, type, business_tag, icon, icon_small, description_url, categories_json, service_category_id, status

                FROM

                zcy_item_service_protocol

                WHERE status != -3
                 AND status = 3
                 AND type = 1
                 AND business_tag = '网超';`
	fmt.Println(MaoHaoChange(sql))
}

func TestRound(t *testing.T) {
	fmt.Println(Round(9.123932, 3))
}
