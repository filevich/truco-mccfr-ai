package cfr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/filevich/truco-ai/abs"
	"github.com/filevich/truco-ai/info"
	"github.com/filevich/truco-ai/utils"
	"github.com/truquito/truco/pdt"
)

type Trainer struct {
	Builder *info.Builder

	CurrentIter int
	TotalIter   int
	InfosetMap  map[string]*RNode
	NumPlayers  int
	// multi
	Mu      *sync.Mutex
	Wg      *sync.WaitGroup
	Working int
}

func (trainer *Trainer) addRootUtils(new_utils []float32) {
	root := trainer.GetRnode("__ROOT__", 2)
	trainer.Lock()
	root.CumulativeRegrets = utils.SumFloat32Slices(root.CumulativeRegrets, new_utils)
	trainer.Unlock()
}

func (trainer *Trainer) inc_t() {
	trainer.Mu.Lock()
	defer trainer.Mu.Unlock()
	trainer.CurrentIter++
}

func (trainer *Trainer) inc_T() {
	trainer.Mu.Lock()
	defer trainer.Mu.Unlock()
	trainer.TotalIter++
}

func (trainer *Trainer) Lock() {
	trainer.Mu.Lock()
}

func (trainer *Trainer) Unlock() {
	trainer.Mu.Unlock()
}

func (trainer *Trainer) SetWorkers(n int) {
	trainer.Working = n
	trainer.Wg.Add(n)
}

func (trainer *Trainer) WorkerDone() {
	trainer.Mu.Lock()
	defer trainer.Mu.Unlock()
	// aviso que este thread termino
	trainer.Working--
}

func (trainer *Trainer) AllDones() bool {
	trainer.Mu.Lock()
	defer trainer.Mu.Unlock()
	return trainer.Working == 0
}

func (trainer *Trainer) Get_t() int {
	trainer.Mu.Lock()
	defer trainer.Mu.Unlock()
	return trainer.CurrentIter
}

func (trainer *Trainer) get_T() int {
	return trainer.TotalIter
}

func (trainer *Trainer) set_T(T int) {
	trainer.TotalIter = T
}

func (trainer *Trainer) getNumPlayers() int {
	return 2 // trainer.Num_players
}

func (t *Trainer) GetAbs() abs.IAbstraction {
	return t.Builder.Abs
}

func (t *Trainer) Reset() {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	for _, rnode := range t.InfosetMap {
		rnode.Reset()
	}
	t.CurrentIter = 0
}

func (t *Trainer) GetRnode(hash string, chiLen int) *RNode {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	if _, ok := t.InfosetMap[hash]; !ok {
		t.InfosetMap[hash] = NewRNode(chiLen)
	}
	return t.InfosetMap[hash]
}

func (t *Trainer) samplePartida() *pdt.Partida {
	A := []string{"Alice", "Ariana", "Anna"}
	B := []string{"Bob", "Ben", "Bill"}
	n := t.NumPlayers / 2
	limEnvite := 4
	verbose := true
	p, _ := pdt.NuevaPartida(pdt.A20, A[:n], B[:n], limEnvite, verbose)
	return p
}

func (trainer *Trainer) CountInfosets() int {
	trainer.Mu.Lock()
	defer trainer.Mu.Unlock()
	return len(trainer.InfosetMap)
}

func (trainer *Trainer) GetAvgStrategy(hash string, chiLen int) []float32 {
	rnode := trainer.GetRnode(hash, chiLen)
	return rnode.GetAverageStrategy()
}

func (t *Trainer) MaxAvgGameValue() float32 {
	r := t.GetRnode("__ROOT__", 0).CumulativeRegrets[0]
	agm := r / float32(t.TotalIter)
	if agm > 0 {
		return agm
	}
	return -agm
}

func (t *Trainer) FinalReport(profile IProfile) {
	if profile.IsFullySilent() {
		return
	}

	log.Println()
	for player := 0; player < t.getNumPlayers(); player++ {
		r := t.GetRnode("__ROOT__", 0).CumulativeRegrets[player]
		log.Printf("Computed average game value for player %d: %.3f\n",
			player+1,
			r/float32(t.TotalIter), // <-- OJO CON ESTO!!! todo: si es un perfil de tiempo el tital iters debe ser actualizado
		)
	}
	log.Println()
}

