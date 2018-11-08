package storage

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/robfig/cron"
    "strings"
    "errors"
    "bytes"
    "github.com/jackielihf/golib/log"
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
    available bool  // 是否可用
}

// to do: sql cache

func (that *PgClient) init() {
    that.formatConnStr()
    that.available = false
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
        log.Errorf("%v", err)
    }
    that.heartbeating()
}

func (that *PgClient) connect() {
    log.Info("connecting db: " + that.ConnStr)
    var err error
    if that.Db, err = sql.Open("postgres", that.ConnStr); err != nil {
        log.Errorf("%v", err)
    }
}

func (that *PgClient) check() error{
    if err := that.Db.Ping(); err != nil {
        that.available = false
        return err
    }
    // exec a simple query
    if rows, err2 := that.Db.Query("select 1"); rows != nil && err2 == nil { // ok
        if !that.available {
            log.Info("connected db: " + that.ConnStr)
            that.available = true
        }
        defer rows.Close()
        return nil    
    }else{
        that.available = false
        return err2
    }
}

// check the connection available or not
func (that *PgClient) heartbeating() {
    that.sched = cron.New()
    that.sched.AddFunc("@every 5s", func(){
        if err := that.check(); err != nil {
            log.Warnf("%v", err)
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
func (that *PgClient) BuildInsertSql(table string, fields []string, returning string) string{
    n := len(fields)
    placeholders := make([]string, n)
    for i := 0; i < n; i++ {
        placeholders[i] = fmt.Sprintf("$%d", i+1)
    }
    sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(fields, ","), strings.Join(placeholders, ","))
    if returning != "" {
        sql += fmt.Sprintf(" returning %s", returning)
    }
    return sql
}

func replaceQuestion(clause string, start int) string {
    var sb bytes.Buffer
    index := start
    for _, char := range clause {
        if char == '?' {
            sb.WriteString(fmt.Sprintf("$%d", index))
            index ++
        }else{
            sb.WriteRune(char)
        }
    }
    return sb.String()
}

// build update sql
// where clauses:  find placeholder "?", then replace it with $n
func (that *PgClient) BuildUpdateSql(table string, fields []string, where string, returning string) string{
    n := len(fields)
    for i, c := range fields {
        fields[i] = fmt.Sprintf("%s=$%d", c, i + 1)
    }
    whereClause := replaceQuestion(where, n + 1)
    sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(fields,","), whereClause)
    if returning != "" {
        sql += fmt.Sprintf(" returning %s", returning)
    }
    return sql
}

// build sql
func (that *PgClient) BuildSql(sql string) string{
    return replaceQuestion(sql, 1)
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
func (that *PgClient) Insert(table string, fields map[string]interface{}, returning string, src interface{}) (error){
    var keys []string
    var values []interface{}
    for key, value := range fields {
        keys = append(keys, key)
        values = append(values, value)
    }
    sql := that.BuildInsertSql(table, keys, returning)
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

// update
func (that *PgClient) Update(table string, fields map[string]interface{}, where string, vars ...interface{}) (int64, error){
    var keys []string
    var values []interface{}
    for key, value := range fields {
        keys = append(keys, key)
        values = append(values, value)
    }
    for _, value := range vars {
        values = append(values, value)
    }
    sql := that.BuildUpdateSql(table, keys, where, "")
    if stmt, err := that.Db.Prepare(sql); err == nil {
        defer stmt.Close()  // close when this function returns
        if res, err2 := stmt.Exec(values...); err2 != nil {
            return 0, err2
        }else{
            return res.RowsAffected()
        }
    }else{
        return 0, err
    }
}

type RowPage struct {
    Total int64
    Page int
    Limit int
    Rows *sql.Rows    
}

type PageInfo struct {
    Total int64
    Page int64
    Limit int64
}

func (that *PgClient) SelectRowPage(sql string, page int, limit int, values ...interface{}) (*RowPage, error){
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

func (that *PgClient) SelectPage(sql string, page int64, limit int64, doMapping func()(interface{}, map[string]interface {}), values ...interface{}) (PageInfo, []interface{}, error){
    pageInfo := PageInfo{0, page, limit}
    // count
    countSql := fmt.Sprintf("select count(1) as total from (%s) _alias", sql)
    if row, err := that.QueryRow(countSql, values...); err != nil {
        return pageInfo, nil, err
    }else{
        if err2 := row.Scan(&pageInfo.Total); err2 != nil {
            return pageInfo, nil, err2
        }
    }
    // no result
    if pageInfo.Total == 0 {
        return pageInfo, nil, nil
    }
    // get page
    offset := (page - 1) * limit
    pageSql := fmt.Sprintf("%s limit %d offset %d", sql, limit, offset)
    if rows, err3 := that.Query(pageSql, values...); err3 != nil {
        return pageInfo, nil, err3
    }else{
        if result, err4 := that.ListScan(rows, doMapping); err4 != nil {
            return pageInfo, nil, err4  
        }else{
            return pageInfo, result, nil    
        }
    }
}

// 查询一个结果，存放到src中
func (that *PgClient) SelectOne(sql string, src FieldMapping, values ...interface{}) (int, error) {
    if rows, err := that.Query(sql, values...); err == nil {
        defer rows.Close()
        // 扫描一次
        if rows.Next() {
            return 1, that.FieldScan(rows, src)    
        }
        // no available row
        return 0, nil
    }else{
        return 0, err
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
// 扫描rows，根据字段映射，获取对象列表
func (that *PgClient) ListScan(rows *sql.Rows, doMapping func()(interface{}, map[string]interface {})) ([]interface{}, error) {
    var list []interface{}
    if rows == nil {
        return list, nil
    }
    for rows.Next() {
        ptr, mapping := doMapping()
        if ptr == nil || mapping == nil {
            return list, errors.New("ListScan err: ptr or mapping is nil")
        }
        if err := that.FieldScan(rows, mapping); err != nil {
            return list, err
        }
        list = append(list, ptr)
    }
    return list, nil
}





