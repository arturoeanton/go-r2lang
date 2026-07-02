let dateObj = date.Date();
let failures = 0;

func check(name, condition) {
    if (condition) {
        std.print("OK  ", name);
    } else {
        std.print("FAIL", name);
        failures = failures + 1;
    }
}

// isLeapYear
check("isLeapYear(2024) true", dateObj.isLeapYear(2024) == true);
check("isLeapYear(1900) false", dateObj.isLeapYear(1900) == false);
check("isLeapYear(2000) true", dateObj.isLeapYear(2000) == true);

// daysInMonth (0-indexed month)
check("daysInMonth(2024, 1) == 29", dateObj.daysInMonth(2024, 1) == 29);
check("daysInMonth(2023, 1) == 28", dateObj.daysInMonth(2023, 1) == 28);
check("daysInMonth(2024, 11) == 31", dateObj.daysInMonth(2024, 11) == 31);

// startOfDay / endOfDay
let mid = dateObj.create(2024, 6, 15, 14, 30, 25);
let start = dateObj.startOfDay(mid);
let end = dateObj.endOfDay(mid);
check("startOfDay hour == 0", dateObj.getHours(start) == 0);
check("startOfDay minute == 0", dateObj.getMinutes(start) == 0);
check("endOfDay hour == 23", dateObj.getHours(end) == 23);
check("endOfDay minute == 59", dateObj.getMinutes(end) == 59);
check("original date not mutated by startOfDay", dateObj.getHours(mid) == 14);

// startOfMonth / endOfMonth (December rollover, month/year handling)
let dec = dateObj.create(2024, 11, 10);
let decStart = dateObj.startOfMonth(dec);
let decEnd = dateObj.endOfMonth(dec);
check("startOfMonth day == 1", dateObj.getDate(decStart) == 1);
check("endOfMonth day == 31 (December)", dateObj.getDate(decEnd) == 31);
check("endOfMonth stays in same month", dateObj.getMonth(decEnd) == 11);

let feb = dateObj.create(2024, 1, 5);
let febEnd = dateObj.endOfMonth(feb);
check("endOfMonth day == 29 (leap February)", dateObj.getDate(febEnd) == 29);

// isWeekend
let saturday = dateObj.create(2024, 6, 13);
let monday = dateObj.create(2024, 6, 15);
check("isWeekend(Saturday) true", dateObj.isWeekend(saturday) == true);
check("isWeekend(Monday) false", dateObj.isWeekend(monday) == false);

// clone (independent from the original, unlike setters which mutate in place)
let original = dateObj.create(2024, 6, 15);
let cloned = dateObj.clone(original);
dateObj.setFullYear(cloned, 2030);
check("clone is independent of original", dateObj.getFullYear(original) == 2024);
check("clone reflects its own mutation", dateObj.getFullYear(cloned) == 2030);

// comparison operators already work on DateValue (no isBefore/isAfter/isSame needed)
check("comparison operators work on dates", saturday < monday);
check("equality operator works on dates", dateObj.create(2024, 6, 15) == dateObj.create(2024, 6, 15));

if (failures == 0) {
    std.print("DATE_SMOKE_OK");
} else {
    std.print("DATE_SMOKE_FAIL", failures);
}
