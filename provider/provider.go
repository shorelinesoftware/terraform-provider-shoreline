// Copyright 2021, Shoreline Software Inc.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CanonicalizeUrl(url string) (urlOut string, err error) {
	urlRegexStr := "^(http(s)?://)?(?P<customer>[^\\.]*).(?P<region>[^\\.]*).ap[ip].shoreline-(?P<cluster>[^\\.]*).io(/)?$"
	urlBaseStr := "https://${customer}.${region}.api.shoreline-${cluster}.io"
	urlRegex := regexp.MustCompile(urlRegexStr)
	match := urlRegex.FindStringSubmatch(url)
	if len(match) < 4 {
		return "", fmt.Errorf("URL -- %s -- couldn't be mapped to canonical form -- %s -- (%d)\n", url, CanonicalUrl, len(match))
	}
	for i, name := range urlRegex.SubexpNames() {
		if i > 0 && i <= len(match) {
			urlBaseStr = strings.Replace(urlBaseStr, "${"+name+"}", match[i], 1)
		}
	}
	return urlBaseStr, nil
}

func StringToJsonArray(data string) ([]interface{}, error) {
	//jsObj := map[string]interface{}{}
	jsObj := []interface{}{}
	jsErr := json.Unmarshal([]byte(data), &jsObj)
	return jsObj, jsErr
}

func Base64ToJsonArray(data string) ([]interface{}, error) {
	jsStr, bError := base64.StdEncoding.DecodeString(data)
	if bError != nil {
		return []interface{}{}, bError
	}
	return StringToJsonArray(string(jsStr))
}

func StringToJson(data string) (map[string]interface{}, error) {
	jsObj := map[string]interface{}{}
	jsErr := json.Unmarshal([]byte(data), &jsObj)
	return jsObj, jsErr
}

func Base64ToJson(data string) (map[string]interface{}, error) {
	// NOTE: there are different encoding styles, i.e.:
	//   b64Str := base64.URLEncoding.EncodeToString([]byte(data))
	//   b64Str := base64.StdEncoding.EncodeToString([]byte(data))
	jsStr, bError := base64.StdEncoding.DecodeString(data)
	if bError != nil {
		return map[string]interface{}{}, bError
	}
	return StringToJson(string(jsStr))
}

func OmitJsonObjectFields(val map[string]interface{}, omitList []interface{}) map[string]interface{} {
	appendActionLog(fmt.Sprintf("Omitting (obj) keys: %+v\n", omitList))
	for _, o := range omitList {
		oStr, isStr := o.(string)
		if isStr {
			delete(val, oStr)
		}
	}
	return val
}

func OmitJsonArrayFields(val *[]interface{}, omitList []interface{}) {
	//appendActionLog(fmt.Sprintf("Omitting (array) keys: %+v\n", omitList))
	for idx, elem := range *val {
		eMap, isMap := elem.(map[string]interface{})
		if isMap {
			(*val)[idx] = OmitJsonObjectFields(eMap, omitList)
		}
	}
}

func timeSuffixToIntSec(tv string) int {
	sz := len(tv)
	l := sz - 1
	if sz < 1 {
		return 0
	}
	mult := 0
	switch tv[sz-1] {
	case 's':
		mult = 1
	case 'm':
		mult = 60
	case 'h':
		mult = 60 * 60
	case 'd':
		mult = 60 * 60 * 24
	default:
		l = sz
		mult = 1
	}
	val := tv[:l]
	i, err := strconv.Atoi(val)
	if err != nil {
		return -1
	}
	return i * mult
}

func appendActionLog(msg string) {
	if !DoDebugLog {
		return
	}
	filename := "/tmp/tf-shoreline.log"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		//panic(err)
		return
	}
	defer f.Close()
	id := fmt.Sprintf("gid:%d - ", curGoroutineID())
	if _, err = f.WriteString(id + msg); err != nil {
		//panic(err)
		return
	}
}

func runOpCommand(command string, checkResult bool) (string, error) {
	//var GlobalOpts = CliOpts{}
	//if !LoadAuthConfig(&GlobalOpts) {
	//	return "", fmt.Errorf("Failed to load auth credentials")
	//}
	result := ""
	err := error(nil)
	for r := 0; r <= RetryLimit; r += 1 {
		appendActionLog(fmt.Sprintf("Running OpLang command (retries %d/%d)   ---   command:(( %s ))\n", r, RetryLimit, command))
		result, err = ExecuteOpCommand(&GlobalOpts, command)
		if err == nil {
			if !checkResult {
				return result, err
			}
			err = CheckUpdateResult(result)
			if err == nil {
				return result, err
			} else {
				appendActionLog(fmt.Sprintf("Failed OpLang update (retries %d/%d)   ---   error:(( %s ))\n", r, RetryLimit, err.Error()))
			}
		} else {
			appendActionLog(fmt.Sprintf("Failed OpLang command (retries %d/%d)   ---   error:(( %s ))\n", r, RetryLimit, err.Error()))
		}
	}
	return result, err
}

func runOpCommandToJson(command string) (map[string]interface{}, error) {
	result, err := runOpCommand(command, false)
	if err != nil {
		errOut := fmt.Errorf("Failed to execute op '%s': %s", command, err.Error())
		return nil, errOut
	}

	js := map[string]interface{}{}
	// Parsing/Unmarshalling JSON encoding/json
	err = json.Unmarshal([]byte(result), &js)
	if err != nil {
		errOut := fmt.Errorf("Failed to parse json from command '%s': %s", command, err.Error())
		return nil, errOut
	}
	return js, nil
}

func getNamedObjectFromClassDef(name string, typ string, classJs map[string]interface{}) map[string]interface{} {
	baseKey := fmt.Sprintf("get_%s_class.%s_classes", typ, typ)
	baseArray, isArray := GetNestedValueOrDefault(classJs, ToKeyPath(baseKey), []interface{}{}).([]interface{})
	if isArray {
		for _, curJs := range baseArray {
			extrName := GetNestedValueOrDefault(curJs, ToKeyPath("name"), "")
			if name == extrName {
				return curJs.(map[string]interface{})
			}
		}
	}
	return map[string]interface{}{}
}

func CheckUpdateResult(result string) error {
	js := map[string]interface{}{}
	err := json.Unmarshal([]byte(result), &js)
	if err != nil {
		return fmt.Errorf("Failed parse json result from resource update %s", err.Error())
	}

	actions := []string{"define", "delete", "update"}
	types := []string{"resource", "metric", "alarm", "action", "bot", "file", "notebook"}
	for _, act := range actions {
		for _, typ := range types {
			key := act + "_" + typ
			def := GetNestedValueOrDefault(js, ToKeyPath(key), nil)
			if def != nil {
				errKey := key + ".error.message"
				err := GetNestedValueOrDefault(js, ToKeyPath(errKey), nil)
				if typ == "notebook" && (err == nil || err == "") {
					// have to special-case for notebooks
					err = ""
					errArray := []string{}
					ve, isArray := GetNestedValueOrDefault(js, ToKeyPath(key+".error.validation_errors"), nil).([]interface{})
					if isArray {
						for i, _ := range ve {
							errn, isStr := GetNestedValueOrDefault(js, ToKeyPath(fmt.Sprintf(key+".error.validation_errors.[%d].message", i)), nil).(string)
							if isStr && errn != "" {
								errArray = append(errArray, errn)
							}
							err = strings.Join(errArray, "\n")
						}
					}
				}
				if err == nil || err == "" {
					// success ...
					return nil
				} else {
					errStr := GetInnerErrorStr(CastToString(err))
					// error ...
					return fmt.Errorf("ERROR: %s.\n", errStr)
				}
			}
		}
	}
	return nil
}

