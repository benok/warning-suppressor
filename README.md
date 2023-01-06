# warning-suppressor
A wrapper to suppress unneeded output (e.g. spurious warning) without modifying original command.

(Currently, only Windows environment is supported.)

## How to use
1. Rename original executable from ExeName.exe to  ExeName_orig.exe.
2. Copy warning-suppressor.exe as ExeName.exe
3. Create ExeName.exe.yml and configure

## Features

### Suppress output 

Suppress output of matched patterns under `suppress`.

```yaml
filter-config:
  suppress:
    - "LINN32:" # Warning: LINN32: Last line 3868 (F1Ch) is less than or equal to first line 3868 (F1Ch) for symbol "chparse()" in module xxx.cpp
    - "Warning: public.*clashes with prior module" # [TLIB Warning] Warning: public 'glui_img_checkbox_0' in module '..\_obj\rs100\x64\Debug\glui\glui_img_checkbox_0.o' clashes with prior module 'glui_bitmap_img_data.o'
```

### Colorize output

Colorize output of matched patterns to specified color under `colorize`.

```yaml
filter-config:
  colorize:
    - "Turbo Incremental Link": "cyan"
    - "(Fatal|Error)": "red"
    - "Warning": "yellow"
```



