# Propuesta: Soporte Nativo de Fechas y Tiempo en R2Lang

**Versión:** 1.0  
**Fecha:** 2025-07-15  
**Estado:** Propuesta  

## Resumen Ejecutivo

Esta propuesta introduce un sistema completo de manejo de fechas y tiempo en R2Lang, incluyendo tipos nativos de fecha, funciones de manipulación, formateo, zona horaria y operaciones aritméticas temporales.

## Problema Actual

R2Lang actualmente carece de soporte nativo para fechas y tiempo:

- **Sin tipo Date:** Solo strings y números para representar fechas
- **Operaciones complejas:** Cálculos manuales con timestamps
- **Formateo limitado:** Sin formateo estándar de fechas
- **Zonas horarias:** Sin soporte para diferentes husos horarios
- **Validación:** Sin validación de fechas

### Ejemplos de Problemas Actuales

```r2
// Actualmente en R2Lang - complejo y propenso a errores
let fechaStr = "2024-12-25";
let timestamp = 1703462400; // ¿Qué fecha es esta?
let año = parseInt(substr(fechaStr, 0, 4));

// Comparar fechas es complicado
func esFechaMayor(fecha1, fecha2) {
    // Lógica compleja manual...
}

// Agregar días requiere cálculos manuales
let mañana = timestamp + (24 * 60 * 60 * 1000);
```

## Solución Propuesta

### 1. Tipo Date Nativo

#### 1.1 Implementación del Tipo Date
```go
// pkg/r2core/date.go
type DateValue struct {
    Time     time.Time
    Location *time.Location
}

func (dv *DateValue) String() string {
    return dv.Time.Format("2006-01-02T15:04:05Z07:00")
}

func (dv *DateValue) Type() string {
    return "Date"
}

func (dv *DateValue) Equals(other interface{}) bool {
    if otherDate, ok := other.(*DateValue); ok {
        return dv.Time.Equal(otherDate.Time)
    }
    return false
}
```

#### 1.2 Parsing en el Lexer
```go
// pkg/r2core/lexer.go
func (l *Lexer) readDateLiteral() Token {
    // Reconocer literales de fecha como @2024-12-25 o @"2024-12-25T10:30:00"
    l.nextChar() // saltar @
    
    if l.ch == '"' || l.ch == '\'' {
        // Fecha con formato específico
        dateStr := l.readString()
        return Token{Type: TOKEN_DATE, Value: dateStr}
    } else {
        // Fecha simple YYYY-MM-DD
        dateStr := l.readSimpleDate()
        return Token{Type: TOKEN_DATE, Value: dateStr}
    }
}

func (l *Lexer) readSimpleDate() string {
    position := l.position
    // Leer hasta encontrar un delimitador válido
    for l.ch != ' ' && l.ch != '\n' && l.ch != '\t' && l.ch != 0 {
        l.nextChar()
    }
    return l.input[position:l.position]
}
```

### 2. Sintaxis de Fechas en R2Lang

#### 2.1 Literales de Fecha
```r2
// Diferentes formas de crear fechas
let fecha1 = @2024-12-25;                    // Solo fecha
let fecha2 = @"2024-12-25 14:30:00";         // Fecha y hora
let fecha3 = @"2024-12-25T14:30:00Z";        // ISO 8601
let fecha4 = @"2024-12-25T14:30:00-05:00";   // Con zona horaria

// Fecha actual
let ahora = Date.now();
let hoy = Date.today();

// Crear desde componentes
let navidad = Date.create(2024, 12, 25);
let reunion = Date.create(2024, 12, 25, 14, 30, 0);
```

#### 2.2 Constructores de Date
```r2
// Diferentes constructores
let fecha = new Date();                      // Fecha actual
let fecha = new Date("2024-12-25");          // Desde string
let fecha = new Date(2024, 12, 25);          // Desde componentes
let fecha = new Date(1703462400000);         // Desde timestamp
```

### 3. Biblioteca de Funciones de Fecha

