// Copyright 2017 ZhongAn Information Technology Services Co.,Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package basesql

// GetInitSQLs get database initialize sqls
// opt sqls to create operation tables
// opi sqls to create operation table-indexs
// qt  sqls to create query tables
// qi  sqls to create query table-indexs
func (bs *Basesql) GetInitSQLs() (opt, opi, qt, qi []string) {
	opt = []string{
		createAccDataSQL,
	}
	opi = createOpIndexs

	qt = []string{
		creatActionSQL,
	}
	qi = createQIndex

	return
}

var (
	createOpIndexs = []string{
		// index for table accdata
		"CREATE INDEX IF NOT EXISTS accdateaccid ON accdata (accountid)",
	}

	createQIndex = []string{
		// indexs for table actions
		"CREATE INDEX IF NOT EXISTS actionstxhash ON actions (txhash)",
		"CREATE INDEX IF NOT EXISTS actionsfromaccount ON actions (fromaccount)",
		"CREATE INDEX IF NOT EXISTS actionstoaccount ON actions (toaccount)",
		"CREATE INDEX IF NOT EXISTS actionscreateat ON actions (createat)",
	}
)

const (
	createAccDataSQL = `CREATE TABLE IF NOT EXISTS accdata
    (
		dataid			INTEGER 		PRIMARY KEY	AUTOINCREMENT,
		accountid       VARCHAR(66)		NOT NULL,
		datakey			VARCHAR(256)	NOT NULL,
		datavalue		VARCHAR(256)	NOT NULL,
		category		VARCHAR(256)	NOT NULL
	);`

	//========================================================================//

	creatActionSQL = `CREATE TABLE IF NOT EXISTS actions
	(
		actionid			INTEGER	PRIMARY KEY	AUTOINCREMENT,
		typei				INT			NOT NULL,
		type				VARCHAR(32)	NOT NULL,
		height				INT			NOT NULL,
		txhash				VARCHAR(64)	NOT NULL,
		fromaccount			VARCHAR(66),			-- only used in payment
		toaccount			VARCHAR(66),			-- only used in payment
		createat			INT			NOT NULL,
		jdata				TEXT		NOT NULL
	);`
)
