// Keep a global variable of world and cities data, to not repeat parsing
let oldWorld = undefined;
let oldCities = undefined;

const buildCitiesForm = (cities) => {
    const div = document.getElementById("tc_left");

    cities.forEach(city => {
        const form = document.createElement("form");
        form.id = city.city;
        form.className = 'left-form';
        form.innerHTML = `
            <h2>${city.city}</h2>
            <div class="form-control">
                <label for="${city.city}_ratio">Requests Ratio</label>
                <input class="requests-total" type="range" min="0" max="100" id="${city.city}_ratio" />
            </div>
            <div class="form-control">
                <label for="${city.city}_device">Device (%)</label>
                <input class="device-mobile" type="range" min="0" max="100" id="${city.city}_device_mobile" />
                <input class="device-web" type="range" min="0" max="100" id="${city.city}_device_web" />                
            </div>
            <div class="form-control">
                <label for="${city.city}_user">User (%)</label>
                <input class="user-new" type="range" min="0" max="100" id="${city.city}_user_new" />
                <input class="user-registered" type="range" min="0" max="100" id="${city.city}_user_registered" />                
            </div>
            <div class="form-control">
                <label for="${city.city}_user">Travel Type (%)</label>
                <input class="travel-t1" type="range" min="0" max="100" id="${city.city}_travel_t1" />
                <input class="travel-t2" type="range" min="0" max="100" id="${city.city}_travel_t2" />
                <input class="travel-t3" type="range" min="0" max="100" id="${city.city}_travel_t3" />                            
            </div>
        `;
        div.append(form);
    });

};

const drawMap = (world, cities) => {
    const div = document.getElementById("tc_map");

    const width = div.offsetWidth - 10,    // Adjust the margin
        height = div.offsetHeight - 10;  // Adjust the margin

    const svg = d3.select("#tc_map")
        .append("svg")
        .attr("id", "svg_map")
        .attr("viewBox", [0, 0, width, height]);

    // Define map projection
    const projection = d3.geoMercator()
        .center([13, 52])                    // comment centrer la carte, longitude, latitude
        .translate([width / 2, height / 2])  // centrer l'image obtenue dans le svg
        .scale([width / 2.5]);                     // zoom, plus la valeur est petit plus le zoom est gros

    // Define path generator
    const path = d3.geoPath()
        .projection(projection);

    const countries = topojson.feature(world, world.objects.countries);
    const filtered = cities.map(c => c.country);
    svg.selectAll("path")
        .data(countries.features)
        .enter()
        .append("path")
        .attr("d", path)
        .attr("stroke", "var(--map-stroke-color)")
        .attr("fill", (d) => {
            if (filtered.includes(d.properties.name)) {
                return "var(--map-fill-selected-color)";
            }
            return "var(--map-fill-color)";
        });
};

const drawChartsTotal = (cities) => {
    const div = document.getElementById("tc_charts_total");
    const width = div.offsetWidth - 10;
    const height = div.offsetHeight - 10;

    const margin = ({top: 30, right: 30, bottom: 10, left: 50})

    const svg = d3.select("#tc_charts_total")
        .append("svg")
        .attr("id", "svg_charts_total")
        .attr("viewBox", [0, 0, width, height]);

    const y = d3.scaleBand()
        .domain(d3.range(cities.length))
        .rangeRound([margin.top, height - margin.bottom])
        .padding(0.1);

    const yAxis = g => g
        .attr("transform", `translate(${margin.left},0)`)
        .call(d3.axisLeft(y).tickFormat(i => cities[i].city).tickSizeOuter(0));

    const x = d3.scaleLinear()
        .domain([0, d3.max(cities, city => city.status.requests.total)])
        .range([margin.left, width - margin.right]);

    const xAxis = g => g
        .attr("transform", `translate(0,${margin.top})`)
        .call(d3.axisTop(x).ticks(width / 80))
        .call(g => g.select(".domain").remove())
        .call(g => g.select(".tick:first-of-type text").clone()
            .attr("y", -18)
            .attr("x", -3)
            .attr("text-anchor", "start")
            .attr("font-size", 16)
            .text("Total Requests per City"));

    const format = x.tickFormat(20);

    svg.append("g")
        .attr("fill", "var(--total-color)")
        .selectAll("rect")
        .data(cities)
        .join("rect")
        .attr("x", x(0))
        .attr("y", (d, i) => y(i))
        .attr("width", d => x(d.status.requests.total) - x(0))
        .attr("height", y.bandwidth());

    svg.append("g")
        .attr("fill", "white")
        .attr("text-anchor", "end")
        .attr("font-family", "sans-serif")
        .attr("font-size", 12)
        .selectAll("text")
        .data(cities)
        .join("text")
        .attr("x", d => x(d.status.requests.total))
        .attr("y", (d, i) => y(i) + y.bandwidth() / 2)
        .attr("dy", "0.35em")
        .attr("dx", -4)
        .text(d => format(d.status.requests.total))
        .call(text => text.filter(d => x(d.status.requests.total) - x(0) < 20) // short bars
            .attr("dx", +4)
            .attr("fill", "black")
            .attr("text-anchor", "start"));

    svg.append("g")
        .call(xAxis);

    svg.append("g")
        .call(yAxis);

};