// io
func (t *Trainer) Save(filename string) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	// esto es debido a que el save se hace alfinal del for, antes de que se de
	// el incremento de la variable `t`
	t.CurrentIter++
	// falta el numero de las iteraciones
	utils.Write(t, filename, true)
	t.CurrentIter--
}

const CURRENT_MODEL_VERSION float64 = 3.0

// changes from v2.2 to v3.0:
//   - add `hash`, `info` and `abs` fields
//   - rm `abstractor` field
// changes from v2.1 to v2.2:
//   - ?
// changes from v2.0 to v2.1:
//   - field `id` is now `trainer`
//   - no use of hyphens in trainers. eg, `es-vmccfr` is now `esvmccfr`
// changes from v1.0 to v2.0
//   - field `Prot` is now `version`
//   - everything else is lowerCamelCased

func (t *Trainer) SaveModel(

	filename string,
	report_interval int,
	id string,
	extras []string,

) {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	// esto es debido a que el save se hace alfinal del for, antes de que se de
	// el incremento de la variable `t`
	t.CurrentIter++

	// falta el numero de las iteraciones
	// utils.Write(t, filename, true)

	n := len(filename)
	if ok := filename[n-6:] == ".model"; !ok {
		msg := fmt.Sprintf("la extension del archivo debe ser `.model`. (%s)", filename)
		panic(msg)
	}

	// creo el archivo
	f := utils.Touch(filename)
	defer f.Close()

	verbose := report_interval > 0

	if verbose {
		log.Printf("Saving: 0%%")
	}

	// Infoset_map  map[string]*RNode

	// agrego los campos de interes:
	// campos extras: como el tipo, o valor de epsilon de OS-MCCFR
	f.Write([]byte(fmt.Sprintf("version %.1f\n", CURRENT_MODEL_VERSION)))
	f.Write([]byte(fmt.Sprintf("trainer %s\n", id)))

	f.Write([]byte(fmt.Sprintf("hash %s\n", t.Builder.Data.Hash)))   // sha1, sha256, sha512
	f.Write([]byte(fmt.Sprintf("info %s\n", t.Builder.Data.Info)))   // InfosetRondaBase
	f.Write([]byte(fmt.Sprintf("abs %s\n", t.Builder.Abs.String()))) // a1, a2, a3...

	f.Write([]byte(fmt.Sprintf("currentiter %d\n", t.CurrentIter)))
	f.Write([]byte(fmt.Sprintf("totaliter %d\n", t.TotalIter)))
	f.Write([]byte(fmt.Sprintf("numplayers %d\n", t.NumPlayers)))

	for _, field := range extras {
		f.Write([]byte(fmt.Sprintf("%s\n", field)))
	}

	// doble salto de linea indica que se acabaron los campos
	f.Write([]byte("\n\n"))

	// agrego los inofsets/rnodes
	i := 0
	for hash, rnode := range t.InfosetMap {
		if verbose && utils.Mod(i+1, report_interval) == 0 {
			progress := float32(i+1) / float32(len(t.InfosetMap))
			fmt.Printf(" | %d%%", int(progress*100))
			runtime.GC()
		}

		bs, _ := json.Marshal(rnode)
		s := fmt.Sprintf("%s %s\n", hash, string(bs))
		if _, err := f.Write([]byte(s)); err != nil {
			panic(err)
		}

		i++
	}

	// retorno Current_Iter a su estado orig
	t.CurrentIter--
}

func (t *Trainer) GetBuilder() *info.Builder {
	return t.Builder
}

