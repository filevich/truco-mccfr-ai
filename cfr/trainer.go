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

	"github.com/filevich/truco-cfr/abs"
	"github.com/filevich/truco-cfr/utils"
	"github.com/truquito/truco/pdt"
)

type Trainer struct {
	Current_Iter int
	Total_Iter   int
	Infoset_map  map[string]*RNode
	Num_players  int
	Abstractor   abs.IAbstraction
	// multi
	Mu      *sync.Mutex
	Wg      *sync.WaitGroup
	Working int
}

func (trainer *Trainer) add_root_utils(new_utils []float32) {
	root := trainer.GetRnode("__ROOT__", 2)
	trainer.Lock()
	root.CumulativeRegrets = utils.SumFloat32Slices(root.CumulativeRegrets, new_utils)
	trainer.Unlock()
}

func (trainer *Trainer) inc_t() {
	trainer.Mu.Lock()
	defer trainer.Mu.Unlock()
	trainer.Current_Iter++
}

func (trainer *Trainer) inc_T() {
	trainer.Mu.Lock()
	defer trainer.Mu.Unlock()
	trainer.Total_Iter++
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
	return trainer.Current_Iter
}

func (trainer *Trainer) get_T() int {
	return trainer.Total_Iter
}

func (trainer *Trainer) set_T(T int) {
	trainer.Total_Iter = T
}

func (trainer *Trainer) getNumPlayers() int {
	return 2 // trainer.Num_players
}

func (t *Trainer) GetAbs() abs.IAbstraction {
	return t.Abstractor
}

func (t *Trainer) Reset() {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	for _, rnode := range t.Infoset_map {
		rnode.Reset()
	}
	t.Current_Iter = 0
}

func (t *Trainer) GetRnode(hash string, chi_len int) *RNode {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	if _, ok := t.Infoset_map[hash]; !ok {
		t.Infoset_map[hash] = NewRNode(chi_len)
	}
	return t.Infoset_map[hash]
}

func (t *Trainer) samplePartida() *pdt.Partida {
	A := []string{"Alice", "Ariana", "Anna"}
	B := []string{"Bob", "Ben", "Bill"}
	n := t.Num_players / 2
	limEnvite := 4
	verbose := true
	p, _ := pdt.NuevaPartida(pdt.A20, A[:n], B[:n], limEnvite, verbose)
	return p
}

func (trainer *Trainer) CountInfosets() int {
	trainer.Mu.Lock()
	defer trainer.Mu.Unlock()
	return len(trainer.Infoset_map)
}

func (trainer *Trainer) GetAvgStrategy(hash string, chi_len int) []float32 {
	rnode := trainer.GetRnode(hash, chi_len)
	return rnode.Get_average_strategy()
}

func (t *Trainer) MaxAvgGameValue() float32 {
	r := t.GetRnode("__ROOT__", 0).CumulativeRegrets[0]
	agm := r / float32(t.Total_Iter)
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
			r/float32(t.Total_Iter), // <-- OJO CON ESTO!!! todo: si es un perfil de tiempo el tital iters debe ser actualizado
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
	t.Current_Iter++
	// falta el numero de las iteraciones
	utils.Write(t, filename, true)
	t.Current_Iter--
}

