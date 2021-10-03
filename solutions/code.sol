@Package Solution
@Entry Start

Start : Void -> Void =
	Println Append (
		GEP (["A", "B", "C"], Sum [-2, 5, 43, -44]),
		GEP ([['d', 'c', 'a'], "abc"], 0));
