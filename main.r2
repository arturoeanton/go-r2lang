

func main(){
    let arr = [1, 2, 3];
    arr = arr.add( 5, 4);
    arr = arr.sort(func (a,b){ b < a});
    print("arr = arr.sort(func (a,b){ b < a});", arr); // [5, 4, 3, 2, 1]
    arr = arr.insert_at( 2, 4); // [5, 4, 4, 3, 2, 1]
    arr = arr.insert_at( 2, 0); // [5, 4, 0, 4, 3, 2, 1]
    print("arr = arr.insert_at( 2, 0); >>", arr); // [5, 4, 0, 4, 3, 2, 1]
    print("arr.indexes(func (v){v==4}).len() >>", arr.indexes(func (v){v==4}).len()); // 2
}