#### 3.1 Funciones Básicas
```go
// pkg/r2libs/r2date.go
func RegisterDate(env *r2core.Environment) {
    // Constructor y utilidades
    env.Set("Date", createDateConstructor())
    
    // Funciones estáticas
    env.Set("now", r2core.BuiltinFunction(dateNow))
    env.Set("today", r2core.BuiltinFunction(dateToday))
    env.Set("parse", r2core.BuiltinFunction(dateParse))
    env.Set("isValid", r2core.BuiltinFunction(dateIsValid))
    
    // Funciones de comparación
    env.Set("isBefore", r2core.BuiltinFunction(dateIsBefore))
    env.Set("isAfter", r2core.BuiltinFunction(dateIsAfter))
    env.Set("isSame", r2core.BuiltinFunction(dateIsSame))
    env.Set("diff", r2core.BuiltinFunction(dateDiff))
    
    // Funciones de manipulación
    env.Set("add", r2core.BuiltinFunction(dateAdd))
    env.Set("subtract", r2core.BuiltinFunction(dateSubtract))
    env.Set("startOf", r2core.BuiltinFunction(dateStartOf))
    env.Set("endOf", r2core.BuiltinFunction(dateEndOf))
}

func dateNow(args ...interface{}) interface{} {
    return &r2core.DateValue{
        Time:     time.Now(),
        Location: time.Local,
    }
}

func dateParse(args ...interface{}) interface{} {
    if len(args) < 1 {
        panic("Date.parse() requiere al menos 1 argumento")
    }
    
    dateStr, ok := args[0].(string)
    if !ok {
        panic("Date.parse() requiere un string")
    }
    
    format := time.RFC3339
    if len(args) > 1 {
        if formatStr, ok := args[1].(string); ok {
            format = convertToGoFormat(formatStr)
        }
    }
    
    parsedTime, err := time.Parse(format, dateStr)
    if err != nil {
        panic("Fecha inválida: " + dateStr)
    }
    
    return &r2core.DateValue{
        Time:     parsedTime,
        Location: time.Local,
    }
}
```

#### 3.2 Métodos de Date
```go
// Métodos disponibles en objetos Date
func (env *Environment) setupDateMethods() {
    dateObj := map[string]interface{}{
        "getYear":        r2core.BuiltinFunction(getYear),
        "getMonth":       r2core.BuiltinFunction(getMonth),
        "getDay":         r2core.BuiltinFunction(getDay),
        "getHour":        r2core.BuiltinFunction(getHour),
        "getMinute":      r2core.BuiltinFunction(getMinute),
        "getSecond":      r2core.BuiltinFunction(getSecond),
        "getDayOfWeek":   r2core.BuiltinFunction(getDayOfWeek),
        "getDayOfYear":   r2core.BuiltinFunction(getDayOfYear),
        "getWeekOfYear":  r2core.BuiltinFunction(getWeekOfYear),
        "format":         r2core.BuiltinFunction(dateFormat),
        "toString":       r2core.BuiltinFunction(dateToString),
        "toISOString":    r2core.BuiltinFunction(dateToISO),
        "valueOf":        r2core.BuiltinFunction(dateValueOf),
        "setTimezone":    r2core.BuiltinFunction(setTimezone),
    }
    
    env.Set("DatePrototype", dateObj)
}

func getYear(args ...interface{}) interface{} {
    date := getDateFromArgs(args, "getYear")
    return float64(date.Time.Year())
}

func dateFormat(args ...interface{}) interface{} {
    if len(args) < 2 {
        panic("format() requiere al menos 2 argumentos")
    }
    
    date := getDateFromArgs(args[:1], "format")
    format, ok := args[1].(string)
    if !ok {
        panic("format() requiere un string como formato")
    }
    
    goFormat := convertToGoFormat(format)
    return date.Time.Format(goFormat)
}
```

### 4. Operaciones Aritméticas con Fechas

#### 4.1 Aritmética de Fechas
```r2
let inicio = @2024-01-01;
let fin = @2024-12-31;

// Diferencia entre fechas
let diferencia = fin - inicio;           // Retorna Duration
print(diferencia.days());                // 365 días
print(diferencia.hours());               // 8760 horas

// Agregar tiempo
let mañana = inicio + Duration.days(1);
let proximaSemana = inicio + Duration.weeks(1);
let proximoMes = inicio + Duration.months(1);

// Operadores de comparación
if (fin > inicio) {
    print("Fin es después de inicio");
}
```

