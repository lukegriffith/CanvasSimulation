const canvas = document.getElementById('simulationCanvas');
const ctx = canvas.getContext('2d');

// Function to update the HTML table with team counts
function updateTeamTable(teamCounts) {
    const tableBody = document.getElementById('teamTableBody');
    tableBody.innerHTML = ''; // Clear the table before updating

    for (const [teamID, count] of Object.entries(teamCounts)) {
        // Create a new row for each team
        const row = document.createElement('tr');

        // Create the team cell
        const teamCell = document.createElement('td');
        teamCell.textContent = `Team ${teamID}`;

        // Create the count cell
        const countCell = document.createElement('td');
        countCell.textContent = count;

        // Append the cells to the row
        row.appendChild(teamCell);
        row.appendChild(countCell);

        // Append the row to the table body
        tableBody.appendChild(row);
    }
}


function resizeCanvas(socket) {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;

    // Optional: Clear the canvas or redraw content after resizing
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    // You can add code here to redraw content if needed
    data = {
     Type: 'resize',
     Width: canvas.width,
     Height: canvas.height,
    }
    console.log(data)
    socket.send(JSON.stringify(data))
}

// Connect to the WebSocket server
const socket = new WebSocket('ws://localhost:8080/ws');

socket.onopen = () => {
    console.log('WebSocket connection established');
    resizeCanvas(socket)
};

socket.onerror = (error) => {
    console.error('WebSocket error:', error);
};

socket.onclose = () => {
    console.log('WebSocket connection closed');
};
// Initial resize to fill the screen

// Event listener to resize the canvas when the window is resized
window.addEventListener('resize', function() { resizeCanvas(socket) });

// Handle incoming messages from the WebSocket
socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    updateCanvas(data);
};

canvas.addEventListener('click', (event) => {
    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    console.log('Canvas clicked at:', x, y);
    // Optionally, send this data to the backend
    // Example: Sending a message to the server
    socket.send(JSON.stringify({ type: 'click', x: x, y: y }));
});



// Example: Draw a background color to test
ctx.fillStyle = 'lightblue';
ctx.fillRect(0, 0, canvas.width, canvas.height);


function updateCanvas(data) {
    //console.log(data);
    // Clear the canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    const activeColor = '#000000';
    const inactiveColor = '#D3D3D3';
    const invulnColor = '#0000FF';

    // Separate active and inactive entities
    const inactiveEntities = [];
    const activeEntities = [];

    data.Entities.forEach((entity) => {
        if (entity.Active) {
            activeEntities.push(entity);
        } else {
            inactiveEntities.push(entity);
        }
    });

    // Draw inactive entities first
    inactiveEntities.forEach((entity) => {
        ctx.beginPath(); // Start a new path
        ctx.arc(entity.X, entity.Y, entity.Width, 0, Math.PI * 2); // Create a circle
        ctx.fillStyle = getTeamColor(entity.TeamID, data.TeamCount, false, true);
        ctx.fill(); // Fill the circle
    });

    // Separate active and inactive food items
    let activeFood = [];
    let inactiveFood = [];

    if (data.Foods != null) { 
        data.Foods.forEach((food) => {
            if (food.Active) {
                activeFood.push(food);
            } else {
                inactiveFood.push(food);
            }
        });

        // Draw inactive food items as gray diamonds
        inactiveFood.forEach((food) => {
            ctx.beginPath(); // Start a new path
            ctx.moveTo(food.X, food.Y - food.Size); // Move to the top point of the diamond
            ctx.lineTo(food.X + food.Size, food.Y); // Draw to the right point
            ctx.lineTo(food.X, food.Y + food.Size); // Draw to the bottom point
            ctx.lineTo(food.X - food.Size, food.Y); // Draw to the left point
            ctx.closePath(); // Close the path to form a diamond
            ctx.fillStyle = "gray"; // Set the fill color for inactive food
            ctx.fill(); // Fill the diamond
        });

        // Draw active food items as gold diamonds
        activeFood.forEach((food) => {
            ctx.beginPath(); // Start a new path
            ctx.moveTo(food.X, food.Y - food.Size); // Move to the top point of the diamond
            ctx.lineTo(food.X + food.Size, food.Y); // Draw to the right point
            ctx.lineTo(food.X, food.Y + food.Size); // Draw to the bottom point
            ctx.lineTo(food.X - food.Size, food.Y); // Draw to the left point
            ctx.closePath(); // Close the path to form a diamond
            ctx.fillStyle = "lightgreen"; // Set the fill color for active food
            ctx.fill(); // Fill the diamond
        });
    }
    // Then draw active entities
    activeEntities.forEach((entity) => {
        ctx.beginPath(); // Start a new path
        ctx.arc(entity.X, entity.Y, entity.Width, 0, Math.PI * 2); // Create a circle
        teamColor = getTeamColor(entity.TeamID, data.TeamCount);
        ctx.fillStyle = teamColor; // Set the fill color for active
        if (entity.Invulnerable) {
            
            ctx.fillStyle = getTeamColor(entity.TeamID, data.TeamCount, true);
        }

        ctx.fill(); // Fill the circle
    });

    // Assuming 'activeEntities' is an array of objects, each with a 'teamID' property
    activeCount = activeEntities.length;

    // Calculate the count of active entities for each team
    const teamCounts = {};
    activeEntities.forEach(entity => {
        const teamID = entity.TeamID;
        if (teamCounts[teamID]) {
            teamCounts[teamID]++;
        } else {
            teamCounts[teamID] = 1;
        }
    });

    // Update the HTML table
    updateTeamTable(teamCounts);

}

