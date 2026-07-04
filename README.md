# rpgmap

CLI tool to generate leaflet markers or polygons, with tag-based layerGroups for autotoggles.

```bash
Usage of ./rpgmap:
  -c, --config string   Config file to read
```
## Point Markers
Point markers match the format below. `x` and `y` must be numeric. Tags must be one word, no punctuation.

`x,y,Text Comment,tag1,tag2`

## Polygon Markers
Polygon markers match the format below. `x` and `y` must be numeric, and must be sequential tuples. The last tuple will automatically close to the first tuple. Tags must be one word, no punctuation.

`p,x,y,x,y,x,y,x,y,Text Comment,tag1,tag2`

## Circle Markers
Circle markers match the format below. `x` and `y` must be numeric and reflect the center of the circle. 'radius' must be numeric. Tags must be one word, no punctuation.

`c,x,y,radius,Text Comment,tag1,tag2`

## Comments
Lines that start with `#` or `//` are ignored.
Lines that start with `/*` start a comment block, and all lines are ignored until *after* a line starts with `*/`