// Takes a regex like: "if (?P<if_expr>.*?) then (?P<then_expr>.*?) fi"
// and parses out the named captures (e.g. 'if_expr', 'then_expr')
// into the returned map, with the name as a key, and the match as the value.
func ExtractRegexToMap(expr string, regex string) map[string]interface{} {
	result := map[string]interface{}{}
	re := regexp.MustCompile(regex)
	vals := re.FindStringSubmatch(expr)
	keys := re.SubexpNames()

	vlen := len(vals)
	// skip index 0, which is the entire expression
	for i := 1; i < len(keys); i++ {
		if i < vlen {
			result[keys[i]] = vals[i]
		}
		// XXX else error, capture group didn't match (but need diags passed in)
	}
	return result
}

func ValidateVariableName(name string) bool {
	// match valid variable string names
	matched, _ := regexp.MatchString(`^[_a-zA-Z][_a-zA-Z0-9]*$`, name)
	return matched
}
func ValidateResourceType(name string) bool {
	switch name {
	case "HOST":
		return true
	case "POD":
		return true
	case "CONTAINER":
		return true
	}
	return false
}
func ForceToBool(val interface{}) bool {
	valBool, _ := CastToBoolMaybe(val)
	return valBool
}
func ConvertBoolInt(val interface{}) int {
	valBool, _ := CastToBoolMaybe(val)
	if valBool {
		return 1
	}
	return 0
}

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func ExtractVersionData(verStr string) (major int64, minor int64, patch int64, err *error) {
	verRe := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	if verRe.MatchString(verStr) {
		match := verRe.FindStringSubmatch(verStr)
		return CastToInt(match[1]), CastToInt(match[2]), CastToInt(match[3]), nil
	}
	major, minor, patch = 0, 0, 0
	erf := fmt.Errorf("Couldn't find version number in string '%s'", verStr)
	err = &erf
	return
}

func GetBackendVersionInfo() (build string, version string, major int64, minor int64, patch int64, err *error) {
	err = nil
	build = "unknown"
	version = "unknown"
	major, minor, patch = 0, 0, 0
	// op> backend_version
	// ... "get_backend_version": "{ \"tag\": \"release-1.2.3-stuff\", \"build_date\": \"Wed_May_18_00:07:11_UTC_2022\" }", ...
	js, opErr := runOpCommandToJson("backend_version")
	if opErr != nil {
		return
	}
	build = GetNestedValueOrDefault(js, ToKeyPath("get_backend_version"), "unknown").(string)
	buildJs := CastToObject(build)
	if buildJs == nil {
		// TODO set error
		return
	}
	version = GetNestedValueOrDefault(buildJs, ToKeyPath("tag"), "unknown").(string)
	if strings.HasPrefix(version, "stable") || strings.HasPrefix(version, "release") {
		// parse out '\d+\.\d+.\d+' suffix
		major, minor, patch, err = ExtractVersionData(version)
	} else {
		// dev build, special case
		major, minor, patch = 9999, 9999, 9999
	}
	return
}

func dataSourceVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	build, version, major, minor, patch, err := GetBackendVersionInfo()
	if err != nil {
		diags = diag.Errorf("Failed to read backend_version: %s", (*err).Error())
		return diags
	}

	d.Set("build_info", CastToString(build))
	d.Set("version", CastToString(version))
	d.Set("major", major)
	d.Set("minor", minor)
	d.Set("patch", patch)
	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			//DataSourcesMap: map[string]*schema.Resource{
			//	"shoreline_datasource": dataSourceShoreline(),
			//},
			ResourcesMap: map[string]*schema.Resource{
				"shoreline_action":          resourceShorelineObject(ObjectConfigJsonStr, "action"),
				"shoreline_alarm":           resourceShorelineObject(ObjectConfigJsonStr, "alarm"),
				"shoreline_bot":             resourceShorelineObject(ObjectConfigJsonStr, "bot"),
				"shoreline_circuit_breaker": resourceShorelineObject(ObjectConfigJsonStr, "circuit_breaker"),
				"shoreline_metric":          resourceShorelineObject(ObjectConfigJsonStr, "metric"),
				"shoreline_resource":        resourceShorelineObject(ObjectConfigJsonStr, "resource"),
				"shoreline_file":            resourceShorelineObject(ObjectConfigJsonStr, "file"),
				"shoreline_notebook":        resourceShorelineObject(ObjectConfigJsonStr, "notebook"),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"shoreline_version": &schema.Resource{
					ReadContext: dataSourceVersionRead,
					Schema: map[string]*schema.Schema{
						"build_info": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"major": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"minor": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"patch": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			Schema: map[string]*schema.Schema{
				"url": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						if !ValidateApiUrl(val.(string)) {
							errs = append(errs, fmt.Errorf("%q must be of the form %s,\n but got: %s", key, CanonicalUrl, val.(string)))
						}
						return
					},
					DefaultFunc: schema.EnvDefaultFunc("SHORELINE_URL", nil),
					Description: "Customer-specific URL for the Shoreline API server. It should be of the form ```" + CanonicalUrl + "``` .",
				},
				"token": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SHORELINE_TOKEN", nil),
					Description: "Customer/user-specific authorization token for the Shoreline API server. May be provided via `SHORELINE_TOKEN` env variable.",
				},
				"retries": {
					Type:        schema.TypeInt,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("SHORELINE_RETRIES", nil),
					Description: "Number of retries for API calls, in case of e.g. transient network failures.",
				},
				"debug": {
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("SHORELINE_DEBUG", nil),
					Description: "Debug logging to `/tmp/tf-shoreline.log`.",
				},
				"min_version": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Minimum version required on the Shoreline backend (API server).",
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
}

