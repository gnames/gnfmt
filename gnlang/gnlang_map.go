package gnlang

import "strings"

// https://iso639-3.sil.org/code_tables/download_tables

type iso639_3 string

var langMap = func() map[string]iso639_3 {
	res := make(map[string]iso639_3)
	for k, v := range langCodeMap {
		res[strings.ToLower(v)] = k
	}
	return res
}()

var langCodeMap = map[iso639_3]string{
	"afr": "Afrikaans",
	"aka": "Akan",     // Macro-language (includes Twi and Fante).                     |
	"sqi": "Albanian", // Macro-language (includes Gheg, Tosk, etc.).                  |
	"amh": "Amharic",  //                                                                      |
	"ara": "Arabic",
	"arw": "Arawak",    // Also known as Lokono, spoken in South America.               |
	"hye": "Armenian",  //                                                                      |
	"asm": "Assamese",  //                                                                      |
	"aym": "Aymara",    //                                                                      |
	"btq": "Batek",     // A dialect of Temoq, spoken by the Batek people in Malaysia.  |
	"ben": "Bengali",   //                                                                      |
	"bis": "Bislama",   // A creole language spoken in Vanuatu.                         |
	"mya": "Burmese",   // Also known as Myanmar language.                              |
	"khm": "Cambodian", // Also known as Khmer.                                         |
	"cat": "Catalan",   // "Catalán" is the same language.                              |
	"ceb": "Cebuano",   //                                                                      |
	"cha": "Chamorro",  // Spoken in Guam and the Northern Mariana Islands.             |
	"zho": "Chinese",
	"hrv": "Croat", // Also known as Croatian.                                      |
	"ces": "Czech", //                                                                      |
	"dan": "Danish",
	"nld": "Dutch",            //                                                                      |
	"aer": "Eastern Arrernte", // An Indigenous Australian language.                           |
	"bin": "Edo",              // Also known as Bini, spoken in Nigeria.                       |
	"efi": "Efik",             // Spoken in Nigeria.                                           |
	"eng": "English",
	"fij": "Fijian",   //                                                                      |
	"fil": "Filipino", // Also known as Tagalog (see below).                           |
	"fin": "Finnish",  //                                                                      |
	"fra": "French",
	"gla": "Gaelic", // Likely Scottish Gaelic (Irish Gaelic would be "gle").        |
	"grt": "Garo",   // Spoken in India and Bangladesh.                              |
	"deu": "German",
	"ell": "Greek",
	"grn": "Guarani",  // "Guaraní" is the same language; macro-language.              |
	"gyn": "Guyanese", // Likely Guyanese Creole English.                              |
	"hau": "Hausa",
	"haw": "Hawaiian",
	"hin": "Hindi",     //                                                                      |
	"auc": "Huarani",   // Likely Waorani (also spelled Huaorani), spoken in Ecuador.   |
	"hun": "Hungarian", //                                                                      |
	"ibo": "Igbo",      //                                                                      |
	"ilo": "Ilocano",   // Spoken in the Philippines.                                   |
	"ind": "Indonesian",
	"ita": "Italian",
	"jpn": "Japanese",
	"jav": "Javanese", //                                                                      |
	"kaz": "Kazakh",   //                                                                      |
	"kha": "Khasi",    // Spoken in India.                                             |
	"kik": "Kikuyu",   // Also known as Gikuyu, spoken in Kenya.                       |
	"kor": "Korean",
	"kur": "Kurdish", // Macro-language (includes Kurmanji, Sorani, etc.).            |
	"lao": "Lao",     //                                                                      |
	"lug": "Luganda", // Also known as Ganda, spoken in Uganda.                       |
	"mas": "Maasai",  // Spoken in Kenya and Tanzania.                                |
	"mlg": "Malagasy",
	"msa": "Malay",         // Macro-language (includes Indonesian).                        |
	"mal": "Malayalam",     //                                                                      |
	"mlt": "Maltese",       //                                                                      |
	"mri": "Maori",         //                                                                      |
	"arn": "Mapuche",       // Also known as Mapudungun, spoken in Chile and Argentina.     |
	"mar": "Marathi",       //                                                                      |
	"yua": "Mayan",         // Family of languages (e.g., K’iche’: "quc", Yucatec: "yua").  |
	"nah": "Nahuatl",       // Macro-language, spoken in Mexico.                            |
	"nep": "Nepali",        // "Nepalese" is the same language.                             |
	"ntj": "Ngaanyatjarra", // An Indigenous Australian language.                           |
	"nor": "Norwegian",     // Macro-language (includes Bokmål: "nob", Nynorsk: "nno").     |
	"orm": "Oromo",         // Macro-language, spoken in Ethiopia and Kenya.                |
	"pol": "Polish",        //                                                                      |
	"por": "Portuguese",
	"que": "Quechua", // Macro-language (includes many dialects); "Quichua" is a variant name. |
	"ron": "Romanian",
	"rus": "Russian",  //                                                                      |
	"smo": "Samoan",   //                                                                      |
	"san": "Sanskrit", // Classical language of India.                                 |
	"sas": "Sasak",    // Spoken in Indonesia (Lombok).                                |
	"jiv": "Shuar",    // Spoken in Ecuador and Peru.                                  |
	"sin": "Sinhala",  // Also known as Sinhalese, spoken in Sri Lanka.                |
	"slv": "Slovenian",
	"som": "Somali", //                                                                      |
	"spa": "Spanish",
	"sun": "Sundanese", // Spoken in Indonesia (West Java).                             |
	"swa": "Swahili",   // Macro-language, spoken in East Africa.                       |
	"swe": "Swedish",
	"tgl": "Tagalog", // "Filipino" is often used interchangeably.                    |
	"tnq": "Taino",   // Extinct language of the Caribbean (reconstructed).           |
	"tgk": "Tajik",   //                                                                      |
	"tam": "Tamil",   //                                                                      |
	"tha": "Thai",
	"tpi": "Tok Pisin",  // A creole language spoken in Papua New Guinea.                |
	"ton": "Tonga",      // Likely Tongan, spoken in Tonga.                              |
	"top": "Totonac",    // Refers to Totonac languages (e.g., Papantla Totonac: "top"). |
	"tur": "Turkish",    //                                                                      |
	"tuk": "Turkmen",    //                                                                      |
	"urd": "Urdu",       //                                                                      |
	"urh": "Urhobo",     // Spoken in Nigeria.                                           |
	"vie": "Vietnamese", //                                                                      |
	"wol": "Wolof",      // Spoken in Senegal, Gambia, and Mauritania.                   |
	"yor": "Yoruba",     //                                                                      |
	"zul": "Zulu",
}

