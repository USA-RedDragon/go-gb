package cartridge

type NewPublisher struct {
	Name string
}

//nolint:goconst,gochecknoglobals
var newPublishers = map[string]NewPublisher{
	"00": {"None"},
	"01": {"Nintendo Research & Development 1"},
	"08": {"Capcom"},
	"13": {"Electronic Arts"},
	"18": {"Hudson Soft"},
	"19": {"B-AI"},
	"20": {"KSS"},
	"22": {"Planning Office WADA"},
	"24": {"PCM Complete"},
	"25": {"San-X"},
	"28": {"Kemco"},
	"29": {"SETA Corporation"},
	"30": {"Viacom"},
	"31": {"Nintendo"},
	"32": {"Bandai"},
	"33": {"Ocean Software"},
	"34": {"Konami"},
	"35": {"HectorSoft"},
	"37": {"Taito"},
	"38": {"Hudson Soft"},
	"39": {"Banpresto"},
	"41": {"Ubi Soft"},
	"42": {"Atlus"},
	"44": {"Malibu Interactive"},
	"46": {"Angel"},
	"47": {"Bullet-Proof Software"},
	"49": {"Irem"},
	"50": {"Absolute"},
	"51": {"Acclaim Entertainment"},
	"52": {"Activision"},
	"53": {"Sammy USA Corporation"},
	"54": {"Konami"},
	"55": {"Hi Tech Expressions"},
	"56": {"LJN"},
	"57": {"Matchbox"},
	"58": {"Mattel"},
	"59": {"Milton Bradley Company"},
	"60": {"Titus Interactive"},
	"61": {"Virgin Games Ltd."},
	"64": {"Lucasfilm Games"},
	"67": {"Ocean Software"},
	"69": {"Electronic Arts"},
	"70": {"Infogrames"},
	"71": {"Interplay Entertainment"},
	"72": {"Broderbund"},
	"73": {"Sculptured Software"},
	"75": {"The Sales Curve Limited"},
	"78": {"THQ"},
	"79": {"Accolade"},
	"80": {"Misawa Entertainment"},
	"83": {"lozc"},
	"86": {"Tokuma Shoten"},
	"87": {"Tsukuda Original"},
	"91": {"Chunsoft Co."},
	"92": {"Video System"},
	"93": {"Ocean Software"},
	"95": {"Varie"},
	"96": {"Yonezawa"},
	"97": {"Kaneko"},
	"99": {"Pack-In-Video"},
	"9H": {"Bottom Up"},
	"A4": {"Konami (Yu-Gi-Oh!)"},
	"BL": {"MTO"},
	"DK": {"Kodansha"},
}

func GetNewPublisher(code string) string {
	pub, ok := newPublishers[code]
	if ok {
		return pub.Name
	}
	return "Unknown Publisher"
}
