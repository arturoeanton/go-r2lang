let true = 1==1
let false = 1==0

func testSuma() {
    let x = 2 + 2;
    assertEq(x, 4, "2+2 debería ser 4");
    assertTrue( x == 4, "x == 4");
}

func testBooleanos() {
    assertTrue(1 < 5, "1 < 5");
    assertEq(true, 1<5, "true vs cond");
}

func testFalla() {
    // Intencional
    assertEq("hola", "mundo", "esto fallará");
}

// Llamamos runAllTests() al final
func main() {
    runAllTests();
}