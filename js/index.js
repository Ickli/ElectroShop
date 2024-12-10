const ID_OUTER_LAYER = ".right-layer1";
const CLASS_BTN_BUY  = '.preview-btn-buy';
const CLASS_BTN_INFO = '.preview-btn-info';
const URL_CATEGORY = '/';
const URL_DETAILS = 'details';
const URL_BUY = 'buy';
const URL_CREDITS = "credits";
const URL_ACCOUNT = "account";
const CLASS_HOTBAR_BTN_CATEGORY = ".hotbar-btn-category";
const ID_HOTBAR_BTN_ACCOUNT = "#hotbar-btn-account";
const ID_HOTBAR_BTN_CREDITS = "#hotbar-btn-credits";

$(document).ready(function() {
    setUpAuthLogic();
    setUpBuyModalLogic();
    $(document).on('click', CLASS_BTN_INFO, function() {
        let data = {
            id: Number($(this).attr('data-id')),
        };
        $.ajax({
            type: "GET",
            url: URL_DETAILS,
            dataType: "html",
            data: data,
            error: function(data, status, error) {
                alert(data.error);
                alert(status);
                alert(error);
            }
        }).done(function(data) {
            $(ID_OUTER_LAYER).html(data);
            if(cookieExists()) {
                disableAuthForm();
                let login = getCookie("login");
                showAuthMsg(login, "");
            } else {
                enableAuthForm();
            }
        });
    });

    $(document).on('click', CLASS_HOTBAR_BTN_CATEGORY, function() {
        let data = {
            category: $(this).attr('data-category'),
        };
        $.ajax({
            type: "GET",
            url: URL_CATEGORY,
            dataType: "html",
            data: data,
            error: function(data, status, error) {
                alert(data.error);
                alert(status);
                alert(error);
            }
        }).done(function(data) {
            $(ID_OUTER_LAYER).html(data);
        });
    });

    $(document).on('click', ID_HOTBAR_BTN_ACCOUNT, function() {
        $.ajax({
            type: "GET",
            url: URL_ACCOUNT,
            dataType: "html",
            error: function(data, status, error) {
                alert(data.error);
                alert(status);
                alert(error);
            }
        }).done(function(data) {
            $(ID_OUTER_LAYER).html(data);
            if(cookieExists()) {
                disableAuthForm();
                let login = getCookie("login");
                showAuthMsg(login, "");
            } else {
                enableAuthForm();
            }
        });
    });
    $(document).on('click', ID_HOTBAR_BTN_CREDITS, function() {
        $.ajax({
            type: "GET",
            url: URL_CREDITS,
            dataType: "html",
            error: function(data, status, error) {
                alert(data.error);
                alert(status);
                alert(error);
            }
        }).done(function(data) {
            $(ID_OUTER_LAYER).html(data);
        });
    });
});

/**************************Auth logic**************************/

const URL_AUTH = "auth";
const ID_AUTH_FORM = "#auth-form";
const ID_AUTH_LOGIN_INPUT = "#auth-login";
const ID_AUTH_PASSWORD_INPUT = "#auth-password";
const ID_AUTH_LOGIN_BTN = "#auth-login-submit";
const ID_AUTH_REGISTER = "#auth-register";
const ID_AUTH_LOGOUT = "#auth-logout";

function setUpAuthLogic() {
    /*
        1. Cookie exists?
            yes -> form inactive, logout active
            no -> vice versa
        2. login -> send get request
        3. register -> send post request
        4. logout -> destroy cookie and make form active
     */
    // set starting activeness of form and logout button
    $(document).on('click', ID_AUTH_LOGIN_BTN, function() {
        if(cookieExists()) {
            alert("Вы уже зашли в аккаунт");
            return;
        }

        let values = getAuthValues();
        $.ajax({
            type: "GET",
            url: URL_AUTH,
            dataType: "json",
            data: values,
            /*
            error: function(data, status, error) {
                processAuthFail(data, status, error);
            },
            success: function(data, status, error) {
                processAuthSuccess(data, status, error);
            }
            */
        })
        .fail(processAuthFail)
        .done(processAuthSuccess);
    });
    $(document).on('click', ID_AUTH_REGISTER, function() {
        if(cookieExists()) {
            alert("Вы уже зашли в аккаунт");
            return;
        }

        // sending form data according to https://stackoverflow.com/questions/5392344/sending-multipart-formdata-with-jquery-ajax/5976031#5976031
        let values = new FormData();
        for(const [key, value] of Object.entries(getAuthValues())) {
            values.append(key, value);
        }
        $.ajax({
            type: "POST",
            url: URL_AUTH,
            dataType: "json",
            contentType: false,
            processData: false,
            data: values,
        })
        .fail(processAuthFail)
        .done(processAuthSuccess);
    });
    $(document).on('click', ID_AUTH_LOGOUT, function() {
        document.cookie = 'login=; Max-Age=-99999999;';
        document.cookie = 'password=; Max-Age=-99999999;';
        showAuthMsg("", "");
        enableAuthForm();
    });
}