func configure(version string, p *schema.Provider) func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		AuthUrl = d.Get("url").(string)
		token, hasToken := d.GetOk("token")

		var diags diag.Diagnostics = nil

		canonUrl, err := CanonicalizeUrl(AuthUrl)
		if err != nil {
			//return nil, diag.Errorf("Couldn't map URL to canonical form.\n" + err.Error())
			diags = diag.FromErr(err)
			diags[0].Severity = diag.Warning
			canonUrl = AuthUrl
			appendActionLog(fmt.Sprintf("Non-standard url: %s -- to -- %s\n", AuthUrl, canonUrl))
		} else {
			appendActionLog(fmt.Sprintf("Mapped url: %s -- to -- %s\n", AuthUrl, canonUrl))
		}

		if hasToken {
			SetAuth(&GlobalOpts, canonUrl, token.(string))
		} else {
			GlobalOpts.Url = canonUrl
			if !LoadAuthConfig(&GlobalOpts) {
				return nil, diag.Errorf("Failed to load auth credentials file.\n" + GetManualAuthMessage(&GlobalOpts))
			}
			if !selectAuth(&GlobalOpts, canonUrl) {
				return nil, diag.Errorf("Failed to load auth credentials for %s\n"+GetManualAuthMessage(&GlobalOpts), canonUrl)
			}
		}

		retries, hasRetry := d.GetOk("retries")
		if hasRetry {
			RetryLimit = retries.(int)
		} else {
			RetryLimit = 0
		}

		debugLog, hasDebugLog := d.GetOk("debug")
		if hasDebugLog {
			DoDebugLog = debugLog.(bool)
		}

		minVer, hasMinVer := d.GetOk("min_version")
		if hasMinVer {
			var diags diag.Diagnostics
			_, version, major, minor, patch, err := GetBackendVersionInfo()
			if err != nil {
				diags = diag.Errorf("Failed to read backend_version: %s", (*err).Error())
				return nil, diags
			}
			minMajor, minMinor, minPatch, err := ExtractVersionData(minVer.(string))
			if err != nil {
				diags = diag.Errorf("Failed to parse min_version: %s", (*err).Error())
				return nil, diags
			}
			wantVer := []int64{minMajor, minMinor, minPatch}
			haveVer := []int64{major, minor, patch}
			verOk := true
			for i, want := range wantVer {
				if haveVer[i] < want {
					verOk = false
					break
				}
				if haveVer[i] > want {
					break
				}
			}
			if !verOk {
				diags = diag.Errorf("Backend version '%s' (%d, %d, %d) does not meet min_version: '%s' (%d, %d, %d)", version, major, minor, patch, minVer.(string), minMajor, minMinor, minPatch)
				return nil, diags
			}
		}

		return &apiClient{}, diags
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

