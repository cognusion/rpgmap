var layerControl = L.control.layers().addTo(map)


var camp = L.layerGroup([
	L.polygon([[992.000000,2872.000000],[1700.000000,2872.000000],[1700.000000,3752.000000],[992.000000,3752.000000]],{color: 'red'}).bindPopup("<i>Camp!</i>"),
	L.marker([1256.000000,3064.000000]).bindPopup("Stump"),
	L.marker([1144.000000,3328.000000]).bindPopup("Tent"),
	L.marker([1584.000000,3240.000000]).bindPopup("Fire")])
layerControl.addOverlay(camp, "Camp");

var fish = L.layerGroup([
	L.marker([3048.000000,2672.000000]).bindPopup("Fish"),
	L.marker([2144.000000,2424.000000]).bindPopup("Fish"),
	L.marker([2040.000000,2312.000000]).bindPopup("Fish"),
	L.marker([1648.000000,2192.000000]).bindPopup("Fish"),
	L.marker([984.000000,368.000000]).bindPopup("<b>Big Fish</b>"),
	L.marker([736.000000,240.000000]).bindPopup("<B>Big Fish</b>")])
layerControl.addOverlay(fish, "Fish");

var things = L.layerGroup([
	L.circle([2816.000000,1416.000000],{radius: 300.000000,color: 'blue',fillColor: '#30f',fillOpacity: 0.200000}).bindPopup("Fountain"),
	L.marker([1512.000000,2424.000000]).bindPopup("<h2>BOAT!</h2>"),
	L.marker([2303.000000,133.000000]).bindPopup("Bird Nest"),
	L.polygon([[992.000000,2872.000000],[1700.000000,2872.000000],[1700.000000,3752.000000],[992.000000,3752.000000]],{color: 'red'}).bindPopup("<i>Camp!</i>")])
layerControl.addOverlay(things, "Things");


map.addLayer(things); // we want the 'things' to be checked by default
