// Ejemplo del módulo "graph": grafo dirigido simple para modelar
// relaciones/jerarquías (por ejemplo, un árbol genealógico) desde un script.

std.print("🌐 R2Lang: módulo graph")
std.print("=========================")

let family = graph.new()
family.addEdge("Federico", "Elias")
family.addEdge("Federico", "Eugenia")
family.addEdge("Elias", "Sara")
family.addEdge("Elias", "Arturo")
family.addEdge("Sara", "Telma")

std.print("\nAncestros de Sara:", family.getAncestors("Sara"))
std.print("Descendientes de Elias:", family.getDescendants("Elias"))
std.print("Nivel Federico -> Telma:", family.getRelationshipLevel("Federico", "Telma"))
std.print("Camino más corto Federico -> Telma:", family.getShortestPath("Federico", "Telma"))
std.print("Sin camino Telma -> Federico:", family.getShortestPath("Telma", "Federico"))

// addBidirectionalEdge agrega la arista en ambos sentidos de una sola vez.
let social = graph.new()
social.addBidirectionalEdge("Ana", "Beto")
social.addBidirectionalEdge("Beto", "Cami")

std.print("\nAna conoce (directo o indirecto) a:", social.getDescendants("Ana"))
std.print("Nivel Ana -> Cami:", social.getRelationshipLevel("Ana", "Cami"))

std.print("\n✅ Ejemplo de graph completado")
