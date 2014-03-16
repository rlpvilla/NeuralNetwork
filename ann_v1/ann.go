package ann

import (
	
)

type SynapseDef struct {
	Input chan float
	Output chan float
	Feedput chan float
	StartWeight float
}