var ObjectConfigJsonStr = `
{
	"action": {
		"attributes": {
			"type":                    { "type": "string",   "computed": true, "value": "ACTION" },
			"name":                    { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"command":                 { "type": "command",  "required": true, "primary": true },
			"description":             { "type": "string",   "optional": true },
			"enabled":                 { "type": "intbool",  "optional": true, "default": false },
			"params":                  { "type": "string[]", "optional": true },
			"resource_tags_to_export": { "type": "string[]", "optional": true },
			"res_env_var":             { "type": "string",   "optional": true },
			"resource_query":          { "type": "command",  "optional": true },
			"shell":                   { "type": "string",   "optional": true },
			"timeout":                 { "type": "unsigned", "optional": true, "default": 60000 },
			"file_deps":               { "type": "string[]", "optional": true },
			"start_short_template":    { "type": "string",   "optional": true, "step": "start_step_class.short_template" },
			"start_long_template":     { "type": "string",   "optional": true, "step": "start_step_class.long_template" },
			"start_title_template":    { "type": "string",   "optional": true, "step": "start_step_class.title_template" },
			"error_short_template":    { "type": "string",   "optional": true, "step": "error_step_class.short_template" },
			"error_long_template":     { "type": "string",   "optional": true, "step": "error_step_class.long_template" },
			"error_title_template":    { "type": "string",   "optional": true, "step": "error_step_class.title_template" },
			"complete_short_template": { "type": "string",   "optional": true, "step": "complete_step_class.short_template" },
			"complete_long_template":  { "type": "string",   "optional": true, "step": "complete_step_class.long_template" },
			"complete_title_template": { "type": "string",   "optional": true, "step": "complete_step_class.title_template" },
			"#user":                   { "type": "string",   "optional": true }
		}
	},

	"alarm": {
		"attributes": {
			"type":                   { "type": "string",   "computed": true, "value": "ALARM" },
			"name":                   { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"fire_query":             { "type": "command",  "required": true, "primary": true },
			"clear_query":            { "type": "command",  "optional": true },
			"description":            { "type": "string",   "optional": true },
			"resource_query":         { "type": "command",  "optional": true },
			"enabled":                { "type": "intbool",  "optional": true, "default": false },
			"mute_query":             { "type": "string",   "optional": true },
			"resolve_short_template": { "type": "string",   "optional": true, "step": "clear_step_class.short_template" },
			"resolve_long_template":  { "type": "string",   "optional": true, "step": "clear_step_class.long_template" },
			"resolve_title_template": { "type": "string",   "optional": true, "step": "clear_step_class.title_template", "suppress_null_regex": "^cleared \\w*$" },
			"fire_short_template":    { "type": "string",   "optional": true, "step": "fire_step_class.short_template" },
			"fire_long_template":     { "type": "string",   "optional": true, "step": "fire_step_class.long_template" },
			"fire_title_template":    { "type": "string",   "optional": true, "step": "fire_step_class.title_template", "suppress_null_regex": "^fired \\w*$" },
			"condition_type":         { "type": "command",  "optional": true, "step": "condition_details.[0].condition_type" },
			"condition_value":        { "type": "string",   "optional": true, "step": "condition_details.[0].condition_value" },
			"metric_name":            { "type": "string",   "optional": true, "step": "condition_details.[0].metric_name" },
			"raise_for":              { "type": "command",  "optional": true, "step": "condition_details.[0].raise_for", "default": "local" },
			"check_interval_sec":     { "type": "command",  "optional": true, "step": "check_interval_sec", "default": 1 },
			"compile_eligible":       { "type": "bool",     "optional": true, "step": "compile_eligible", "default": true },
			"resource_type":          { "type": "resource", "optional": true, "step": "resource_type" },
			"family":                 { "type": "command",  "optional": true, "step": "config_data.family", "default": "custom" }
		}
	},

	"bot": {
		"attributes": {
			"type":             { "type": "string",   "computed": true, "value": "BOT" },
			"name":             { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"command":          { "type": "command",  "required": true, "primary": true,
				"compound_in": "^\\s*if\\s*(?P<alarm_statement>.*?)\\s*then\\s*(?P<action_statement>.*?)\\s*fi\\s*$",
				"compound_out": "if ${alarm_statement} then ${action_statement} fi"
			},
			"description":      { "type": "string",   "optional": true },
			"enabled":          { "type": "intbool",  "optional": true, "default": false },
			"family":           { "type": "command",  "optional": true, "step": "config_data.family", "default": "custom" },
			"action_statement": { "type": "command",  "internal": true },
			"alarm_statement":  { "type": "command",  "internal": true }
		}
	},

	"circuit_breaker": {
		"attributes": {
			"type": { "type": "string", "computed": true, "value": "CIRCUIT_BREAKER" },
			"name": { "type": "label", "required": true, "forcenew": true, "skip": true },
			"command": { "type": "command", "required": true, "primary": true, "forcenew": true,
				"compound_in": "^\\s*(?P<resource_query>.+)\\s*\\|\\s*(?P<action_name>[a-zA-Z_][a-zA-Z_]*)\\s*$",
				"compound_out": "${resource_query} | ${action_name}"
			},
			"breaker_type": { "type": "string", "required": true },
			"hard_limit": { "type": "int", "required": true },
			"soft_limit": { "type": "int", "optional": true, "default": -1 },
			"duration": { "type": "time_s", "required": true },
			"fail_over": { "type": "string", "required": true },
			"enabled": { "type": "bool", "optional": true, "default": false },
			"action_name": { "type": "command", "internal": true },
			"resource_query": { "type": "command", "internal": true }
		}
	},

	"metric": {
		"attributes": {
			"type":           { "type": "string",   "computed": true, "value": "METRIC" },
			"name":           { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"value":          { "type": "command",   "required": true, "primary": true, "alias_out": "val" },
			"description":    { "type": "string",   "optional": true },
			"units":          { "type": "string",   "optional": true },
			"resource_type":  { "type": "resource", "optional": true },
			"#enabled":        { "type": "intbool",  "optional": true, "default": false },
			"#user":           { "type": "string",   "optional": true }
		}
	},

	"resource": {
		"attributes": {
			"type":           { "type": "string",   "computed": true, "value": "RESOURCE" },
			"name":           { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"value":          { "type": "command",  "required": true, "primary": true },
			"description":    { "type": "string",   "optional": true },
			"#units":          { "type": "string",   "optional": true },
			"#resource_type":  { "type": "resource", "optional": true },
			"#user":           { "type": "string",   "optional": true },
			"#read_only":      { "type": "bool",     "optional": true }
		}
	},

	"file": {
		"attributes": {
			"type":             { "type": "string",   "computed": true, "value": "FILE" },
			"name":             { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"destination_path": { "type": "string",   "required": true, "primary": true },
			"description":      { "type": "string",   "optional": true },
			"resource_query":   { "type": "string",   "required": true },
			"enabled":          { "type": "intbool",  "optional": true, "default": false },
			"input_file":       { "type": "string",   "required": true, "skip": true, "not_stored": true },
			"file_data":        { "type": "string",   "computed": true },
			"file_length":      { "type": "int",      "computed": true },
			"checksum":         { "type": "string",   "computed": true },
			"md5":              { "type": "string",   "optional": true, "proxy": "file_length,checksum,file_data" },
			"#resource_type":    { "type": "resource", "optional": true },
			"#last_modified_timestamp": { "type": "string",   "optional": true }
		}
	},

	"notebook": {
		"attributes": {
			"type":                   { "type": "string",   "computed": true, "value": "NOTEBOOK" },
			"name":                   { "type": "label",    "required": true, "forcenew": true, "skip": true },
			"data":                   { "type": "b64json",  "required": true, "step": ".", "primary": true,
				                          "omit": {"cells": "dynamic_cell_fields", ".": "dynamic_fields" },
				                          "cast": { "params": "string[]", "params_values": "string[]" }
			                          },
			"description":            { "type": "string",   "optional": true },
			"#enabled":                { "type": "intbool",  "optional": true, "default": false },
			"timeout_ms":             { "type": "unsigned", "optional": true, "default": 60000 }
		}
	},

	"docs": {
		"objects": {
				"action":   "A command that can be run.\n\nSee the Shoreline [Actions Documentation](https://docs.shoreline.io/actions) for more info.",
				"alarm":    "A condition that triggers Alerts or Actions.\n\nSee the Shoreline [Alarms Documentation](https://docs.shoreline.io/alarms) for more info.",
				"bot":      "An automation that ties an Action to an Alert.\n\nSee the Shoreline [Bots Documentation](https://docs.shoreline.io/bots) for more info.",
				"circuit_breaker":      "An automatic rate limit on actions.\n\nSee the Shoreline [CircuitBreakers Documentation](https://docs.shoreline.io/circuit_breakers) for more info.",
				"metric":   "A periodic measurement of a system property.\n\nSee the Shoreline [Metrics Documentation](https://docs.shoreline.io/metrics) for more info.",
				"resource": "A server or compute resource in the system (e.g. host, pod, container).\n\nSee the Shoreline [Resources Documentation](https://docs.shoreline.io/platform/resources) for more info.",
				"file":     "A datafile that is automatically copied/distributed to defined Resources.\n\nSee the Shoreline [OpCp Documentation](https://docs.shoreline.io/op/commands/cp) for more info.",
				"notebook":     "An interactive notebook of Op commands and user documentation .\n\nSee the Shoreline [Notebook Documentation](https://docs.shoreline.io/ui/notebooks) for more info."
		},

		"attributes": {
				"type":                    "The type of object (i.e., Alarm, Action, Bot, Metric, Resource, or File).",
				"cells":                   "The data cells inside a notebook.",
				"check_interval":          "Interval (in seconds) between Alarm evaluations.",
				"checksum":                "Cryptographic hash (e.g. md5) of a File Resource.",
				"clear_query":             "The Alarm's resolution condition.",
				"command":                 "A specific action to run.",
				"compile_eligible":        "If the Alarm can be effectively optimized.",
				"complete_short_template": "The short description of the Action's completion.",
				"complete_long_template":  "The long description of the Action's completion.",
				"complete_title_template": "UI title of the Action's completion.",
				"condition_type":          "Kind of check in an Alarm (e.g. above or below) vs a threshold for a Metric.",
				"condition_value":         "Switching value (threshold) for a Metric in an Alarm.",
				"data":                    "The downloaded (JSON) representation of a Notebook.",
				"description":             "A user-friendly explanation of an object.",
				"destination_path":        "Target location for a copied distributed File object.  See [Op: cp](https://docs.shoreline.io/op/commands/cp).",
				"enabled":                 "If the object is currently enabled or disabled.",
				"error_short_template":    "The short description of the Action's error condition.",
				"error_long_template":     "The long description of the Action's error condition.",
				"error_title_template":    "UI title of the Action's error condition.",
				"family":                  "General class for an Action or Bot (e.g., custom, standard, metric, or system check).",
				"file_data":               "Internal representation of a distributed File object's data (computed).",
				"file_length":             "Length, in bytes, of a distributed File object (computed)",
				"fire_query":              "The Alarm's trigger condition.",
				"fire_short_template":     "The short description of the Alarm's triggering condition.",
				"fire_long_template":      "The long description of the Alarm's triggering condition.",
				"fire_title_template":     "UI title of the Alarm's triggering condition.",
				"input_file":              "The local source of a distributed File object.",
				"metric_name":             "The Alarm's triggering Metric.",
				"mute_query":              "The Alarm's mute condition.",
				"md5":                     "The md5 checksum of a file, e.g. filemd5(\"${path.module}/data/example-file.txt\")",
				"name":                    "The name of the object (must be unique).",
				"params":                  "Named variables to pass to an object (e.g. an Action).",
				"raise_for":               "Where an Alarm is raised (e.g., local to a resource, or global to the system).",
				"res_env_var":             "Result environment variable ... an environment variable used to output values through.",
				"resolve_short_template":  "The short description of the Alarm's resolution.",
				"resolve_long_template":   "The long description of the Alarm's resolution.",
				"resolve_title_template":  "UI title of the Alarm's' resolution.",
				"resource_query":          "A set of Resources (e.g. host, pod, container), optionally filtered on tags or dynamic conditions.",
				"shell":                   "The commandline shell to use (e.g. /bin/sh).",
				"start_short_template":    "The short description when starting the Action.",
				"start_long_template":     "The long description when starting the Action.",
				"start_title_template":    "UI title of the start of the Action.",
				"timeout":                 "Maximum time to wait, in milliseconds.",
				"file_deps":               "file object dependencies.",
				"units":                   "Units of a Metric (e.g., bytes, blocks, packets, percent).",
				"value":                   "The Op statement that defines a Metric or Resource."
		}
	}
}
`