func lineCounter(filename string) (int, error) {
	file, err := os.Open(filename)

	if err != nil {
		panic(fmt.Sprintf("failed to open model `%s`", filename))
	}

	defer file.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func LoadModel(
	filename string,
	verbose bool,
	report_interval int,
	lazy bool,
) ITrainer {

	var t Trainer_T
	base := &Trainer{
		CurrentIter: 0,
		TotalIter:   0,
		InfosetMap:  make(map[string]*RNode),
		NumPlayers:  0,
		Builder:     nil,
		Mu:          &sync.Mutex{},
		Wg:          &sync.WaitGroup{},
	}

	file, err := os.Open(filename)

	if err != nil {
		panic(fmt.Sprintf("failed to open model `%s`", filename))
	}

	if verbose {
		log.Printf("Fetching model size...\n")
	}

	locs, _ := lineCounter(filename)

	if verbose {
		log.Printf("%d lines red\n", locs)
		log.Printf("Loading model: 0%%")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	return_counter := 0
	i := 0

	_hash := ""
	_info := ""
	_abs := ""

	for scanner.Scan() {
		line := scanner.Text()

		// caso 1: atributos
		if return_counter < 2 {

			if line == "" {
				return_counter++
				if return_counter == 2 && lazy {
					break
				} else {
					continue
				}
			}

			words := strings.Fields(line)
			val := words[1]

			switch words[0] {
			case "version":
				words[1] = strings.TrimSuffix(words[1], "-distil")
				v, err := strconv.ParseFloat(words[1], 64)
				if err != nil {
					panic(fmt.Sprintf("couldn't parse version: %s", err))
				}
				if isOldVersion := v < CURRENT_MODEL_VERSION; isOldVersion {
					tmpl := "cannot load deprecated .model version; has:%v want:%v"
					panic(fmt.Sprintf(tmpl, v, CURRENT_MODEL_VERSION))
				}
			case "trainer":
				t = Trainer_T(words[1])
			case "currentiter":
				base.CurrentIter, _ = strconv.Atoi(val)
			case "totaliter":
				base.TotalIter, _ = strconv.Atoi(val)
			case "numplayers":
				base.NumPlayers, _ = strconv.Atoi(val)
			case "hash":
				_hash = val
			case "info":
				_info = val
			case "abs":
				_abs = val
			default:
				continue
			}

			// caso 2: rnode
		} else if return_counter == 2 {

			if utils.Mod(i+1, report_interval) == 0 {
				runtime.GC()
				if verbose {
					progress := float32(i+1) / float32(locs)
					log.Printf(" | %d%%", int(progress*100))
				}
			}

			i++

			ix := strings.Index(line, " ")
			hash, data := line[:ix], line[ix+1:]
			rnode := &RNode{}
			if err := json.Unmarshal([]byte(data), rnode); err != nil {
				panic(fmt.Sprintf("no se pudo parsear la linea %s", hash))
			}

			// lo agrego
			base.InfosetMap[hash] = rnode
		}
	}

	if ok := len(_hash)*len(_info)*len(_abs) > 0; !ok {
		panic("either `hash`, `info` or `abs` field is missing")
	}

	base.Builder = info.BuilderFactory(_hash, _info, _abs)

	if verbose {
		log.Println()
	}

	return Embed(t, base)
}

type Trainer_T string

const (
	// vanilla
	CFR_T      Trainer_T = "cfr"
	CFRP_T     Trainer_T = "cfrp"
	DCFR_T     Trainer_T = "dcfr"
	ESLMCCFR_T Trainer_T = "eslmccfr"
	ESVMCCFR_T Trainer_T = "esvmccfr"
	OSMCCFR_T  Trainer_T = "oslmccfr"
	// exploit
	BR_T Trainer_T = "bestresponse"
)

func NewTrainer(

	t Trainer_T,
	numPlayers int,
	_hash string,
	_info string,
	_abs string,

) ITrainer {

	base := Trainer{
		Builder:     info.BuilderFactory(_hash, _info, _abs),
		CurrentIter: 0,
		TotalIter:   0,
		InfosetMap:  make(map[string]*RNode),
		NumPlayers:  numPlayers,
		Mu:          &sync.Mutex{},
		Wg:          &sync.WaitGroup{},
	}

	return Embed(t, &base)
}

func Embed(t Trainer_T, base *Trainer) ITrainer {
	switch t {
	// vanilla
	case CFR_T:
		return &CFR{base}

	// variantes
	// case CFRP_T:
	// 	return &CFRP{base}

	// case DCFR_T:
	// 	return &DCFR{base, 1.5, 0, 2}

	case ESLMCCFR_T:
		return &ESLMCCFR{base}

	case ESVMCCFR_T:
		return &ESVMCCFR{base}

	// case OSMCCFR_T:
	// 	return &OSMCCFR{base, 0.1} // epsilon

	// // exploitability
	// case BR_T:
	// 	return &BestResponse{base}

	default:
		panic("trainer unknown")
	}
}
