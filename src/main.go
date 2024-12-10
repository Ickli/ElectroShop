package main

import (
    "net/http"
    "log"
    fp "path/filepath"
    rg "regexp"
    "strconv"
    "strings"
    "fmt"
//    "os"
)

const (
    // filepaths
    PATH_HTML = "html/";
    PATH_MAIN_TMPL = PATH_HTML + "main_tmpl.html";

    // filenames
    NAME_CREDITS_TMPL = "credits_tmpl.html";

    // URLs
    URL_MAIN = "/{$}";
    URL_STATIC = "/static/{staticFile}";
    URL_DETAILS = "/details";
    URL_CREDITS = "/credits";
    URL_ACCOUNT = "/account";
    URL_AUTH = "/auth";
    URL_BUY = "/buy";

    // other
    DEFAULT_CATEGORY = "NoCategory";
    COOKIE_MAX_AGE_SECS = 60*20;
)

// handles both first visit and ajax calls
func mainHandle(w http.ResponseWriter, r *http.Request) {
    category := r.URL.Query().Get("category");

    if category == "" {
        category = DEFAULT_CATEGORY;
    }

    firstVisit := len(r.Header["Referer"]) == 0 || !strings.Contains(r.Header["Referer"][0], r.Host);

    log.Printf("REQUEST: main page (cat = %s, first-visit = %t)\n", category, firstVisit);
    err := WritePreviewsFilled(category, w, firstVisit);

    if err != nil {
        log.Printf("Failed to fill main page with previews: %s", err.Error());
        w.WriteHeader(http.StatusInternalServerError);
    }
}

func staticHandle(w http.ResponseWriter, r *http.Request) {
    staticFile := r.PathValue("staticFile");
    log.Printf("REQUEST: file %s\n", staticFile);

    if staticFile == "" {
        w.WriteHeader(http.StatusNotFound);
        return;
    }

    switch ext := fp.Ext(staticFile); ext {
    case ".js":
        http.ServeFile(w, r, fp.Join("js", staticFile));
    case ".html":
        http.ServeFile(w, r, fp.Join("html", staticFile));
    case ".css":
        http.ServeFile(w, r, fp.Join("css", staticFile));
    default:
        http.ServeFile(w, r, fp.Join("static", staticFile));
    }
}

func detailsHandle(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id");

    if idStr == "" {
        w.WriteHeader(http.StatusNotFound);
        return;
    }

    id, err := strconv.Atoi(idStr);
    
    if err != nil {
        w.WriteHeader(http.StatusNotFound);
        return;
    }

    log.Printf("REQUEST: details (id = %d)", id);
    
    w.Header().Add("Content-Type", "text/html");
    err = WriteDetailsFilled(id, w);

    if err != nil {
        w.WriteHeader(http.StatusNotFound);
        return;
    }
}

func creditsHandle(w http.ResponseWriter, r *http.Request) {
    log.Printf("REQUEST: credits");
    r.SetPathValue("staticFile", NAME_CREDITS_TMPL);
    staticHandle(w, r);
}

func authHandlePost(w http.ResponseWriter, r *http.Request) {
    log.Printf("REQUEST: auth (method = POST)");
    /*
        1. Check if login is free
            if not free -> return 400 and json with error msg
            (assuming password is not empty)
        2. Check if password is free :)
            if not free -> return 400 and json with error msg, containing login
        3. Post that user
        4. Set cookie
    */
    err := r.ParseMultipartForm(4096);

    if err != nil {
        log.Printf("ERROR: couldn't parse form: %s", err.Error());
        w.WriteHeader(http.StatusInternalServerError);
        return;
    }

    login := r.FormValue("login");
    pwd := r.FormValue("password");

    log.Printf("Login: %s; Password: %s", login, pwd);
    if pwd == "" || login == "" {
        log.Printf("ERROR: empty login or password");
        w.WriteHeader(http.StatusUnprocessableEntity);
        return;
    }

    user, found := UserGetByLoginOrPassword(login, pwd);

    w.Header().Add("Content-Type", "application/json");
    if !found {
        err = UserPost(login, pwd);
        if err != nil {
            log.Printf("ERROR: %s", err.Error());
            w.WriteHeader(http.StatusInternalServerError);
            return;
        }
        user.Login = login;
        user.Pass = pwd;
        setAuthCookie(w, user); // success
        fmt.Fprintf(w, "{}");
        return;
    }

    w.WriteHeader(http.StatusBadRequest);
    jsonErrorFmt := "{ \"match\": \"%s\", \"matched_login\": \"%s\" }";

    fmt.Printf("%v\n", user);
    if user.Login == login {
        log.Printf("ERROR: user found, login match");
        fmt.Fprintf(w, jsonErrorFmt, "login", "");
    } else if user.Pass == pwd {
        log.Printf("ERROR: user found, password match");
        fmt.Fprintf(w, jsonErrorFmt, "password", user.Login);
    } else {
        log.Printf("ERROR: user found, but no match");
    }
}

