// Ejemplo de uso de fechas en r2lang

let hoy = Date.today();
print("Hoy es: " + hoy);

let navidad = Date.create(2024, 12, 25);
print("Navidad es: " + navidad);

print("El año de navidad es: " + Date.getYear(navidad));
print("El mes de navidad es: " + Date.getMonth(navidad));
print("El dia de navidad es: " + Date.getDay(navidad));

let ahora = Date.now();
print("Ahora mismo es: " + ahora);

// Formateo de fechas
let formato = "YYYY-MM-DD HH:mm:ss";
let fecha_formateada = Date.format(ahora, formato);
print("Fecha formateada: " + fecha_formateada);

// Zonas horarias
let tokio = Date.timezone("Asia/Tokyo", 2025, 1, 1, 10, 0, 0);
print("Año nuevo en Tokio: " + tokio);

let tokio_en_ny = Date.toTimezone(tokio, "America/New_York");
print("Año nuevo de Tokio en Nueva York: " + tokio_en_ny);

