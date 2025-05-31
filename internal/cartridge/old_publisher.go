package cartridge

type OldPublisher uint8

const (
	OldPublisherNone                     OldPublisher = 0x00
	OldPublisherNintendo                 OldPublisher = 0x01
	OldPublisherCapcom                   OldPublisher = 0x08
	OldPublisherHOTB                     OldPublisher = 0x09
	OldPublisherJaleco                   OldPublisher = 0x0A
	OldPublisherCoconutsJapan            OldPublisher = 0x0B
	OldPublisherEliteSystems             OldPublisher = 0x0C
	OldPublisherElectronicArts           OldPublisher = 0x13
	OldPublisherHudsonSoft               OldPublisher = 0x18
	OldPublisherITCEntertainment         OldPublisher = 0x19
	OldPublisherYanoman                  OldPublisher = 0x1A
	OldPublisherJapanClary               OldPublisher = 0x1D
	OldPublisherVirginGamesLtd           OldPublisher = 0x1F
	OldPublisherPCMComplete              OldPublisher = 0x24
	OldPublisherSanX                     OldPublisher = 0x25
	OldPublisherKemco                    OldPublisher = 0x28
	OldPublisherSETACorporation          OldPublisher = 0x29
	OldPublisherInfogrames               OldPublisher = 0x30
	OldPublisherNintendo2                OldPublisher = 0x31
	OldPublisherBandai                   OldPublisher = 0x32
	OldPublisherSeeOther                 OldPublisher = 0x33
	OldPublisherKonami                   OldPublisher = 0x34
	OldPublisherHectorSoft               OldPublisher = 0x35
	OldPublisherCapcom2                  OldPublisher = 0x38
	OldPublisherBanpresto                OldPublisher = 0x39
	OldPublisherEntertainmentInteractive OldPublisher = 0x3C
	OldPublisherGremlin                  OldPublisher = 0x3E
	OldPublisherUbiSoft                  OldPublisher = 0x41
	OldPublisherAtlus                    OldPublisher = 0x42
	OldPublisherMalibuInteractive        OldPublisher = 0x44
	OldPublisherAngel                    OldPublisher = 0x46
	OldPublisherSpectrumHoloByte         OldPublisher = 0x47
	OldPublisherIrem                     OldPublisher = 0x49
	OldPublisherVirginGamesLtd2          OldPublisher = 0x4A
	OldPublisherMalibuInteractive2       OldPublisher = 0x4D
	OldPublisherUSGold                   OldPublisher = 0x4F
	OldPublisherAbsolute                 OldPublisher = 0x50
	OldPublisherAcclaimEntertainment     OldPublisher = 0x51
	OldPublisherActivision               OldPublisher = 0x52
	OldPublisherSammyUSACorporation      OldPublisher = 0x53
	OldPublisherGameTek                  OldPublisher = 0x54
	OldPublisherParkPlace                OldPublisher = 0x55
	OldPublisherLJN                      OldPublisher = 0x56
	OldPublisherMatchbox                 OldPublisher = 0x57
	OldPublisherMiltonBradleyCompany     OldPublisher = 0x59
	OldPublisherMindscape                OldPublisher = 0x5A
	OldPublisherRomstar                  OldPublisher = 0x5B
	OldPublisherNaxatSoft                OldPublisher = 0x5C
	OldPublisherTradewest                OldPublisher = 0x5D
	OldPublisherTitusInteractive         OldPublisher = 0x60
	OldPublisherVirginGamesLtd3          OldPublisher = 0x61
	OldPublisherOceanSoftware            OldPublisher = 0x67
	OldPublisherElectronicArts2          OldPublisher = 0x69
	OldPublisherEliteSystems2            OldPublisher = 0x6E
	OldPublisherElectroBrain             OldPublisher = 0x6F
	OldPublisherInfogrames2              OldPublisher = 0x70
	OldPublisherInterplayEntertainment   OldPublisher = 0x71
	OldPublisherBroderbund               OldPublisher = 0x72
	OldPublisherSculpturedSoftware       OldPublisher = 0x73
	OldPublisherTheSalesCurveLimited     OldPublisher = 0x75
	OldPublisherTHQ                      OldPublisher = 0x78
	OldPublisherAccolade                 OldPublisher = 0x79
	OldPublisherTriffixEntertainment     OldPublisher = 0x7A
	OldPublisherMicroProse               OldPublisher = 0x7C
	OldPublisherKemco2                   OldPublisher = 0x7F
	OldPublisherMisawaEntertainment      OldPublisher = 0x80
	OldPublisherLOZCG                    OldPublisher = 0x83
	OldPublisherTokumaShoten             OldPublisher = 0x86
	OldPublisherBulletProofSoftware      OldPublisher = 0x8B
	OldPublisherVicTokaiCorp             OldPublisher = 0x8C
	OldPublisherApeInc                   OldPublisher = 0x8E
	OldPublisherIMax                     OldPublisher = 0x8F
	OldPublisherChunsoftCo               OldPublisher = 0x91
	OldPublisherVideoSystem              OldPublisher = 0x92
	OldPublisherTsubarayaProductions     OldPublisher = 0x93
	OldPublisherVarie                    OldPublisher = 0x95
	OldPublisherYonezawa                 OldPublisher = 0x96
	OldPublisherKemco3                   OldPublisher = 0x97
	OldPublisherArc                      OldPublisher = 0x99
	OldPublisherNihonBussan              OldPublisher = 0x9A
	OldPublisherTecmo                    OldPublisher = 0x9B
	OldPublisherImagineer                OldPublisher = 0x9C
	OldPublisherBanpresto2               OldPublisher = 0x9D
	OldPublisherNova                     OldPublisher = 0x9F
	OldPublisherHoriElectric             OldPublisher = 0xA1
	OldPublisherBandai2                  OldPublisher = 0xA2
	OldPublisherKonami2                  OldPublisher = 0xA4
	OldPublisherKawada                   OldPublisher = 0xA6
	OldPublisherTakara                   OldPublisher = 0xA7
	OldPublisherTechnosJapan             OldPublisher = 0xA9
	OldPublisherBroderbund2              OldPublisher = 0xAA
	OldPublisherToeiAnimation            OldPublisher = 0xAC
	OldPublisherToho                     OldPublisher = 0xAD
	OldPublisherNamco                    OldPublisher = 0xAF
	OldPublisherAcclaimEntertainment2    OldPublisher = 0xB0
	OldPublisherASCIICorporation         OldPublisher = 0xB1
	OldPublisherBandai3                  OldPublisher = 0xB2
	OldPublisherSquareEnix               OldPublisher = 0xB4
	OldPublisherHALLaboratory            OldPublisher = 0xB6
	OldPublisherSNK                      OldPublisher = 0xB7
	OldPublisherPonyCanyon               OldPublisher = 0xB9
	OldPublisherCultureBrain             OldPublisher = 0xBA
	OldPublisherSunsoft                  OldPublisher = 0xBB
	OldPublisherSonyImagesoft            OldPublisher = 0xBD
	OldPublisherSammyCorporation         OldPublisher = 0xBF
	OldPublisherTaito                    OldPublisher = 0xC0
	OldPublisherKemco4                   OldPublisher = 0xC2
	OldPublisherSquare                   OldPublisher = 0xC3
	OldPublisherTokumaShoten2            OldPublisher = 0xC4
	OldPublisherDataEast                 OldPublisher = 0xC5
	OldPublisherTonkinHouse              OldPublisher = 0xC6
	OldPublisherKoei                     OldPublisher = 0xC8
	OldPublisherUFL                      OldPublisher = 0xC9
	OldPublisherUltraGames               OldPublisher = 0xCA
	OldPublisherVAPInc                   OldPublisher = 0xCB
	OldPublisherUseCorporation           OldPublisher = 0xCC
	OldPublisherMeldac                   OldPublisher = 0xCD
	OldPublisherPonyCanyon2              OldPublisher = 0xCE
	OldPublisherAngel2                   OldPublisher = 0xCF
	OldPublisherTaito2                   OldPublisher = 0xD0
	OldPublisherSOFEL                    OldPublisher = 0xD1
	OldPublisherQuest                    OldPublisher = 0xD2
	OldPublisherSigmaEnterprises         OldPublisher = 0xD3
	OldPublisherASKKodanshaCo            OldPublisher = 0xD4
	OldPublisherNaxatSoft2               OldPublisher = 0xD6
	OldPublisherCopyaSystem              OldPublisher = 0xD7
	OldPublisherBanpresto3               OldPublisher = 0xD9
	OldPublisherTomy                     OldPublisher = 0xDA
	OldPublisherLJN2                     OldPublisher = 0xDB
	OldPublisherNipponComputerSystems    OldPublisher = 0xDD
	OldPublisherHumanEnt                 OldPublisher = 0xDE
	OldPublisherAltron                   OldPublisher = 0xDF
	OldPublisherJaleco2                  OldPublisher = 0xE0
	OldPublisherTowaChiki                OldPublisher = 0xE1
	OldPublisherYutaka                   OldPublisher = 0xE2
	OldPublisherVarie2                   OldPublisher = 0xE3
	OldPublisherEpoch                    OldPublisher = 0xE5
	OldPublisherAthena                   OldPublisher = 0xE7
	OldPublisherAsmikAceEntertainment    OldPublisher = 0xE8
	OldPublisherNatsume                  OldPublisher = 0xE9
	OldPublisherKingRecords              OldPublisher = 0xEA
	OldPublisherAtlus2                   OldPublisher = 0xEB
	OldPublisherEpicRecords              OldPublisher = 0xEC
	OldPublisherIGS                      OldPublisher = 0xEE
	OldPublisherAWave                    OldPublisher = 0xF0
	OldPublisherExtremeEntertainment     OldPublisher = 0xF3
	OldPublisherLJN3                     OldPublisher = 0xFF
)

