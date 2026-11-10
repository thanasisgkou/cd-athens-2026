# Game Design Document (GDD): 🩸 ΤΙΤΑΝΟΜΑΧΙΑ ΣΤΗΝ ΑΘΗΝΑ — WRATH OF THE TITANS (ATHENS RAMPAGE)

> **Game Title**: 🩸 ΤΙΤΑΝΟΜΑΧΙΑ ΣΤΗΝ ΑΘΗΝΑ: WRATH OF KRONOS (Athens Rampage)
> **Target 2D Genre**: Classic Arcade Titan Destruction / Gore Platform-Climber (*Rampage* 1986 + *God of War* / *Metal Slug* Pixel Gore)
> **Camera Projection & Spatial Model**: 2D Side-Profile / Street-Level & Skyscraper Elevation View (`1280x720`)
> **Simulation Model**: Real-Time Continuous (60 TPS fixed timestep in Ebitengine v2)
> **Target Platform**: Desktop (Windows/macOS/Linux) & WebAssembly (Browser / Game Dev Studio Playtest Tab / Google Cloud Run)
> **Virtual Canvas & Aspect Ratio**: 16:9 Widescreen (`1280x720` virtual canvas)
> **Author / Lead Designer**: User & Antigravity

---

## 1. Executive Summary & Elevator Pitch
- **Elevator Pitch**: Οι αλυσίδες του Ταρτάρου έσπασαν κάτω από τον Κηφισό! Διάλεξε έναν από τους **3 Μυθικούς Τιτάνες** (**Κρόνος ο Αιμοβόρος**, **Τυφώνας ο Δρακο-Τιτάνας** ή **Πολύφημος ο Κύκλωπας**), σκαρφάλωσε στις πολυκατοικίες της Αθήνας με τις πορτοκαλί τέντες και στους κίονες της Ακρόπολης, γκρέμισε τα κτίρια με σεισμικές μπουνιές και **κατασπάραξε ζωντανούς τουρίστες, κράχτες, influencers και ταξιτζήδες** με over-the-top arcade gore (πίδακες αίματος που βάφουν τις πολυκατοικίες, κομμένα κεφάλια, κόκαλα, σπλάχνα, σανδάλια και Freddo Espresso) σε 5 ξεκαρδιστικές πίστες της Αθήνας!
- **Core Inspiration**: *Rampage* (Midway 1986) + *God of War* Mythic Brutality + *Metal Slug* Pixel Gore + Αθηναϊκή σάτιρα.
- **Target Mood & Pacing**: Καταιγιστικό, ξεκαρδιστικό, αιμοβόρο arcade action με διαμελισμούς, βάψιμο κτιρίων με αίμα, screen-shake, hitstop και επικές ελληνικές ατάκες.

---

## 2. Οι 3 Μυθικοί Τιτάνες (Playable Roster — Classic Rampage Trio)
1. **🩸 ΚΡΟΝΟΣ Ο ΑΙΜΟΒΟΡΟΣ (KRONOS THE DEVOURER)**: Ο Βασιλιάς των Τιτάνων που έτρωγε τα παιδιά του! Τεράστιος, με φλεγόμενα μάτια, σπασμένες αλυσίδες του Ταρτάρου και χάλκινη πανοπλία. **Special**: *Δρεπάνι του Ταρτάρου & Κεραυνός Αιωνιότητας* (extra Gore & Heal όταν τρώει ανθρώπους!).
2. **🐍 ΤΥΦΩΝΑΣ Ο ΚΑΤΑΣΤΡΟΦΕΑΣ (TYPHON THE HYDRA-TITAN)**: Το πιο τρομακτικό μυθικό τέρας της αρχαιότητας με δρακόμορφο σώμα, φίδια και πύρινο στόμα. **Special**: *Πύρινη Καταιγίδα Λάβας* (λιώνει ελικόπτερα και ολόκληρους ορόφους!).
3. **👁️ ΠΟΛΥΦΗΜΟΣ Ο ΓΙΓΑΝΤΑΣ ΚΥΚΛΩΠΑΣ (POLYPHEMUS THE MAN-EATER)**: Γιγάντιος ανθρωποφάγος Κύκλωπας που ξεριζώνει μπαλκόνια. **Special**: *Σεισμός του Εγκέλαδου* (διπλή ζημιά στις πολυκατοικίες και λιώσιμο εχθρών!).

