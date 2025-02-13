package sql

import (
	"fmt"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/tgorm"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

// ExampleUTXOSelector_buildQueryForInputs_sqlite demonstrates what would be the query used to select inputs for a transaction.
func ExampleUTXOSelector_buildQueryForInputs_sqlite() {
	db := tgorm.GormDBForPrintingSQL(tgorm.SQLite)

	// and:
	selector := givenInputsSelector(db)

	query := db.ToSQL(func(db *gorm.DB) *gorm.DB {
		query := selector.buildQueryForInputs(db, "someuserid", 1, 10)
		query.Find(&database.UserUTXO{})
		return query
	})

	fmt.Println(query)

	// Output: SELECT `tx_id`,`vout`,`custom_instructions`,`satoshis`,`estimated_input_size` FROM `xapi_user_utxos` WHERE (tx_id, vout) in (SELECT tx_id,vout FROM (SELECT tx_id,vout,change,min(case when change >= 0 then change end) over () as min_change FROM (SELECT tx_id,vout,case when remaining_value - fee_no_change_output <= 0 then remaining_value - fee_no_change_output else remaining_value - fee_with_change_output end as change FROM (SELECT `tx_id`,`vout`,sum(satoshis) over (order by touched_at ASC, created_at ASC, tx_id ASC, vout ASC) - 1 as remaining_value,ceil((sum(estimated_input_size) over (order by touched_at ASC, created_at ASC, tx_id ASC, vout ASC) + 10) / cast(1000 as float)) * 1 as fee_no_change_output,ceil((sum(estimated_input_size) over (order by touched_at ASC, created_at ASC, tx_id ASC, vout ASC) + 10 + 34) / cast(1000 as float)) * 1 as fee_with_change_output FROM `xapi_user_utxos` WHERE user_id = "someuserid") as utxo) as utxoWithChange) as utxoWithMinChange WHERE change <= min_change)
}

// ExampleUTXOSelector_buildQueryForInputs_postgresql demonstrates what would be the query used to select inputs for a transaction.
func ExampleUTXOSelector_buildQueryForInputs_postgresql() {
	db := tgorm.GormDBForPrintingSQL(tgorm.PostgreSQL)

	// and:
	selector := givenInputsSelector(db)

	query := db.ToSQL(func(db *gorm.DB) *gorm.DB {
		query := selector.buildQueryForInputs(db, "someuserid", 1, 10)
		query.Find(&database.UserUTXO{})
		return query
	})

	fmt.Println(query)

	// Output: SELECT "tx_id","vout","custom_instructions","satoshis","estimated_input_size" FROM "xapi_user_utxos" WHERE (tx_id, vout) in (SELECT tx_id,vout FROM (SELECT tx_id,vout,change,min(case when change >= 0 then change end) over () as min_change FROM (SELECT tx_id,vout,case when remaining_value - fee_no_change_output <= 0 then remaining_value - fee_no_change_output else remaining_value - fee_with_change_output end as change FROM (SELECT "tx_id","vout",sum(satoshis) over (order by touched_at ASC, created_at ASC, tx_id ASC, vout ASC) - 1 as remaining_value,ceil((sum(estimated_input_size) over (order by touched_at ASC, created_at ASC, tx_id ASC, vout ASC) + 10) / cast(1000 as float)) * 1 as fee_no_change_output,ceil((sum(estimated_input_size) over (order by touched_at ASC, created_at ASC, tx_id ASC, vout ASC) + 10 + 34) / cast(1000 as float)) * 1 as fee_with_change_output FROM "xapi_user_utxos" WHERE user_id = 'someuserid') as utxo) as utxoWithChange) as utxoWithMinChange WHERE change <= min_change)
}

// ExampleUTXOSelector_buildUpdateTouchedAtQuery_sqlite demonstrates what would be the SQL statement used to update inputs after selecting them.
func ExampleUTXOSelector_buildUpdateTouchedAtQuery_sqlite() {
	db := tgorm.GormDBForPrintingSQL(tgorm.SQLite)

	selector := givenInputsSelector(db)

	utxos := []*database.UserUTXO{
		{UserID: "id_of_user_1", TxID: "tx_id_1", Vout: 0, Satoshis: 10, EstimatedInputSize: 148, Bucket: "bsv", CreatedAt: time.Now(), TouchedAt: time.Now()},
		{UserID: "id_of_user_1", TxID: "tx_id_1", Vout: 1, Satoshis: 10, EstimatedInputSize: 148, Bucket: "bsv", CreatedAt: time.Now(), TouchedAt: time.Now()},
		{UserID: "id_of_user_1", TxID: "tx_id_2", Vout: 0, Satoshis: 10, EstimatedInputSize: 148, Bucket: "bsv", CreatedAt: time.Now(), TouchedAt: time.Now()},
	}

	query := db.ToSQL(func(db *gorm.DB) *gorm.DB {
		query := selector.buildUpdateTouchedAtQuery(db, utxos)
		query.UpdateColumn("touched_at", time.Date(2006, 02, 01, 15, 4, 5, 7, time.UTC))
		return query
	})

	fmt.Println(query)

	// Output: UPDATE `xapi_user_utxos` SET `touched_at`="2006-02-01 15:04:05" WHERE (tx_id, vout) in (("tx_id_1",0),("tx_id_1",1),("tx_id_2",0))
}

// ExampleUTXOSelector_buildUpdateTouchedAtQuery_postgres demonstrates what would be the SQL statement used to update inputs after selecting them.
func ExampleUTXOSelector_buildUpdateTouchedAtQuery_postgres() {
	db := tgorm.GormDBForPrintingSQL(tgorm.PostgreSQL)

	selector := givenInputsSelector(db)

	utxos := []*database.UserUTXO{
		{UserID: "id_of_user_1", TxID: "tx_id_1", Vout: 0, Satoshis: 10, EstimatedInputSize: 148, Bucket: "bsv", CreatedAt: time.Now(), TouchedAt: time.Now()},
		{UserID: "id_of_user_1", TxID: "tx_id_1", Vout: 1, Satoshis: 10, EstimatedInputSize: 148, Bucket: "bsv", CreatedAt: time.Now(), TouchedAt: time.Now()},
		{UserID: "id_of_user_1", TxID: "tx_id_2", Vout: 0, Satoshis: 10, EstimatedInputSize: 148, Bucket: "bsv", CreatedAt: time.Now(), TouchedAt: time.Now()},
	}

	query := db.ToSQL(func(db *gorm.DB) *gorm.DB {
		query := selector.buildUpdateTouchedAtQuery(db, utxos)
		query.UpdateColumn("touched_at", time.Date(2006, 02, 01, 15, 4, 5, 7, time.UTC))
		return query
	})

	fmt.Println(query)

	// Output: UPDATE "xapi_user_utxos" SET "touched_at"='2006-02-01 15:04:05' WHERE (tx_id, vout) in (('tx_id_1',0),('tx_id_1',1),('tx_id_2',0))
}

func givenInputsSelector(db *gorm.DB) *UTXOSelector {
	selector := NewUTXOSelector(db, bsv.FeeUnit{Satoshis: 1, Bytes: 1000})
	return selector
}
