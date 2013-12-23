package mssql

import (
    "testing"
    "time"
)

func TestSelect(t *testing.T) {
    conn := open(t)
    defer conn.Close()

    type testStruct struct {
        sql string
        val interface{}
    }

    values := []testStruct{
        {"1", int32(1)},
        {"cast(1 as tinyint)", uint8(1)},
        {"cast(1 as smallint)", int16(1)},
        {"cast(1 as bigint)", int64(1)},
        {"cast(1 as bit)", true},
        {"cast(0 as bit)", false},
        {"'abc'", string("abc")},
        {"cast(0.5 as float)", float64(0.5)},
        {"cast(0.5 as real)", float32(0.5)},
        {"cast(1 as decimal)", Decimal{[...]uint32{1, 0, 0, 0}, true, 18, 0}},
        {"cast(0.5 as decimal(18,1))", Decimal{[...]uint32{5, 0, 0, 0}, true, 18, 1}},
        {"cast(-0.5 as decimal(18,1))", Decimal{[...]uint32{5, 0, 0, 0}, false, 18, 1}},
        {"cast(-0.5 as numeric(18,1))", Decimal{[...]uint32{5, 0, 0, 0}, false, 18, 1}},
        {"N'abc'", string("abc")},
        {"NULL", nil},
        {"cast('2000-01-01' as datetime)", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
        {"cast('2000-01-01T12:13:14.12' as datetime)",
         time.Date(2000, 1, 1, 12, 13, 14, 120000000, time.UTC)},
        {"cast(NULL as datetime)", nil},
        {"cast('2000-01-01T12:13:00' as smalldatetime)",
         time.Date(2000, 1, 1, 12, 13, 0, 0, time.UTC)},
        {"cast('2000-01-01' as date)",
         time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
        {"cast(NULL as date)", nil},
        {"cast(0x6F9619FF8B86D011B42D00C04FC964FF as uniqueidentifier)",
         [...]byte{0x6F, 0x96, 0x19, 0xFF, 0x8B, 0x86, 0xD0, 0x11, 0xB4, 0x2D, 0x00, 0xC0, 0x4F, 0xC9, 0x64, 0xFF}},
        {"cast(NULL as uniqueidentifier)", nil},
    }

    for _, test := range values {
        stmt, err := conn.Prepare("select " + test.sql)
        if err != nil {
            t.Error("Prepare failed:", test.sql, err.Error())
            return
        }
        defer stmt.Close()

        row := stmt.QueryRow()
        var retval interface{}
        err = row.Scan(&retval)
        if err != nil {
            t.Error("Scan failed:", test.sql, err.Error())
            return
        }
        if retval != test.val {
            t.Error("Values don't match", test.sql, retval, test.val)
            return
        }
    }
}