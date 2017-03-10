package storage

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/robfig/cron"
    "strings"
    // "reflect"
    "errors"
)

// postgresql client
type PgClient struct {
    Host string
    Port string
    User string
    Password string
    Dbname string
    Db *sql.DB
    ConnStr string
    sched *cron.Cron
}

func (that *PgClient) init() {
    that.formatConnStr()
}

func (that *PgClient) formatConnStr() {
    if that.Password != "" {
        that.ConnStr = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", that.Host, that.Port, that.User, that.Password, that.Dbname)    
    }else{
        that.ConnStr = fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=disable", that.Host, that.Port, that.User, that.Dbname)    
    }
}
// open a db connection
func (that *PgClient) Open() {
    that.init()
    that.connect()
    if err := that.check(); err != nil {
        fmt.Println(err)
    }
    that.heartbeating()
}

func (that *PgClient) connect() {
    fmt.Println("connecting db: " + that.ConnStr)
    var err error
    if that.Db, err = sql.Open("postgres", that.ConnStr); err != nil {
        fmt.Println(err)
    }
}

func (that *PgClient) check() error{
    // fmt.Println("check")
    if err := that.Db.Ping(); err == nil {
        // exec a simple query
        if rows, err2 := that.Db.Query("select 1"); rows != nil && err2 == nil {
            defer rows.Close()
            return nil    
        }else{
            return err2
        }    
    }else{
        return err    
    }
}

// check the connection available or not
func (that *PgClient) heartbeating() {
    that.sched = cron.New()
    that.sched.AddFunc("@every 5s", func(){
        if err := that.check(); err != nil {
            fmt.Println(err)
            tempDb := that.Db
            that.connect()
            tempDb.Close()        
        }
    })
    that.sched.Start()
}

func (that *PgClient) Close() {
    that.sched.Stop()
    that.Db.Close()
}


// tools
// sql builder
// build insert sql
func (that *PgClient) BuildInsertSql(table string, fields string, returning string) string{
    n := strings.Count(fields, ",") + 1
    placeholders := make([]string, n)
    for i := 1; i <= n; i++ {
        placeholders[i-1] = fmt.Sprintf("$%d", i)
    }
    sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, fields, strings.Join(placeholders, ","))
    if returning != "" {
        sql += fmt.Sprintf(" returning %s", returning)
    }
    return sql
}

// build update sql
// where clauses:  find placeholder "?", then replace it with $n
func (that *PgClient) BuildUpdateSql(table string, fields []string, where string, returning string) string{
    n := len(fields)
    for i, c := range fields {
        fields[i] = fmt.Sprintf("%s=$%d", c, i + 1)
    }
    where += " " // add a blank
    clauses := strings.Split(where, "?")
    m := len(clauses)
    for i := 0; i < m-1; i++ {
        clauses[i] = fmt.Sprintf("%s$%d", clauses[i], n+i+1)
    }
    sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(fields,","), strings.Join(clauses, ""))
    if returning != "" {
        sql += fmt.Sprintf(" returning %s", returning)
    }
    return sql
}

// build sql
func (that *PgClient) BuildSql(sql string) string{
    sql += " " // add a blank
    strs := strings.Split(sql, "?")
    n := len(strs)
    for i := 0; i < n-1; i++ {
        strs[i] = fmt.Sprintf("%s$%d", strs[i], i+1)
    }
    return strings.Join(strs, "")
}

// sql command
// prepare
func (that *PgClient) Prepare(sql string) (*sql.Stmt, error){
    return that.Db.Prepare(sql)
}

// query
func (that *PgClient) Query(sql string, values ...interface{}) (*sql.Rows, error){
    sql = that.BuildSql(sql)
    if stmt, err := that.Db.Prepare(sql); err == nil {
        defer stmt.Close()  // close when this function returns
        return stmt.Query(values...)
    }else{
        return nil, err
    }
}

func (that *PgClient) QueryRow(sql string, values ...interface{}) (*sql.Row, error){
    sql = that.BuildSql(sql)
    if stmt, err := that.Db.Prepare(sql); err == nil {
        defer stmt.Close()  // close when this function returns
        return stmt.QueryRow(values...), nil
    }else{
        return nil, err
    }
}

