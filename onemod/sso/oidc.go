/*
 *@Description
 *@author          lirui
 *@create          2021-06-23 14:22
 */
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func main() {

	configURL := "http://172.20.34.203:1989/auth/realms/lr"
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, configURL)
	if err != nil {
		panic(err)
	}

	clientID := "lroidc"
	clientSecret := "e528327f-f065-4e50-8b7a-3125c669cd25"

	redirectURL := "http://localhost:1990/demo/callback"
	// Configure an OpenID Connect aware OAuth2 client.

	// 初始化 oauth2 配置
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	state := "somestate"

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	// 初始化验证器
	verifier := provider.Verifier(oidcConfig)

	// 正常业务服务
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rawAccessToken := r.Header.Get("Authorization")

		// 第一次登陆么有token 进行跳转到
		if rawAccessToken == "" {
			// 跳转的时候带上一个标记
			// State is a token to protect the user from CSRF attacks. You must
			// always provide a non-empty string and validate that it matches the
			// the state query parameter on your redirect callback.
			// 向keylocak发送请求的时候会在URL后面带上 &state=somestate
			http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
			// 跳转完成后进行认证，认证完成后keyloak就会访问下面的回调url /demo/callback
			return
		}

		// 获取token
		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != 2 {
			w.WriteHeader(400)
			return
		}

		// 验证 token
		_, err := verifier.Verify(ctx, parts[1])

		if err != nil {
			http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
			return
		}

		// 验证成功进行返回
		w.Write([]byte("hello world"))
	})

	// 配置回调服务
	http.HandleFunc("/demo/callback", func(w http.ResponseWriter, r *http.Request) {

		// 首先判断state是不是发送请求的时候注入的
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		/**
			oauth2Config.Exchange
			oauth2Token.Extra
			verifier.Verify
			下面的三个方法是协议实现，用于获取验证token
		**/
		// 回调url
		// http://localhost:1990/demo/callback?
		//state=somestate
		//&session_state=7342990d-bf59-4bd7-9e57-64b2373c1f2d
		//&code=7957a4ef-71b8-4ba2-a6ea-eba8b55be186.7342990d-bf59-4bd7-9e57-64b2373c1f2d.372c64a3-a8b4-494a-b449-c03fc9422fed

		// Exchange converts an authorization code into a token
		oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回 id_token 字段的值
		// 实际上从 token的raw里面获取
		// raw optionally contains extra metadata from the server
		// when updating a token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
		}{
			oauth2Token,
			new(json.RawMessage),
		}

		// Claims unmarshals the raw JSON payload of the ID Token into a provided struct
		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 美化json数据
		data, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 传送给客户端
		w.Write(data)

		// 给了客户端token就可以通过http发送请求
		// curl -H "Authorization: Bearer $TOKEN" localhost:1990
		// 进入 http.HandleFunc("/",  方法

	})

	log.Fatal(http.ListenAndServe("localhost:1990", nil))
}
