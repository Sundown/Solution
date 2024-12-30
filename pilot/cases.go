package pilot

var Cases = []Test{
	{Name: "Plus reduce", Code: `Print +/1 1 1 2;`, Result: "5"},
	{Name: "Addition #1", Code: `Print 1 + -1;`, Result: "0"},
	{Name: "Map #1", Code: `Print ¨ 0 1;`, Result: "01"},
	{Name: "Train #1", Code: `Print 1.92 (⌊÷⌈) 5.6;`, Result: "0.342857"},
	{Name: "Train #2", Code: `Print 1.2(+++++)11.3;`, Result: "37.500000"},
	{Name: "Train #3", Code: `Print 1.2(⌈+⌊+⌊)11.3;`, Result: "13.700000"},
	{Name: "Vectored-op #1", Code: `Print+/13.4 4.1+5.1 4.1;`, Result: "26.700000"},
}
