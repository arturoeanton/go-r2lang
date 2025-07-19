function main(){

    a = [1,5,3,2,4];
    std.print("a",a)
    a = a.sort(func (a,b){ b < a})
    std.print("a = a.sort(func (a,b){ b < a})")
    std.print("a",a)// tiene que ser 5,4,3,2,1
    a = a.sort() // es lo mismo que a.sort(func (a,b){ a < b})
    std.print("a = a.sort()  // es lo mismo que a.sort(func (a,b){ a < b})")
    std.print("a",a)// tiene que ser 1,2,3,4,5
    std.print('a.join("-") >>>',a.join("-")); // tiene que ser 1,2,3,4,5


    a = a.add(6);
    std.print("a.add(6) >>>",a);
    a = a.del(a.length()-1);
    std.print("a.del(a.length()-1) >>>",a);

    std.print("a.find(3)  >>>",a.find(3));
    std.print("a.find(func(v){ v==3 }) >>>",a.find(func(v){ v==3 }));
    std.print("a.find(func(v,p){ v==p },3) >>>",a.find(func(v,p){ v==p },3));



    std.print("a.reverse() >>>",a.reverse());
    std.print("a",a)// tiene que ser 5,4,3,2,1
    std.print("a.length >>>",a.length); // tiene que ser 1,2,3,4,5
    a = a.map(func(v){v*2}).filter(func(v){v<10}).reduce(func(v,c){v+c;});
    std.print("map -> filter -> reduce:",a); // tiene que ser 20  -> de map 2,4,6,8,10 -> de filter 2,4,6,8 -> de reduce 20


    let arr = ["Hola", "Mundo", "R2"];
    let joined = arr.join("-");
    std.print("arr.join('-') =>", joined); // "Hola-Mundo-R2"

}