// old bot
//			"#action_statement": { "type": "command",  "required": true, "primary": true },
//			"#alarm_statement":  { "type": "command",  "required": true }

func resourceShorelineObject(configJsStr string, key string) *schema.Resource {
	params := map[string]*schema.Schema{}

	objects := map[string]interface{}{}
	// Parsing/Unmarshalling JSON encoding/json
	err := json.Unmarshal([]byte(configJsStr), &objects)
	if err != nil {
		WriteMsg("WARNING: Failed to parse JSON config from resourceShorelineObject().\n")
		return nil
	}
	object := GetNestedValueOrDefault(objects, ToKeyPath(key), nil)
	if object == nil {
		WriteMsg("WARNING: Failed to parse JSON config from resourceShorelineObject(%s).\n", key)
		return nil
	}
	attributes := GetNestedValueOrDefault(object, ToKeyPath("attributes"), map[string]interface{}{}).(map[string]interface{})
	primary := "name"
	for k, attrs := range attributes {
		if strings.HasPrefix(k, "#") {
			continue
		}

		// internal objects, i.e. components of compound fields
		internal := GetNestedValueOrDefault(attrs, ToKeyPath("internal"), false).(bool)
		if internal {
			continue
		}

		sch := &schema.Schema{}

		description := CastToString(GetNestedValueOrDefault(objects, ToKeyPath("docs.attributes."+k), ""))
		sch.Description = description

		attrMap := attrs.(map[string]interface{})
		typ := GetNestedValueOrDefault(attrMap, ToKeyPath("type"), "string")
		switch typ {
		case "command":
			sch.Type = schema.TypeString
			sch.DiffSuppressFunc = func(k, old, nu string, d *schema.ResourceData) bool {
				// ignore whitespace changes in command strings
				if strings.ReplaceAll(old, " ", "") == strings.ReplaceAll(nu, " ", "") {
					return true
				}
				return false
			}
		case "time_s":
			sch.Type = schema.TypeString
			// special case for notebook cells
			sch.DiffSuppressFunc = func(k, old, nu string, d *schema.ResourceData) bool {
				oldT := timeSuffixToIntSec(old)
				nuT := timeSuffixToIntSec(nu)
				appendActionLog(fmt.Sprintf("time_s DiffSuppressFunc: diffing (%s)=(%d) and (%s)=(%d)", old, oldT, nu, nuT))
				if oldT == nuT {
					return true
				}
				return false
			}
		case "b64json":
			sch.Type = schema.TypeString
			// special case for notebook JSON data
			sch.DiffSuppressFunc = func(k, old, nu string, d *schema.ResourceData) bool {
				if old == "" && nu == "" {
					return true
				}
				oldJs, oldErr := StringToJson(old)
				nuJs, nuErr := StringToJson(nu)
				if oldErr != nil || nuErr != nil {
					return false
				}
				// special case top-level notebook "enabled" which may be returned by old backends
				delete(nuJs, "enabled")
				delete(oldJs, "enabled")
				if reflect.DeepEqual(oldJs, nuJs) {
					return true
				}
				return false
			}
		case "string":
			sch.Type = schema.TypeString
		case "string[]":
			sch.Type = schema.TypeList
			sch.Elem = &schema.Schema{
				Type: schema.TypeString,
			}
		case "bool":
			sch.Type = schema.TypeBool
		case "intbool":
			// special handling to/from backend ("1"/"0")
			sch.Type = schema.TypeBool
		case "float":
			sch.Type = schema.TypeFloat
		case "int":
			sch.Type = schema.TypeInt
		case "unsigned":
			sch.Type = schema.TypeInt
			// non-negative validator
			sch.ValidateFunc = func(val interface{}, key string) (warns []string, errs []error) {
				v := val.(int)
				if v <= 0 {
					errs = append(errs, fmt.Errorf("%q must be > 0, got: %d", key, v))
				}
				return
			}
		case "label":
			sch.Type = schema.TypeString
			// TODO ValidateVariableName()
		case "resource":
			sch.Type = schema.TypeString
			// TODO ValidateResourceType() "^(HOST|POD|CONTAINER)$"
		}
		sch.Optional = GetNestedValueOrDefault(attrMap, ToKeyPath("optional"), false).(bool)
		sch.Required = GetNestedValueOrDefault(attrMap, ToKeyPath("required"), false).(bool)
		sch.Computed = GetNestedValueOrDefault(attrMap, ToKeyPath("computed"), false).(bool)
		sch.ForceNew = GetNestedValueOrDefault(attrMap, ToKeyPath("forcenew"), false).(bool)
		//WriteMsg("WARNING: JSON config from resourceShorelineObject(%s) %s.Optional = %+v.\n", key, k, sch.Optional)
		//WriteMsg("WARNING: JSON config from resourceShorelineObject(%s) %s.Required = %+v.\n", key, k, sch.Required)
		//WriteMsg("WARNING: JSON config from resourceShorelineObject(%s) %s.Computed = %+v.\n", key, k, sch.Computed)
		//defowlt := GetNestedValueOrDefault(attrMap, ToKeyPath("value"), nil)
		defowlt := GetNestedValueOrDefault(attrMap, ToKeyPath("default"), nil)
		if defowlt != nil {
			sch.Default = defowlt
		}
		suppressNullDiffRegex, isStr := GetNestedValueOrDefault(attrMap, ToKeyPath("suppress_null_regex"), nil).(string)
		if isStr {
			sch.DiffSuppressFunc = func(k, old, nu string, d *schema.ResourceData) bool {
				if old == nu {
					return true
				}
				if nu == "" {
					matched, _ := regexp.MatchString(suppressNullDiffRegex, old)
					if matched {
						return true
					}
				}
				return false
			}
		}

		// NOTE: This actually messes up the file objects. Need a suppress function that's just for acceptance test comparisions.
		//notStored, isBool := GetNestedValueOrDefault(attrMap, ToKeyPath("not_stored"), nil).(bool)
		//if isBool && notStored {
		//	sch.DiffSuppressFunc = func(k, old, nu string, d *schema.ResourceData) bool {
		//		//appendActionLog(fmt.Sprintf("Not Stored Value: '%s': '%s' -- vs -- '%s'\n", k, old, nu))
		//		if old == nu || nu == "" || old == "" {
		//			return true
		//		}
		//		return false
		//	}
		//}

		if GetNestedValueOrDefault(attrMap, ToKeyPath("primary"), false).(bool) {
			primary = k
		}
		params[k] = sch

	}

	objDescription := CastToString(GetNestedValueOrDefault(objects, ToKeyPath("docs.objects."+key), ""))

	return &schema.Resource{
		Description: "Shoreline " + key + ". " + objDescription,

		CreateContext: resourceShorelineObjectCreate(key, primary, attributes),
		ReadContext:   resourceShorelineObjectRead(key, attributes),
		UpdateContext: resourceShorelineObjectUpdate(key, attributes),
		DeleteContext: resourceShorelineObjectDelete(key),
		Importer:      &schema.ResourceImporter{State: schema.ImportStatePassthrough},

		Schema: params,
	}

}

