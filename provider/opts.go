package provider

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	zstd "github.com/klauspost/compress/zstd"
	"github.com/spf13/viper"
	"io/ioutil"
	prand "math/rand"
	"os/user"
	"regexp"
	"strings"
	"time"
)

type CliOpts struct {
	Execute     bool
	Verbose     bool
	StdIn       bool
	Quiet       bool
	CfgFile     string
	HasAuth     bool
	AuthChanged bool
	Url         string
	Token       string
}

const CanonicalUrl = "https://<customer>.<region>.api.shoreline-<cluster>.io"

var AuthUrl string
var AuthToken string
var RetryLimit int
var DoDebugLog = false
var GlobalOpts = CliOpts{}

var clientAuth *ClientAuth

var AuthConfig = viper.New()

func GetHomeDir() string {
	user, err := user.Current()
	homeDir := "/"
	if err == nil {
		homeDir = user.HomeDir
	}
	return homeDir
}

func getAuthFilename() string {
	return ".ops_auth.yaml"
}

func getAuthUrls() []string {
	urls := []string{}
	auth := AuthConfig.Get("Auth")
	if auth != nil {
		authArr, isArr := auth.([]interface{})
		if isArr {
			for _, obj := range authArr {
				objMap, isMap := obj.(map[interface{}]interface{})
				if isMap {
					url, uOk := objMap["Url"]
					if uOk {
						urlStr, isStr := url.(string)
						if isStr {
							urls = append(urls, urlStr)
						}
					}
				}
			}
		}
	}
	return urls
}

func LoadAuthConfig(GlobalOpts *CliOpts) bool {
	AuthConfig.SetConfigName(getAuthFilename())
	AuthConfig.SetConfigType("yaml")
	AuthConfig.AddConfigPath(GetHomeDir())
	AuthConfig.ReadInConfig() // ignore errors...

	URL := AuthConfig.GetString("Url")
	TOKEN := AuthConfig.GetString("Token")
	//WriteMsg("*** URL: %s, TOKEN:%.16s\n", URL, TOKEN)
	if URL == "" || TOKEN == "" {
		GlobalOpts.HasAuth = false
	} else {
		GlobalOpts.HasAuth = true
		GlobalOpts.Url = URL
		GlobalOpts.Token = TOKEN
	}
	return GlobalOpts.HasAuth
}

func SetAuth(GlobalOpts *CliOpts, Url string, Token string) {
	AuthUrl = Url
	AuthToken = Token
	// set default
	GlobalOpts.Url = Url
	GlobalOpts.Token = Token
	GlobalOpts.HasAuth = true
	GlobalOpts.AuthChanged = true

	//AddAuthEntry(GlobalOpts, Url, Token, false)
}

func selectAuth(GlobalOpts *CliOpts, toUrl string) bool {
	auth := AuthConfig.Get("Auth")
	if auth == nil {
		WriteMsg("No 'Auth' object in config!\n")
		return false
	}
	authArr, isArr := auth.([]interface{})
	if !isArr {
		WriteMsg("Config 'Auth' object is not an array!\n")
		return false
	}
	for i, obj := range authArr {
		objMap, isMap := obj.(map[interface{}]interface{})
		if !isMap {
			WriteMsg("Config 'Auth' element %d is not an object (%T)!\n", i, obj)
			continue
		}
		urlOb, uOk := objMap["Url"]
		url, uStr := urlOb.(string)
		if !(uOk && uStr) {
			continue
		}
		tokenOb, tOk := objMap["Token"]
		token, tStr := tokenOb.(string)
		if !(tOk && tStr) {
			continue
		}
		if url == toUrl {
			SetAuth(GlobalOpts, toUrl, token)
			return true
		}
	}
	return false
}

func PrintAuthWarning() {
	//WriteMsg("Make sure URL and TOKEN are set in the config file: " + "~/.ops_auth.yaml\n")
	WriteMsg("Missing URL and Authorization Token!\n")
	WriteMsg("Get your customer URL from an administrator and enter the command:\n")
	WriteMsg("   'auth " + CanonicalUrl + "' \n")
}

