// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pluginreceiver

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Plugin is a templated pipeline of receivers and processors
type Plugin struct {
	Title       string      `yaml:"title,omitempty"`
	Template    string      `yaml:"template,omitempty"`
	Version     string      `yaml:"version,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Parameters  []Parameter `yaml:"parameters,omitempty"`
}

// LoadPlugin loads a plugin from a file path
func LoadPlugin(path string) (*Plugin, error) {
	cleanPath := filepath.Clean(path)
	bytes, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var plugin Plugin
	if err := yaml.Unmarshal(bytes, &plugin); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin from yaml: %w", err)
	}

	return &plugin, nil
}

// Render renders the plugin's template as a config
func (p *Plugin) Render(values map[string]interface{}) (*RenderedConfig, error) {
	template, err := template.New(p.Title).Parse(p.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin template: %w", err)
	}

	templateValues := p.ApplyDefaults(values)

	var writer bytes.Buffer
	if err := template.Execute(&writer, templateValues); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	renderedCfg, err := NewRenderedConfig(writer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to create rendered config: %w", err)
	}

	return renderedCfg, nil
}

// ApplyDefaults returns a copy of the values map with parameter defaults applied.
// If a value is already present in the map, it supercedes the default.
func (p *Plugin) ApplyDefaults(values map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, parameter := range p.Parameters {
		if parameter.Default == nil {
			continue
		}

		result[parameter.Name] = parameter.Default
	}

	for key, value := range values {
		result[key] = value
	}

	return result
}

// CheckParameters checks the supplied values against the defined parameters of the plugin
func (p *Plugin) CheckParameters(values map[string]interface{}) error {
	if err := p.checkDefined(values); err != nil {
		return fmt.Errorf("definition failure: %w", err)
	}

	if err := p.checkRequired(values); err != nil {
		return fmt.Errorf("required failure: %w", err)
	}

	if err := p.checkType(values); err != nil {
		return fmt.Errorf("type failure: %w", err)
	}

	if err := p.checkSupported(values); err != nil {
		return fmt.Errorf("supported value failure: %w", err)
	}

	return nil
}

// checkDefined checks if any of the supplied values are not defined by the plugin
func (p *Plugin) checkDefined(values map[string]interface{}) error {
	parameterMap := make(map[string]struct{})
	for _, parameter := range p.Parameters {
		parameterMap[parameter.Name] = struct{}{}
	}

	for key := range values {
		if _, ok := parameterMap[key]; !ok {
			return fmt.Errorf("parameter %s is not defined in plugin", key)
		}
	}

	return nil
}

// checkRequired checks if required values are defined
func (p *Plugin) checkRequired(values map[string]interface{}) error {
	for _, parameter := range p.Parameters {
		_, ok := values[parameter.Name]
		if parameter.Required && !ok {
			return fmt.Errorf("parameter %s is missing but required in plugin", parameter.Name)
		}
	}

	return nil
}

// checkType checks if the values match their parameter type
func (p *Plugin) checkType(values map[string]interface{}) error {
	for _, parameter := range p.Parameters {
		value, ok := values[parameter.Name]
		if !ok {
			continue
		}

		switch parameter.Type {
		case stringType:
			if _, ok := value.(string); !ok {
				return fmt.Errorf("parameter %s must be a string", parameter.Name)
			}
		case stringArrayType:
			raw, ok := value.([]interface{})
			if !ok {
				return fmt.Errorf("parameter %s must be a []string", parameter.Name)
			}

			for _, v := range raw {
				if _, ok := v.(string); !ok {
					return fmt.Errorf("parameter %s: expected string, but got %v", parameter.Name, v)
				}
			}
		case intType:
			if _, ok := value.(int); !ok {
				return fmt.Errorf("parameter %s must be an int", parameter.Name)
			}
		case boolType:
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("parameter %s must be a bool", parameter.Name)
			}
		case timezoneType:
			tzlist := []string{"UTC", "Africa/Abidjan", "Africa/Accra", "Africa/Addis_Ababa", "Africa/Algiers", "Africa/Asmara", "Africa/Bamako", "Africa/Bangui", "Africa/Banjul", "Africa/Bissau", "Africa/Blantyre", "Africa/Brazzaville", "Africa/Bujumbura", "Africa/Cairo", "Africa/Casablanca", "Africa/Ceuta", "Africa/Conakry", "Africa/Dakar", "Africa/Dar_es_Salaam", "Africa/Djibouti", "Africa/Douala", "Africa/El_Aaiun", "Africa/Freetown", "Africa/Gaborone", "Africa/Harare", "Africa/Johannesburg", "Africa/Juba", "Africa/Kampala", "Africa/Khartoum", "Africa/Kigali", "Africa/Kinshasa", "Africa/Lagos", "Africa/Libreville", "Africa/Lome", "Africa/Luanda", "Africa/Lubumbashi", "Africa/Lusaka", "Africa/Malabo", "Africa/Maputo", "Africa/Maseru", "Africa/Mbabane", "Africa/Mogadishu", "Africa/Monrovia", "Africa/Nairobi", "Africa/Ndjamena", "Africa/Niamey", "Africa/Nouakchott", "Africa/Ouagadougou", "Africa/Porto-Novo", "Africa/Sao_Tome", "Africa/Tripoli", "Africa/Tunis", "Africa/Windhoek", "America/Adak", "America/Anchorage", "America/Anguilla", "America/Antigua", "America/Araguaina", "America/Argentina/Buenos_Aires", "America/Argentina/Catamarca", "America/Argentina/Cordoba", "America/Argentina/Jujuy", "America/Argentina/La_Rioja", "America/Argentina/Mendoza", "America/Argentina/Rio_Gallegos", "America/Argentina/Salta", "America/Argentina/San_Juan", "America/Argentina/San_Luis", "America/Argentina/Tucuman", "America/Argentina/Ushuaia", "America/Aruba", "America/Asuncion", "America/Atikokan", "America/Bahia", "America/Bahia_Banderas", "America/Barbados", "America/Belem", "America/Belize", "America/Blanc-Sablon", "America/Boa_Vista", "America/Bogota", "America/Boise", "America/Cambridge_Bay", "America/Campo_Grande", "America/Cancun", "America/Caracas", "America/Cayenne", "America/Cayman", "America/Chicago", "America/Chihuahua", "America/Costa_Rica", "America/Creston", "America/Cuiaba", "America/Curacao", "America/Danmarkshavn", "America/Dawson", "America/Dawson_Creek", "America/Denver", "America/Detroit", "America/Dominica", "America/Edmonton", "America/Eirunepe", "America/El_Salvador", "America/Fort_Nelson", "America/Fortaleza", "America/Glace_Bay", "America/Goose_Bay", "America/Grand_Turk", "America/Grenada", "America/Guadeloupe", "America/Guatemala", "America/Guayaquil", "America/Guyana", "America/Halifax", "America/Havana", "America/Hermosillo", "America/Indiana/Indianapolis", "America/Indiana/Knox", "America/Indiana/Marengo", "America/Indiana/Petersburg", "America/Indiana/Tell_City", "America/Indiana/Vevay", "America/Indiana/Vincennes", "America/Indiana/Winamac", "America/Inuvik", "America/Iqaluit", "America/Jamaica", "America/Juneau", "America/Kentucky/Louisville", "America/Kentucky/Monticello", "America/Kralendijk", "America/La_Paz", "America/Lima", "America/Los_Angeles", "America/Lower_Princes", "America/Maceio", "America/Managua", "America/Manaus", "America/Marigot", "America/Martinique", "America/Matamoros", "America/Mazatlan", "America/Menominee", "America/Merida", "America/Metlakatla", "America/Mexico_City", "America/Miquelon", "America/Moncton", "America/Monterrey", "America/Montevideo", "America/Montserrat", "America/Nassau", "America/New_York", "America/Nipigon", "America/Nome", "America/Noronha", "America/North_Dakota/Beulah", "America/North_Dakota/Center", "America/North_Dakota/New_Salem", "America/Nuuk", "America/Ojinaga", "America/Panama", "America/Pangnirtung", "America/Paramaribo", "America/Phoenix", "America/Port-au-Prince", "America/Port_of_Spain", "America/Porto_Velho", "America/Puerto_Rico", "America/Punta_Arenas", "America/Rainy_River", "America/Rankin_Inlet", "America/Recife", "America/Regina", "America/Resolute", "America/Rio_Branco", "America/Santarem", "America/Santiago", "America/Santo_Domingo", "America/Sao_Paulo", "America/Scoresbysund", "America/Sitka", "America/St_Barthelemy", "America/St_Johns", "America/St_Kitts", "America/St_Lucia", "America/St_Thomas", "America/St_Vincent", "America/Swift_Current", "America/Tegucigalpa", "America/Thule", "America/Thunder_Bay", "America/Tijuana", "America/Toronto", "America/Tortola", "America/Vancouver", "America/Whitehorse", "America/Winnipeg", "America/Yakutat", "America/Yellowknife", "Antarctica/Casey", "Antarctica/Davis", "Antarctica/DumontDUrville", "Antarctica/Macquarie", "Antarctica/Mawson", "Antarctica/McMurdo", "Antarctica/Palmer", "Antarctica/Rothera", "Antarctica/Syowa", "Antarctica/Troll", "Antarctica/Vostok", "Arctic/Longyearbyen", "Asia/Aden", "Asia/Almaty", "Asia/Amman", "Asia/Anadyr", "Asia/Aqtau", "Asia/Aqtobe", "Asia/Ashgabat", "Asia/Atyrau", "Asia/Baghdad", "Asia/Bahrain", "Asia/Baku", "Asia/Bangkok", "Asia/Barnaul", "Asia/Beirut", "Asia/Bishkek", "Asia/Brunei", "Asia/Chita", "Asia/Choibalsan", "Asia/Colombo", "Asia/Damascus", "Asia/Dhaka", "Asia/Dili", "Asia/Dubai", "Asia/Dushanbe", "Asia/Famagusta", "Asia/Gaza", "Asia/Hebron", "Asia/Ho_Chi_Minh", "Asia/Hong_Kong", "Asia/Hovd", "Asia/Irkutsk", "Asia/Jakarta", "Asia/Jayapura", "Asia/Jerusalem", "Asia/Kabul", "Asia/Kamchatka", "Asia/Karachi", "Asia/Kathmandu", "Asia/Khandyga", "Asia/Kolkata", "Asia/Krasnoyarsk", "Asia/Kuala_Lumpur", "Asia/Kuching", "Asia/Kuwait", "Asia/Macau", "Asia/Magadan", "Asia/Makassar", "Asia/Manila", "Asia/Muscat", "Asia/Nicosia", "Asia/Novokuznetsk", "Asia/Novosibirsk", "Asia/Omsk", "Asia/Oral", "Asia/Phnom_Penh", "Asia/Pontianak", "Asia/Pyongyang", "Asia/Qatar", "Asia/Qostanay", "Asia/Qyzylorda", "Asia/Riyadh", "Asia/Sakhalin", "Asia/Samarkand", "Asia/Seoul", "Asia/Shanghai", "Asia/Singapore", "Asia/Srednekolymsk", "Asia/Taipei", "Asia/Tashkent", "Asia/Tbilisi", "Asia/Tehran", "Asia/Thimphu", "Asia/Tokyo", "Asia/Tomsk", "Asia/Ulaanbaatar", "Asia/Urumqi", "Asia/Ust-Nera", "Asia/Vientiane", "Asia/Vladivostok", "Asia/Yakutsk", "Asia/Yangon", "Asia/Yekaterinburg", "Asia/Yerevan", "Atlantic/Azores", "Atlantic/Bermuda", "Atlantic/Canary", "Atlantic/Cape_Verde", "Atlantic/Faroe", "Atlantic/Madeira", "Atlantic/Reykjavik", "Atlantic/South_Georgia", "Atlantic/St_Helena", "Atlantic/Stanley", "Australia/Adelaide", "Australia/Brisbane", "Australia/Broken_Hill", "Australia/Currie", "Australia/Darwin", "Australia/Eucla", "Australia/Hobart", "Australia/Lindeman", "Australia/Lord_Howe", "Australia/Melbourne", "Australia/Perth", "Australia/Sydney", "Europe/Amsterdam", "Europe/Andorra", "Europe/Astrakhan", "Europe/Athens", "Europe/Belgrade", "Europe/Berlin", "Europe/Bratislava", "Europe/Brussels", "Europe/Bucharest", "Europe/Budapest", "Europe/Busingen", "Europe/Chisinau", "Europe/Copenhagen", "Europe/Dublin", "Europe/Gibraltar", "Europe/Guernsey", "Europe/Helsinki", "Europe/Isle_of_Man", "Europe/Istanbul", "Europe/Jersey", "Europe/Kaliningrad", "Europe/Kiev", "Europe/Kirov", "Europe/Lisbon", "Europe/Ljubljana", "Europe/London", "Europe/Luxembourg", "Europe/Madrid", "Europe/Malta", "Europe/Mariehamn", "Europe/Minsk", "Europe/Monaco", "Europe/Moscow", "Europe/Oslo", "Europe/Paris", "Europe/Podgorica", "Europe/Prague", "Europe/Riga", "Europe/Rome", "Europe/Samara", "Europe/San_Marino", "Europe/Sarajevo", "Europe/Saratov", "Europe/Simferopol", "Europe/Skopje", "Europe/Sofia", "Europe/Stockholm", "Europe/Tallinn", "Europe/Tirane", "Europe/Ulyanovsk", "Europe/Uzhgorod", "Europe/Vaduz", "Europe/Vatican", "Europe/Vienna", "Europe/Vilnius", "Europe/Volgograd", "Europe/Warsaw", "Europe/Zagreb", "Europe/Zaporozhye", "Europe/Zurich", "Indian/Antananarivo", "Indian/Chagos", "Indian/Christmas", "Indian/Cocos", "Indian/Comoro", "Indian/Kerguelen", "Indian/Mahe", "Indian/Maldives", "Indian/Mauritius", "Indian/Mayotte", "Indian/Reunion", "Pacific/Apia", "Pacific/Auckland", "Pacific/Bougainville", "Pacific/Chatham", "Pacific/Chuuk", "Pacific/Easter", "Pacific/Efate", "Pacific/Enderbury", "Pacific/Fakaofo", "Pacific/Fiji", "Pacific/Funafuti", "Pacific/Galapagos", "Pacific/Gambier", "Pacific/Guadalcanal", "Pacific/Guam", "Pacific/Honolulu", "Pacific/Kiritimati", "Pacific/Kosrae", "Pacific/Kwajalein", "Pacific/Majuro", "Pacific/Marquesas", "Pacific/Midway", "Pacific/Nauru", "Pacific/Niue", "Pacific/Norfolk", "Pacific/Noumea", "Pacific/Pago_Pago", "Pacific/Palau", "Pacific/Pitcairn", "Pacific/Pohnpei", "Pacific/Port_Moresby", "Pacific/Rarotonga", "Pacific/Saipan", "Pacific/Tahiti", "Pacific/Tarawa", "Pacific/Tongatapu", "Pacific/Wake", "Pacific/Wallis"}
			raw, ok := value.(string)
			if !ok {
				return fmt.Errorf("parameter %s must be a string", parameter.Name)
			}
			for _, tz := range tzlist {
				if raw == tz {
					return nil
				}
			}
			return fmt.Errorf("parameter %s must be a valid timezone", parameter.Name)
		default:
			return fmt.Errorf("unsupported parameter type: %s", parameter.Type)
		}
	}

	return nil
}

// checkSupported checks the values against the plugin's supported values
func (p *Plugin) checkSupported(values map[string]interface{}) error {
OUTER:
	for _, parameter := range p.Parameters {
		if parameter.Supported == nil {
			continue
		}

		value, ok := values[parameter.Name]
		if !ok {
			continue
		}

		for _, v := range parameter.Supported {
			if v == value {
				continue OUTER
			}
		}

		return fmt.Errorf("parameter %s does not match the list of supported values: %v", parameter.Name, parameter.Supported)
	}

	return nil
}

// Parameter is the parameter of plugin
type Parameter struct {
	Name      string        `yaml:"name,omitempty"`
	Type      ParameterType `yaml:"type,omitempty"`
	Default   interface{}   `yaml:"default,omitempty"`
	Supported []interface{} `yaml:"supported,omitempty"`
	Required  bool          `yaml:"required,omitempty"`
}

// ParameterType is the type of a parameter
type ParameterType string

const (
	stringType      ParameterType = "string"
	stringArrayType ParameterType = "[]string"
	boolType        ParameterType = "bool"
	intType         ParameterType = "int"
	timezoneType    ParameterType = "timezone"
)
