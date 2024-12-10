package main

import (
//    "encoding/json"
//    "path/filepath"
//    "strconv"
//    "strings"
//    "unsafe"
//    "bytes"
    "log"
//    "os"
)

const IMG_PATH = "static/";
const PRODUCT_DIR = "products/"; // temporary

type Product struct {
    Id int;
    Price float32;
    Name string;
    Company string;
    Category string;
    ShortDesc string;
    FullDesc string;
    Standards string;
    ImgPath string;
}

func Post(p Product) (int, error) {
    const queryString = "INSERT INTO electro (price, name, company, category, short_desc, full_desc, standards, img_path) VALUES (?, ?,?,?,?,?,?,?)";
    res, err := DB.Exec(queryString, p.Price, p.Name, p.Company, p.Category, p.ShortDesc, p.FullDesc, p.Standards, p.ImgPath);

    if err != nil {
        return -1, err;
    }
    id, err := res.LastInsertId();
    if err != nil {
        return -1, err;
    }

    log.Printf("Inserted %d", id);
    return int(id), nil;
}

func Get(id int) (p Product, err error) {
    const queryString = "SELECT * FROM product WHERE id=?";
    row := DB.QueryRow(queryString, id);

    err = row.Scan(&p.Id, &p.Price, &p.Name, &p.Company, &p.Category, &p.ShortDesc, &p.FullDesc, &p.Standards, &p.ImgPath);
    log.Printf("Made get query, product (id = %d)", id);

    return;
}

func GetByCategory(cat string) ([]Product, error) {
    const queryString = "SELECT * FROM product WHERE category=?";
    ps := make([]Product, 0, 16);
    rows, err := DB.Query(queryString, cat);
    defer rows.Close();

    log.Printf("Made get query, product (category = %s)", cat);
    if err != nil {
        return ps, err;
    }

    var p Product;
    for rows.Next() {
        err = rows.Scan(&p.Id, &p.Price, &p.Name, &p.Company, &p.Category, &p.ShortDesc, &p.FullDesc, &p.Standards, &p.ImgPath);
        if err != nil {
            return ps, err;
        }

        ps = append(ps, p);
    }

    return ps, err;
}

func GetCategories() (cts []string, err error) {
    const queryString = "SELECT DISTINCT category FROM product";
    cts = make([]string, 0, 16);
    rows, err := DB.Query(queryString);

    if err != nil {
        return;
    }

    for rows.Next() {
        var ct string;
        err = rows.Scan(&ct);
        if err != nil {
            return;
        }
        cts = append(cts, ct);
    }

    return;
}
