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

## Example

```
// Circle!
c,-38.85682, -24.98291,2250000,Whirlpool,nature,dangers

// Point!
39.027719, 52.382813,That Place,places

// Polygon!
p,15.45368, 3.339844,27.722436, 29.003906,38.410558, 42.1875,36.491973, 55.063477,31.728167, 58.31543,30.826781, 63.588867,14.689881, 82.045898,13.795406, 91.625977,3.908099, 95.317383,3.294082, 104.589844,-6.8828, 103.447266,-14.85985, 106.391602,-24.806681, 107.841797,-33.687782, 89.692383,-33.100745, 71.411133,-5.790897, 22.324219,-6.489983, 4.130859,That Region,regions
```

## Thoughts/TODOish

### Styles
Styles are hardcoded, although the underlying types have it abstracted for futureproofing. PointMarkers could use icons, poly and circles should allow custom colors and opacity.

Should that be per tag? If so, what if there are multiple tags with conflicting styles? First tag wins (e.g. primary vs secondaries)? Even that could be nightmarish with poly's/circles as the stacking opacity would be obscuring, I believe (CONFIRMED).