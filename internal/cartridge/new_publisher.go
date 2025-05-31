package cartridge

type NewPublisher string

const (
	NewPublisherNone                            NewPublisher = "00"
	NewPublisherNintendoResearchAndDevelopment1 NewPublisher = "01"
	NewPublisherCapcom                          NewPublisher = "08"
	NewPublisherElectronicArts                  NewPublisher = "13"
	NewPublisherHudsonSoft                      NewPublisher = "18"
	NewPublisherBAI                             NewPublisher = "19"
	NewPublisherKSS                             NewPublisher = "20"
	NewPublisherPlanningOfficeWADA              NewPublisher = "22"
	NewPublisherPCMComplete                     NewPublisher = "24"
	NewPublisherSanX                            NewPublisher = "25"
	NewPublisherKemco                           NewPublisher = "28"
	NewPublisherSETACorporation                 NewPublisher = "29"
	NewPublisherViacom                          NewPublisher = "30"
	NewPublisherNintendo                        NewPublisher = "31"
	NewPublisherBandai                          NewPublisher = "32"
	NewPublisherOceanSoftware                   NewPublisher = "33"
	NewPublisherKonami                          NewPublisher = "34"
	NewPublisherHectorSoft                      NewPublisher = "35"
	NewPublisherTaito                           NewPublisher = "37"
	NewPublisherHudsonSoft2                     NewPublisher = "38"
	NewPublisherBanpresto                       NewPublisher = "39"
	NewPublisherUbiSoft                         NewPublisher = "41"
	NewPublisherAtlus                           NewPublisher = "42"
	NewPublisherMalibuInteractive               NewPublisher = "44"
	NewPublisherAngel                           NewPublisher = "46"
	NewPublisherBulletProofSoftware             NewPublisher = "47"
	NewPublisherIrem                            NewPublisher = "49"
	NewPublisherAbsolute                        NewPublisher = "50"
	NewPublisherAcclaimEntertainment            NewPublisher = "51"
	NewPublisherActivision                      NewPublisher = "52"
	NewPublisherSammyUSACorporation             NewPublisher = "53"
	NewPublisherKonami2                         NewPublisher = "54"
	NewPublisherHiTechExpressions               NewPublisher = "55"
	NewPublisherLJN                             NewPublisher = "56"
	NewPublisherMatchbox                        NewPublisher = "57"
	NewPublisherMattel                          NewPublisher = "58"
	NewPublisherMiltonBradleyCompany            NewPublisher = "59"
	NewPublisherTitusInteractive                NewPublisher = "60"
	NewPublisherVirginGamesLtd                  NewPublisher = "61"
	NewPublisherLucasfilmGames                  NewPublisher = "64"
	NewPublisherOceanSoftware2                  NewPublisher = "67"
	NewPublisherElectronicArts2                 NewPublisher = "69"
	NewPublisherInfogrames                      NewPublisher = "70"
	NewPublisherInterplayEntertainment          NewPublisher = "71"
	NewPublisherBroderbund                      NewPublisher = "72"
	NewPublisherSculpturedSoftware              NewPublisher = "73"
	NewPublisherTheSalesCurveLimited            NewPublisher = "75"
	NewPublisherTHQ                             NewPublisher = "78"
	NewPublisherAccolade                        NewPublisher = "79"
	NewPublisherMisawaEntertainment             NewPublisher = "80"
	NewPublisherlozc                            NewPublisher = "83"
	NewPublisherTokumaShoten                    NewPublisher = "86"
	NewPublisherTsukudaOriginal                 NewPublisher = "87"
	NewPublisherChunsoftCo                      NewPublisher = "91"
	NewPublisherVideoSystem                     NewPublisher = "92"
	NewPublisherOceanSoftware3                  NewPublisher = "93"
	NewPublisherVarie                           NewPublisher = "95"
	NewPublisherYonezawa                        NewPublisher = "96"
	NewPublisherKaneko                          NewPublisher = "97"
	NewPublisherPackInVideo                     NewPublisher = "99"
	NewPublisherBottomUp                        NewPublisher = "9H"
	NewPublisherKonamiYuGiOh                    NewPublisher = "A4"
	NewPublisherMTO                             NewPublisher = "BL"
	NewPublisherKodansha                        NewPublisher = "DK"
)

