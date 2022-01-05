package id

import (
	"fmt"
	"testing"
)

func TestGenID(t *testing.T) {
	type testCase struct {
		name string
		id   string
	}
	tbID, _ := ShortIDWithDic(-1, "ABCDFGHJKLMNPQRSTVWXZ012456789", 26)
	testCases := []*testCase{
		&testCase{"StringUUID", StringUUID()},
		&testCase{"String", String(DefaultShortNameLen)},

		&testCase{"HexUUID_l", HexUUID(false)},
		&testCase{"HexUUID_u", HexUUID(true)},
		&testCase{"BaseUUID", BaseUUID()},
		&testCase{"ShortID", ShortID(-1)},
		&testCase{"UpperShortID", UpperShortID(-1)},
		&testCase{"WithPrefix", WithPrefix(ShortID(12), "req_")},
		&testCase{"WithPrefix", WithPrefix(ShortID(-1), "app_")},
		&testCase{"WithPrefix", WithPrefix(ShortID(-1), "t_")},
		&testCase{"ShortIDWithDic", WithPrefix(tbID, "TB_")},
	}
	for _, v := range testCases {
		fmt.Printf("%-16s: (%-2d)%s\n", v.name, len(v.id), v.id)
	}
}
