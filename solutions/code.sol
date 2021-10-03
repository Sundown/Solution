@Package Solution
@Entry Start

Rec : Int -> Int = Return Rec @;

Start : Void -> Void =
	Rec 0;
