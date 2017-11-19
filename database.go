package tealtech

import (
	"sync"
	"io/ioutil"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/net/context"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
)

type Database interface {
	GetPatient(uid string) *Patient
	SavePatient(uid string, p *Patient)
	GetAllPatients() []*Patient
	Reset()
}

type DumpDatabase struct {
	mtx *sync.RWMutex
	m   map[string]*Patient
}

func (d *DumpDatabase) GetAllPatients() []*Patient {
	return nil
}

func NewDumpDatabase() *DumpDatabase {
	d := DumpDatabase{
		mtx: &sync.RWMutex{},
		m:   make(map[string]*Patient),
	}
	return &d
}

func (d *DumpDatabase) GetPatient(uid string) *Patient {
	d.mtx.RLock()
	defer d.mtx.RUnlock()

	u, ok := d.m[uid]
	if !ok {
		return nil
	}
	return u
}

func (d *DumpDatabase) SavePatient(uid string, p *Patient) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	d.m[uid] = p
}

func (d *DumpDatabase) Reset() {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	d.m = make(map[string]*Patient)
}

type Patient struct {
	ResourceType string      `json:"resourceType,omitempty"`
	ID           string      `json:"id,omitempty"`
	Active       bool        `json:"active,omitempty"`
	Name         []HumanName `json:"name,omitempty"`
	Gender       string      `json:"gender,omitempty"`
	BirthDate    string      `json:"birthDate,omitempty"`
	Address      []Address   `json:"address,omitempty"`
	Telecom      []Telecom   `json:"telecom,omitempty"`
}

type HumanName struct {
	Use    string   `json:"use"`
	Family []string `json:"family"`
	Given  []string `json:"given"`
}

type Address struct {
	Type     string   `json:"type"`
	Use      string   `json:"use"`
	City     string   `json:"city"`
	Line     []string `json:"line"`
	District string   `json:"district"`
}

type Telecom struct {
	Use    string `json:"use"`
	Rank   int    `json:"rank"`
	Value  string `json:"value"`
	System string `json:"system"`
}

type AidBoxDatabase struct {
	*http.Client
}

func NewAidBoxDatabase() *AidBoxDatabase {
	c := clientcredentials.Config{
		ClientID:     "golang",
		ClientSecret: "doesnotmatter",
		TokenURL:     "https://panacea.aidbox.io/oauth/token",
	}
	return &AidBoxDatabase{c.Client(context.Background())}
}

func (a *AidBoxDatabase) GetAllPatients() []*Patient {
	resp, err := a.Get("https://panacea.aidbox.io/fhir/Patient")
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	patients := new(AllPatients)
	if err = json.Unmarshal(b, patients); err != nil {
		panic(err)
	}

	var p []*Patient
	for _, v := range patients.Entry {
		p = append(p, v.Resource)
	}
	return p
}

type AllPatients struct {
	Entry []struct {
		Resource *Patient `json:"resource"`
	} `json:"entry"`
}

func (a *AidBoxDatabase) Reset() {
	return
}

func (a *AidBoxDatabase) GetPatient(uid string) *Patient {
	resp, err := a.Get("https://panacea.aidbox.io/fhir/Patient/" + uid)
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	patient := new(Patient)
	_ = json.Unmarshal(body, patient)
	return patient
}

func (a *AidBoxDatabase) SavePatient(uid string, p *Patient) {
	jsonPatient, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("PUT", "https://panacea.aidbox.io/fhir/Patient/"+uid, bytes.NewBuffer(jsonPatient))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.Do(req)
	if err != nil {
		panic(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
	if resp != nil {
		resp.Body.Close()
	}
}

//createPatient(client, "pat-viktor1111", "Виктор", "Плис", "Александрович", "male", "1996-06-26", "Россия", "Сибай", "ул. Революционная, д. 10", "89631338188")

func createPatient(ID, firstName, lastName, patronymic, gender, birthDate, district, city, street, telephone string) *Patient {
	return &Patient{
		ID: ID,
		Name: []HumanName{
			{
				Family: []string{lastName, patronymic},
				Given:  []string{firstName},
				Use:    "official",
			},
		},
		Gender:    gender,
		BirthDate: birthDate,
		Address: []Address{
			{
				District: district,
				City:     city,
				Line:     []string{street},
				Type:     "both",
				Use:      "home",
			},
		},
		Telecom: []Telecom{
			{
				Value:  telephone,
				Use:    "mobile",
				Rank:   1,
				System: "phone",
			},
		},

		ResourceType: "Patient",
		Active:       true,
	}
}