func (p NewPublisher) String() string {
	switch p {
	case NewPublisherNone:
		return "None"
	case NewPublisherNintendoResearchAndDevelopment1:
		return "Nintendo Research and Development 1"
	case NewPublisherCapcom:
		return "Capcom"
	case NewPublisherElectronicArts:
		return "Electronic Arts"
	case NewPublisherHudsonSoft:
		return "Hudson Soft"
	case NewPublisherBAI:
		return "B-AI"
	case NewPublisherKSS:
		return "KSS"
	case NewPublisherPlanningOfficeWADA:
		return "Planning Office WADA"
	case NewPublisherPCMComplete:
		return "PCM Complete"
	case NewPublisherSanX:
		return "San-X"
	case NewPublisherKemco:
		return "Kemco"
	case NewPublisherSETACorporation:
		return "SETA Corporation"
	case NewPublisherViacom:
		return "Viacom"
	case NewPublisherNintendo:
		return "Nintendo"
	case NewPublisherBandai:
		return "Bandai"
	case NewPublisherOceanSoftware:
		return "Ocean Software"
	case NewPublisherKonami:
		return "Konami"
	case NewPublisherHectorSoft:
		return "Hector Soft"
	case NewPublisherTaito:
		return "Taito"
	case NewPublisherHudsonSoft2:
		return "Hudson Soft"
	case NewPublisherBanpresto:
		return "Banpresto"
	case NewPublisherUbiSoft:
		return "Ubi Soft"
	case NewPublisherAtlus:
		return "Atlus"
	case NewPublisherMalibuInteractive:
		return "Malibu Interactive"
	case NewPublisherAngel:
		return "Angel"
	case NewPublisherBulletProofSoftware:
		return "Bullet-Proof Software"
	case NewPublisherIrem:
		return "Irem"
	case NewPublisherAbsolute:
		return "Absolute"
	case NewPublisherAcclaimEntertainment:
		return "Acclaim Entertainment"
	case NewPublisherActivision:
		return "Activision"
	case NewPublisherSammyUSACorporation:
		return "Sammy USA Corporation"
	case NewPublisherKonami2:
		return "Konami"
	case NewPublisherHiTechExpressions:
		return "Hi Tech Expressions"
	case NewPublisherLJN:
		return "LJN"
	case NewPublisherMatchbox:
		return "Matchbox"
	case NewPublisherMattel:
		return "Mattel"
	case NewPublisherMiltonBradleyCompany:
		return "Milton Bradley Company"
	case NewPublisherTitusInteractive:
		return "Titus Interactive"
	case NewPublisherVirginGamesLtd:
		return "Virgin Games Ltd."
	case NewPublisherLucasfilmGames:
		return "Lucasfilm Games"
	case NewPublisherOceanSoftware2:
		return "Ocean Software"
	case NewPublisherElectronicArts2:
		return "Electronic Arts"
	case NewPublisherInfogrames:
		return "Infogrames"
	case NewPublisherInterplayEntertainment:
		return "Interplay Entertainment"
	case NewPublisherBroderbund:
		return "Broderbund"
	case NewPublisherSculpturedSoftware:
		return "Sculptured Software"
	case NewPublisherTheSalesCurveLimited:
		return "The Sales Curve Limited"
	case NewPublisherTHQ:
		return "THQ"
	case NewPublisherAccolade:
		return "Accolade"
	case NewPublisherMisawaEntertainment:
		return "Misawa Entertainment"
	case NewPublisherlozc:
		return "lozc"
	case NewPublisherTokumaShoten:
		return "Tokuma Shoten"
	case NewPublisherTsukudaOriginal:
		return "Tsukuda Original"
	case NewPublisherChunsoftCo:
		return "Chunsoft Co."
	case NewPublisherVideoSystem:
		return "Video System"
	case NewPublisherOceanSoftware3:
		return "Ocean Software"
	case NewPublisherVarie:
		return "Varie"
	case NewPublisherYonezawa:
		return "Yonezawa"
	case NewPublisherKaneko:
		return "Kaneko"
	case NewPublisherPackInVideo:
		return "Pack-In-Video"
	case NewPublisherBottomUp:
		return "Bottom-Up"
	case NewPublisherKonamiYuGiOh:
		return "Konami (Yu-Gi-Oh!)"
	case NewPublisherMTO:
		return "MTO"
	case NewPublisherKodansha:
		return "Kodansha"
	default:
		return "Unknown Publisher"
	}
}
