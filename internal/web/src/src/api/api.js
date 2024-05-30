const AuthPath = "http://localhost:6600"
const AccountPath = "/auth/account"
const LoginPath = "/auth/login"

const ProfilePath = "http://localhost:6700"
const ProfileUpsertPath = "/profile/upsert"
const ProfileGetPath = "/profile/get"


const CreateAccount = (name, login, password, onSuccess, onFailure) => {
    const body = JSON.stringify({
            "name": name,
            "login": login,
            "password": password,
            "password_confirm": password
        });

    fetch(AuthPath + AccountPath, bodyRequestOptions(body, headers(), "PUT"))
        .then(response => response.json())
        .then(b => {
            if (b.status >= 400) {
                onFailure(b.data)
            } else {
                onSuccess(b.data)
            }
        })
        .catch(onFailure);
}

const Login = (login, password, onSuccess, onFailure) => {
    const body = JSON.stringify({
        "login":login,
        "password": password,
    });

    fetch(AuthPath + LoginPath, bodyRequestOptions(body, headers(),"POST"))
        .then(response => response.json())
        .then(b => {
            if (b.status >= 400) {
                onFailure(b.data)
            } else {
                onSuccess(b.data)
            }
        })
        .catch(onFailure);
}

const GetProfile = (accessToken, onSuccess, onFailure) => {
    fetch(ProfilePath + ProfileGetPath, getRequestOptions(headersWithBearer(accessToken)))
        .then(response => response.json())
        .then(b => {
            if (b.status >= 400) {
                onFailure(b.data)
            } else {
                onSuccess(b.data)
            }
        })
        .catch(onFailure);
}

const UpsertProfile = (accessToken,nickname, about_me, onSuccess, onFailure) => {
    const body = JSON.stringify({
        "nickname":nickname,
        "about_me": about_me,
    });

    fetch(ProfilePath + ProfileUpsertPath,
        bodyRequestOptions(body, headersWithBearer(accessToken),"POST"))
        .then(response => response.json())
        .then(b => {
            if (b.status >= 400) {
                onFailure(b.data)
            } else {
                onSuccess(b.data)
            }
        })
        .catch(onFailure);
}

export {
    CreateAccount,
    Login,
    GetProfile,
    UpsertProfile
}

const headers = () => {
    const headers = new Headers();
    headers.append("Content-Type", "application/json");

    return headers
}

const headersWithBearer = (accessToken) => {
    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("Bearer", accessToken);

    return headers
}

const bodyRequestOptions = (body, headers ,method) => {
    return {
        method: method,
        headers: headers,
        body: body,
        redirect: 'follow'
    };
}

const getRequestOptions = (headers) => {
    return {
        method: "GET",
        headers: headers,
        redirect: 'follow'
    };
}