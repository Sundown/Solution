@Package Solution
@Entry Start

F :: (Int, Int) -> Bool
F = Return Equals @;

P :: Bool -> Void
P = Println @;

Start :: Void -> Void
Start = Map (F, [(1, 1)]);