#### 4.2 Objeto Duration
```go
// pkg/r2core/duration.go
type DurationValue struct {
    Duration time.Duration
}

func (dv *DurationValue) String() string {
    return dv.Duration.String()
}

func (dv *DurationValue) Days() float64 {
    return dv.Duration.Hours() / 24
}

func (dv *DurationValue) Hours() float64 {
    return dv.Duration.Hours()
}

func (dv *DurationValue) Minutes() float64 {
    return dv.Duration.Minutes()
}
```

### 5. Formateo Avanzado de Fechas

#### 5.1 Formatos Predefinidos
```r2
let fecha = @"2024-12-25 14:30:00";

// Formatos predefinidos
print(fecha.format("short"));        // 25/12/24
print(fecha.format("medium"));       // 25 dic 2024
print(fecha.format("long"));         // 25 de diciembre de 2024
print(fecha.format("full"));         // miércoles, 25 de diciembre de 2024

// Formatos personalizados
print(fecha.format("YYYY-MM-DD"));   // 2024-12-25
print(fecha.format("DD/MM/YYYY"));   // 25/12/2024
print(fecha.format("MMM DD, YYYY")); // dic 25, 2024
print(fecha.format("dddd, MMMM Do, YYYY")); // miércoles, diciembre 25°, 2024
```

#### 5.2 Localización de Fechas
```r2
// Configurar idioma
Date.setLocale("es-ES");
print(fecha.format("long")); // 25 de diciembre de 2024

Date.setLocale("en-US");
print(fecha.format("long")); // December 25, 2024

Date.setLocale("fr-FR");
print(fecha.format("long")); // 25 décembre 2024
```

### 6. Soporte de Zonas Horarias

#### 6.1 Manejo de Zonas Horarias
```r2
// Crear fecha en zona específica
let utc = Date.utc(2024, 12, 25, 14, 30);
let madrid = Date.timezone("Europe/Madrid", 2024, 12, 25, 14, 30);
let tokyo = Date.timezone("Asia/Tokyo", 2024, 12, 25, 14, 30);

// Convertir entre zonas
let fechaLocal = utc.toTimezone("America/New_York");
print(fechaLocal.format("YYYY-MM-DD HH:mm:ss Z"));

// Información de zona horaria
print(madrid.getTimezone());     // "Europe/Madrid"
print(madrid.getOffset());       // +01:00
print(madrid.isDST());           // true/false
```

#### 6.2 Lista de Zonas Horarias
```r2
// Obtener zonas disponibles
let zonas = Date.getTimezones();
print(zonas); // ["UTC", "America/New_York", "Europe/Madrid", ...]

// Información detallada de zona
let info = Date.getTimezoneInfo("Europe/Madrid");
print(info.offset);     // "+01:00"
print(info.dst);        // true
print(info.name);       // "Central European Time"
```

### 7. Funciones de Utilidad

#### 7.1 Validación y Cálculos
```r2
// Validación
print(Date.isValid("2024-02-29"));  // true (año bisiesto)
print(Date.isValid("2023-02-29"));  // false

print(Date.isLeapYear(2024));       // true
print(Date.daysInMonth(2024, 2));   // 29

// Cálculos útiles
let fecha = @2024-12-25;
print(fecha.isWeekend());           // true/false
print(fecha.isHoliday("US"));       // true si es navidad en US
print(fecha.getQuarter());          // 4

// Fechas relativas
let inicio = Date.startOfWeek();    // Lunes de esta semana
let fin = Date.endOfMonth();        // Último día del mes
let ayer = Date.yesterday();
let mañana = Date.tomorrow();
```

#### 7.2 Parsing Inteligente
```r2
// Parsing flexible
let f1 = Date.parse("mañana");
let f2 = Date.parse("en 3 días");
let f3 = Date.parse("el próximo lunes");
let f4 = Date.parse("hace 2 semanas");
let f5 = Date.parse("2024-12-25");
let f6 = Date.parse("25 de diciembre de 2024");
```

### 8. Integración con el Sistema

