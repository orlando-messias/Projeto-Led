package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"math"
	"testing"

	"io"
	"os"
	"sort"

	"strconv"

	"github.com/stretchr/testify/assert"
)

func main() {
	err := rideAtPark()
	if err != nil {
		log.Println(err)
		return
	}
}

func rideAtPark() error {
	jurassic := initJurassic()

	// Load first file
	err := jurassic.loadBase()
	if err != nil {
		return err
	}

	// Load second file
	err = jurassic.loadAdd()
	if err != nil {
		return err
	}

	// Calc Speed
	err = jurassic.calcSpeed()
	if err != nil {
		return err
	}

	// Export Dinos
	jurassic.exportFastestDinos("output.txt")
	return err
}

type dino struct {
	name         string
	legLenght    float64
	strideLenght float64
	speed        float64
	stance       string
}

type jurassicPark struct {
	dinos map[string]dino
}

func initJurassic() *jurassicPark {
	return &jurassicPark{
		dinos: make(map[string]dino),
	}
}

// Considerando que o arquivo dataset2 é base do nosso problema, pois nele está a informação se o nosso dinossauro é bípede, carregamos apenas estes dinossauros
func (jp *jurassicPark) loadBase() error {
	fileLoad, err := parseCSV("dataset2.csv", "NAME", "STRIDE_LENGTH", "STANCE")
	if err != nil {
		return err
	}

	for i := 0; i < len(fileLoad["NAME"]); i++ {
		if fileLoad["STANCE"][i] != "bipedal" {
			continue
		}

		sLenght, err := strconv.ParseFloat(fileLoad["STRIDE_LENGTH"][i], 64)
		if err != nil {
			return err
		}

		jp.dinos[fileLoad["NAME"][i]] = dino{
			name:         fileLoad["NAME"][i],
			strideLenght: sLenght,
			stance:       fileLoad["STANCE"][i],
		}
	}

	return nil
}

// Nesta função trazemos a informação adicional do arquivo dataset1
func (jp *jurassicPark) loadAdd() error {
	fileLoad, err := parseCSV("dataset1.csv", "NAME", "LEG_LENGTH")
	if err != nil {
		return err
	}

	for i := 0; i < len(fileLoad["NAME"]); i++ {
		if dinossaur, ok := jp.dinos[fileLoad["NAME"][i]]; ok {
			lLenght, err := strconv.ParseFloat(fileLoad["LEG_LENGTH"][i], 64)
			if err != nil {
				return err
			}

			dinossaur.legLenght = lLenght

			jp.dinos[fileLoad["NAME"][i]] = dinossaur
		}
	}

	return nil
}

// Calcula velocidade
func (jp *jurassicPark) calcSpeed() error {
	for key, dino := range jp.dinos {
		if dino.strideLenght == 0 || dino.legLenght == 0 {
			continue
		}
		speed := ((dino.strideLenght / dino.legLenght) - 1) * math.Sqrt(dino.legLenght*math.Pow(9.8, 2))

		dino.speed = speed

		jp.dinos[key] = dino
	}

	return nil
}

// Ordena e exporta dinossauros
func (jp *jurassicPark) exportFastestDinos(file string) error {
	var Dinos []dino

	for _, value := range jp.dinos {
		if value.speed == 0 {
			continue
		}

		Dinos = append(Dinos, value)
	}

	sort.Slice(Dinos, func(i, j int) bool {
		return Dinos[i].speed > Dinos[j].speed
	})

	err := createFile(Dinos, file)

	return err
}

// Cria arquivo com os dinos informados
func createFile(dinos []dino, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	for i, d := range dinos {
		name := d.name

		if i != 0 {
			name = "\n" + name
		}

		_, err := f.Write([]byte(name))
		if err != nil {
			return err
		}
	}

	return nil
}

// parseCSV é uma função que recebe um CSV e um slice de colunas
// A função retorna uma matriz com as informações solicitadas
func parseCSV(file string, colsToReturn ...string) (map[string][]string, error) {
	csvFile, err := os.Open(file)
	if err != nil {
		return map[string][]string{}, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))

	mapCSV := make(map[string][]string)
	nrow := 0

	keys, err := reader.Read() // using header as keys
	if err != nil {
		return map[string][]string{}, err
	}

	mapKeyPosition := make(map[string]int)
	for i, v := range keys {
		mapKeyPosition[v] = i
	}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return map[string][]string{}, err
		}

		mapCSV["row"] = append(mapCSV["row"], strconv.Itoa(nrow))
		for _, key := range colsToReturn {
			mapCSV[key] = append(mapCSV[key], line[mapKeyPosition[key]])

		}
		nrow++
	}
	return mapCSV, nil
}

func TestFiles(t *testing.T) {
	jurassic := initJurassic()

	// Load first file
	err := jurassic.loadBase()
	assert.Nil(t, err)

	// Check if any dinossaur exists
	var ok string
	var dinossaur dino
	for ok, dinossaur = range jurassic.dinos {
		break
	}

	assert.True(t, ok != "")

	// Check if Leg_Lenght exists
	assert.NotNil(t, dinossaur.legLenght)

	// Load aditional file
	err = jurassic.loadAdd()
	assert.Nil(t, err)

	for ok, dinossaur = range jurassic.dinos {
		if dinossaur.strideLenght != 0 {
			break
		}
	}

	assert.True(t, ok != "")
	assert.NotNil(t, dinossaur.strideLenght)

}

func TestCalcSpeed(t *testing.T) {
	jurassic := initJurassic()

	jurassic.loadBase()
	jurassic.loadAdd()

	jurassic.calcSpeed()

	var ok string
	var dinossaur dino

	// Check if any dinossaur has speed calculated
	for ok, dinossaur = range jurassic.dinos {
		if dinossaur.speed != 0 {
			break
		}
	}

	assert.True(t, ok != "")
	assert.NotNil(t, dinossaur.speed)

	// Check if speed is right

	speedValidade := ((dinossaur.strideLenght / dinossaur.legLenght) - 1) * math.Sqrt(dinossaur.legLenght*math.Pow(9.8, 2))

	assert.Equal(t, dinossaur.speed, speedValidade)

}

func TestExportDinos(t *testing.T) {
	jurassic := initJurassic()

	jurassic.loadBase()
	jurassic.loadAdd()

	jurassic.calcSpeed()

	err := jurassic.exportFastestDinos("teste_file.txt")
	assert.Nil(t, err)

	// Check if file is created
	_, err = os.Stat("teste_file.txt")
	assert.False(t, os.IsNotExist(err))

}