func attrValueString(typ string, key string, val interface{}, attrs map[string]interface{}) string {
	strVal := ""
	attrTyp := GetNestedValueOrDefault(attrs, ToKeyPath(key+".type"), "string").(string)
	switch attrTyp {
	case "command":
		strVal = fmt.Sprintf("%s", val)
	case "time_s":
		strVal = fmt.Sprintf("%s", val)
	case "b64json":
		jsStr, isStr := val.(string)
		if !isStr {
			jsStr = ""
		}
		strVal = fmt.Sprintf("\"%s\"", base64.StdEncoding.EncodeToString([]byte(jsStr)))
	case "string":
		strVal = fmt.Sprintf("\"%s\"", val)
	case "string[]":
		valArr, isArr := val.([]interface{})
		listStr := ""
		sep := ""
		if isArr {
			for _, v := range valArr {
				listStr = listStr + fmt.Sprintf("%s\"%s\"", sep, v)
				sep = ", "
			}
		}
		return "[ " + listStr + " ]"
	case "bool":
		if ForceToBool(val) {
			strVal = fmt.Sprintf("true")
		} else {
			strVal = fmt.Sprintf("false")
		}
	case "intbool": // special handling to/from backend ("1"/"0")
		strVal = fmt.Sprintf("%d", ConvertBoolInt(val))
	case "float":
		strVal = fmt.Sprintf("%f", val)
	case "int":
		strVal = fmt.Sprintf("%d", val)
	case "unsigned":
		strVal = fmt.Sprintf("%d", val)
	case "label":
		strVal = fmt.Sprintf("\"%s\"", val)
	case "resource":
		strVal = fmt.Sprintf("\"%s\"", val)
	}
	return strVal
}

func setFieldViaOp(typ string, attrs map[string]interface{}, name string, key string, val interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	valStr := attrValueString(typ, key, val, attrs)
	appendActionLog(fmt.Sprintf("Setting %s field: '%s'.'%s' :: %+v\n", typ, name, key, val))

	op := fmt.Sprintf("%s.%s = %s", name, key, valStr)

	// TODO Let alias to be a list of fallbacks for versioning,
	//   or have alternate ObjectConfigJsonStr based on backend version,
	//   or let backend return ObjectConfigJsonStr to use.
	alias, isStr := GetNestedValueOrDefault(attrs, ToKeyPath(key+".alias_out"), nil).(string)
	if isStr {
		op = fmt.Sprintf("%s.%s = %s", name, alias, valStr)
	}

	appendActionLog(fmt.Sprintf("Setting with op statement... '%s'\n", op))
	result, err := runOpCommand(op, true)
	if err != nil {
		diags = diag.Errorf("Failed to set %s %s.%s: %s", typ, name, key, err.Error())
		appendActionLog(fmt.Sprintf("Failed to set %s %s.%s: %s\nval: (( %+v ))\nop-statement: %s\n", typ, name, key, val, err.Error(), op))
		return diags
	}
	err = CheckUpdateResult(result)
	if err != nil {
		diags = diag.Errorf("Failed to update %s %s.%s: %s", typ, name, key, err.Error())
		return diags
	}
	return nil
}

