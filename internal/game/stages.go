package game

import "image/color"

type TitanClass struct {
	ID          int
	Name        string
	Title       string
	Description string
	SpecialName string
	VoiceIntro  string
	SpeedMul    float64
	ClimbMul    float64
	SmashMul    float64
	GoreHealMul float64
	TintColor   color.RGBA
	AuraColor   color.RGBA
}

var TitanRoster = []TitanClass{
	{
		ID:          0,
		Name:        "ΚΡΟΝΟΣ",
		Title:       "Ο ΑΙΜΟΒΟΡΟΣ ΤΙΤΑΝΑΣ",
		Description: "+40% Αίμα & Ζωή όταν καταβροχθίζει ανθρώπους! Κεραυνοί του Ταρτάρου.",
		SpecialName: "ΔΡΕΠΑΝΙ & ΚΕΡΑΥΝΟΣ ΤΑΡΤΑΡΟΥ",
		VoiceIntro:  "Είμαι ο Κρόνος! Θα φάω όλη την Αθήνα χωρίς τζατζίκι!",
		SpeedMul:    1.0,
		ClimbMul:    1.05,
		SmashMul:    1.1,
		GoreHealMul: 1.45,
		TintColor:   color.RGBA{255, 230, 220, 255},
		AuraColor:   color.RGBA{220, 30, 50, 220},
	},
	{
		ID:          1,
		Name:        "ΤΥΦΩΝΑΣ",
		Title:       "Ο ΔΡΑΚΟ-ΤΙΤΑΝΑΣ ΤΗΣ ΦΩΤΙΑΣ",
		Description: "+35% Ταχύτητα Σκαρφαλώματος & Πύρινη Ανάσα που λιώνει ελικόπτερα και ορόφους!",
		SpecialName: "ΚΑΤΑΙΓΙΔΑ ΛΑΒΑΣ ΤΟΥ ΟΛΥΜΠΟΥ",
		VoiceIntro:  "Ο Τυφώνας ξύπνησε! Κάψτε τις πολυκατοικίες και τα κίτρινα ταξί!",
		SpeedMul:    1.25,
		ClimbMul:    1.35,
		SmashMul:    1.0,
		GoreHealMul: 1.1,
		TintColor:   color.RGBA{210, 255, 210, 255},
		AuraColor:   color.RGBA{40, 220, 90, 220},
	},
	{
		ID:          2,
		Name:        "ΠΟΛΥΦΗΜΟΣ",
		Title:       "Ο ΑΝΘΡΩΠΟΦΑΓΟΣ ΚΥΚΛΩΠΑΣ",
		Description: "+50% Ζημιά Κατεδάφισης στις πολυκατοικίες! Σεισμικές γροθιές Εγκέλαδου.",
		SpecialName: "ΣΕΙΣΜΟΣ ΤΟΥ ΕΓΚΕΛΑΔΟΥ (10 ΡΙΧΤΕΡ)",
		VoiceIntro:  "Πολύφημος πεινάει! Φέρτε μου τουρίστες με κάλτσα και σανδάλι!",
		SpeedMul:    0.92,
		ClimbMul:    0.95,
		SmashMul:    1.50,
		GoreHealMul: 1.25,
		TintColor:   color.RGBA{225, 220, 255, 255},
		AuraColor:   color.RGBA{255, 165, 30, 220},
	},
}

type BuildingSpec struct {
	SignTitle   string
	Subtitle    string
	Floors      int
	Cols        int
	WallColor   color.RGBA
	AwningColor color.RGBA
	IsMarble    bool
}

type VictimDialog struct {
	Archetype  string
	PanicQuote string
	EatenQuote string
	VoiceLang  string
	ShirtColor color.RGBA
}

type StageSpec struct {
	Number     int
	Name       string
	Subtitle   string
	IntroQuip  string
	SkyTop     color.RGBA
	SkyBottom  color.RGBA
	HasBoss    bool
	Buildings  []BuildingSpec
	Victims    []VictimDialog
}

