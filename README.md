# rpgmap

CLI tool to generate [leaflet](https://leafletjs.com/) markers, polygons, and circles; With tag-based layerGroups for autotoggles, and tag-based custom icons.

```bash
Usage of ./rpgmap:
  -c, --config string   Config file to read
```

The output of `rpgmap` is a valid JS file for use with a predefined map var named `map`. 

Below is a terse, non-working snippet to demonstrate. (The header `link` and `script` need to be corrected. `var imageUrl` needs to point to a valid image for the map, the `bounds` need to be correct, etc. etc.)

You can use tiles, projections, whatever.

```html
<html>
<head>
    <link rel="stylesheet" href="https://unpkg.com/leaflet.css>
    <script src="https://unpkg.com/leaflet.js"></script>
</head>
<body>


<div id="map" style="width: 100%; height: 100%;"></div>

<script>
    var map = L.map('map', {
        crs: L.CRS.Simple,
        minZoom: -3,
        maxZoom: 3 
    });

    var bounds = [[0, 0], [6144, 8192]];

    var imageUrl = 'yermap.png';
    L.imageOverlay(imageUrl, bounds).addTo(map);
    map.fitBounds(bounds);
</script>

<script src="rpgmap.out.js"></script>

</body>
</html>
```


## Point Markers
Point markers match the format below. `x` and `y` must be numeric. Tags must be one word, no punctuation. The first tag is used to set the icon, if any.

`x,y,Text Comment,tag1,tag2`

### Point Marker icons
To define an icon for a tag type, use the format below. `x` and `y` must be numeric, and must be sequential tuples: The first tuple is *required* and is the width and height of the icon; The second tuple is optional, and references which point on the icon will be anchored to the spot. The default is `0,0`. Only the first tag in a Point Marker will be used when determining which icon to set.

**WARNING**: Icons need to be defined before a PointMarker that uses the icon, or it will be missed. (Yes, I could do a first-pass to find all of the icon lines, but I don't want to. Put them up top. Now. Stop reading. Go do it.)

`i,tag,uri,x,y,x,y`

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

If you put the configuration stanzas below into a file called `ex`

```
// Circle!
c,-38.85682, -24.98291,2250000,Whirlpool,nature,dangers

// icon!
i,places,icons/places.png,50,50,25,25

// Point!
39.027719, 52.382813,That Place,places

// Polygon!
p,15.45368, 3.339844,27.722436, 29.003906,38.410558, 42.1875,36.491973, 55.063477,31.728167, 58.31543,30.826781, 63.588867,14.689881, 82.045898,13.795406, 91.625977,3.908099, 95.317383,3.294082, 104.589844,-6.8828, 103.447266,-14.85985, 106.391602,-24.806681, 107.841797,-33.687782, 89.692383,-33.100745, 71.411133,-5.790897, 22.324219,-6.489983, 4.130859,That Region,regions
```

and then process it:

```bash
$ rpgmap -c ex > rpgmap.out.js
$ cat rpgmap.out.js
var layerControl = L.control.layers().addTo(map)
var placesIcon = L.icon({ iconUrl: 'icons/places.png', iconSize: [50.000000,50.000000], iconAnchor: [25.000000,25.000000] });
var dangers = L.layerGroup([
	L.circle([-38.856820,-24.982910],{color: 'red',fillColor: '#f03',radius: 2250000.000000,fillOpacity: 0.200000}).bindPopup("Whirlpool")])
layerControl.addOverlay(dangers, "Dangers");
var nature = L.layerGroup([
	L.circle([-38.856820,-24.982910],{color: 'red',fillColor: '#f03',radius: 2250000.000000,fillOpacity: 0.200000}).bindPopup("Whirlpool")])
layerControl.addOverlay(nature, "Nature");
var places = L.layerGroup([
	L.marker([39.027719,52.382813],{icon: placesIcon}).bindPopup("That Place")])
layerControl.addOverlay(places, "Places");
var regions = L.layerGroup([
	L.polygon([[15.453680,3.339844],[27.722436,29.003906],[38.410558,42.187500],[36.491973,55.063477],[31.728167,58.315430],[30.826781,63.588867],[14.689881,82.045898],[13.795406,91.625977],[3.908099,95.317383],[3.294082,104.589844],[-6.882800,103.447266],[-14.859850,106.391602],[-24.806681,107.841797],[-33.687782,89.692383],[-33.100745,71.411133],[-5.790897,22.324219],[-6.489983,4.130859]],{color: 'red',fillColor: '#f03',fillOpacity: 0.200000}).bindPopup("That Region")])
```

You get a valid [leaflet](https://leafletjs.com/) JS file to reference, that autobuilds the stanzas necessary for layers, icons, etc.

## Thoughts/TODOish

### Styles
Styles are hardcoded, although the underlying types have it abstracted for futureproofing. ~~PointMarkers could use icons~~, poly and circles should allow custom colors and opacity.

~~Should that be per tag? If so, what if there are multiple tags with conflicting styles?~~ First tag wins (e.g. primary vs secondaries)? Even that could be nightmarish with poly's/circles as the stacking opacity would be obscuring, I believe (CONFIRMED).