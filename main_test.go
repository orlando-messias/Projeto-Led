package main

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testes unitários

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