func (p OldPublisher) String() string {
	switch p {
	case OldPublisherNone:
		return "None"
	case OldPublisherNintendo:
		return "Nintendo"
	case OldPublisherCapcom:
		return "Capcom"
	case OldPublisherHOTB:
		return "HOT-B"
	case OldPublisherJaleco:
		return "Jaleco"
	case OldPublisherCoconutsJapan:
		return "Coconuts Japan"
	case OldPublisherEliteSystems:
		return "Elite Systems"
	case OldPublisherElectronicArts:
		return "Electronic Arts"
	case OldPublisherHudsonSoft:
		return "Hudson Soft"
	case OldPublisherITCEntertainment:
		return "ITC Entertainment"
	case OldPublisherYanoman:
		return "Yanoman"
	case OldPublisherJapanClary:
		return "Japan Clary"
	case OldPublisherVirginGamesLtd:
		return "Virgin Games Ltd."
	case OldPublisherPCMComplete:
		return "PCM Complete"
	case OldPublisherSanX:
		return "San-X"
	case OldPublisherKemco:
		return "Kemco"
	case OldPublisherSETACorporation:
		return "SETA Corporation"
	case OldPublisherInfogrames:
		return "Infogrames"
	case OldPublisherNintendo2:
		return "Nintendo"
	case OldPublisherBandai:
		return "Bandai"
	case OldPublisherSeeOther:
		return "See Other"
	case OldPublisherKonami:
		return "Konami"
	case OldPublisherHectorSoft:
		return "Hector Soft"
	case OldPublisherCapcom2:
		return "Capcom"
	case OldPublisherBanpresto:
		return "Banpresto"
	case OldPublisherEntertainmentInteractive:
		return "Entertainment Interactive"
	case OldPublisherGremlin:
		return "Gremlin"
	case OldPublisherUbiSoft:
		return "Ubi Soft"
	case OldPublisherAtlus:
		return "Atlus"
	case OldPublisherMalibuInteractive:
		return "Malibu Interactive"
	case OldPublisherAngel:
		return "Angel"
	case OldPublisherSpectrumHoloByte:
		return "Spectrum HoloByte"
	case OldPublisherIrem:
		return "Irem"
	case OldPublisherVirginGamesLtd2:
		return "Virgin Games Ltd."
	case OldPublisherMalibuInteractive2:
		return "Malibu Interactive"
	case OldPublisherUSGold:
		return "US Gold"
	case OldPublisherAbsolute:
		return "Absolute"
	case OldPublisherAcclaimEntertainment:
		return "Acclaim Entertainment"
	case OldPublisherActivision:
		return "Activision"
	case OldPublisherSammyUSACorporation:
		return "Sammy USA Corporation"
	case OldPublisherGameTek:
		return "GameTek"
	case OldPublisherParkPlace:
		return "Park Place"
	case OldPublisherLJN:
		return "LJN"
	case OldPublisherMatchbox:
		return "Matchbox"
	case OldPublisherMiltonBradleyCompany:
		return "Milton Bradley Company"
	case OldPublisherMindscape:
		return "Mindscape"
	case OldPublisherRomstar:
		return "Romstar"
	case OldPublisherNaxatSoft:
		return "Naxat Soft"
	case OldPublisherTradewest:
		return "Tradewest"
	case OldPublisherTitusInteractive:
		return "Titus Interactive"
	case OldPublisherVirginGamesLtd3:
		return "Virgin Games Ltd."
	case OldPublisherOceanSoftware:
		return "Ocean Software"
	case OldPublisherElectronicArts2:
		return "Electronic Arts"
	case OldPublisherEliteSystems2:
		return "Elite Systems"
	case OldPublisherElectroBrain:
		return "Electro Brain"
	case OldPublisherInfogrames2:
		return "Infogrames"
	case OldPublisherInterplayEntertainment:
		return "Interplay Entertainment"
	case OldPublisherBroderbund:
		return "Broderbund"
	case OldPublisherSculpturedSoftware:
		return "Sculptured Software"
	case OldPublisherTheSalesCurveLimited:
		return "The Sales Curve Limited"
	case OldPublisherTHQ:
		return "THQ"
	case OldPublisherAccolade:
		return "Accolade"
	case OldPublisherTriffixEntertainment:
		return "Triffix Entertainment"
	case OldPublisherMicroProse:
		return "MicroProse"
	case OldPublisherKemco2:
		return "Kemco"
	case OldPublisherMisawaEntertainment:
		return "Misawa Entertainment"
	case OldPublisherLOZCG:
		return "LOZC G."
	case OldPublisherTokumaShoten:
		return "Tokuma Shoten"
	case OldPublisherBulletProofSoftware:
		return "Bullet-Proof Software"
	case OldPublisherVicTokaiCorp:
		return "Vic Tokai Corp."
	case OldPublisherApeInc:
		return "Ape Inc."
	case OldPublisherIMax:
		return "I-Max"
	case OldPublisherChunsoftCo:
		return "Chunsoft Co."
	case OldPublisherVideoSystem:
		return "Video System"
	case OldPublisherTsubarayaProductions:
		return "Tsubaraya Productions"
	case OldPublisherVarie:
		return "Varie"
	case OldPublisherYonezawa:
		return "Yonezawa"
	case OldPublisherKemco3:
		return "Kemco"
	case OldPublisherArc:
		return "Arc"
	case OldPublisherNihonBussan:
		return "Nihon Bussan"
	case OldPublisherTecmo:
		return "Tecmo"
	case OldPublisherImagineer:
		return "Imagineer"
	case OldPublisherBanpresto2:
		return "Banpresto"
	case OldPublisherNova:
		return "Nova"
	case OldPublisherBandai2:
		return "Bandai"
	case OldPublisherKonami2:
		return "Konami"
	case OldPublisherKawada:
		return "Kawada"
	case OldPublisherTakara:
		return "Takara"
	case OldPublisherTechnosJapan:
		return "Technos Japan"
	case OldPublisherBroderbund2:
		return "Broderbund"
	case OldPublisherToeiAnimation:
		return "Toei Animation"
	case OldPublisherToho:
		return "Toho"
	case OldPublisherNamco:
		return "Namco"
	case OldPublisherAcclaimEntertainment2:
		return "Acclaim Entertainment"
	case OldPublisherASCIICorporation:
		return "ASCII Corporation"
	case OldPublisherBandai3:
		return "Bandai"
	case OldPublisherHoriElectric:
		return "Hori Electric"
	case OldPublisherSquareEnix:
		return "Square Enix"
	case OldPublisherHALLaboratory:
		return "HAL Laboratory"
	case OldPublisherSNK:
		return "SNK"
	case OldPublisherPonyCanyon:
		return "Pony Canyon"
	case OldPublisherCultureBrain:
		return "Culture Brain"
	case OldPublisherSunsoft:
		return "Sunsoft"
	case OldPublisherSonyImagesoft:
		return "Sony Imagesoft"
	case OldPublisherSammyCorporation:
		return "Sammy Corporation"
	case OldPublisherTaito:
		return "Taito"
	case OldPublisherKemco4:
		return "Kemco"
	case OldPublisherSquare:
		return "Square"
	case OldPublisherTokumaShoten2:
		return "Tokuma Shoten"
	case OldPublisherDataEast:
		return "Data East"
	case OldPublisherTonkinHouse:
		return "Tonkin House"
	case OldPublisherKoei:
		return "Koei"
	case OldPublisherUFL:
		return "UFL"
	case OldPublisherUltraGames:
		return "Ultra Games"
	case OldPublisherVAPInc:
		return "VAP Inc."
	case OldPublisherUseCorporation:
		return "Use Corporation"
	case OldPublisherMeldac:
		return "Meldac"
	case OldPublisherPonyCanyon2:
		return "Pony Canyon"
	case OldPublisherAngel2:
		return "Angel"
	case OldPublisherTaito2:
		return "Taito"
	case OldPublisherSOFEL:
		return "SOFEL"
	case OldPublisherQuest:
		return "Quest"
	case OldPublisherSigmaEnterprises:
		return "Sigma Enterprises"
	case OldPublisherASKKodanshaCo:
		return "ASK Kodansha Co."
	case OldPublisherNaxatSoft2:
		return "Naxat Soft"
	case OldPublisherCopyaSystem:
		return "Copya System"
	case OldPublisherBanpresto3:
		return "Banpresto"
	case OldPublisherTomy:
		return "Tomy"
	case OldPublisherLJN2:
		return "LJN"
	case OldPublisherNipponComputerSystems:
		return "Nippon Computer Systems"
	case OldPublisherHumanEnt:
		return "Human Entertainment"
	case OldPublisherAltron:
		return "Altron"
	case OldPublisherJaleco2:
		return "Jaleco"
	case OldPublisherTowaChiki:
		return "Towa Chiki"
	case OldPublisherYutaka:
		return "Yutaka"
	case OldPublisherVarie2:
		return "Varie"
	case OldPublisherEpoch:
		return "Epoch"
	case OldPublisherAthena:
		return "Athena"
	case OldPublisherAsmikAceEntertainment:
		return "Asmik Ace Entertainment"
	case OldPublisherNatsume:
		return "Natsume"
	case OldPublisherKingRecords:
		return "King Records"
	case OldPublisherAtlus2:
		return "Atlus"
	case OldPublisherEpicRecords:
		return "Epic Records"
	case OldPublisherIGS:
		return "IGS"
	case OldPublisherAWave:
		return "A-Wave"
	case OldPublisherExtremeEntertainment:
		return "Extreme Entertainment"
	case OldPublisherLJN3:
		return "LJN"
	default:
		return "Unknown Publisher"
	}
}
