@Package Solution
@Entry Start

Printandspace :: Int -> Void
Printandspace =
	Print @;
	Print ", ";

Printintvec :: [Int] -> Void
Printintvec =
	Print "[";
	Map (Printandspace, @);
	Print "]";


Start :: Void -> Void
Start =
	Printintvec [1, 2, 3, 4];
	Println "";