function getTeamColor(teamID, totalTeams, isInvulnerable=false, isInactive=false) {
    // Scale the hue based on the team ID
    let hue = (360 / totalTeams) * teamID; // Evenly distribute hues across 360 degrees

    // Set saturation and lightness for regular and invulnerable states
    let saturation = isInvulnerable ? 65 : 70; // Lower saturation to make it closer to gray when invulnerable
    let lightness = isInvulnerable ? 60 : 50; // Slightly higher lightness for a muted look when invulnerable

    saturation = isInactive ? 10 : saturation
    lightness = isInactive ? 30 : lightness

    // Convert HSL to hex and return the color
    return hslToHex(hue, saturation, lightness);
}

function hslToHex(h, s, l) {
    s /= 100;
    l /= 100;

    let c = (1 - Math.abs(2 * l - 1)) * s;
    let x = c * (1 - Math.abs((h / 60) % 2 - 1));
    let m = l - c / 2;

    let r = 0, g = 0, b = 0;
    if (0 <= h && h < 60) { r = c; g = x; b = 0; }
    else if (60 <= h && h < 120) { r = x; g = c; b = 0; }
    else if (120 <= h && h < 180) { r = 0; g = c; b = x; }
    else if (180 <= h && h < 240) { r = 0; g = x; b = c; }
    else if (240 <= h && h < 300) { r = x; g = 0; b = c; }
    else if (300 <= h && h < 360) { r = c; g = 0; b = x; }

    r = Math.round((r + m) * 255).toString(16).padStart(2, '0');
    g = Math.round((g + m) * 255).toString(16).padStart(2, '0');
    b = Math.round((b + m) * 255).toString(16).padStart(2, '0');

    return `#${r}${g}${b}`;
}


// Show the form
document.getElementById('showFormButton').addEventListener('click', () => {
    document.getElementById('formModal').style.display = 'block';
});

// Hide the form
function hideForm() {
    document.getElementById('formModal').style.display = 'none';
}

// Save simulation settings
function saveSimulationSettings() {
    const population = document.getElementById('population').value;
    const teamCount = document.getElementById('teamCount').value;
    const foodCount = document.getElementById('foodCount').value;
    const minSize = document.getElementById('MinSize').value;
    const startMaxSize = document.getElementById('StartMaxSize').value;
    const maxSize = document.getElementById('MaxSize').value;
    const baseSpeed = document.getElementById('BaseSpeed').value;

    // Example of handling the settings
    //
    data = {
        Type: 'settings',
        TeamCount: Number(teamCount),
        Population: Number(population),
        FoodCount: Number(foodCount),
        MinSize: Number(minSize),
        StartMaxSize: Number(startMaxSize),
        MaxSize: Number(maxSize),
        baseSpeed: Number(baseSpeed),
    }
    socket.send(JSON.stringify(data))


    // Add your simulation control logic here

    // Hide the form after saving
    hideForm();
}