---

## 2. Core Gameplay Loop, Simulation & Mechanics
- **Primary Gameplay Loop**:
  1. **Climb & Smash**: Περπάτα στον δρόμο ή **σκαρφάλωσε στις προσόψεις των πολυκατοικιών** (`W`/`S` ή `Up`/`Down`). Ρίξε μπουνιές (`J` / `Z` / `Left Click`) στα παράθυρα και τις κολώνες για να ρίξεις τη στατικότητα του κτιρίου στο `0%` και να το κατεδαφίσεις!
  2. **Grab & Devour Humans (Arcade Splatter)**: Όταν σπας ένα παράθυρο ή βλέπεις κόσμο στα μπαλκόνια/δρόμο, πάτα **Αρπαγή & Φάγωμα (`K` / `X` / `Right Click`)**! Ο Αθηνόσαυρος τους καταβροχθίζει με **πίδακες αίματος (blood particles), ιπτάμενα κόκαλα/κρανία, βάψιμο των τοίχων με αίμα**, +HP heal και +Rage!
  3. **Pickups & Traps**:
     - **Πιτόγυρο 4,80€**: +25 HP!
     - **Freddo Espresso Μέτριο**: +35 Οργή (Rage) & Speed Boost!
     - **Παγίδες στα Μπαλκόνια**: Προσοχή στις **Γιαγιάδες με Παντόφλα/Γλάστρα** και στα **Γκαζάκια/Μετασχηματιστές ΔΕΗ** που σε χτυπάνε με ηλεκτροσόκ!
  4. **Special — Φωτιά & Κεραυνός του Κηφισού (`Space` / `L` / `Middle Click`)**: Όταν έχεις Rage, ο Αθηνόσαυρος βρυχάται και εξαπολύει πύρινη/ηλεκτρική ανάσα που λιώνει ελικόπτερα, κλούβες ΜΑΤ και ολόκληρους ορόφους!

---

## 3. Οι 5 Πίστες της Αθήνας (Stages & Inside Jokes)
1. **ΠΙΣΤΑ 1: ΜΟΝΑΣΤΗΡΑΚΙ & ΠΛΑΚΑ (Tourist Trap)**
   - **Κτίρια**: Σουβλατζίδικο «Ο ΤΖΙΜΗΣ (Πιτόγυρο 4,80€)», Ταβέρνα «MOUSAKA SPECIAL», Airbnb «Acropolis View 250€/βραδιά».
   - **Θύματα & Ατάκες**: Τουρίστες με κάλτσα-σανδάλι (*«Is this dinosaur included in the Free Walking Tour?!»*), Κράχτες (*«My friend best Tzatziki... ΩΧ ΜΕ ΤΡΩΕΙ!»*), Influencers (*«Chat, με τρώει ένας Τυραννόσαυρος στο Μοναστηράκι!»*).
2. **ΠΙΣΤΑ 2: ΣΥΝΤΑΓΜΑ & ΕΡΜΟΥ (Yellow Taxi Fury)**
   - **Κτίρια**: Πολυκαταστήματα Ερμού, Ξενοδοχείο Μεγάλη Βρεταννία, Στάση Μετρό Συντάγματος, Υπουργείο.
   - **Θύματα & Εχθροί**: Κίτρινα Ταξί (*«Δεν πάω προς τα κει, αλλάζω βάρδια!»*, *«Διπλή ταρίφα ρε τέρας!»*), Περιπολικά, Κλούβες ΜΑΤ.
3. **ΠΙΣΤΑ 3: ΕΞΑΡΧΕΙΑ & ΟΜΟΝΟΙΑ (Freddo & Delivery Chaos)**
   - **Κτίρια**: Πολυκατοικίες του '70 με πορτοκαλί τέντες, κεραίες, ηλιακούς θερμοσίφωνες και καφετέριες.
   - **Θύματα & Εχθροί**: Μηχανάκια Delivery με κομμένη εξάτμιση (*«Άσε με να περάσω, θα νερώσει το Freddo!»*), Διαχειριστής Πολυκατοικίας (*«Δεν πλήρωσες τα κοινόχρηστα!»*), Γιαγιάδες που πετάνε γλάστρες (*«Μη μου πατάς τα γεράνια ρε αλήτη!»*).