func authHandleGet(w http.ResponseWriter, r *http.Request) {
    /*
        1. Fetch user with login provided
        2. If user doesn't exist -> 404
        3. If user's password doesn't match -> 401
        4. Set cookie 
    */
    log.Printf("REQUEST: auth (method = GET)");
    login := r.URL.Query().Get("login");
    pwd := r.URL.Query().Get("password");

    if login == "" || pwd == "" {
        log.Printf("ERROR: invalid login or password");
        w.WriteHeader(http.StatusUnprocessableEntity);
        return;
    }

    userWithSamePwd, found := UserGetByLoginOrPassword(login, pwd);
    actualUser, errFindActualUser := UserGet(login);

    if !found {
        log.Printf("ERROR: no user found");
        w.WriteHeader(http.StatusNotFound);
        return;
    }

    w.Header().Set("Content-Type", "application/json");
    const jsonErrorFmt = "{ \"match\": \"%s\", \"matched_login\": \"%s\" }";

    if errFindActualUser == nil {
        if actualUser.Pass == pwd {
            log.Printf("SUCCESS: password match"); // success
            setAuthCookie(w, actualUser);
            fmt.Fprintf(w, "{}");
            return;
        } else {
            log.Printf("ERROR: password mismatch");
            w.WriteHeader(http.StatusBadRequest);
            fmt.Fprintf(w, jsonErrorFmt, "login", "");
            return;
        }
    }

    // if pwd match but login mismatch
    w.WriteHeader(http.StatusBadRequest);
    fmt.Fprintf(w, jsonErrorFmt, "password", userWithSamePwd.Login);
}

func authHandle(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case "POST": authHandlePost(w, r);
        case "GET": authHandleGet(w, r);
        default: w.WriteHeader(http.StatusMethodNotAllowed);
    }
}

func setAuthCookie(w http.ResponseWriter, u User) error {
    cookieLogin, err := http.ParseSetCookie(fmt.Sprintf("login=%s", u.Login));
    if err != nil {
        return err;
    }

    cookiePwd, err := http.ParseSetCookie(fmt.Sprintf("password=%s", u.Pass));
    if err != nil {
        return err;
    }

    cookieLogin.MaxAge = COOKIE_MAX_AGE_SECS;
    cookiePwd.MaxAge = COOKIE_MAX_AGE_SECS;
    http.SetCookie(w, cookieLogin);
    http.SetCookie(w, cookiePwd);
    return nil;
}

func getAuthCookies(r *http.Request) (u User, err error) {
    cLogin, err := r.Cookie("login");
    
    if err != nil {
        return;
    }

    cPass, err := r.Cookie("password");

    if err != nil {
        return;
    }

    u.Login = cLogin.Value;
    u.Pass = cPass.Value;
    return;
}

func accountHandle(w http.ResponseWriter, r *http.Request) {
    log.Printf("REQUEST: account");
    if err := WriteAccount(w); err != nil {
        log.Printf("ERROR: %s", err.Error());
    }
}

func buyHandle(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case "POST": buyHandlePost(w, r);
        case "GET": buyHandleGet(w, r);
        default: w.WriteHeader(http.StatusMethodNotAllowed);
    }
}

func buyHandlePost(w http.ResponseWriter, r *http.Request) {
    // regex taken from: https://stackoverflow.com/questions/9315647/regex-credit-card-number-tests
    // const rgCardNumber = `^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`;
    var patternNonDigit = rg.MustCompile("[^0-9]");
    const rgCardNumber = "[0-9]{16}";
    const rgCVV = "^[0-9]{3}$";
    const rgExpDate = "^[0-9]{2}.[0-9]{2}.[0-9]{4}$";

    log.Printf("REQUEST: buy (method = POST)");

    sentUser, err := getAuthCookies(r);

    if err != nil {
        w.WriteHeader(http.StatusForbidden);
        return;
    }

    actualUser, err := UserGet(sentUser.Login);

    if err != nil || actualUser.Pass != sentUser.Pass {
        log.Printf("ERROR: user not found: %s", err.Error());
        w.WriteHeader(http.StatusNotFound);
        return;
    }

    err = r.ParseMultipartForm(4096);

    if err != nil {
        log.Printf("ERROR: couldn't parse form: %s", err.Error());
        w.WriteHeader(http.StatusInternalServerError);
        return;
    }

    CVV := r.FormValue("cvv");
    cardNumber := r.FormValue("card_number");
    cardNumber = patternNonDigit.ReplaceAllString(cardNumber, "");
    expDate := r.FormValue("expiration_date");

    log.Printf("%s %s %s", CVV, cardNumber, expDate);
    if match, err := rg.MatchString(rgCVV, CVV); !match || err != nil {
        log.Printf("ERROR: invalid CVV");
        w.WriteHeader(http.StatusBadRequest);
    } else if match, err := rg.MatchString(rgCardNumber, cardNumber); !match || err != nil {
        log.Printf("ERROR: invalid card number");
        w.WriteHeader(http.StatusBadRequest);
    } else if match, err := rg.MatchString(rgExpDate, expDate); !match || err != nil {
        log.Printf("ERROR: invalid expiration date");
        w.WriteHeader(http.StatusBadRequest);
    }
}

func buyHandleGet(w http.ResponseWriter, r *http.Request) {
    log.Printf("REQUEST: buy (method = GET)");

    _, err := getAuthCookies(r);

    if err != nil {
        w.WriteHeader(http.StatusForbidden);
        return;
    }

    idStr := r.URL.Query().Get("id");

    if idStr == "" {
        w.WriteHeader(http.StatusNotFound);
        return;
    }

    id, err := strconv.Atoi(idStr);
    
    if err != nil {
        w.WriteHeader(http.StatusNotFound);
        return;
    }

    err = WriteBuyProductFilled(id, w);

    if err != nil {
        log.Printf("ERROR: couldn't fill template: %s", err.Error());
        w.WriteHeader(http.StatusInternalServerError);
        return;
    }
}

func main() {
    log.Printf("Hello cruel world!");

    http.HandleFunc(URL_STATIC, staticHandle);
    http.HandleFunc(URL_MAIN, mainHandle);
    http.HandleFunc(URL_DETAILS, detailsHandle);
    http.HandleFunc(URL_CREDITS, creditsHandle);
    http.HandleFunc(URL_ACCOUNT, accountHandle);
    http.HandleFunc(URL_AUTH, authHandle);
    http.HandleFunc(URL_BUY, buyHandle);
    
    log.Fatal(http.ListenAndServe(":8080", nil));
}
