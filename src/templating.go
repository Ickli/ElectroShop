package main

import (
    "text/template"
    "log"
//    "fmt"
//    "os"
    "io"
)

type mainPageInfo struct {
    ToWrap bool;
    CategoryName string;
    Products []Product;
    Categories map[string]int;
}

const (
    NAME_ACCOUNT_TEMPLATE = "auth";
    NAME_BASE_TEMPLATE = "main";

    PATH_PREVIEW = PATH_HTML + "preview_tmpl.html";
    PATH_MAIN = PATH_HTML + "main_tmpl.html";
    PATH_DETAILS = PATH_HTML + "details_tmpl.html";
    PATH_BUY = PATH_HTML + "buy_tmpl.html";
    PATH_ACCOUNT = PATH_HTML + "auth_tmpl.html";
)

var (
    detailsTemplate = template.Must(template.ParseFiles(PATH_DETAILS, PATH_ACCOUNT));
    mainTemplate = template.Must(template.ParseFiles(PATH_MAIN, PATH_PREVIEW));
    buyTemplate = template.Must(template.ParseFiles(PATH_BUY));
    accountTemplate = template.Must(template.ParseFiles(PATH_ACCOUNT));
)

func WriteDetailsFilled(productId int, wr io.Writer) error {
    product, err := Get(productId);

    if err != nil {
        log.Printf("Couldn't get info about product (id = %d): %s", productId, err.Error());
        return err;
    }

    err = detailsTemplate.ExecuteTemplate(wr, NAME_BASE_TEMPLATE, product);
    
    if err != nil {
        log.Printf("Couldn't fill preview template: %s", err.Error());
        return err;
    }
    return nil;
}

func WritePreviewsFilled(category string, w io.Writer, toWrap bool) error {
    info := mainPageInfo {
        ToWrap: toWrap,
        CategoryName: category,
        Products: make([]Product, 0, 16),
        Categories: make(map[string]int),
    };

    products, err := GetByCategory(category);
    if err != nil {
        return err;
    }

    for i := 0; i < len(products); i++ {
        product := products[i];
        if product.Category == category {
            info.Products = append(info.Products, product);
        }
    }

    if toWrap {
        cts, err := GetCategories();
        if err != nil {
            return err;
        }

        for _, ct := range cts {
            info.Categories[ct] = 1;
        }
    }

    return mainTemplate.ExecuteTemplate(w, NAME_BASE_TEMPLATE, info);
}

func WriteBuyProductFilled(id int, w io.Writer) error {
    product, err := Get(id);

    if err != nil {
        return err;
    }

    return buyTemplate.Execute(w, product);
}

func WriteAccount(w io.Writer) error {
    return accountTemplate.ExecuteTemplate(w, NAME_ACCOUNT_TEMPLATE,  nil);
}
