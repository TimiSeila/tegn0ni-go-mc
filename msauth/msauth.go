package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Authorize Device Req n Res
type AuthorizeDeviceRequest struct {
	ClientId string `json:"client_id"`
	Scope string `json:"scope"`
}

type AuthorizeDeviceResponse struct {
	DeviceCode string `json:"device_code"`
	UserCode string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	ExpiresIn int `json:"expires_in"`
	Interval int `json:"interval"`
	Message string `json:"message"`
}

// Authenicate User Req n Res
type AuthenticateUserRequest struct {
	GrantType string `json:"grant_type"`
	ClientId string `json:"client_id"`
	DeviceCode string `json:"device_code"`
}

type AuthenticateUserResponse struct {
	TokenType string `json:"token_type"`
	Scope string `json:"scope"`
	ExpiresIn int `json:"expires_in"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken string `json:"id_token"`
}

// Authenticate to Xbox Live Req n Res
type XBOXAuthRequest struct {
	Properties struct {
		AuthMethod string `json:"AuthMethod"`
		SiteName string `json:"SiteName"`
		RpsTicket string `json:"RpsTicket"`
	} `json:"Properties"`
	RelyingParty string `json:"RelyingParty"`
	TokenType string `json:"TokenType"`
}

type XBOXAuthResponse struct {
	IssueInstant string `json:"IssueInstant"`
	NotAfter string `json:"NotAfter"`
	Token string `json:"Token"`
	DisplayClaims struct {
		XUI []struct {
			UHS string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

// Obtain XSTS token Req n Res
type XSTSRequest struct {
	Properties struct {
		SandboxID string `json:"SandboxId"`
		UserTokens []string `json:"UserTokens"`
	} `json:"Properties"`
	RelyingParty string `json:"RelyingParty"`
	TokenType string `json:"TokenType"`
}

type XSTSResponse struct {
	IssueInstant string `json:"IssueInstant"`
	NotAfter string `json:"NotAfter"`
	Token string `json:"Token"`
	DisplayClaims struct {
		XUI []struct {
			UHS string `json:"uhs"`
		} `json:"xui"`
	} `json:"DisplayClaims"`
}

// Authenticate to Minecraft Req n Res
type MinecraftAuthRequest struct {
	IdentityToken string `json:"identityToken"`
}

type MinecraftAuthResponse struct {
	Username string `json:"username"`
	Roles []string `json:"roles"`
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
}

// Minecraft profile Res
type MinecraftProfileResponse struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Skins []struct {
		Id string `json:"id"`
		State string `json:"state"`
		Url string `json:"url"`
		Variant string `json:"variant"`
		Alias string `json:"alias"`
	} `json:"skins"`
	Capes []struct {
		Id string `json:"id"`
		State string `json:"state"`
		Url string `json:"url"`
		Alias string `json:"alias"`
	} `json:"capes"`
}

const client_id = "6a2c1746-4dc3-4e5a-bbc5-fffa060c3006"

func main(){
	ctx := context.Background()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}


	device_code, err := Authorize_Device(ctx, client)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	var access_token string
	for {
		fmt.Println("Authenticating Microsoft user")
		access_token, err = Authenticate_User(ctx, client, device_code)
		if err == nil && access_token != "" {
			break
		}

		select {
		case <- ctx.Done():
			return
		case <- time.After(5 * time.Second):
		}
	}
	fmt.Println("Authenticated to Microsoft")

	fmt.Println("Authenticating to Xbox Live")
	xbl_token, err := Authenticate_Xbox_Live(ctx, client, access_token)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println("Authenticated to Xbox Live")

	fmt.Println("Otaining XSTS token")
	xsts_token, hash, err := Obtain_XSTS_Token(ctx, client, xbl_token)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println("XSTS token obtained")

	fmt.Println("Authenticating to Minecraft")
	minecraft_token, err := Authenticate_Minecraft(ctx, client, xsts_token, hash)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println("Authenticated to Minecraft")
	fmt.Printf("Token: %v\n", minecraft_token)
}

func Authorize_Device(ctx context.Context, client *http.Client) (string, error) {
	var payload AuthorizeDeviceRequest
	payload.ClientId = client_id
	payload.Scope = "XboxLive.signin"

	url_format_request := url.Values{}
	url_format_request.Set("client_id", payload.ClientId)
	url_format_request.Set("scope", payload.Scope)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://login.microsoftonline.com/consumers/oauth2/v2.0/devicecode",
		strings.NewReader(url_format_request.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Received status: %v", resp.StatusCode)
	}

	var response AuthorizeDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	fmt.Printf("Log in at: %s with the code: %s\n", response.VerificationUri, response.UserCode)

	return response.DeviceCode, nil
}

func Authenticate_User(ctx context.Context, client *http.Client, device_code string) (string, error){
	var payload AuthenticateUserRequest
	payload.GrantType = "urn:ietf:params:oauth:grant-type:device_code"
	payload.ClientId = client_id
	payload.DeviceCode = device_code

	url_format_request := url.Values{}
	url_format_request.Set("grant_type", payload.GrantType)
	url_format_request.Set("client_id", payload.ClientId)
	url_format_request.Set("device_code", payload.DeviceCode)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://login.microsoftonline.com/consumers/oauth2/v2.0/token",
		strings.NewReader(url_format_request.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Received status: %v", resp.StatusCode)
	}

	var response AuthenticateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.AccessToken, nil
}

func Authenticate_Xbox_Live(ctx context.Context, client *http.Client, access_token string) (string, error) {
	var payload XBOXAuthRequest
	payload.Properties.AuthMethod = "RPS"
	payload.Properties.SiteName = "user.auth.xboxlive.com"
	payload.Properties.RpsTicket = "d=" + access_token
	payload.RelyingParty = "http://auth.xboxlive.com"
	payload.TokenType = "JWT"

	json_request, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://user.auth.xboxlive.com/user/authenticate",
		bytes.NewBuffer(json_request),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Received status: %v", resp.StatusCode)
	}

	var response XBOXAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Token, nil
}

func Obtain_XSTS_Token(ctx context.Context, client *http.Client, xbl_token string) (string, string, error) {
	var payload XSTSRequest
	payload.Properties.SandboxID = "RETAIL"
	payload.Properties.UserTokens = []string{xbl_token}
	payload.RelyingParty = "rp://api.minecraftservices.com/"
	payload.TokenType = "JWT"

	json_request, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://xsts.auth.xboxlive.com/xsts/authorize",
		bytes.NewReader(json_request),
	)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Received status: %v", resp.StatusCode)
	}

	var response XSTSResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", err
	}

	if len(response.DisplayClaims.XUI) <= 0 {
		return "", "", fmt.Errorf("Expected UHS but did not receive")
	}

	return response.Token, response.DisplayClaims.XUI[0].UHS, nil
}

func Authenticate_Minecraft(ctx context.Context, client *http.Client, xsts_token string, hash string) (string, error) {
	var payload MinecraftAuthRequest
	payload.IdentityToken = "XBL3.0 x=" + hash + ";" + xsts_token

	json_request, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.minecraftservices.com/authentication/login_with_xbox",
		bytes.NewReader(json_request),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Received status: %v", resp.StatusCode)
	}

	var response XSTSResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err 
	}

	return response.Token, nil
}