const drawChartsDetails = (cities, divId, svgId, keys, colors, dataMap, title) => {
    const div = document.getElementById(divId);

    const width = div.offsetWidth - 10,
        height = div.offsetHeight - 10;

    const svg = d3.select("#" + divId)
        .append("svg")
        .attr("id", svgId)
        .attr("viewBox", [0, 0, width, height]);

    const margin = {top: 10, right: 10, bottom: 20, left: 40};

    const color = d3.scaleOrdinal()
        .range(colors);

    const x0 = d3.scaleBand()
        .domain(cities.map(d => d.city))
        .rangeRound([margin.left, width - margin.right])
        .paddingInner(0.1);

    const x1 = d3.scaleBand()
        .domain(keys)
        .rangeRound([0, x0.bandwidth()])
        .padding(0.05)

    const xAxis = g => g
        .attr("transform", `translate(0,${height - margin.bottom})`)
        .call(d3.axisBottom(x0).tickSizeOuter(0))
        .call(g => g.select(".domain").remove())

    // Note that we know that d.status.requests.total will represent max request per city
    const y = d3.scaleLinear()
        .domain([0, d3.max(cities, d => d.status.requests.total)]).nice()
        .rangeRound([height - margin.bottom, margin.top])

    const yAxis = g => g
        .attr("transform", `translate(${margin.left},0)`)
        .call(d3.axisLeft(y).ticks(null, "s"))
        .call(g => g.select(".domain").remove())
        .call(g => g.select(".tick:last-of-type text").clone()
            .attr("x", 3)
            .attr("text-anchor", "start")
            .attr("font-size", 16)
            .text(title));

    const legend = svg => {
        const g = svg
            .attr("transform", `translate(${width},0)`)
            .attr("text-anchor", "end")
            .attr("font-family", "sans-serif")
            .attr("font-size", 16)
            .selectAll("g")
            .data(color.domain().slice())
            .join("g")
            .attr("transform", (d, i) => `translate(0,${i * 20})`);

        g.append("rect")
            .attr("x", -19)
            .attr("width", 19)
            .attr("height", 19)
            .attr("fill", color);

        g.append("text")
            .attr("x", -24)
            .attr("y", 9.5)
            .attr("dy", "0.35em")
            .text(d => d);
    };

    svg.append("g")
        .call(xAxis);

    svg.append("g")
        .call(yAxis);

    svg.append("g")
        .selectAll("g")
        .data(cities)
        .join("g")
        .attr("transform", d => `translate(${x0(d.city)},0)`)
        .selectAll("rect")
        .data(dataMap)
        .join("rect")
        .attr("x", d => x1(d.key))
        .attr("y", d => y(d.value))
        .attr("width", x1.bandwidth())
        .attr("height", d => y(0) - y(d.value))
        .attr("fill", d => color(d.key));

    svg.append("g")
        .call(legend);

};

const drawChartsDevices = (cities) => {
    const keys = ["mobile", "web"];
    const colors = ["var(--device-mobile-color)", "var(--device-web-color)"];
    const dataMap = d => keys.map(key => {
        let value = 0;
        switch (key) {
            case "web":
                value = d.status.requests.devices.web;
                break;
            case "mobile":
                value = d.status.requests.devices.mobile;
                break;
            default:
            // Nothing to do in default
        }
        return {
            key,
            value: value,
        }
    });
    drawChartsDetails(
        cities,
        "tc_charts_devices",
        "svg_charts_devices",
        keys,
        colors,
        dataMap,
        "Requests per Device"
    );
}

const drawChartsUsers = (cities) => {
    const keys = ["new", "registered"];
    const colors = ["var(--user-new-color)", "var(--user-registered-color)"];
    const dataMap = d => keys.map(key => {
        let value = 0;
        switch (key) {
            case "registered":
                value = d.status.requests.users.registered;
                break;
            case "new":
                value = d.status.requests.users.new;
                break;
            default:
            // Nothing to do in default
        }
        return {
            key,
            value: value,
        }
    });
    drawChartsDetails(
        cities,
        "tc_charts_users",
        "svg_charts_users",
        keys,
        colors,
        dataMap,
        "Requests per User Type"
    );
};

const drawChartsTypes = (cities) => {
    const keys = ["t1", "t2", "t3"];
    const colors = ["var(--travel-t1-color)", "var(--travel-t2-color)", "var(--travel-t3-color)"];
    const dataMap = d => keys.map(key => {
        let value = 0;
        switch (key) {
            case "t1":
                value = d.status.requests.travel_type.t1;
                break;
            case "t2":
                value = d.status.requests.travel_type.t2;
                break;
            case "t3":
                value = d.status.requests.travel_type.t3;
                break;
            default:
            // Nothing to do in default
        }
        return {
            key,
            value: value,
        }
    });
    drawChartsDetails(
        cities,
        "tc_charts_travel_types",
        "svg_charts_travel_types",
        keys,
        colors,
        dataMap,
        "Requests per Travel Type"
    );
}

// From https://github.com/topojson/world-atlas
// Load json data, note that d3 v5 uses promises, so function callbaks are not supported
Promise.all([
    d3.json("data/countries-110m.json"),
    d3.json("data/example_status.json")
])
    .then(([world, cities]) => {
        oldWorld = world;
        oldCities = cities;

        buildCitiesForm(cities);

        drawMap(oldWorld, oldCities);
        drawChartsTotal(oldCities);
        drawChartsDevices(oldCities);
        drawChartsUsers(oldCities);
        drawChartsTypes(oldCities);
    });

window.onresize = () => {
    if (oldWorld !== undefined && oldCities !== undefined) {

        d3.select("#svg_map").remove();
        d3.select("#svg_charts_total").remove();
        d3.select("#svg_charts_devices").remove();
        d3.select("#svg_charts_users").remove();
        d3.select("#svg_charts_travel_types").remove();

        drawMap(oldWorld, oldCities);
        drawChartsTotal(oldCities);
        drawChartsDevices(oldCities);
        drawChartsUsers(oldCities);
        drawChartsTypes(oldCities);
    }
};

