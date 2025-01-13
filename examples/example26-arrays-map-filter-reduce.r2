function main(){

    a = [1,5,3,2,4];
    println("a",a)
    a = a.sort(func (a,b){ b < a})
    println("a = a.sort(func (a,b){ b < a})")
    println("a",a)// tiene que ser 5,4,3,2,1
    a = a.sort() // es lo mismo que a.sort(func (a,b){ a < b})
    println("a = a.sort()  // es lo mismo que a.sort(func (a,b){ a < b})")
    println("a",a)// tiene que ser 1,2,3,4,5
    println('a.join("-") >>>',a.join("-")); // tiene que ser 1,2,3,4,5


    a = a.add(6);
    println("a.add(6) >>>",a);
    a = a.del(a.length()-1);
    println("a.del(a.length()-1) >>>",a);

    println("a.find(3)  >>>",a.find(3));
    println("a.find(func(v){ v==3 }) >>>",a.find(func(v){ v==3 }));
    println("a.find(func(v,p){ v==p },3) >>>",a.find(func(v,p){ v==p },3));



    println("a.reverse() >>>",a.reverse());
    println("a",a)// tiene que ser 5,4,3,2,1
    println("a.length >>>",a.length); // tiene que ser 1,2,3,4,5
    a = a.map(func(v){v*2}).filter(func(v){v<10}).reduce(func(v,c){v+c;});
    print("map -> filter -> reduce:",a); // tiene que ser 20  -> de map 2,4,6,8,10 -> de filter 2,4,6,8 -> de reduce 20


    let arr = ["Hola", "Mundo", "R2"];
    let joined = arr.join("-");
    print("arr.join('-') =>", joined); // "Hola-Mundo-R2"

}