func (t *Trainer) SaveModel(filename string, report_interval int, id string, extras []string) {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	// esto es debido a que el save se hace alfinal del for, antes de que se de
	// el incremento de la variable `t`
	t.Current_Iter++

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
	f.Write([]byte("Prot 1.0\n"))
	f.Write([]byte(fmt.Sprintf("ID %s\n", id)))
	f.Write([]byte(fmt.Sprintf("Current_Iter %d\n", t.Current_Iter)))
	f.Write([]byte(fmt.Sprintf("Total_Iter %d\n", t.Total_Iter)))
	f.Write([]byte(fmt.Sprintf("Num_players %d\n", t.Num_players)))
	f.Write([]byte(fmt.Sprintf("Abstractor %s\n", t.Abstractor.String())))
	for _, field := range extras {
		f.Write([]byte(fmt.Sprintf("%s\n", field)))
	}

	// doble salto de linea indica que se acabaron los campos
	f.Write([]byte("\n\n"))

	// agrego los inofsets/rnodes
	i := 0
	for hash, rnode := range t.Infoset_map {
		if verbose && utils.Mod(i+1, report_interval) == 0 {
			progress := float32(i+1) / float32(len(t.Infoset_map))
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
	t.Current_Iter--
}

func (t *Trainer) Load(filename string) {
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(jsonFile, t); err != nil {
		panic(err)
	}
}

/*

	1. Tenemos pdt.GetA [15]bool o pdt.Chi []IJugada para obtener las jugadas
	   posibles.
	2. Una abstraccion puede modificar esos slices.
	3. El resultado de CFR va depender de la abstraccion usada.
	4. El orden de los regrets sera segun el que la abstraccion diga, luego:
		 - realidad: 2-4-5
		 - abstraccion 1: {"t1": 3, "t2": 0, "t3": 0}
		 - cfr: [ 1.234 ] <- significa que puedo puedo jugar una carta de abs 1
		 la pregunta es: ¿Cuál carta elijo? ~> queda delegado a la abstraccion.
	5. En cfr itero sobre el slice de todas las jugadas posibles.
	   luego a cada jugada se le asigna un indice (cartas), o bien
		 problema: una misma jugada se actualiza dos veces ->
		 [!] en esa iteracion se elige una sola carta por bucket [!]
	6. Conclusion: una instancia de CFR depende de que abstraccion se haya usado.

	opciones para almacenar los regrets:
	- slice de floats [1.11, 2.22, 3.33 ... 9.99]
	- diccionario {"t1": 1.11, "t2":2.22, "mazo":777}

	si estoy en `p`:

	tengo pesos:
	[1.11,      2.22, 3.33 ... 9.99]
	   ^        ^--->
	carta_t1    el resto de aca pa alla no son cartas

	entonces las primeras [:n] son de tirar carta, en orden
	[bucket_12, bucket_5, bucket_1]
	[bucket_3, bucket_1]
	[bucket_1]
	[]

	sinemabrgo en tengo chi:
	{"t1": 3, "t2":0, "t1":0, ... , "mazo":777}
	? como itero? tengo que poder iterar de tal forma que pueda ejecutar la accion.

	----------------------------

	si tengo el slice de pesos
	se lo puedo "meter" como parametro a un infoset
	el cual se encargue de saber a que accion le corresponde.


*/

func Load(t Trainer_T, filename string) ITrainer {

	data := struct {
		Trainer
		// sobre-escribo los atributos de interfaces
		Abstractor abs.Abstractor_ID
	}{}

	bs, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bs, &data); err != nil {
		panic(err)
	}

	base := Trainer{
		Current_Iter: data.Current_Iter,
		Total_Iter:   data.Total_Iter,
		Infoset_map:  data.Infoset_map,
		Num_players:  data.Num_players,
		Mu:           &sync.Mutex{},
		Wg:           &sync.WaitGroup{},
	}

	switch data.Abstractor {
	case abs.A1_ID:
		base.Abstractor = &abs.A1{}
	case abs.A2_ID:
		base.Abstractor = &abs.A2{}
	default:
		panic("abstraccion desconocida")
	}

	return Embed(t, &base)
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

func LoadModel(filename string, verbose bool, report_interval int) ITrainer {

	var t Trainer_T
	base := &Trainer{
		Current_Iter: 0,
		Total_Iter:   0,
		Infoset_map:  make(map[string]*RNode),
		Num_players:  0,
		Abstractor:   nil,
		Mu:           &sync.Mutex{},
		Wg:           &sync.WaitGroup{},
	}

	file, err := os.Open(filename)

	if err != nil {
		panic(fmt.Sprintf("failed to open model `%s`", filename))
	}

	if verbose {
		log.Printf("Leyendo tamano del modelo...\n")
	}

	locs, _ := lineCounter(filename)

	if verbose {
		log.Printf("%d lineas leidas\n", locs)
		log.Printf("Cargando modelo: 0%%")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	return_counter := 0
	i := 0

	for scanner.Scan() {
		line := scanner.Text()

		// caso 1: atributos
		if return_counter < 2 {

			if line == "" {
				return_counter++
				continue
			}

			words := strings.Fields(line)
			val := words[1]

			switch words[0] {
			case "ID":
				t = Trainer_T(words[1])
			case "Current_Iter":
				base.Current_Iter, _ = strconv.Atoi(val)
			case "Total_Iter":
				base.Total_Iter, _ = strconv.Atoi(val)
			case "Num_players":
				base.Num_players, _ = strconv.Atoi(val)
			case "Abstractor":
				if val == string(abs.A1_ID) {
					base.Abstractor = &abs.A1{}
				} else if val == string(abs.A2_ID) {
					base.Abstractor = &abs.A2{}
				} else if val == string(abs.NULL_ID) {
					base.Abstractor = &abs.Null{}
				} else {
					panic("abstraccion desconocida")
				}
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
			base.Infoset_map[hash] = rnode
		}
	}

	if verbose {
		log.Println()
	}

	return Embed(t, base)
}

// *****************************************************************************

type Trainer_T string

const (
	// vanilla
	CFR_T      Trainer_T = "cfr"
	CFRP_T     Trainer_T = "cfr+"
	DCFR_T     Trainer_T = "dcfr"
	ESLMCCFR_T Trainer_T = "es-lmccfr"
	ESVMCCFR_T Trainer_T = "es-vmccfr"
	OSMCCFR_T  Trainer_T = "os-lmccfr"
	// exploit
	BR_T Trainer_T = "best-response"
	// LCFR_T     Trainer_T = "lcfr" // <- DEPRECATED
)

func NewTrainer(t Trainer_T, num_players int, abs abs.IAbstraction) ITrainer {
	base := Trainer{
		Current_Iter: 0,
		Total_Iter:   0,
		Infoset_map:  make(map[string]*RNode),
		Num_players:  num_players,
		Abstractor:   abs,
		Mu:           &sync.Mutex{},
		Wg:           &sync.WaitGroup{},
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

	// case ESLMCCFR_T:
	// 	return &ESLMCCFR{base}

	case ESVMCCFR_T:
		return &ESVMCCFR{base}

	// case OSMCCFR_T:
	// 	return &OSMCCFR{base, 0.1} // epsilon

	// // exploitability
	// case BR_T:
	// 	return &BestResponse{base}

	default:
		panic("tainer unknown")
	}
}
