package main

import (
//    "database/sql";
//    "log"
)

type User struct {
    Login string;
    Pass string;
}

func UserPost(login string, pwd string) error {
    const queryString = "INSERT INTO user (login, password) VALUES (?, ?)";
    _, err := DB.Exec(queryString, login, pwd);

    return err;
}

func UserGet(login string) (u User, err error) {
    const queryString = "SELECT * FROM user WHERE login=?";

    row := DB.QueryRow(queryString, login);

    if err = row.Err(); err != nil {
        return;
    }

    err = row.Scan(&u.Login, &u.Pass);
    return;
}

func UserPut(login string, pwd string) error {
    const queryString = "UPDATE user SET password=? WHERE login=? LIMIT 1";
    _, err := DB.Exec(queryString, pwd, login);

    return err;
}
func UserGetByLoginOrPassword(login string, pwd string) (u User, found bool) {
    const queryString = "SELECT * FROM user WHERE login=? OR password=?";

    row := DB.QueryRow(queryString, login, pwd);
    if err := row.Err(); err != nil {
        return;
    }

    if row.Scan(&u.Login, &u.Pass) != nil {
        return;
    }

    found = true;
    return;
}