func resourceShorelineObjectSetFields(typ string, attrs map[string]interface{}, ctx context.Context, d *schema.ResourceData, meta interface{}, doDiff bool, isCreate bool) diag.Diagnostics {
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	// valid-variable-name check (and non-null)
	//appendActionLog(fmt.Sprintf("RESOURCE TYPE IS: %s\n", typ))

	forcedChangeKeys := map[string]bool{}
	forcedChangeVals := map[string]interface{}{}

	for key, _ := range attrs {
		proxy := GetNestedValueOrDefault(attrs, ToKeyPath(key+".proxy"), "").(string)
		if proxy != "" {
			proxyKeys := strings.Split(proxy, ",")
			for _, k := range proxyKeys {
				forcedChangeKeys[k] = true
			}
		}
	}

	if typ == "file" {
		infile, exists := d.GetOk("input_file")
		if exists {
			base64Data, ok, fileSize, md5sum := FileToBase64(infile.(string))
			if ok {
				appendActionLog(fmt.Sprintf("file_length is %d (%v)\n", int(fileSize), fileSize))
				if forcedChangeKeys["file_data"] {
					forcedChangeVals["file_length"] = int(fileSize)
					forcedChangeVals["checksum"] = md5sum
					forcedChangeVals["file_data"] = base64Data
				}
				d.Set("file_length", int(fileSize))
				d.Set("checksum", md5sum)
				d.Set("file_data", base64Data)
			} else {
				diags = diag.Errorf("Failed to read file object %s", infile)
				return diags
			}
		}
	}

	// TODO handle intbool type (aside from enable)
	writeEnable := false
	enableVal := false
	anyChange := false
	for key, _ := range attrs {

		skip := GetNestedValueOrDefault(attrs, ToKeyPath(key+".skip"), false).(bool)
		if skip {
			continue
		}

		internal := GetNestedValueOrDefault(attrs, ToKeyPath(key+".internal"), false).(bool)
		if internal {
			continue
		}
		proxy := GetNestedValueOrDefault(attrs, ToKeyPath(key+".proxy"), "").(string)
		if proxy != "" {
			continue
		}

		forceSet := false
		// CS-336 workaround: Force explicit set of action_statement/alarm_statement to patch quoting issue
		_, botEnvDefined := os.LookupEnv("BOT_SKIP_PATCH")
		isPrimary := GetNestedValueOrDefault(attrs, ToKeyPath(key+".primary"), false).(bool)
		if isCreate && isPrimary && typ == "bot" {
			if botEnvDefined {
				appendActionLog(fmt.Sprintf("Bot skipping post-ctor set: %s: '%s'.'%s' HasChange(%v)\n", typ, name, key, d.HasChange(key)))
				// primary value is set on creation, and redundant set currently triggers an issue with bots
				continue
			} else {
				appendActionLog(fmt.Sprintf("Bot running post-ctor set: %s: '%s'.'%s'  HasChange(%v)\n", typ, name, key, d.HasChange(key)))
				forceSet = true
			}
		}

		val, exists := d.GetOk(key)
		// NOTE: Terraform reports !exists when a value is explicitly supplied, but matches the 'default'
		if !exists && !d.HasChange(key) && !forceSet && !forcedChangeKeys[key] {
			appendActionLog(fmt.Sprintf("FieldDoesNotExist: %s: '%s'.'%s' val(%v) HasChange(%v), forceSet(%v)\n", typ, name, key, val, d.HasChange(key), forceSet))
			continue
		}

		// Because OpLang auto-toggles some objects to "disabled" on *any* property change,
		// we have to restore the value as needed.
		if key == "enabled" {
			enableVal, _ = CastToBoolMaybe(val)
			if d.HasChange(key) || !doDiff {
				writeEnable = true
			}
			appendActionLog(fmt.Sprintf("CheckEnableState: %s: '%s' write(%v) val(%v) change(%v) hasChange:(%v) doDiff(%v)\n", typ, name, writeEnable, enableVal, anyChange, d.HasChange(key), doDiff))
			continue
		}
		if doDiff && !d.HasChange(key) && !forcedChangeKeys[key] {
			continue
		}

		compoundRegex, isStr := GetNestedValueOrDefault(attrs, ToKeyPath(key+".compound_in"), nil).(string)
		if isStr {
			curMap := ExtractRegexToMap(CastToString(val), compoundRegex)
			appendActionLog(fmt.Sprintf("CompoundSet: %s: '%s'.'%s' map(%v) from (( %v ))\n", typ, name, key, curMap, val))

			unchanged := map[string]bool{}
			if doDiff {
				old, _ := d.GetChange(key)
				oldMap := ExtractRegexToMap(CastToString(old), compoundRegex)
				for k, v := range oldMap {
					nu, exists := curMap[k]
					if exists && v == nu {
						unchanged[k] = true
					}
				}
			}

			for k, v := range curMap {
				_, skip := unchanged[k]
				if skip {
					continue
				}
				result := setFieldViaOp(typ, attrs, name, k, v)
				if result != nil {
					return result
				}
			}
			anyChange = true
			continue
		}

		result := diag.Diagnostics(nil)
		if forcedChangeKeys[key] {
			result = setFieldViaOp(typ, attrs, name, key, forcedChangeVals[key])
		} else {
			result = setFieldViaOp(typ, attrs, name, key, val)
		}
		if result != nil {
			return result
		}
		anyChange = true
	}

	appendActionLog(fmt.Sprintf("EnableState: %s: '%s' write(%v) val(%v) change(%v)\n", typ, name, writeEnable, enableVal, anyChange))
	// Enabled is automatically toggled to "false" by oplang on any other attribute change.
	// So, it requires special handling.
	if writeEnable || (enableVal && anyChange) {
		act := "enable"
		if !enableVal {
			act = "disable"
		}
		op := fmt.Sprintf("%s %s", act, name)
		appendActionLog(fmt.Sprintf("EnableState: %s: '%s' Op:'%s'\n", typ, name, op))
		result, err := runOpCommand(op, true)
		if err != nil {
			diags = diag.Errorf("Failed to %s (1) %s: %s", act, typ, err.Error())
			return diags
		}
		err = CheckUpdateResult(result)
		if err != nil {
			diags = diag.Errorf("Failed to %s (2) %s: %s", act, typ, err.Error())
			return diags
		}
	}
	return nil
}

func resourceShorelineObjectCreate(typ string, primary string, attrs map[string]interface{}) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// use the meta value to retrieve your client from the provider configure method
		// client := meta.(*apiClient)

		var diags diag.Diagnostics
		name := d.Get("name").(string)
		primaryVal := d.Get(primary)
		idFromAPI := name
		appendActionLog(fmt.Sprintf("Creating %s: '%s' (%v) :: %+v\n", typ, idFromAPI, name, d))

		primaryValStr := attrValueString(typ, primary, primaryVal, attrs)
		//appendActionLog(fmt.Sprintf("primaryValStr is ((( %+v )))\n", primaryValStr))
		//op := fmt.Sprintf("%s %s = \"%s\"", typ, name, primaryVal)
		op := fmt.Sprintf("%s %s = %s", typ, name, primaryValStr)
		//if typ == "bot" {
		//	// special handling for BOT creation statement "bot <name>=
		//	action := d.Get("action_statement").(string)
		//	alarm := d.Get("alarm_statement").(string)
		//	op = fmt.Sprintf("%s %s = if %s then %s fi", typ, name, alarm, action)
		//}
		result, err := runOpCommand(op, true)
		if err != nil {
			// TODO check if already exists
			diags = diag.Errorf("Failed to create (1) %s: %s", typ, err.Error())
			return diags
		}
		err = CheckUpdateResult(result)
		if err != nil {
			diags = diag.Errorf("Failed to create (2) %s: %s", typ, err.Error())
			return diags
		}

		diags = resourceShorelineObjectSetFields(typ, attrs, ctx, d, meta, false, true)
		if diags != nil {
			// delete incomplete object
			resourceShorelineObjectDelete(typ)(ctx, d, meta)
			return diags
		}

		// once the object is ok, set the ID to tell terraform it's valid...
		d.SetId(name)
		// update the data in terraform
		return resourceShorelineObjectRead(typ, attrs)(ctx, d, meta)
	}
}

