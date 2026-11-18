
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	apiBase       = "https://prim.iledefrance-mobilites.fr/marketplace/v2/navitia"
	campusCoord   = "2.20891;48.90162"
	cacheDuration = 5 * time.Minute
)

var client = &http.Client{
	Timeout: 15 * time.Second,
}

var cache struct {
	sync.Mutex
	data      []byte
	expiresAt time.Time
}

type Reponse struct {
	Places []Lieu `json:"places"`
}

type Lieu struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	StopArea StopArea `json:"stop_area"`
}

type StopArea struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Coord Coord   `json:"coord"`
	Lines []Ligne `json:"lines"`
}

type Coord struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type Ligne struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	CommercialMode Mode   `json:"commercial_mode"`
	PhysicalModes  []Mode `json:"physical_modes"`
}

type Mode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ReponseHoraires struct {
	Departures []Depart `json:"departures"`
}

type Depart struct {
	Informations Informations `json:"display_informations"`
	Passage       Passage      `json:"stop_date_time"`
}

type Informations struct {
	Code           string `json:"code"`
	Direction      string `json:"direction"`
	CommercialMode string `json:"commercial_mode"`
}

type Passage struct {
	DepartureDateTime string `json:"departure_date_time"`
	BaseDateTime      string `json:"base_departure_date_time"`
	DataFreshness     string `json:"data_freshness"`
}

type ReponseProximite struct {
	PlacesNearby []ArretProximite `json:"places_nearby"`
	Disruptions  []json.RawMessage `json:"disruptions"`
}

type ArretProximite struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Distance string   `json:"distance"`
	StopArea StopArea `json:"stop_area"`
}

type ReponseLignes struct {
	Lines []Ligne `json:"lines"`
}

type ReponsePerturbations struct {
	Disruptions []json.RawMessage `json:"disruptions"`
}

type Erreur struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/api/transports/recherche", rechercherTransports)
	http.HandleFunc("/api/transports/horaires", horairesTransports)
	http.HandleFunc("/api/transports/proximite", transportsProximite)

	http.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Println("CampusLink est lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func obtenirCleAPI() (string, error) {
	apiKey := os.Getenv("PRIM_API_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("la clé PRIM_API_KEY est introuvable")
	}

	return apiKey, nil
}

func effectuerRequete(endpoint string) ([]byte, error) {
	apiKey, err := obtenirCleAPI()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", apiBase+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contenu, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"erreur de l'API PRIM : code %d - %s",
			resp.StatusCode,
			string(contenu),
		)
	}

	return contenu, nil
}

func envoyerJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(
			w,
			"Erreur lors de la création de la réponse JSON",
			http.StatusInternalServerError,
		)
	}
}

func envoyerErreur(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Erreur{
		Message: message,
	})
}

func rechercherTransports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		envoyerErreur(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	nom := strings.TrimSpace(r.URL.Query().Get("q"))

	if nom == "" {
		envoyerErreur(w, "Veuillez saisir le nom d'un arrêt", http.StatusBadRequest)
		return
	}

	params := url.Values{}
	params.Set("q", nom)
	params.Add("type[]", "stop_area")

	endpoint := "/places?" + params.Encode()

	body, err := effectuerRequete(endpoint)
	if err != nil {
		envoyerErreur(w, err.Error(), http.StatusBadGateway)
		return
	}

	var resultat Reponse

	if err := json.Unmarshal(body, &resultat); err != nil {
		envoyerErreur(w, "Impossible de lire les résultats de recherche", http.StatusInternalServerError)
		return
	}

	envoyerJSON(w, resultat)
}

func horairesTransports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		envoyerErreur(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	idArret := strings.TrimSpace(r.URL.Query().Get("id"))

	if idArret == "" {
		envoyerErreur(w, "Identifiant de l'arrêt manquant", http.StatusBadRequest)
		return
	}

	fromDateTime := time.Now().Format("20060102T150405")

	endpoint := "/stop_areas/" + url.PathEscape(idArret) +
		"/departures?from_datetime=" + fromDateTime +
		"&count=10&duration=3600"

	body, err := effectuerRequete(endpoint)
	if err != nil {
		envoyerErreur(w, err.Error(), http.StatusBadGateway)
		return
	}

	var resultat ReponseHoraires

	if err := json.Unmarshal(body, &resultat); err != nil {
		envoyerErreur(w, "Impossible de lire les prochains passages", http.StatusInternalServerError)
		return
	}

	envoyerJSON(w, resultat)
}

func transportsProximite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		envoyerErreur(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	cache.Lock()

	if len(cache.data) > 0 && time.Now().Before(cache.expiresAt) {
		donnees := append([]byte(nil), cache.data...)
		cache.Unlock()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(donnees)
		return
	}

	cache.Unlock()

	params := url.Values{}
	params.Set("distance", "1000")
	params.Add("type[]", "stop_area")
	params.Set("count", "20")
	params.Set("depth", "2")

	endpoint := "/coverage/fr-idf/coords/" + campusCoord +
		"/places_nearby?" + params.Encode()

	body, err := effectuerRequete(endpoint)
	if err != nil {
		envoyerErreur(w, err.Error(), http.StatusBadGateway)
		return
	}

	var resultat ReponseProximite

	if err := json.Unmarshal(body, &resultat); err != nil {
		envoyerErreur(w, "Impossible de lire les arrêts à proximité", http.StatusInternalServerError)
		return
	}

	for i := range resultat.PlacesNearby {
		idArret := resultat.PlacesNearby[i].ID

		if idArret == "" {
			continue
		}

		if i > 0 {
			time.Sleep(220 * time.Millisecond)
		}

		endpointLignes := "/stop_areas/" +
			url.PathEscape(idArret) + "/lines?depth=2"

		bodyLignes, err := effectuerRequete(endpointLignes)
		if err != nil {
			log.Printf("Impossible de récupérer les lignes de %s : %v", idArret, err)
			continue
		}

		var resultatLignes ReponseLignes

		if err := json.Unmarshal(bodyLignes, &resultatLignes); err != nil {
			log.Printf("Erreur de lecture des lignes de %s : %v", idArret, err)
			continue
		}

		resultat.PlacesNearby[i].StopArea.ID = idArret
		resultat.PlacesNearby[i].StopArea.Name = resultat.PlacesNearby[i].Name
		resultat.PlacesNearby[i].StopArea.Lines = resultatLignes.Lines
	}

	endpointPerturbations := "/line_reports/line_reports?count=100"

	bodyPerturbations, err := effectuerRequete(endpointPerturbations)

	if err != nil {
		log.Printf("Impossible de récupérer les perturbations : %v", err)
	} else {
		var resultatPerturbations ReponsePerturbations

		if err := json.Unmarshal(bodyPerturbations, &resultatPerturbations); err != nil {
			log.Printf("Erreur de lecture des perturbations : %v", err)
		} else {
			resultat.Disruptions = resultatPerturbations.Disruptions
		}
	}

	if resultat.Disruptions == nil {
		resultat.Disruptions = []json.RawMessage{}
	}

	reponse, err := json.Marshal(resultat)
	if err != nil {
		envoyerErreur(w, "Impossible de préparer les résultats", http.StatusInternalServerError)
		return
	}

	cache.Lock()
	cache.data = append([]byte(nil), reponse...)
	cache.expiresAt = time.Now().Add(cacheDuration)
	cache.Unlock()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(reponse)
}