#### 8.1 Operadores Sobrecargados
```go
// pkg/r2core/binary_expression.go
func (be *BinaryExpression) evalDateOperations(left, right interface{}, env *Environment) interface{} {
    leftDate, leftIsDate := left.(*DateValue)
    rightDate, rightIsDate := right.(*DateValue)
    rightDuration, rightIsDuration := right.(*DurationValue)
    
    switch be.Operator {
    case "+":
        if leftIsDate && rightIsDuration {
            return &DateValue{
                Time:     leftDate.Time.Add(rightDuration.Duration),
                Location: leftDate.Location,
            }
        }
    case "-":
        if leftIsDate && rightIsDate {
            diff := leftDate.Time.Sub(rightDate.Time)
            return &DurationValue{Duration: diff}
        }
        if leftIsDate && rightIsDuration {
            return &DateValue{
                Time:     leftDate.Time.Add(-rightDuration.Duration),
                Location: leftDate.Location,
            }
        }
    case "<", "<=", ">", ">=", "==", "!=":
        if leftIsDate && rightIsDate {
            return compareDates(leftDate, rightDate, be.Operator)
        }
    }
    
    return nil // No es operación de fecha
}
```

### 9. Plan de Implementación

#### Fase 1: Tipo Date Básico (3-4 días)
- [ ] Implementar DateValue y DurationValue
- [ ] Soporte en lexer para literales @fecha
- [ ] Parser integration
- [ ] Tests básicos

#### Fase 2: Funciones Core (4-5 días)
- [ ] Date.now(), Date.parse(), etc.
- [ ] Métodos básicos (getYear, getMonth, etc.)
- [ ] Operaciones aritméticas básicas
- [ ] Tests de funcionalidad

#### Fase 3: Formateo y Localización (3-4 días)
- [ ] Sistema de formateo
- [ ] Soporte de locales
- [ ] Formatos predefinidos
- [ ] Tests de formateo

#### Fase 4: Zonas Horarias (4-5 días)
- [ ] Soporte de timezone
- [ ] Conversión entre zonas
- [ ] DST handling
- [ ] Tests de timezone

#### Fase 5: Funciones Avanzadas (3-4 días)
- [ ] Parsing inteligente
- [ ] Validaciones
- [ ] Funciones de utilidad
- [ ] Documentación completa

## Beneficios

1. **Simplicidad:** API intuitiva para fechas
2. **Potencia:** Funcionalidad completa para aplicaciones modernas
3. **Estándares:** Compatibilidad con ISO 8601 y RFC 3339
4. **Internacionalización:** Soporte multi-idioma
5. **Rendimiento:** Implementación eficiente basada en Go

## Consideraciones

- **Tamaño:** Aumento moderado en tamaño del binario
- **Complejidad:** Manejo cuidadoso de zonas horarias
- **Compatibilidad:** Mantener funciones de tiempo existentes
- **Rendimiento:** Optimización para operaciones comunes

## Ejemplos de Uso

```r2
// Aplicación de calendario
func calcularDiasLaborables(inicio, fin) {
    let dias = 0;
    let actual = inicio;
    
    while (actual <= fin) {
        if (!actual.isWeekend()) {
            dias++;
        }
        actual = actual + Duration.days(1);
    }
    
    return dias;
}

// Cálculo de edad
func calcularEdad(fechaNacimiento) {
    let hoy = Date.today();
    let diff = hoy - fechaNacimiento;
    return Math.floor(diff.days() / 365.25);
}

// Recordatorios
func proximoCumpleanos(fechaNacimiento) {
    let hoy = Date.today();
    let cumpleEsteAño = Date.create(hoy.getYear(), fechaNacimiento.getMonth(), fechaNacimiento.getDay());
    
    if (cumpleEsteAño < hoy) {
        cumpleEsteAño = Date.create(hoy.getYear() + 1, fechaNacimiento.getMonth(), fechaNacimiento.getDay());
    }
    
    return cumpleEsteAño - hoy;
}
```

## Conclusión

Esta implementación proporcionará un sistema completo y moderno de manejo de fechas en R2Lang, facilitando el desarrollo de aplicaciones que requieren manipulación temporal compleja mientras mantiene una API simple e intuitiva.