var AthensStages = []StageSpec{
	{
		Number:    1,
		Name:      "ΠΙΣΤΑ 1: ΜΟΝΑΣΤΗΡΑΚΙ & ΠΛΑΚΑ",
		Subtitle:  "Τουρίστες με Κάλτσα-Σανδάλι, Κράχτες & Πιτόγυρο 4,80€!",
		IntroQuip: "Καλώς ήρθατε στο Μοναστηράκι! Φάτε τους τουρίστες και γκρεμίστε τα Airbnb!",
		SkyTop:    color.RGBA{24, 16, 48, 255},
		SkyBottom: color.RGBA{120, 48, 72, 255},
		HasBoss:   false,
		Buildings: []BuildingSpec{
			{
				SignTitle:   "Ο ΤΖΙΜΗΣ - ΠΙΤΟΓΥΡΟ 4.80€",
				Subtitle:    "ΑΥΘΕΝΤΙΚΟ ΣΟΥΒΛΑΚΙ ΑΠ' ΟΛΑ",
				Floors:      4,
				Cols:        3,
				WallColor:   color.RGBA{215, 195, 165, 255},
				AwningColor: color.RGBA{220, 95, 30, 255},
			},
			{
				SignTitle:   "TAVERNA MOUSSAKA SPECIAL",
				Subtitle:    "FREE SHOT OUZO MY FRIEND",
				Floors:      5,
				Cols:        3,
				WallColor:   color.RGBA{225, 215, 195, 255},
				AwningColor: color.RGBA{35, 115, 210, 255},
			},
			{
				SignTitle:   "AIRBNB ACROPOLIS VIEW",
				Subtitle:    "280€ / ΒΡΑΔΙΑ (ΧΩΡΙΣ ΘΕΡΜΑΝΣΗ)",
				Floors:      6,
				Cols:        3,
				WallColor:   color.RGBA{195, 185, 175, 255},
				AwningColor: color.RGBA{210, 75, 25, 255},
			},
			{
				SignTitle:   "SOUVENIR SPARTAN HELMETS",
				Subtitle:    "MADE IN CHINA - GENUINE GREEK",
				Floors:      4,
				Cols:        3,
				WallColor:   color.RGBA{210, 180, 155, 255},
				AwningColor: color.RGBA{185, 45, 55, 255},
			},
		},
		Victims: []VictimDialog{
			{
				Archetype:  "Τουρίστας με Σανδάλι",
				PanicQuote: "OMG is this Titan part of the walking tour?!",
				EatenQuote: "MY AIRBNB DEPOSIT! ΑΑΑΑΧ!",
				VoiceLang:  "en-US",
				ShirtColor: color.RGBA{255, 90, 120, 255},
			},
			{
				Archetype:  "Κράχτης Ταβέρνας",
				PanicQuote: "Come my friend best Moussaka in Plaka!",
				EatenQuote: "ΩΧ ΜΕ ΤΡΩΕΙ ΧΩΡΙΣ ΤΖΑΤΖΙΚΙ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{50, 160, 250, 255},
			},
			{
				Archetype:  "Influencer στο Μοναστηράκι",
				PanicQuote: "Chat, ένας Τιτάνας γκρεμίζει την Πλάκα!",
				EatenQuote: "ΚΑΝΤΕ LIKE ΚΑΙ SUBSCRIBE ΑΑΑΧ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{240, 80, 230, 255},
			},
		},
	},
	{
		Number:    2,
		Name:      "ΠΙΣΤΑ 2: ΣΥΝΤΑΓΜΑ & ΕΡΜΟΥ",
		Subtitle:  "Κίτρινα Ταξί που Αλλάζουν Βάρδια, Ερμού & Έργα Μετρό!",
		IntroQuip: "Πίστα Δύο, Σύνταγμα! Διαλύστε τα ταξί που δεν πάνε προς τα κει!",
		SkyTop:    color.RGBA{18, 24, 58, 255},
		SkyBottom: color.RGBA{95, 55, 105, 255},
		HasBoss:   false,
		Buildings: []BuildingSpec{
			{
				SignTitle:   "ΞΕΝΟΔΟΧΕΙΟ ΜΕΓΑΛΗ ΒΡΕΤΑΝΝΙΑ",
				Subtitle:    "VIP LUXURY SUITES",
				Floors:      6,
				Cols:        3,
				WallColor:   color.RGBA{230, 215, 185, 255},
				AwningColor: color.RGBA{180, 140, 45, 255},
			},
			{
				SignTitle:   "ΕΡΜΟΥ MEGA FASHION STORE",
				Subtitle:    "ΕΚΠΤΩΣΕΙΣ -10% ΣΕ ΟΛΑ",
				Floors:      5,
				Cols:        3,
				WallColor:   color.RGBA{200, 205, 215, 255},
				AwningColor: color.RGBA{215, 45, 95, 255},
			},
			{
				SignTitle:   "ΥΠΟΥΡΓΕΙΟ ΟΙΚΟΝΟΜΙΚΩΝ",
				Subtitle:    "ΕΝΦΙΑ & ΦΟΡΟΛΟΓΗΣΗ ΤΙΤΑΝΩΝ",
				Floors:      6,
				Cols:        3,
				WallColor:   color.RGBA{185, 190, 195, 255},
				AwningColor: color.RGBA{45, 85, 160, 255},
			},
			{
				SignTitle:   "ΜΕΤΡΟ ΣΥΝΤΑΓΜΑΤΟΣ",
				Subtitle:    "ΕΠΟΜΕΝΟΣ ΣΥΡΜΟΣ ΣΕ 48 ΛΕΠΤΑ",
				Floors:      4,
				Cols:        3,
				WallColor:   color.RGBA{175, 180, 190, 255},
				AwningColor: color.RGBA{225, 110, 25, 255},
			},
		},
		Victims: []VictimDialog{
			{
				Archetype:  "Ταξιτζής",
				PanicQuote: "Δεν πάω Τάρταρα φίλε, αλλάζω βάρδια!",
				EatenQuote: "ΒΑΛΕ ΔΙΠΛΗ ΤΑΡΙΦΑ ΡΕ ΚΡΟΝΕΕΕ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{255, 215, 30, 255},
			},
			{
				Archetype:  "Shopper στην Ερμού",
				PanicQuote: "Μη μου πατάς τις σακούλες, πήρα παπούτσια!",
				EatenQuote: "ΚΡΑΤΑ ΤΗΝ ΚΑΡΤΑ ΑΛΛΑΓΗΣ ΑΑΑΧ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{255, 110, 180, 255},
			},
			{
				Archetype:  "Εφοριακός ΑΑΔΕ",
				PanicQuote: "Έχεις άδεια κατεδάφισης απ' την Πολεοδομία;!",
				EatenQuote: "ΠΡΟΣΤΙΜΟ ΓΙΑ ΩΜΗ ΚΑΤΑΝΑΛΩΣΗ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{120, 140, 175, 255},
			},
		},
	},
	{
		Number:    3,
		Name:      "ΠΙΣΤΑ 3: ΕΞΑΡΧΕΙΑ & ΟΜΟΝΟΙΑ",
		Subtitle:  "Πορτοκαλί Τέντες '70s, Απλήρωτα Κοινόχρηστα & Γιαγιάδες με Παντόφλα!",
		IntroQuip: "Πίστα Τρία, Εξάρχεια και Ομόνοια! Προσοχή στις Γιαγιάδες με τις παντόφλες!",
		SkyTop:    color.RGBA{32, 12, 36, 255},
		SkyBottom: color.RGBA{145, 42, 48, 255},
		HasBoss:   false,
		Buildings: []BuildingSpec{
			{
				SignTitle:   "ΠΟΛΥΚΑΤΟΙΚΙΑ ΕΤΟΣ 1974",
				Subtitle:    "ΑΠΛΗΡΩΤΑ ΚΟΙΝΟΧΡΗΣΤΑ 18 ΜΗΝΕΣ",
				Floors:      6,
				Cols:        3,
				WallColor:   color.RGBA{190, 180, 165, 255},
				AwningColor: color.RGBA{230, 95, 20, 255},
			},
			{
				SignTitle:   "COFFEE CULT - FREDDO ESPRESSO",
				Subtitle:    "ΜΕΤΡΙΟΣ ΜΕ ΜΑΥΡΗ ΖΑΧΑΡΗ",
				Floors:      5,
				Cols:        3,
				WallColor:   color.RGBA{175, 165, 155, 255},
				AwningColor: color.RGBA{55, 150, 75, 255},
			},
			{
				SignTitle:   "ΦΩΤΟΤΥΠΙΕΣ & ΦΡΟΝΤΙΣΤΗΡΙΟ",
				Subtitle:    "ΣΗΜΕΙΩΣΕΙΣ ΕΞΕΤΑΣΤΙΚΗΣ",
				Floors:      5,
				Cols:        3,
				WallColor:   color.RGBA{200, 190, 175, 255},
				AwningColor: color.RGBA{220, 105, 30, 255},
			},
			{
				SignTitle:   "24ΩΡΟ ΠΕΡΙΠΤΕΡΟ ΟΜΟΝΟΙΑΣ",
				Subtitle:    "ΚΑΡΤΕΣ, ΤΣΙΧΛΕΣ & ΑΝΑΠΤΗΡΕΣ",
				Floors:      4,
				Cols:        3,
				WallColor:   color.RGBA{185, 175, 170, 255},
				AwningColor: color.RGBA{200, 50, 45, 255},
			},
		},
		Victims: []VictimDialog{
			{
				Archetype:  "Διαχειριστής Πολυκατοικίας",
				PanicQuote: "Μη γκρεμίζεις το ασανσέρ, μόλις βάλαμε πετρέλαιο!",
				EatenQuote: "ΠΟΙΟΣ ΘΑ ΠΛΗΡΩΣΕΙ ΤΑ ΚΟΙΝΟΧΡΗΣΤΑΑΑ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{90, 185, 120, 255},
			},
			{
				Archetype:  "Ντελιβεράς",
				PanicQuote: "Άσε με να περάσω, θα νερώσει το Freddo!",
				EatenQuote: "ΧΥΘΗΚΕ Ο ΚΑΦΕΣ ΡΕ ΤΙΤΑΝΑΑΑ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{255, 130, 35, 255},
			},
			{
				Archetype:  "Γιαγιά στο Μπαλκόνι",
				PanicQuote: "Φύγε απ' τα γεράνια μου σατανά, μόλις σφουγγάρισα!",
				EatenQuote: "ΘΑ ΣΕ ΚΑΤΑΡΑΣΤΩ ΑΛΗΤΗΗΗ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{165, 95, 210, 255},
			},
		},
	},
	{
		Number:    4,
		Name:      "ΠΙΣΤΑ 4: ΚΗΦΙΣΟΣ (ΩΡΑ ΑΙΧΜΗΣ) & ΠΥΡΓΟΣ ΑΘΗΝΩΝ",
		Subtitle:  "Μποτιλιάρισμα στη Λένορμαν, ΛΕΑ & Γυάλινα Γραφεία!",
		IntroQuip: "Πίστα Τέσσερα, Κηφισός σε ώρα αιχμής! Φάτε τους οδηγούς που μπαίνουν από τη ΛΕΑ!",
		SkyTop:    color.RGBA{15, 25, 45, 255},
		SkyBottom: color.RGBA{75, 95, 130, 255},
		HasBoss:   false,
		Buildings: []BuildingSpec{
			{
				SignTitle:   "ΠΥΡΓΟΣ ΑΘΗΝΩΝ (TOWER 1)",
				Subtitle:    "CORPORATE SYNERGY HQ",
				Floors:      7,
				Cols:        3,
				WallColor:   color.RGBA{135, 175, 205, 255},
				AwningColor: color.RGBA{40, 110, 175, 255},
			},
			{
				SignTitle:   "ΓΕΦΥΡΑ ΚΗΦΙΣΟΥ - ΕΞΟΔΟΣ ΛΕΝΟΡΜΑΝ",
				Subtitle:    "ΜΕΣΗ ΤΑΧΥΤΗΤΑ: 2 KM/H",
				Floors:      5,
				Cols:        3,
				WallColor:   color.RGBA{170, 175, 180, 255},
				AwningColor: color.RGBA{225, 140, 25, 255},
			},
			{
				SignTitle:   "STARTUP UNICORN HUB",
				Subtitle:    "AI BLOCKCHAIN FRAPPE",
				Floors:      6,
				Cols:        3,
				WallColor:   color.RGBA{145, 185, 210, 255},
				AwningColor: color.RGBA{85, 60, 190, 255},
			},
			{
				SignTitle:   "ΚΤΕΛ ΚΗΦΙΣΟΥ",
				Subtitle:    "ΔΕΜΑΤΑ & ΤΥΡΟΠΙΤΕΣ",
				Floors:      4,
				Cols:        3,
				WallColor:   color.RGBA{195, 185, 170, 255},
				AwningColor: color.RGBA{205, 70, 35, 255},
			},
		},
		Victims: []VictimDialog{
			{
				Archetype:  "Οδηγός στον Κηφισό",
				PanicQuote: "Πού μπαίνεις ρε Καραμήτρο από τη ΛΕΑ;!",
				EatenQuote: "ΚΑΛΥΤΕΡΑ ΣΤΟ ΣΤΟΜΑΧΙ ΣΟΥ ΠΑΡΑ ΣΤΟΝ ΚΗΦΙΣΟ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{70, 200, 235, 255},
			},
			{
				Archetype:  "Manager σε Zoom Call",
				PanicQuote: "Παιδιά με τρώει ένας Τιτάνας, ας το δούμε offline!",
				EatenQuote: "ΤΑ KPI ΤΟΥ ΤΡΙΜΗΝΟΥΥΥ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{210, 220, 240, 255},
			},
		},
	},
	{
		Number:    5,
		Name:      "ΠΙΣΤΑ 5: ΙΕΡΟΣ ΒΡΑΧΟΣ ΑΚΡΟΠΟΛΗΣ & BOSS FIGHT!",
		Subtitle:  "Παρθενώνας, Προπύλαια & Ο MECHA-ΜΙΝΩΤΑΥΡΟΣ ΤΗΣ ΑΑΔΕ!",
		IntroQuip: "Τελική Πίστα, Ακρόπολη! Καταστρέψτε τον Μηχανικό Μινώταυρο της Εφορίας!",
		SkyTop:    color.RGBA{45, 8, 22, 255},
		SkyBottom: color.RGBA{185, 35, 30, 255},
		HasBoss:   true,
		Buildings: []BuildingSpec{
			{
				SignTitle:   "ΠΡΟΠΥΛΑΙΑ & ΕΚΔΟΤΗΡΙΑ",
				Subtitle:    "ΕΙΣΙΤΗΡΙΟ 30€ (ΜΕΙΩΜΕΝΟ 15€)",
				Floors:      5,
				Cols:        3,
				WallColor:   color.RGBA{242, 235, 220, 255},
				AwningColor: color.RGBA{205, 165, 55, 255},
				IsMarble:    true,
			},
			{
				SignTitle:   "ΕΡΕΧΘΕΙΟ & ΚΑΡΥΑΤΙΔΕΣ",
				Subtitle:    "ΜΗΝ ΑΚΟΥΜΠΑΤΕ ΤΑ ΜΑΡΜΑΡΑ!",
				Floors:      5,
				Cols:        3,
				WallColor:   color.RGBA{245, 238, 225, 255},
				AwningColor: color.RGBA{190, 145, 45, 255},
				IsMarble:    true,
			},
			{
				SignTitle:   "ΠΑΡΘΕΝΩΝΑΣ — ΙΕΡΟΣ ΒΡΑΧΟΣ",
				Subtitle:    "ΘΕΑ ΑΘΗΝΑ VS ΚΡΟΝΟΣ",
				Floors:      6,
				Cols:        4,
				WallColor:   color.RGBA{248, 242, 230, 255},
				AwningColor: color.RGBA{220, 175, 50, 255},
				IsMarble:    true,
			},
		},
		Victims: []VictimDialog{
			{
				Archetype:  "Φύλακας Αρχαιολογικής",
				PanicQuote: "ΦΥΥΥΥΤ! Μην ακουμπάτε τους κίονες κύριε Τιτάνα!",
				EatenQuote: "Η ΣΦΥΡΙΧΤΡΑ ΜΟΥΥΥ!",
				VoiceLang:  "el-GR",
				ShirtColor: color.RGBA{60, 120, 210, 255},
			},
			{
				Archetype:  "VIP Ξεναγός",
				PanicQuote: "As you can see, Kronos is devouring our tour group!",
				EatenQuote: "DON'T FORGET TO TIP YOUR GUIDE!",
				VoiceLang:  "en-US",
				ShirtColor: color.RGBA{245, 200, 70, 255},
			},
		},
	},
}