function processAuthSuccess(data) {
    disableAuthForm();
}

function processAuthFail(data, status, error) {
    switch(data.status) {
        case 422: 
            showAuthMsg("Пустой логин или пароль", "");
            break;
        case 500:
            showAuthMsg("Ошибка сервера", "");
            break;
        case 400:
            let loginMsg;
            let pwdMsg;
            let parsed = JSON.parse(data.responseText);
            console.log(parsed);
            if(parsed.match == "password") {
                loginMsg = "Неверный логин";
                pwdMsg = `Этот пароль уже используется '${parsed.matched_login}'`;
            } else if(parsed.match == "login") {
                loginMsg = "Неверный пароль";
                pwdMsg = "Попробуйте другой пароль";
            } else {
                loginMsg = "WHAT";
                pwdMsg = "WHAT";
            }
            showAuthMsg(loginMsg, pwdMsg);
            break;
        case 404:
            showAuthMsg("Не удалось найти пользователя", "");
            break;
        default:
            alert(`unknown error; data: ${data} status: ${status} error: ${error}`);
            break;
    }
}

function showAuthMsg(loginMsg, pwdMsg) {
    $(ID_AUTH_LOGIN_INPUT).val(loginMsg);
    $(ID_AUTH_PASSWORD_INPUT).val(pwdMsg);
}

function getAuthValues() {
    let values = {}
    $("form" + ID_AUTH_FORM + " :input").each(function() {
        values[this.name] = $(this).val();
    });
    return values;
}

function enableAuthForm() {
    let form = $(ID_AUTH_FORM);
    $("form" + ID_AUTH_FORM + " :input").each(function() {
        $(this).prop("disabled", false);
    })
}

function disableAuthForm() {
    //alert('disable');
    let form = $(ID_AUTH_FORM);
    $("form" + ID_AUTH_FORM + " :input").each(function() {
        $(this).prop("disabled", true);
    })
}

function cookieExists() {
    let exists = getCookie("login") !== undefined && getCookie("password") !== undefined;
    console.log(`exists: ${exists} ${getCookie("login")} ${getCookie("password")}`);
    return exists;
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

/************************** Buy modal **************************/
const ID_BUY_MODAL = "#buy-modal";
const ID_BUY_MODAL_FORM = "#buy-modal-form";
const ID_BUY_MODAL_OVERLAY = "#buy-modal-overlay";
const ID_BUY_MODAL_NAME = "#buy-modal-name";
const ID_BUY_MODAL_PRICE = "#buy-modal-price";
const ID_BUY_MODAL_EXIT = "#buy-modal-exit";
const ID_BUY_MODAL_SUBMIT = "#buy-modal-submit";
const CLASS_BUY_MODAL_SHOW = ".buy-modal-show";

function setUpBuyModalLogic() {
    // TODO: add these buttons to every 'buy' button
    // and details button
    $(document).on('click', CLASS_BUY_MODAL_SHOW, showBuyModal);
    $(document).on('click', CLASS_BTN_BUY, showBuyModal);
    $(document).on('click', ID_BUY_MODAL_SUBMIT, function() {
        let values = new FormData();
        for(const [key, value] of Object.entries(getBuyValues())) {
            values.append(key, value);
        }

        $.ajax({
            type: "POST",
            url: URL_BUY,
            contentType: false,
            processData: false,
            data: values,
        })
        .fail(processBuyFail)
        .done(processBuySuccess);
    });
}

function showBuyModal() {
    let data = {
        id: Number($(this).attr('data-id')),
    };
    $.ajax({
        type: "GET",
        url: URL_BUY,
        dataType: "html",
        data: data,
        error: function(data, status, error) {
            if(data.status == 403) {
                alert("Необходимо войти в аккаунт или зарегистрироваться");
                return;
            }
            alert(data.error);
            alert(status);
            alert(error);
        },
        success: function(data) {
            $(ID_OUTER_LAYER).html(data);
        }
    });
}

function getBuyValues() {
    let values = {}
    $("form" + ID_BUY_MODAL_FORM + " :input").each(function() {
        values[this.name] = $(this).val();
    });
    return values;
}

function processBuyFail(data, status, error) {
    switch(data.status) {
    case 400:
        alert("Неправильно введены данные с карты");
        break;
    case 403:
        alert("Необходимо зайти в аккаунт или зарегистрироваться");
        break;
    case 404:
        alert("Пожалуйста, перезайдите в аккаунт");
        break;
    case 500:
        alert("Ошибка сервера, повторите попытку позже");
        break;
    default:
        alert(`unknown error; data: ${data} status: ${status} error: ${error}`);
        break;
    }
}

function processBuySuccess(data) {
    alert("Запрос успешно выполнен и перенаправлен в банк");
}
