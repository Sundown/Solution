@Package Solution
@Entry Main

Main : Void -> Void =
	Println "should be 15: ";
	Println Sum [1, 2, 3, 4, 5];
	Println "should be 10.0: ";
	Println Sum [1.5, 3.3, 2.7, 2.5];