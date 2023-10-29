# 🧰 toolbox

Toolbox with small CLI utilities.

- cloudlogs-json-format - format JSON dumps from Google Cloud Logs into pretty lines
- uniqcut -  like `uniq` but considering only a part of the line, e.g. only content after a timestamp.
- [format-columns](scripting/format-columns) - format CSV into columns
- [cutcsv](scripting/cutcsv) - Use specific columns from CSV


# `tscalc`

Basic calculations on times.

```bash
% bin/tscalc
2023-10-29T19:40:09+00:00

% echo "2023-10-29T19:39:09+00:00 + 100s" | bin/tscalc
2023-10-29T19:40:49+00:00

% echo "2023-10-29T19:40:39+00:00 - 2022-09-28T18:22:32+00:00" | bin/tscalc
9505h18m7s
```

Works on timestamps as well:
```bash
% echo "12345 - 2345" | ./bin/tscalc
2h46m40s

% echo "2023-10-29T19:42:44+00:00 - 1698603564.000000" | ./bin/tscalc
1h23m20s
```
