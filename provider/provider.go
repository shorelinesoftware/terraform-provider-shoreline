package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func appendActionLog(msg string) {
	filename := "/tmp/tf-json.log"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		//panic(err)
		return
	}
	defer f.Close()
	if _, err = f.WriteString(msg); err != nil {
		//panic(err)
		return
	}
}

func runOpCommand(command string) (string, error) {
	//var GlobalOpts = CliOpts{}
	//if !LoadAuthConfig(&GlobalOpts) {
	//	return "", fmt.Errorf("Failed to load auth credentials")
	//}
	result, err := ExecuteOpCommand(&GlobalOpts, command, false)
	return result, err
}

func runOpCommandToJson(command string) (map[string]interface{}, error) {
	result, err := runOpCommand(command)
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
	types := []string{"resource", "metric", "alarm", "action", "bot", "file"}
	for _, act := range actions {
		for _, typ := range types {
			key := act + "_" + typ
			def := GetNestedValueOrDefault(js, ToKeyPath(key), nil)
			if def != nil {
				errKey := key + ".error.message"
				err := GetNestedValueOrDefault(js, ToKeyPath(errKey), nil)
				if err == nil {
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
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			//DataSourcesMap: map[string]*schema.Resource{
			//	"shoreline_datasource": dataSourceShoreline(),
			//},
			ResourcesMap: map[string]*schema.Resource{
				//"shoreline_resource": resourceShorelineBasic(),
				//"shoreline_action":   resourceShorelineAction(),
				"shoreline_action":   resourceShorelineObject(ObjectConfigJsonStr, "action"),
				"shoreline_alarm":    resourceShorelineObject(ObjectConfigJsonStr, "alarm"),
				"shoreline_bot":      resourceShorelineObject(ObjectConfigJsonStr, "bot"),
				"shoreline_metric":   resourceShorelineObject(ObjectConfigJsonStr, "metric"),
				"shoreline_resource": resourceShorelineObject(ObjectConfigJsonStr, "resource"),
				"shoreline_file":     resourceShorelineObject(ObjectConfigJsonStr, "file"),
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
				},
				"token": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("SHORELINE_TOKEN", nil),
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

		if hasToken {
			SetAuth(&GlobalOpts, AuthUrl, token.(string))
		} else {
			GlobalOpts.Url = AuthUrl
			if !LoadAuthConfig(&GlobalOpts) {
				return nil, diag.Errorf("Failed to load auth credentials file.\n" + GetManualAuthMessage(&GlobalOpts))
			}
			if !selectAuth(&GlobalOpts, AuthUrl) {
				return nil, diag.Errorf("Failed to load auth credentials for %s\n"+GetManualAuthMessage(&GlobalOpts), AuthUrl)
			}
		}

		return &apiClient{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

var ObjectConfigJsonStr = `
{
	"action": {
		"attributes": {
			"type":                    { "type": "string",   "computed": true, "value": "ACTION" },
			"name":                    { "type": "label",    "required": true, "forcenew": true },
			"command":                 { "type": "command",  "required": true, "primary": true },
			"description":             { "type": "string",   "optional": true },
			"enabled":                 { "type": "intbool",  "optional": true, "default": false },
			"params":                  { "type": "string[]", "optional": true },
			"res_env_var":             { "type": "string",   "optional": true },
			"resource_query":          { "type": "command",  "optional": true },
			"shell":                   { "type": "string",   "optional": true },
			"timeout":                 { "type": "unsigned", "optional": true, "default": 60 },
			"start_short_template":    { "type": "string",   "optional": true, "step": "start_step_class.short_template" },
			"start_title_template":    { "type": "string",   "optional": true, "step": "start_step_class.title_template" },
			"error_short_template":    { "type": "string",   "optional": true, "step": "error_step_class.short_template" },
			"error_title_template":    { "type": "string",   "optional": true, "step": "error_step_class.title_template" },
			"complete_short_template": { "type": "string",   "optional": true, "step": "complete_step_class.short_template" },
			"complete_title_template": { "type": "string",   "optional": true, "step": "complete_step_class.title_template" },
			"#user":                   { "type": "string",   "optional": true }
		}
	},

	"alarm": {
		"attributes": {
			"type":                   { "type": "string",   "computed": true, "value": "ALARM" },
			"name":                   { "type": "label",    "required": true, "forcenew": true },
			"fire_query":             { "type": "command",  "required": true, "primary": true },
			"description":            { "type": "string",   "optional": true },
			"enabled":                { "type": "intbool",  "optional": true, "default": false },
			"clear_query":            { "type": "command",  "optional": true },
			"mute_query":             { "type": "string",   "optional": true },
			"resource_query":         { "type": "command",  "optional": true },
			"resolve_short_template": { "type": "string",   "optional": true, "step": "clear_step_class.short_template" },
			"resolve_title_template": { "type": "string",   "optional": true, "step": "clear_step_class.title_template", "suppress_null_regex": "^cleared \\w*$" },
			"fire_short_template":    { "type": "string",   "optional": true, "step": "fire_step_class.short_template" },
			"fire_title_template":    { "type": "string",   "optional": true, "step": "fire_step_class.title_template", "suppress_null_regex": "^fired \\w*$" },
			"condition_type":         { "type": "command",  "optional": true, "step": "condition_details.[0].condition_type" },
			"condition_value":        { "type": "command",  "optional": true, "step": "condition_details.[0].condition_value" },
			"metric_name":            { "type": "command",  "optional": true, "step": "condition_details.[0].metric_name" },
			"raise_for":              { "type": "command",  "optional": true, "step": "condition_details.[0].raise_for" },
			"check_interval":         { "type": "command",  "optional": true, "step": "check_interval" },
			"resource_type":          { "type": "command",  "optional": true, "step": "resource_type" },
			"family":                 { "type": "command",  "optional": true, "step": "config_data.family", "default": "custom" }
		}
	},

	"bot": {
		"attributes": {
			"type":             { "type": "string",   "computed": true, "value": "BOT" },
			"name":             { "type": "label",    "required": true, "forcenew": true },
			"command":          { "type": "command",  "required": true, "primary": true, 
				"compound_in": "if (?P<alarm_statement>.*?) then (?P<action_statement>.*?) fi", 
				"compound_out": "if ${alarm_statement} then ${action_statement} fi"
			},
			"description":      { "type": "string",   "optional": true },
			"enabled":          { "type": "intbool",  "optional": true, "default": false },
			"family":           { "type": "command",  "optional": true, "step": "config_data.family", "default": "custom" }
		}
	},

	"metric": {
		"attributes": {
			"type":           { "type": "string",   "computed": true, "value": "METRIC" },
			"name":           { "type": "label",    "required": true, "forcenew": true },
			"value":          { "type": "command",   "required": true, "primary": true },
			"description":    { "type": "string",   "optional": true },
			"params":         { "type": "string[]", "optional": true },
			"res_env_var":    { "type": "string",   "optional": true },
			"shell":          { "type": "string",   "optional": true },
			"timeout":        { "type": "unsigned", "optional": true },
			"units":          { "type": "string",   "optional": true },
			"#resource_query": { "type": "command",   "optional": true },
			"#enabled":        { "type": "intbool",  "optional": true, "default": false },
			"#resource_type":  { "type": "resource", "optional": true },
			"#user":           { "type": "string",   "optional": true }
		}
	},

	"resource": {
		"attributes": {
			"type":           { "type": "string",   "computed": true, "value": "RESOURCE" },
			"name":           { "type": "label",    "required": true, "forcenew": true },
			"value":          { "type": "command",  "required": true, "primary": true },
			"description":    { "type": "string",   "optional": true },
			"params":         { "type": "string[]", "optional": true },
			"res_env_var":    { "type": "string",   "optional": true },
			"resource_query": { "type": "command",  "optional": true },
			"shell":          { "type": "string",   "optional": true },
			"timeout":        { "type": "unsigned", "optional": true },
			"units":          { "type": "string",   "optional": true },
			"#resource_type":  { "type": "resource", "optional": true },
			"#user":           { "type": "string",   "optional": true },
			"#read_only":      { "type": "bool",     "optional": true }
		}
	},

	"file": {
		"attributes": {
			"type":             { "type": "string",   "computed": true, "value": "FILE" },
			"name":             { "type": "label",    "required": true, "forcenew": true },
			"destination_path": { "type": "string",   "required": true, "primary": true },
			"description":      { "type": "string",   "optional": true },
			"resource_query":   { "type": "string",   "optional": true },
			"enabled":          { "type": "intbool",  "optional": true, "default": false },
			"input_file":       { "type": "string",   "required": true },
			"file_data":        { "type": "string",   "computed": true },
			"file_length":      { "type": "int",      "computed": true },
			"checksum":         { "type": "string",   "computed": true },
			"#resource_type":    { "type": "resource", "optional": true },
			"#last_modified_timestamp": { "type": "string",   "optional": true }
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
		attrMap := attrs.(map[string]interface{})
		sch := &schema.Schema{}
		typ := GetNestedValueOrDefault(attrMap, ToKeyPath("type"), "string")
		switch typ {
		case "command":
			sch.Type = schema.TypeString
			sch.DiffSuppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
				// ignore whitespace changes in command strings
				if strings.ReplaceAll(old, " ", "") == strings.ReplaceAll(new, " ", "") {
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
		if GetNestedValueOrDefault(attrMap, ToKeyPath("primary"), false).(bool) {
			primary = k
		}
		params[k] = sch
	}

	return &schema.Resource{
		Description: "Shoreline " + key + " resource.",

		CreateContext: resourceShorelineObjectCreate(key, primary, attributes),
		ReadContext:   resourceShorelineObjectRead(key, attributes),
		UpdateContext: resourceShorelineObjectUpdate(key, attributes),
		DeleteContext: resourceShorelineObjectDelete(key),

		Schema: params,
	}

}

func attrValueString(typ string, key string, val interface{}, attrs map[string]interface{}) string {
	strVal := ""
	attrTyp := GetNestedValueOrDefault(attrs, ToKeyPath(key+".type"), "string").(string)
	switch attrTyp {
	case "command":
		strVal = fmt.Sprintf("%s", val)
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
		appendActionLog(fmt.Sprintf("Setting with op statement... '%s'\n", op))
		result, err := runOpCommand(op)
		if err != nil {
			diags = diag.Errorf("Failed to set %s %s.%s: %s", typ, name, key, err.Error())
			return diags
		}
		err = CheckUpdateResult(result)
		if err != nil {
			diags = diag.Errorf("Failed to update %s %s.%s: %s", typ, name, key, err.Error())
			return diags
		}
		return nil
}

func resourceShorelineObjectSetFields(typ string, attrs map[string]interface{}, ctx context.Context, d *schema.ResourceData, meta interface{}, doDiff bool) diag.Diagnostics {
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	// valid-variable-name check (and non-null)
	if typ == "file" {
		infile, exists := d.GetOk("input_file")
		if exists {
			base64Data, ok, fileSize, md5sum := FileToBase64(infile.(string))
			if ok {
				appendActionLog(fmt.Sprintf("file_length is %d (%v)\n", int(fileSize), fileSize))
				d.Set("file_length", int(fileSize))
				d.Set("checksum", md5sum)
				d.Set("file_data", base64Data)
			} else {
				diags = diag.Errorf("Failed to read file object %s", infile)
				return diags
			}
		}
	}

	// TODO handle intbool type
	doEnable := false
	enableVal := false
	for key, _ := range attrs {
		val, exists := d.GetOk(key)
		if !exists {
			continue
		}
		if doDiff && !d.HasChange(key) {
			continue
		}
		if key == "enabled" {
			doEnable = true
			enableVal, _ = CastToBoolMaybe(val)
			continue
		}

		compoundRegex, isStr := GetNestedValueOrDefault(attrs, ToKeyPath(key+".compound_out"), nil).(string)
		if isStr {
			re := regexp.MustCompile(compoundRegex)
			vals := re.FindStringSubmatch(CastToString(val))
			keys := re.SubexpNames()
			for i := 1; i < len(keys); i++ {
				result :=  setFieldViaOp(typ, attrs, name, keys[i], vals[i])
				if result != nil {
					return result
				}
			}
			continue
		}

		result :=  setFieldViaOp(typ, attrs, name, key, val)
		if result != nil {
			return result
		}
	}

	if doEnable {
		act := "enable"
		if !enableVal {
			act = "disable"
		}
		op := fmt.Sprintf("%s %s", act, name)
		result, err := runOpCommand(op)
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
		//op := fmt.Sprintf("%s %s = \"%s\"", typ, name, primaryVal)
		op := fmt.Sprintf("%s %s = %s", typ, name, primaryValStr)
		//if typ == "bot" {
		//	// special handling for BOT creation statement "bot <name>=
		//	action := d.Get("action_statement").(string)
		//	alarm := d.Get("alarm_statement").(string)
		//	op = fmt.Sprintf("%s %s = if %s then %s fi", typ, name, alarm, action)
		//}
		result, err := runOpCommand(op)
		if err != nil {
			// TODO check already exists
			diags = diag.Errorf("Failed to create (1) %s: %s", typ, err.Error())
			return diags
		}
		err = CheckUpdateResult(result)
		if err != nil {
			diags = diag.Errorf("Failed to create (2) %s: %s", typ, err.Error())
			return diags
		}

		diags = resourceShorelineObjectSetFields(typ, attrs, ctx, d, meta, false)
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

		if typ=="alarm" || typ=="action" || typ=="bot" {
			// TODO extract fields from step objects
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

			compoundValue, isStr := GetNestedValueOrDefault(attrs, ToKeyPath(key+".compound_out"), nil).(string)
			if isStr {
				fullVal := compoundValue
				re := regexp.MustCompile(`\$\{\w\w*\}`)
				for expr := re.FindString(fullVal); expr != ""; expr = re.FindString(fullVal) {
					l := len(expr)
					varName := expr[2:l-1]
					valStr := CastToString(GetNestedValueOrDefault(record, ToKeyPath("attributes."+varName), ""))
					fullVal = strings.Replace(fullVal, expr, valStr, -1)
				}
				val = fullVal
			} else {
				stepPath, isStr := GetNestedValueOrDefault(attr, ToKeyPath("step"), nil).(string)
				if isStr {
					val = GetNestedValueOrDefault(stepsJs, ToKeyPath(stepPath), nil)
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
		appendActionLog(fmt.Sprintf("Updated action: '%s' :: %+v\n", name, d))

		diags = resourceShorelineObjectSetFields(typ, attrs, ctx, d, meta, true)
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
		result, err := runOpCommand(op)
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