4. **ΠΙΣΤΑ 4: ΚΗΦΙΣΟΣ (ΩΡΑ ΑΙΧΜΗΣ) & ΠΥΡΓΟΣ ΑΘΗΝΩΝ (Gridlock Carnage)**
   - **Κτίρια**: Πύργος Αθηνών, Γυάλινα Γραφεία, Υπερυψωμένη Γέφυρα Κηφισού.
   - **Θύματα & Εχθροί**: Οδηγοί στο μποτιλιάρισμα (*«Πού μπαίνεις ρε Καραμήτρο από τη ΛΕΑ;!»*, *«Καλύτερα στο στομάχι σου παρά στον Κηφισό!»*), Στρατιωτικά Ελικόπτερα & Τανκς.
5. **ΠΙΣΤΑ 5: ΙΕΡΟΣ ΒΡΑΧΟΣ ΑΚΡΟΠΟΛΗΣ & ΠΑΡΘΕΝΩΝΑΣ (Final Boss)**
   - **Κτίρια**: Προπύλαια, Ερέχθειο, Παρθενώνας.
   - **FINAL BOSS**: **MECHA-ΜΙΝΩΤΑΥΡΟΣ ΤΗΣ ΑΑΔΕ / ΑΡΧΑΙΟΛΟΓΙΚΗΣ** (*«ΕΧΕΙΣ ΑΠΟΔΕΙΞΗ ΚΑΙ ΕΙΣΙΤΗΡΙΟ ΓΙΑ ΤΟΝ ΙΕΡΟ ΒΡΑΧΟ;!»*) + Σμήνος Ελικοπτέρων!

---

## 4. Controls & Input Mapping Scheme

| Logical Action | Keyboard | Mouse / Pointer / Touch |
| :--- | :--- | :--- |
| `Κίνηση Αριστερά / Δεξιά` | `A` / `D` ή `Left` / `Right` | Κίνηση προς τον κέρσορα / Virtual D-Pad |
| `Σκαρφάλωμα Πάνω / Κάτω σε Κτίριο` | `W` / `S` ή `Up` / `Down` | Drag Πάνω / Κάτω πάνω σε πολυκατοικία |
| `Άλμα (Jump)` | `W` / `Up` (στον δρόμο) ή `Shift` | Κουμπί JUMP |
| `ΜΠΟΥΝΙΑ / ΚΑΤΕΔΑΦΙΣΗ (Smash)` | `J` ή `Z` | `Left Click` (Αριστερό Κλικ) |
| `ΑΡΠΑΓΗ & ΦΑΓΩΜΑ ΑΝΘΡΩΠΟΥ (Eat & Gore)` | `K` ή `X` ή `E` | `Right Click` (Δεξί Κλικ) |
| `ΠΥΡΙΝΗ ΑΝΑΣΑ / ΚΕΡΑΥΝΟΣ (Rage Breath)` | `Space` ή `L` ή `C` | `Middle Click` / Κουμπί RAGE |
| `Pause / Mute` | `P` / `Esc` (Pause), `M` (Mute) | Κουμπί Pause |

---

## 5. Visual & Audio Pipeline (`sprite-design`, `procedural-sfx`, `procedural-music`)
- **Sprites (`assets/sprites/*.sprite`)**:
  - `athenosaur.sprite`: Ο μεταλλαγμένος Αθηνόσαυρος (`idle`, `walk`, `climb`, `punch`, `eat`, `roar`).
  - `victim_athens.sprite`: Κάτοικοι/Τουρίστες στα μπαλκόνια & στον δρόμο (`idle`, `panic`, `eaten`).
  - `chopper_police.sprite`: Ελικόπτερο Αστυνομίας/Στρατού (`fly`, `hit`).
- **Gore & Destruction Particle System (`ebitengineer`)**:
  - Pixel blood fountain particles, persistent blood splatters on building facades, flying skulls/bones/ribs, flying sandals/frappes, rubble chunks, fire/smoke emitters, and trauma screen shake.
- **Sound Effects (`assets/audio/*.sfx`)**:
  - `punch_building.sfx`, `eat_chomp.sfx`, `scream_victim.sfx`, `building_collapse.sfx`, `atomic_roar.sfx`, `explosion_vehicle.sfx`, `yiayia_slipper.sfx`.
- **Background Music (`assets/audio/athens_rampage.music`)**:
  - 8-channel 16-bit YM2612 FM + SN76489 PSG high-energy action soundtrack.