// type for ISO 3166-1 alpha3
type iso3166 string

var countryMap = func() map[string]iso3166 {
	res := make(map[string]iso3166)
	for k, v := range countryCodeMap {
		res[strings.ToLower(v)] = k
	}
	return res
}()

var countryCodeMap = map[iso3166]string{
	"ABW": "Aruba",
	"AFG": "Afghanistan",
	"AGO": "Angola",
	"AIA": "Anguilla",
	"ALA": "Aland Islands",
	"ALB": "Albania",
	"AND": "Andorra",
	"ARE": "United Arab Emirates",
	"ARG": "Argentina",
	"ARM": "Armenia",
	"ASM": "American Samoa",
	"ATA": "Antarctica",
	"ATF": "French Southern Territories",
	"ATG": "Antigua and Barbuda",
	"AUS": "Australia",
	"AUT": "Austria",
	"AZE": "Azerbaijan",
	"BDI": "Burundi",
	"BEL": "Belgium",
	"BEN": "Benin",
	"BES": "Bonaire, Sint Eustatius and Saba",
	"BFA": "Burkina Faso",
	"BGD": "Bangladesh",
	"BGR": "Bulgaria",
	"BHR": "Bahrain",
	"BHS": "Bahamas",
	"BIH": "Bosnia and Herzegovina",
	"BLM": "Saint Barthélemy",
	"BLR": "Belarus",
	"BLZ": "Belize",
	"BMU": "Bermuda",
	"BOL": "Bolivia, Plurinational State of",
	"BRA": "Brazil",
	"BRB": "Barbados",
	"BRN": "Brunei Darussalam",
	"BTN": "Bhutan",
	"BVT": "Bouvet Island",
	"BWA": "Botswana",
	"CAF": "Central African Republic",
	"CAN": "Canada",
	"CCK": "Cocos (Keeling) Islands",
	"CHE": "Switzerland",
	"CHL": "Chile",
	"CHN": "China",
	"CIV": "Côte d'Ivoire",
	"CMR": "Cameroon",
	"COD": "Congo, Democratic Republic of the",
	"COG": "Congo",
	"COK": "Cook Islands",
	"COL": "Colombia",
	"COM": "Comoros",
	"CPV": "Cabo Verde",
	"CRI": "Costa Rica",
	"CUB": "Cuba",
	"CUW": "Curaçao",
	"CXR": "Christmas Island",
	"CYM": "Cayman Islands",
	"CYP": "Cyprus",
	"CZE": "Czechia",
	"DEU": "Germany",
	"DJI": "Djibouti",
	"DMA": "Dominica",
	"DNK": "Denmark",
	"DOM": "Dominican Republic",
	"DZA": "Algeria",
	"ECU": "Ecuador",
	"EGY": "Egypt",
	"ERI": "Eritrea",
	"ESH": "Western Sahara",
	"ESP": "Spain",
	"EST": "Estonia",
	"ETH": "Ethiopia",
	"FIN": "Finland",
	"FJI": "Fiji",
	"FLK": "Falkland Islands",
	"FRA": "France",
	"FRO": "Faroe Islands",
	"FSM": "Micronesia, Federated States of",
	"GAB": "Gabon",
	"GBR": "United Kingdom of Great Britain and Northern Ireland",
	"GEO": "Georgia",
	"GGY": "Guernsey",
	"GHA": "Ghana",
	"GIB": "Gibraltar",
	"GIN": "Guinea",
	"GLP": "Guadeloupe",
	"GMB": "Gambia",
	"GNB": "Guinea-Bissau",
	"GNQ": "Equatorial Guinea",
	"GRC": "Greece",
	"GRD": "Grenada",
	"GRL": "Greenland",
	"GTM": "Guatemala",
	"GUF": "French Guiana",
	"GUM": "Guam",
	"GUY": "Guyana",
	"HKG": "Hong Kong",
	"HMD": "Heard Island and McDonald Islands",
	"HND": "Honduras",
	"HRV": "Croatia",
	"HTI": "Haiti",
	"HUN": "Hungary",
	"IDN": "Indonesia",
	"IMN": "Isle of Man",
	"IND": "India",
	"IOT": "British Indian Ocean Territory",
	"IRL": "Ireland",
	"IRN": "Iran, Islamic Republic of",
	"IRQ": "Iraq",
	"ISL": "Iceland",
	"ISR": "Israel",
	"ITA": "Italy",
	"JAM": "Jamaica",
	"JEY": "Jersey",
	"JOR": "Jordan",
	"JPN": "Japan",
	"KAZ": "Kazakhstan",
	"KEN": "Kenya",
	"KGZ": "Kyrgyzstan",
	"KHM": "Cambodia",
	"KIR": "Kiribati",
	"KNA": "Saint Kitts and Nevis",
	"KOR": "Korea, Republic of",
	"KWT": "Kuwait",
	"LAO": "Lao People's Democratic Republic",
	"LBN": "Lebanon",
	"LBR": "Liberia",
	"LBY": "Libya",
	"LCA": "Saint Lucia",
	"LIE": "Liechtenstein",
	"LKA": "Sri Lanka",
	"LSO": "Lesotho",
	"LTU": "Lithuania",
	"LUX": "Luxembourg",
	"LVA": "Latvia",
	"MAC": "Macao",
	"MAF": "Saint Martin (French part)",
	"MAR": "Morocco",
	"MCO": "Monaco",
	"MDA": "Moldova, Republic of",
	"MDG": "Madagascar",
	"MDV": "Maldives",
	"MEX": "Mexico",
	"MHL": "Marshall Islands",
	"MKD": "North Macedonia",
	"MLI": "Mali",
	"MLT": "Malta",
	"MMR": "Myanmar",
	"MNE": "Montenegro",
	"MNG": "Mongolia",
	"MNP": "Northern Mariana Islands",
	"MOZ": "Mozambique",
	"MRT": "Mauritania",
	"MSR": "Montserrat",
	"MTQ": "Martinique",
	"MUS": "Mauritius",
	"MWI": "Malawi",
	"MYS": "Malaysia",
	"MYT": "Mayotte",
	"NAM": "Namibia",
	"NCL": "New Caledonia",
	"NER": "Niger",
	"NFK": "Norfolk Island",
	"NGA": "Nigeria",
	"NIC": "Nicaragua",
	"NIU": "Niue",
	"NLD": "Netherlands, Kingdom of the",
	"NOR": "Norway",
	"NPL": "Nepal",
	"NRU": "Nauru",
	"NZL": "New Zealand",
	"OMN": "Oman",
	"PAK": "Pakistan",
	"PAN": "Panama",
	"PCN": "Pitcairn",
	"PER": "Peru",
	"PHL": "Philippines",
	"PLW": "Palau",
	"PNG": "Papua New Guinea",
	"POL": "Poland",
	"PRI": "Puerto Rico",
	"PRK": "Korea, Democratic People's Republic of",
	"PRT": "Portugal",
	"PRY": "Paraguay",
	"PSE": "Palestine, State of",
	"PYF": "French Polynesia",
	"QAT": "Qatar",
	"REU": "Réunion",
	"ROU": "Romania",
	"RUS": "Russian Federation",
	"RWA": "Rwanda",
	"SAU": "Saudi Arabia",
	"SDN": "Sudan",
	"SEN": "Senegal",
	"SGP": "Singapore",
	"SGS": "South Georgia and the South Sandwich Islands",
	"SHN": "Saint Helena, Ascension and Tristan da Cunha",
	"SJM": "Svalbard and Jan Mayen",
	"SLB": "Solomon Islands",
	"SLE": "Sierra Leone",
	"SLV": "El Salvador",
	"SMR": "San Marino",
	"SOM": "Somalia",
	"SPM": "Saint Pierre and Miquelon",
	"SRB": "Serbia",
	"SSD": "South Sudan",
	"STP": "Sao Tome and Principe",
	"SUR": "Suriname",
	"SVK": "Slovakia",
	"SVN": "Slovenia",
	"SWE": "Sweden",
	"SWZ": "Eswatini",
	"SXM": "Sint Maarten (Dutch part)",
	"SYC": "Seychelles",
	"SYR": "Syrian Arab Republic",
	"TCA": "Turks and Caicos Islands",
	"TCD": "Chad",
	"TGO": "Togo",
	"THA": "Thailand",
	"TJK": "Tajikistan",
	"TKL": "Tokelau",
	"TKM": "Turkmenistan",
	"TLS": "Timor-Leste",
	"TON": "Tonga",
	"TTO": "Trinidad and Tobago",
	"TUN": "Tunisia",
	"TUR": "Türkiye",
	"TUV": "Tuvalu",
	"TWN": "Taiwan, Province of China",
	"TZA": "Tanzania, United Republic of",
	"UGA": "Uganda",
	"UKR": "Ukraine",
	"UMI": "United States Minor Outlying Islands",
	"URY": "Uruguay",
	"USA": "United States of America",
	"UZB": "Uzbekistan",
	"VAT": "Holy See",
	"VCT": "Saint Vincent and the Grenadines",
	"VEN": "Venezuela, Bolivarian Republic of",
	"VGB": "Virgin Islands (British)",
	"VIR": "Virgin Islands (U.S.)",
	"VNM": "Viet Nam",
	"VUT": "Vanuatu",
	"WLF": "Wallis and Futuna",
	"WSM": "Samoa",
	"YEM": "Yemen",
	"ZAF": "South Africa",
	"ZMB": "Zambia",
	"ZWE": "Zimbabwe",
}

var countryAbbr = map[string]string{
	"bolivia":       "bolivia, plurinational state of",
	"car":           "central african republic",
	"drc":           "congo, democratic republic of the",
	"gb":            "united kingdom of great britain and northern ireland",
	"great britain": "united kingdom of great britain and northern ireland",
	"iran":          "iran, islamic republic of",
	"kndr":          "korea, democratic people's republic of",
	"laos":          "lao people's democratic republic",
	"malvinas":      "falkland islands",
	"micronesia":    "micronesia, federated states of",
	"moldova":       "moldova, republic of",
	"netherlands":   "netherlands, kingdom of the",
	"russia":        "russian federation",
	"taiwan":        "taiwan, province of china",
	"uae":           "united arab emirates",
	"united states": "united states of america",
	"usa":           "united states of america",
	"vatican":       "holy see",
	"vietnam":       "viet nam",
}
