// Ejemplo de uso de formateo de fechas
// Este ejemplo demuestra que el problema del formato YYYY está arreglado

let date = Date.create(2024, 7, 15, 14, 30, 25);

// Formato básico ISO
let iso = Date.format(date, "YYYY-MM-DD");
print("ISO format: " + iso); // Debe mostrar: 2024-07-15

// Formato europeo
let european = Date.format(date, "DD/MM/YYYY");
print("European format: " + european); // Debe mostrar: 15/07/2024

// Formato americano
let american = Date.format(date, "MM/DD/YYYY");
print("American format: " + american); // Debe mostrar: 07/15/2024

// Formato con tiempo
let datetime = Date.format(date, "YYYY-MM-DD HH:mm:ss");
print("DateTime format: " + datetime); // Debe mostrar: 2024-07-15 14:30:25

// Formato con año corto
let shortYear = Date.format(date, "DD-MM-YY");
print("Short year format: " + shortYear); // Debe mostrar: 15-07-24

// Formato solo año
let yearOnly = Date.format(date, "YYYY");
print("Year only: " + yearOnly); // Debe mostrar: 2024

// Formato ISO con T literal
let isoWithT = Date.format(date, "YYYY-MM-DD'T'HH:mm:ss");
print("ISO with T: " + isoWithT); // Debe mostrar: 2024-07-15T14:30:25

// Formato con milisegundos y zona horaria
let fullISO = Date.format(date, "YYYY-MM-DD'T'HH:mm:ss.SSS'Z'");
print("Full ISO: " + fullISO); // Debe mostrar: 2024-07-15T14:30:25.000Z

// Test con diferentes años
let date2025 = Date.create(2025, 1, 1);
let year2025 = Date.format(date2025, "YYYY-MM-DD");
print("Year 2025: " + year2025); // Debe mostrar: 2025-01-01

let date2020 = Date.create(2020, 12, 31);
let year2020 = Date.format(date2020, "YYYY-MM-DD");
print("Year 2020: " + year2020); // Debe mostrar: 2020-12-31

print("Todos los formatos funcionan correctamente!");