func GetManualAuthMessage(GlobalOpts *CliOpts) string {
	// handle failure with manual copy/paste message (with URL)
	return fmt.Sprintf("ERROR: Automatic authentication token retrieval failed.\n") +
		fmt.Sprintf("Please retry or manually visit the authentication URL:\n") +
		fmt.Sprintf("\t"+GetTokenAuthUrl(GlobalOpts, true)+"\n")
}

func ValidateApiUrl(url string) bool {
	// NOTE: URLs are in the form -- "https://<customer>.<region>.api.shoreline-<cluster>.io"
	urlRegex := regexp.MustCompile(`^https://\w+\.\w+\.api\.shoreline-\w+\.io$`)
	if !urlRegex.MatchString(url) {
		WriteMsg("ERROR: Invalid URL to auth! (%s)\n", url)
		WriteMsg("It should be of the form: '" + CanonicalUrl + "' \n")
		return false
	}
	return true
}

// randomly generated key
func GetIdempotencyKey() string {
	// generate 128 high-quality random bytes, return as a hex string
	data := make([]byte, 16)
	hexBytes := make([]byte, 32)
	_, err := rand.Read(data[:])
	if err != nil {
		// Fall back to low-quality psuedo-random numbers.
		// NOTE: This is not ideal, but should be rare.
		if viper.GetBool("debug") {
			WriteMsg("WARNING: High-Quality random numbers unavailable. -- (%v)\n", err.Error())
		}
		prand.Seed(time.Now().UnixNano())
		for i := 0; i < 16; i++ {
			data[i] = byte(prand.Intn(256))
		}
	}
	hex.Encode(hexBytes, data)
	hexStr := string(hexBytes)
	return hexStr
}

func GetInnerErrorStr(errStr string) string {
	msgRegex := regexp.MustCompile(`message: \\".*\\"`)
	loc := msgRegex.FindStringIndex(errStr)
	if loc == nil {
		errStr = strings.Replace(errStr, "\\n", "\n", -1)
		return errStr
	}
	innerStr := errStr[loc[0]+11 : loc[1]-2]
	innerStr = strings.Replace(innerStr, "\\\"", "\"", -1)
	innerStr = strings.Replace(innerStr, "\\\\", "\\", -1)
	return GetInnerErrorStr(innerStr)
}

func GetInnerError(err error) string {
	outer := error.Error(err)
	innerStr := GetInnerErrorStr(outer)
	return innerStr
}

func ExecuteOpCommand(GlobalOpts *CliOpts, expr string) (string, error) {
	if !GlobalOpts.HasAuth {
		return "", fmt.Errorf("No valid auth credentials.")
	} else {
		if clientAuth == nil || clientAuth.BaseURL != GlobalOpts.Url || GlobalOpts.AuthChanged {
			// Auth data is persisted, so that we don't have to re-authorize for every command
			clientAuth = NewClientAuth(GlobalOpts.Url, GlobalOpts.Token, GetIdempotencyKey())
			GlobalOpts.AuthChanged = false
		} else {
			// Fresh Idempotency key for every command.
			clientAuth.ApiKey = GetIdempotencyKey()
		}
		fullExpr := expr
		new_client := NewClient(clientAuth)
		//fix this to be resolved input
		ret, error := new_client.Execute(fullExpr, false)
		if error != nil {
			inner := GetInnerError(error)
			return "", fmt.Errorf(inner)

		} else {
			retStr := string(ret)
			return retStr, nil
		}
	}
}

// Returns base64 data, success/failure, file size, md5 checksum.
func FileToBase64(filename string) (string, bool, int64, string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		//logging.WriteMsgColor(Red, "ERROR: Couldn't read input file: %s !\n", filename)
		return "", false, 0, ""
	}

	fileLen := int64(len(raw))

	// Create a writer that caches compressors.
	// For this operation type we supply a nil Reader.
	var encoder, _ = zstd.NewWriter(nil)

	// Compress a buffer.
	// If you have a destination buffer, the allocation in the call can also be eliminated.
	compressed := encoder.EncodeAll(raw, make([]byte, 0, len(raw)))

	encoded := base64.StdEncoding.EncodeToString(compressed)
	md5Sum := fmt.Sprintf("%x", md5.Sum(raw))
	return encoded, true, fileLen, md5Sum
}
