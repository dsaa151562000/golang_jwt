package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/google_oauth/db"
	e "github.com/google_oauth/entity"
	s "github.com/google_oauth/service"
)

// JWT jwtを格納する構造体
type JWT struct {
	Token string `json:"token"`
}

// Error Errorを格納する構造体
type Error struct {
	Message string `json:"message"`
}

type Ping struct {
	Status int
	Rssult string
}

type User e.User

// func 関数名 (引数 型, 引数 型)
// JSON 形式で結果を返却
// data interface{} とすると、どのような変数の型でも引数として受け取ることができる
func responseByJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
	return
}

// レスポンスにエラーを突っ込んで、返却するメソッド
func errorInResponse(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status) // 400 とか 500 などの HTTP status コードが入る
	json.NewEncoder(w).Encode(error)
	return
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user e.User
	var error Error

	// https://golang.org/pkg/encoding/json/#NewDecoder
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email は必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "パスワードは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	// dump
	// fmt.Println("---------------------")
	// spew.Dump(user)

	// パスワードのハッシュを生成
	// https://godoc.org/golang.org/x/crypto/bcrypt#GenerateFromPassword
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("パスワード: ", user.Password)
	// fmt.Println("ハッシュ化されたパスワード", hash)
	user.Password = string(hash)

	// できれば引数は&userで渡したい
	//p, err := s.CreateUser(&user)
	var s s.Service
	spew.Dump(&user)
	p, err := s.CreateUser(user.Email, user.Password)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		//Dumpを吐く
		spew.Dump(err)
		//responseByJSON(w, err)
	} else {
		//DB に登録できたらパスワードをからにしておく
		p.Password = ""
		// jsonを返却する
		responseByJSON(w, p)
	}
}

// Token 作成関数
func createToken(user e.User) (string, error) {
	var err error

	// 鍵となる文字列(多分なんでもいい)
	secret := "secret"

	// Token を作成
	// jwt -> JSON Web Token - JSON をセキュアにやり取りするための仕様
	// jwtの構造 -> {Base64 encoded Header}.{Base64 encoded Payload}.{Signature}
	// HS254 -> 証明生成用(https://ja.wikipedia.org/wiki/JSON_Web_Token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "__init__", // JWT の発行者が入る(文字列(__init__)は任意)
	})

	//Dumpを吐く
	spew.Dump(token)

	tokenString, err := token.SignedString([]byte(secret))

	fmt.Println("-----------------------------")
	fmt.Println("tokenString:", tokenString)

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func login(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("login 関数実行")
	var user e.User
	var error Error
	var jwt JWT

	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email は必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "パスワードは、必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
	}

	// 追加(この位置であること)
	// password := user.Password

	//db := db.GetDB()
	//db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
	//db.Where(" email = ? AND password = ?", user.Email ,user.Password).First(&user)

	// if existUser(user.Email) {
	// 	error.Message = "ユーザが存在しません。"
	// 	return
	// } else {
	// 	fmt.Println("ユーザは存在します。")
	// }

	var s s.Service
	p, err := s.GetByPassword(user.Email)

	if err != nil {
		error.Message = "ユーザが存在しません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
		//fmt.Println("ユーザが存在しません。")
	} else {
		// fmt.Println("ユーザは存在します。")
		fmt.Println(p.Password)
		// w.Header().Set("Content-Type", "application/json")
		// responseByJSON(w, p)
	}

	err = bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(user.Password))
	if err != nil {
		error.Message = "無効なパスワードです。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	token, err := createToken(user)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	jwt.Token = token
	responseByJSON(w, jwt)
}

// 認証結果をブラウザに返却
func verifyEndpoint(w http.ResponseWriter, r *http.Request) {
	responseByJSON(w, "認証OK")
}

// verifyEndpoint のラッパーみたいなもの
func tokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var errorObject Error
		ping := Ping{http.StatusOK, "ok"}
		res, err := json.Marshal(ping)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// HTTP リクエストヘッダーを読み取る
		authHeader := r.Header.Get("Authorization")
		// Restlet Client から以下のような文字列を渡す
		// bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3Q5OUBleGFtcGxlLmNvLmpwIiwiaXNzIjoiY291cnNlIn0.7lJKe5SlUbdo2uKO_iLzzeGoxghG7SXsC3w-4qBRLvs
		bearerToken := strings.Split(authHeader, " ")
		fmt.Println("bearerToken: ", bearerToken)
		fmt.Println(len(bearerToken))

		if len(bearerToken) == 3 {
			authToken := bearerToken[2]

			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("エラーが発生しました。")
				}
				return []byte("secret"), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				errorInResponse(w, http.StatusUnauthorized, errorObject)
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Println(claims["email"])
			} else {
				fmt.Println(err)
			}

			if token.Valid {
				// レスポンスを返す
				// next.ServeHTTP(w, r)
				w.Header().Set("Content-Type", "application/json")
				w.Write(res)
			} else {
				errorObject.Message = error.Error()
				errorInResponse(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Token が無効です。"
			errorInResponse(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}

// ユーザー存在チェック
func existUser(email string) bool {
	var user e.User
	var s s.Service
	p, err := s.GetByPassword(user.Email)

	if err != nil {
		spew.Dump(p)
		return true
	} else {
		return false
	}
}

func main() {
	db.Init()
	// db.Close()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//fmt.Println(os.Getenv("googleClientID"))
	//fmt.Println(os.Getenv("googleClientSecret"))

	router := mux.NewRouter()
	router.HandleFunc("/singup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/verify", tokenVerifyMiddleWare(verifyEndpoint)).Methods("GET")

	// console に出力する
	log.Println("サーバー起動 : 8000 port で受信")

	// log.Fatal は、異常を検知すると処理の実行を止めてくれる
	log.Fatal(http.ListenAndServe(":8000", router))
}
