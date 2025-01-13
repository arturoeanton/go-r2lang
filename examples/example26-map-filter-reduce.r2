function main(){
    a = [1,2,3,4,5];
    println("a.length >>>",a.length); // tiene que ser 1,2,3,4,5
    a = a.map(func(v){v*2}).filter(func(v){v<10}).reduce(func(v,c){v+c;});
    print("map -> filter -> reduce:",a); // tiene que ser 20  -> de map 2,4,6,8,10 -> de filter 2,4,6,8 -> de reduce 20
}