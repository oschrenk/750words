# 750words

This tool helps split the monthly text exports of 750words.com into daily markdown files.

## Usage

```
750words split --version old input_file.txt existing_output_dir
750words split --version new input_file.txt existing_output_dir
```

## Format

### Old

https://750words.com/

```
------ ENTRY ------
Date:    2019-11-24
Words:   752
Minutes: 18

...

------ ENTRY ------
Date:    2019-11-25
Words:   756
Minutes: 22

...
```

## New

https://new.750words.com/

```
===== ENTRY =====
=== DATE: 2023-01-01 ===
=== TITLE: Sunday, January  1, 2023 ===
=== WORDS: 750 ===
=== MINUTES: 28481 ===
=== BODY ===

...

===== ENTRY =====
=== DATE: 2023-01-02 ===
=== TITLE: Monday, January  2, 2023 ===
=== WORDS: 752 ===
=== MINUTES: 28 ===
=== BODY ===

...
``` 