// insert
// @params
func (that *PgClient) Insert(table string, fields string, returning string, src interface{}, values ...interface{}) (error){
    sql := that.BuildInsertSql(table, fields, returning)
    if stmt, err := that.Db.Prepare(sql); err == nil {
        defer stmt.Close()  // close when this function returns
        row := stmt.QueryRow(values...)
        if returning != "" && src != nil{
            return row.Scan(src)
        }
        return nil
    }else{
        return err
    }
}


type RowPage struct {
    Total int64
    Page int
    Limit int
    Rows *sql.Rows    
}


func (that *PgClient) SelectPage(sql string, page int, limit int, values ...interface{}) (*RowPage, error){
    resultPage := RowPage{0, page, limit, nil}
    // count
    countSql := fmt.Sprintf("select count(1) as total from (%s) _alias", sql)
    if row, err := that.QueryRow(countSql, values...); err == nil {
        if err2 := row.Scan(&resultPage.Total); err2 != nil {
            return nil, err2
        }
    }else{
        return nil, err
    }
    // get page
    offset := (page - 1) * limit
    pageSql := fmt.Sprintf("%s limit %d offset %d", sql, limit, offset)
    if rows, err3 := that.Query(pageSql, values...); err3 == nil {
        resultPage.Rows = rows
        return &resultPage, nil
    }else{
        return nil, err3
    }
}

// 查询一个结果，存放到src中
func (that *PgClient) SelectOne(sql string, src FieldMapping, values ...interface{}) error{
    if rows, err := that.Query(sql, values...); err == nil {
        defer rows.Close()
        // 扫描一次
        if rows.Next() {
            return that.FieldScan(rows, src)    
        }
        return errors.New("no available row")
    }else{
        return err
    }
}

// 返回用于存放扫描结果的slice和每个元素的指针
func makePointers(colNum int) (pointers []interface{}, values []interface{}) {
    pointers = make([]interface{}, colNum)
    values = make([]interface{}, colNum)
    for i := 0; i < colNum; i++ {
        pointers[i] = &values[i]
    }
    return pointers, values
}

// 字段映射类
type FieldMapping map[string]interface{}

// 扫描行，结果存放在map中
func (that *PgClient) MapScan(rows *sql.Rows) (FieldMapping, error){
    cols, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    pointers, _ := makePointers(len(cols))

    if err2 := rows.Scan(pointers...); err2 != nil {
        return nil, err2
    }
    ret := make(FieldMapping)
    for i, name := range cols {
        v := pointers[i].(*interface{})
        ret[name] = *v   
    }
    return ret, nil
}
// 按指定的字段名进行扫描, 结果存放在map中
func (that *PgClient) CustomMapScan(rows *sql.Rows, ret FieldMapping) error{
    cols, err := rows.Columns()
    if err != nil {
        return err
    }
    // make pointers
    pointers, _ := makePointers(len(cols))

    if err2 := rows.Scan(pointers...); err2 != nil {
        return err2
    }
    // pick fields
    for i, name := range cols {
        if _, ok := ret[name]; ok{
            v := pointers[i].(*interface{})
            ret[name] = *v   
        }
    }
    return nil
}

// field scan
// fieldAddr 存放变量内存地址的map
func (that *PgClient) FieldScan(rows *sql.Rows, fieldAddr FieldMapping) (error){
    cols, err := rows.Columns()
    if err != nil {
        return err
    }
    pointers, _ := makePointers(len(cols))

    // set field address by column name
    for i, name := range cols {
        if addr, ok := fieldAddr[name]; ok{
            pointers[i] = addr 
        }
    }
    // scan
    if err2 := rows.Scan(pointers...); err2 != nil {
        return err2
    }
    return nil
}

// list scan
// func (that *PgClient) ListScan(rows *sql.Rows, list []interface{}, doMapping func()(interface{}, FieldMapping)) error{
//     for rows.Next() {
//         ptr, mapping := doMapping()
//         if mapping == nil {
//             return errors.New("mapping is nil")
//         }
//         if err := that.FieldScan(rows, mapping); err != nil {
//             return err
//         }else{
//             list = append(list, ptr)
//         }
//     }
//     return nil
// }