func resourceShorelineObjectRead(typ string, attrs map[string]interface{}) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// use the meta value to retrieve your client from the provider configure method
		// client := meta.(*apiClient)

		var diags diag.Diagnostics
		name := d.Get("name").(string)
		if name == "" {
			// fallback for 'terraform import'
			name = d.Id()
		}
		// valid-variable-name check
		idFromAPI := name
		appendActionLog(fmt.Sprintf("Reading %s: '%s' (%v) :: %+v\n", typ, idFromAPI, name, d))

		op := fmt.Sprintf("list %ss | name = \"%s\"", typ, name)
		js, err := runOpCommandToJson(op)
		if err != nil {
			diags = diag.Errorf("Failed to read %s - %s: %s", typ, name, err.Error())
			return diags
		}

		stepsJs := map[string]interface{}{}

		if typ == "alarm" || typ == "action" || typ == "bot" || typ == "notebook" {
			// extract fields from step objects
			op := fmt.Sprintf("get_%s_class( %s_name = \"%s\" )", typ, typ, name)
			extraJs, err := runOpCommandToJson(op)
			if err != nil {
				diags = diag.Errorf("Failed to read %s - %s: %s", typ, name, err.Error())
				return diags
			}
			stepsJs = getNamedObjectFromClassDef(name, typ, extraJs)
		}

		found := false
		record := map[string]interface{}{}
		symbols, isArray := GetNestedValueOrDefault(js, ToKeyPath("list_type.symbol"), []interface{}{}).([]interface{})
		if isArray {
			for _, s := range symbols {
				sName, isStr := GetNestedValueOrDefault(s, ToKeyPath("attributes.name"), "").(string)
				if isStr && name == sName {
					record = s.(map[string]interface{})
					found = true
				}
			}
		}

		if !found {
			diags = diag.Errorf("Failed to find %s '%s'", typ, name)
			return diags
		}

		for key, attr := range attrs {
			var val interface{}

			if strings.HasPrefix(key, "#") {
				// skip commented fields
				continue
			}

			internal := GetNestedValueOrDefault(attrs, ToKeyPath(key+".internal"), false).(bool)
			if internal {
				continue
			}

			compoundValue, isStr := GetNestedValueOrDefault(attrs, ToKeyPath(key+".compound_out"), nil).(string)
			if isStr {
				fullVal := compoundValue
				re := regexp.MustCompile(`\$\{\w\w*\}`)
				for expr := re.FindString(fullVal); expr != ""; expr = re.FindString(fullVal) {
					l := len(expr)
					varName := expr[2 : l-1]
					valStr := CastToString(GetNestedValueOrDefault(record, ToKeyPath("attributes."+varName), ""))
					fullVal = strings.Replace(fullVal, expr, valStr, -1)
				}
				val = fullVal
			} else {
				stepPath, isStr := GetNestedValueOrDefault(attr, ToKeyPath("step"), nil).(string)
				if isStr {
					if stepPath == "." {
						val = stepsJs
					} else {
						val = GetNestedValueOrDefault(stepsJs, ToKeyPath(stepPath), nil)
					}

					// special handling (notebooks)... field is base64 outgoing, and json incoming
					attrTyp := GetNestedValueOrDefault(attrs, ToKeyPath(key+".type"), "string").(string)
					if attrTyp == "b64json" {
						// handle cast-map, as get_notebook_class() returns objects with some string-wrapped sub-fields
						castMap := GetNestedValueOrDefault(attr, ToKeyPath("cast"), map[string]interface{}{}).(map[string]interface{})
						for castPath, castType := range castMap {
							cur := GetNestedValueOrDefault(val, ToKeyPath(castPath), nil)
							if cur != nil {
								// TODO add additional types as needed
								switch castType {
								case "string[]":
									SetNestedValue(val, ToKeyPath(castPath), CastToArray(cur))
								case "object":
									SetNestedValue(val, ToKeyPath(castPath), CastToObject(cur))
								}
							}
						}
						// handle omit map, and nested deletions, as get_notebook_class() returns objects with dynamic/temporary fields
						omitMap := GetNestedValueOrDefault(attr, ToKeyPath("omit"), map[string]interface{}{}).(map[string]interface{})
						// "." has to be last, or it will wipe out other objects
						omitPaths := []string{}
						hasDot := false
						for omitPath, _ := range omitMap {
							if omitPath != "." {
								omitPaths = append(omitPaths, omitPath)
							} else {
								hasDot = true
							}
						}
						if hasDot {
							omitPaths = append(omitPaths, ".")
						}
						for _, omitPath := range omitPaths {
							omitTag := omitMap[omitPath]
							appendActionLog(fmt.Sprintf("Omit path:'%+v' tag: '%+v'\n", omitPath, omitTag))
							var cur interface{}
							if omitPath == "." {
								cur = val
							} else {
								cur = GetNestedValueOrDefault(val, ToKeyPath(omitPath), nil)
							}
							omitTagStr, isStr := omitTag.(string)
							if !isStr {
								continue
							}
							omitList, isList := GetNestedValueOrDefault(stepsJs, ToKeyPath(omitTagStr), []interface{}{}).([]interface{})
							//appendActionLog(fmt.Sprintf("Omit-list path:'%+v' tag: '%+v' list:'%+v'\n", omitPath, omitTag, omitList))
							if cur != nil && isList {
								if typ == "notebook" && omitPath == "." {
									// NOTE: The top-level object returned by get_notebook_class contains most/all of the object attributes.
									// So remove them from the inner object
									for akey, _ := range attrs {
										omitList = append(omitList, akey)
									}
									omitList = append(omitList, "enabled")
								}
								switch cur.(type) {
								case map[string]interface{}:
									OmitJsonObjectFields(cur.(map[string]interface{}), omitList)
								case []interface{}:
									curArr := cur.([]interface{})
									OmitJsonArrayFields(&curArr, omitList)
								}
								if omitPath == "." {
									val = cur
								} else {
									SetNestedValue(val, ToKeyPath(omitPath), cur)
								}
							}
						}
						b, err := json.Marshal(val)
						if err != nil {
							diags = diag.Errorf("Failed to marshall JSON %s:%s '%s'", typ, key, name)
							return diags
						}
						//val = base64.URLEncoding.EncodeToString(b)
						val = string(b)
					}
				} else {
					val = GetNestedValueOrDefault(record, ToKeyPath("attributes."+key), nil)
				}
				if val == nil {
					continue
				}
			}
			appendActionLog(fmt.Sprintf("Setting %s field: '%s'.'%s' :: %+v\n", typ, name, key, val))
			//typ := GetNestedValueOrDefault(attrs, ToKeyPath(key+".type"), "string").(string)
			//if typ == "string[]" {
			//}
			attrTyp := GetNestedValueOrDefault(attrs, ToKeyPath(key+".type"), "string").(string)
			switch attrTyp {
			case "float":
				d.Set(key, float64(CastToNumber(val)))
			case "int":
				d.Set(key, CastToInt(val))
			case "unsigned":
				d.Set(key, CastToInt(val))
			case "bool":
				d.Set(key, CastToBool(val))
			case "intbool":
				d.Set(key, CastToBool(val))
			case "string[]":
				d.Set(key, CastToArray(val))
			case "string":
				d.Set(key, CastToString(val))
			case "command":
				d.Set(key, CastToString(val))
			case "label":
				d.Set(key, CastToString(val))
			case "time_s":
				d.Set(key, CastToString(val)+"s")
			default:
				d.Set(key, val)
			}
		}
		return diags
	}
}

func resourceShorelineObjectUpdate(typ string, attrs map[string]interface{}) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// use the meta value to retrieve your client from the provider configure method
		// client := meta.(*apiClient)

		var diags diag.Diagnostics
		name := d.Get("name").(string)
		appendActionLog(fmt.Sprintf("Updated object '%s': '%s' :: %+v\n", typ, name, d))

		diags = resourceShorelineObjectSetFields(typ, attrs, ctx, d, meta, true, false)
		if diags != nil {
			// TODO delete incomplete object?
			return diags
		}

		// update the data in terraform
		return resourceShorelineObjectRead(typ, attrs)(ctx, d, meta)
	}
}

func resourceShorelineObjectDelete(typ string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// use the meta value to retrieve your client from the provider configure method
		// client := meta.(*apiClient)

		var diags diag.Diagnostics
		name := d.Get("name").(string)
		appendActionLog(fmt.Sprintf("deleting %s: '%s' :: %+v\n", typ, name, d))

		op := fmt.Sprintf("delete %s", name)
		result, err := runOpCommand(op, true)
		if err != nil {
			// TODO check already exists
			diags = diag.Errorf("Failed to delete %s: %s", typ, err.Error())
		}
		err = CheckUpdateResult(result)
		if err != nil {
			diags = diag.Errorf("Failed to delete %s: %s", typ, err.Error())
			return diags
		}
		return diags
	}
}
